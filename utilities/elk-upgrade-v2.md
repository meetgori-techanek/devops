# ELK Stack Upgrade Plan: 8.18.3 to 9.5.1 (via 8.19.x and Ubuntu 24.04)

Test environment: AMI-restored clone of central server (172-30-0-157), hostname set to elk-upgrade-test
Starting OS: Ubuntu 18.04.6 LTS
Order executed: Phase 1 (8.18.3 to 8.19.x) → Phase 2 (OS 18.04 to 24.04) → Phase 3 (8.19.x to 9.5.1)

## Current State (as audited)

Central server (172-30-0-157): Elasticsearch 8.18.3, Kibana 8.18.3, Filebeat 8.18.3, APM Server 8.18.3

App servers:

| Server | Filebeat | Metricbeat |
|---|---|---|
| PROD-A | 8.18.1 | 8.18.1 |
| DEV-A | 8.18.3 | 8.18.3 |
| STAGE | 8.18.1 | 8.18.1 |
| PROD-B | 8.18.1 | 8.18.1 |
| PROD-C | 8.18.1 | 8.18.1 |
| PROD-D | 8.18.1 | 8.18.1 |
| dev-process-new | 7.17.29 | 7.17.29 |
| prod-process-new | 7.17.29 | 7.17.29 |

dev-process-new and prod-process-new are on 7.17.29, a major version behind. These two need the two-step 7.17→8.19 path before or alongside Phase 1.

---

# Phase 1: Elasticsearch 8.18.3 to 8.19.x

## 1. Pre-Upgrade Checks

```bash
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:**********
curl -X GET "localhost:9200/" -u elastic:**********
curl -X GET "localhost:9200/_cat/nodes?v" -u elastic:**********
curl -X GET "localhost:9200/_all/_mapping/field/_doc_count?pretty" -u elastic:**********
```

Fix cluster health if not green (yellow from single-node unassigned replicas is expected and not blocking).

**Launch the Kibana Upgrade Assistant** (on the current 8.18.3 Kibana, before touching any packages):
1. Log into Kibana in the browser: `http://<server-ip>:5601`
2. Open the hamburger menu (top left) > scroll to **Management** > **Stack Management**
3. In the left sidebar under **Kibana**, click **Upgrade Assistant**
   Direct path: Stack Management > Upgrade Assistant
4. Two tabs appear: **Elasticsearch** and **Kibana**. Check both.
5. Anything marked **Critical** must be resolved before upgrading, click into each issue, it will offer **Reindex** or **Mark as read-only** for old indices, or show manual config changes needed for deprecated settings.
6. **Warning**-level issues don't block the upgrade but are worth reviewing and fixing where the resolution is "Automated".

Resolve everything Critical here, on 8.18.3, before starting Phase 1. This is the same check to repeat again after Phase 1 (now on 8.19.x Kibana) before starting Phase 3.

Review 8.19 breaking changes: https://www.elastic.co/guide/en/elasticsearch/reference/8.19/migrating-8.19.html

## 2. Prepare Cluster

```bash
curl -X POST "localhost:9200/_flush?pretty" -u elastic:**********
curl -X POST "localhost:9200/_ml/set_upgrade_mode?enabled=true&pretty" -u elastic:**********
```

## 3. Stop Services

```bash
sudo systemctl stop kibana apm-server filebeat elasticsearch
```

## 4. Upgrade Elasticsearch

```bash
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-8.x.list

sudo apt-get update
sudo apt-get install elasticsearch=8.19.*

sudo systemctl start elasticsearch
curl -X GET "localhost:9200/" -u elastic:**********
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:**********
```

**Config prompt encountered:** `jvm.options` — new default enables auto heap sizing. Chose to keep the new default (accepted the auto-sizing), did not create a manual override, since the previous fixed heap (2g/6g) is smaller than what auto-sizing gives on this instance's RAM.

## 5. Upgrade Kibana

```bash
sudo apt-get install kibana=8.19.*
sudo systemctl start kibana
```

Log in, confirm version in Stack Management > About. Fleet-related log warnings (`failed to get message signing key pair`, `Failed to decrypt attribute "passphrase"`) are expected noise since Fleet isn't in use, not a failure indicator, confirmed via `curl localhost:5601/api/status` showing `overall: available`.

## 6. Upgrade APM Server

```bash
sudo apt-get install apm-server=8.19.*
sudo systemctl start apm-server
```

**Config prompt encountered:** `apm-server.yml` — kept current file (`N`). The package default would reset `host` to `127.0.0.1:8200` (blocking remote app servers) and comment out the `secret_token` (removing auth). Keeping the existing file preserved both.

## 7. Upgrade Filebeat (Central Node)

```bash
sudo apt-get install filebeat=8.19.*
sudo systemctl start filebeat
```

**Config prompt encountered:** `modules.d/aws.yml.disabled` — took the package default (`Y`), module is unused and disabled either way.

## 8. Disable ML Upgrade Mode

```bash
curl -X POST "localhost:9200/_ml/set_upgrade_mode?enabled=false&pretty" -u elastic:**********
```

## 9. Upgrade Ingest Agents on App Servers

PROD-A, PROD-B, PROD-C, PROD-D, STAGE, DEV-A (already 8.18.x, straight bump):

```bash
sudo systemctl stop filebeat metricbeat

wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-8.x.list

sudo apt-get update
sudo apt-get install filebeat=8.19.* metricbeat=8.19.*

sudo systemctl start filebeat metricbeat
```

dev-process-new and prod-process-new (7.17.29, do not skip straight to 8.19):

```bash
sudo systemctl stop filebeat metricbeat

# Step A: latest 7.17.x first
sudo apt-get update
sudo apt-get install filebeat=7.17.* metricbeat=7.17.*

# Step B: switch to 8.x repo and upgrade
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-8.x.list
sudo rm -f /etc/apt/sources.list.d/elastic-7.x.list

sudo apt-get update
sudo apt-get install filebeat=8.19.* metricbeat=8.19.*

sudo systemctl start filebeat metricbeat
```

Back up `filebeat.yml` and `metricbeat.yml` on these two before the jump, config syntax changed enough between 7.17 and 8.x to warrant a manual diff.

## 10. Validation

```bash
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:**********
curl -X GET "localhost:9200/_cat/nodes?v&h=name,node.role,master,ip,version" -u elastic:**********
curl -X GET "localhost:5601/api/status" -u elastic:**********
```

Confirmed on test run: all four components on 8.19.20, cluster health yellow (single-node replica shards, expected, not a fault, all primaries active), Kibana available, saved objects migrated cleanly, dashboards and Discover render existing data. A pre-existing mapping conflict on `azure.activitylogs.properties` (two indices from 2025.05.14, 7.11/7.12 era) was investigated and confirmed unrelated to this upgrade, safe to ignore or clean up separately.

---

# Phase 2: Ubuntu OS Upgrade (18.04.6 LTS to 24.04 LTS)

Elasticsearch 9.5.1 requires Ubuntu 20.04+, and 18.04 is past its standard support window. Done in three hops since `do-release-upgrade` cannot skip LTS releases: 18.04 → 20.04 → 22.04 → 24.04. Run this after Phase 1 and before Phase 3.

## 1. Stop ELK Services and Disable Auto-Start

```bash
sudo systemctl stop kibana apm-server filebeat elasticsearch
sudo systemctl disable kibana apm-server filebeat elasticsearch
```

## 2. Snapshot the AMI

Take a fresh AMI snapshot before touching the OS.

## 3. Fix Slow/Broken Mirror (if hit)

The regional EC2 mirror (`us-east-1.ec2.archive.ubuntu.com`) returned repeated `503` errors and near-zero throughput during testing. If `apt update && apt upgrade` stalls or throws `503 Service Unavailable` / connection failures, switch to the global mirror before continuing:

```bash
sudo sed -i 's|http://us-east-1.ec2.archive.ubuntu.com/ubuntu|http://archive.ubuntu.com/ubuntu|g' /etc/apt/sources.list
sudo apt update
```

Safe to `Ctrl+C` and retry at this stage, since this is plain `apt upgrade`, not `do-release-upgrade` yet.

## 4. Hop 1: 18.04 to 20.04

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install update-manager-core -y
sudo do-release-upgrade
```

If `do-release-upgrade` refuses with "Please install all available updates for your release before upgrading" even after `apt upgrade -y`, the remaining holdbacks are usually Ubuntu Pro/ESM-only packages that can't be upgraded without a Pro subscription, expected on unregistered AWS instances. Force the remaining dependency-driven updates through instead:

```bash
sudo apt-get dist-upgrade -y
sudo do-release-upgrade
```

Reboot when prompted. Confirm:

```bash
lsb_release -a
```

## 5. Hop 2: 20.04 to 22.04

```bash
sudo apt update && sudo apt upgrade -y
sudo do-release-upgrade
```

Reboot when prompted. Confirm 22.04 with `lsb_release -a`.

## 6. Hop 3: 22.04 to 24.04

```bash
sudo apt update && sudo apt upgrade -y
sudo do-release-upgrade
```

Reboot when prompted. Confirm 24.04 with `lsb_release -a`.

## 7. Re-enable and Start ELK Services

```bash
sudo systemctl enable kibana apm-server filebeat elasticsearch
sudo systemctl start elasticsearch
sudo systemctl start kibana apm-server filebeat
```

## 8. Validate Packages Survived the OS Hops

```bash
dpkg -l | grep -E 'elasticsearch|kibana|filebeat|apm-server'
curl -X GET "localhost:9200/" -u elastic:**********
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:**********
cat /etc/apt/sources.list.d/elastic-8.x.list
```

Confirmed on test run: all four components remained 8.19.20 across all three hops, cluster came back healthy each time. `do-release-upgrade` did not disturb the Elastic repo file in this run, but re-verify it each time regardless.

Note: repeat this OS upgrade on each app server too, since Filebeat/Metricbeat 8.19.x and 9.5.1 also require 20.04+.

---

# Phase 3: Elasticsearch 8.19.x to 9.5.1

## 1. Pre-Upgrade Checks — Run the Upgrade Assistant Fully

This step is not optional. Elasticsearch 9.x will refuse to even boot if any pre-8.0 index isn't marked read-only or reindexed, this crashed the test run with an `IllegalStateException` on `.apm-custom-link` (created in 7.10.1), which the Upgrade Assistant would have flagged in advance.

**Launch the Kibana Upgrade Assistant** (on the 8.19.x Kibana from Phase 1, before starting any 9.5.1 package installs):
1. Log into Kibana: `http://<server-ip>:5601`
2. Hamburger menu > **Stack Management** > **Upgrade Assistant** (left sidebar, under Kibana)
3. Open the **Elasticsearch** tab first:
   - Every row marked **Critical** must show "Reindex complete" or a read-only resolution before proceeding.
   - Click into any row still showing "Recommended: reindex" (not yet actioned) and run the reindex from there.
   - On the test run this list included several old `apm-*`, `metricbeat-*`, `filebeat-*`, and `.transform-notifications-*` indices from 2025.05.14 (7.11/7.12 era). All 8 critical entries needed "Reindex complete" before Elasticsearch 9.5.1 would boot.
4. Open the **Kibana** tab:
   - Confirm **Critical: 0**. Warnings (config deprecations, legacy OpenSSL provider, etc.) don't block the upgrade but note them for later cleanup.
   - Any row with an **Automated** resolution can be clicked and applied directly from here.
5. **Important gap found in testing**: `.apm-custom-link` (a 7.10.1-era index) did **not** appear in the Elasticsearch tab's list but still blocked Elasticsearch 9.5.1 from starting. Don't treat a clean Upgrade Assistant screen as a full guarantee, manually check for any other old indices the assistant might have missed:

```bash
curl -X GET "localhost:9200/_cat/indices?v" -u elastic:********** | awk '{print $3}'
```
Cross-check any suspicious old index (especially ones tied to features you barely use, like `.apm-custom-link`) against the fix below.

Manually verify no other pre-8.0 indices are missed:

```bash
curl -X GET "localhost:9200/.apm-custom-link/_settings?pretty" -u elastic:**********
```
Should show `"blocks":{"write":"true"}`. If not:

```bash
curl -X PUT "localhost:9200/.apm-custom-link/_settings" -u elastic:********** -H 'Content-Type: application/json' -d '{"index.blocks.write": true}'
```

Also remove any deprecated `elasticsearch.yml` settings no longer recognized in 9.x. On the test run, `cluster.routing.allocation.disk.watermark.enable_for_single_data_node` (a 7.x-era prerequisite setting) caused a fatal boot error and had to be removed:

```bash
sudo grep -n "enable_for_single_data_node" /etc/elasticsearch/elasticsearch.yml
sudo sed -i '/cluster.routing.allocation.disk.watermark.enable_for_single_data_node/d' /etc/elasticsearch/elasticsearch.yml
```

Review 9.x breaking changes: https://www.elastic.co/guide/en/elasticsearch/reference/9.5/migrating-9.5.html

If Enterprise Search is in use anywhere, remove it first, it's dropped in 9.0.

## 2. Stop Services

```bash
sudo systemctl stop kibana apm-server filebeat elasticsearch
```

## 3. Upgrade Elasticsearch

```bash
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/9.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-9.x.list
sudo rm -f /etc/apt/sources.list.d/elastic-8.x.list

sudo apt-get update
sudo apt-get install elasticsearch=9.5.*

sudo systemctl start elasticsearch
curl -X GET "localhost:9200/" -u elastic:**********
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:**********
```

**If Elasticsearch fails to start**, check the real error, not just the systemd summary:

```bash
sudo tail -100 /var/log/elasticsearch/elasticsearch.log
```

**If you need to roll back to 8.19.20 to fix a blocking issue** (as happened on the test run for the read-only index problem), the 8.x repo file was already removed, restore it temporarily:

```bash
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-8.x.list
sudo mv /etc/apt/sources.list.d/elastic-9.x.list /etc/apt/sources.list.d/elastic-9.x.list.disabled

sudo apt-get update
sudo apt-get install elasticsearch=8.19.20
sudo systemctl start elasticsearch
```

Fix the blocking issue, then swap back and retry:

```bash
sudo mv /etc/apt/sources.list.d/elastic-9.x.list.disabled /etc/apt/sources.list.d/elastic-9.x.list
sudo rm -f /etc/apt/sources.list.d/elastic-8.x.list
sudo apt-get update
sudo apt-get install elasticsearch=9.5.*
sudo systemctl start elasticsearch
```

## 4. Upgrade Kibana

```bash
sudo apt-get install kibana=9.5.*
sudo systemctl start kibana
```

Kibana takes noticeably longer to become available on a major version jump (saved object migrations). The same Fleet-related log warnings from Phase 1 will reappear and are expected. Give it a few minutes and confirm:

```bash
curl -X GET "localhost:5601/api/status" -u elastic:**********
```

## 5. Upgrade APM Server

```bash
sudo apt-get install apm-server=9.5.*
sudo systemctl start apm-server
```

Same config prompt as Phase 1, keep the current `apm-server.yml` (`N`).

## 6. Upgrade Filebeat (Central Node)

```bash
sudo apt-get install filebeat=9.5.*
sudo systemctl start filebeat
```

## 7. Ingest Agents on App Servers

8.19.x Filebeat/Metricbeat are fully compatible with a 9.5.1 cluster, upgrading them is optional. To upgrade anyway, repeat the same stop/repo-switch/install/start pattern from Phase 1 step 9, pointed at the 9.x repo.

## 8. Validation

```bash
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:**********
curl -X GET "localhost:9200/_cat/nodes?v&h=name,node.role,master,ip,version" -u elastic:**********
dpkg -l | grep -E 'elasticsearch|kibana|filebeat|apm-server'
```

Confirmed on test run: all four components on 9.5.1, cluster healthy at yellow (same known single-node reason), Kibana available with all plugins reporting available, APM Server active and processing, metrics/traces/logs all confirmed flowing normally in Kibana after the full upgrade.

---

# What May Break

- **Pre-8.0 indices will crash Elasticsearch 9.x on boot** if not marked read-only or reindexed first. Confirmed on this test run (`.apm-custom-link`). Always run the Upgrade Assistant fully and don't rely on package install succeeding as confirmation, the crash happens at Elasticsearch startup, after the package is already installed.
- **Deprecated `elasticsearch.yml` settings cause a fatal boot error** on 9.x, not a warning. Confirmed with `cluster.routing.allocation.disk.watermark.enable_for_single_data_node`, a leftover from the original 7.17 prerequisites doc. Review the full `elasticsearch.yml` against the 9.x breaking changes doc before attempting the install, not after it fails.
- **Regional EC2 apt mirrors can silently degrade** (503s, near-zero throughput) without warning. Confirmed on this test run. Switching to the global `archive.ubuntu.com` mirror resolved it immediately.
- **`do-release-upgrade` refuses to proceed if Ubuntu Pro/ESM-only security updates are held back**, even though they're not truly blocking. `apt-get dist-upgrade -y` clears this in practice.
- Kibana saved objects (old TSVB visualizations, Timelion) can fail to render after a major jump, check dashboards manually post-upgrade.
- Custom index templates or ILM policies referencing settings deprecated in 8.x or 9.x will need updates, Upgrade Assistant flags these.
- APM intake API and data stream naming changed between major versions, confirm APM agents in the app code still connect after the server upgrade.
- Filebeat and Metricbeat config syntax (especially module definitions and output blocks) changed meaningfully between 7.17 and 8.x, the two 7.17.29 servers need a manual config diff, not just a package swap.
- Custom Kibana plugins or dashboards built against old APIs may break, none known here but worth confirming.
- Disk space, both old and new package versions plus logs briefly coexist during `apt-get install`, confirm free space before each upgrade step.
- `do-release-upgrade` can silently disable or re-point third-party apt repos (including the Elastic repo), re-verify `/etc/apt/sources.list.d/elastic-*.list` after each OS hop even though it wasn't disturbed on this test run.
