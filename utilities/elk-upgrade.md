# ELK Stack Upgrade Plan: 8.18.3 to 9.5.1

Test environment: AMI-restored clone of central server (172-30-0-157)
OS: Ubuntu 18.04.6 LTS (needs upgrade too, see Phase 3)

## Current State

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

dev-process-new and prod-process-new are on 7.17.29, a major version behind the cluster. These two need attention first, before or alongside Phase 1.

---

# Phase 1: 8.18.3 to 8.19.x

## 1. Pre-Upgrade Checks
> Cluster health: shows overall status (green/yellow/red).\
> Root /: shows version info. \
> Cat nodes: lists nodes and resource usage. \
> Doc count mapping: checks for a known field mapping issue that can block upgrades.
```
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:<password>
curl -X GET "localhost:9200/" -u elastic:<password>
curl -X GET "localhost:9200/_cat/nodes?v" -u elastic:<password>

curl -X GET "localhost:9200/_all/_mapping/field/_doc_count?pretty" -u elastic:**********
```

Fix cluster health if not green. Run Kibana Upgrade Assistant (Stack Management > Upgrade Assistant) and resolve all critical issues.

Review 8.19 breaking changes:
https://www.elastic.co/guide/en/elasticsearch/reference/8.19/migrating-8.19.html

## 2. Prepare Cluster

Flush indices and enable ML upgrade mode:
>Flush saves in-memory data to disk.\
>ML upgrade mode pauses ML jobs during the upgrade.
```
curl -X POST "localhost:9200/_flush?pretty" -u elastic:**********
curl -X POST "localhost:9200/_ml/set_upgrade_mode?enabled=true&pretty" -u elastic:**********
```

## 3. Stop Services

```bash
sudo systemctl stop kibana
sudo systemctl stop apm-server
sudo systemctl stop filebeat
sudo systemctl stop elasticsearch
```

## 4. Upgrade Elasticsearch

```bash
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-8.x.list

sudo apt-get update
sudo apt-get install elasticsearch=8.19.*

sudo systemctl start elasticsearch
```
valiate 
```
curl -X GET "localhost:9200/" -u elastic:**********
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:**********
```

## 5. Upgrade Kibana

```bash
sudo apt-get install kibana=8.19.*
sudo systemctl start kibana
```

Log in, check version in Stack Management > About, re-check Upgrade Assistant.

## 6. Upgrade APM Server

```bash
sudo apt-get install apm-server=8.19.*
sudo systemctl start apm-server
```

## 7. Upgrade Filebeat (Central Node)

```bash
sudo apt-get install filebeat=8.19.*
sudo systemctl start filebeat
```

## 8. Disable ML Upgrade Mode

```
curl -X POST "localhost:9200/_ml/set_upgrade_mode?enabled=false&pretty" -u elastic:**********
```


## 9. Upgrade Ingest Agents on App Servers

For PROD-A, PROD-B, PROD-C, PROD-D, STAGE (currently 8.18.1, straight package bump):

```bash
sudo systemctl stop filebeat metricbeat

wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-8.x.list

sudo apt-get update
sudo apt-get install filebeat=8.19.* metricbeat=8.19.*

sudo systemctl start filebeat metricbeat
```

For DEV-A (already 8.18.3), same commands, smaller jump.

For dev-process-new and prod-process-new (7.17.29, major version behind), upgrade in two steps, do not skip straight to 8.19:

```bash
sudo systemctl stop filebeat metricbeat

# Step A: move to latest 7.17.x first if not already there
sudo apt-get update
sudo apt-get install filebeat=7.17.* metricbeat=7.17.*

# Step B: switch repo to 8.x and upgrade
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-8.x.list
sudo rm -f /etc/apt/sources.list.d/elastic-7.x.list

sudo apt-get update
sudo apt-get install filebeat=8.19.* metricbeat=8.19.*

sudo systemctl start filebeat metricbeat
```

Back up `filebeat.yml` and `metricbeat.yml` on these two servers before the jump, config syntax changed enough between 7.17 and 8.x that modules and output blocks are worth a manual diff.

## 10. Validation

```
curl -X GET "localhost:9200/_cluster/health?pretty" -u elastic:**********
curl -X GET "localhost:9200/_cat/nodes?v&h=name,node.role,master,ip,version" -u elastic:**********
curl -X GET "localhost:5601/api/status" -u elastic:**********
```

Confirm:
- Cluster status green
- All nodes report 8.19.x
- Kibana loads, dashboards render existing data
- New log events from each app server are landing
- New APM traces are landing
- `journalctl -u elasticsearch -u kibana -u apm-server -u filebeat --since "10 min ago"` clean

---

# Phase 2: 8.19.x to 9.5.1

Run only after Phase 1 is validated and stable in production for a period.

## 1. Pre-Upgrade Checks

Run Kibana Upgrade Assistant again (8.19 version this time), resolve all critical issues, including reindexing or marking read-only any indices created before 8.0.

Review 9.x breaking changes:
https://www.elastic.co/guide/en/elasticsearch/reference/9.5/migrating-9.5.html

If Enterprise Search is in use anywhere, remove it first, it is dropped in 9.0.

## 2. Stop Services

```bash
sudo systemctl stop kibana apm-server filebeat elasticsearch
```

## 3. Upgrade Elasticsearch

```bash
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/9.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-9.x.list
sudo rm -f /etc/apt/sources.list.d/elastic-8.x.list

sudo apt-get update
sudo apt-get install elasticsearch=9.5.*

sudo systemctl start elasticsearch
GET /
GET /_cluster/health?pretty
```

## 4. Upgrade Kibana

```bash
sudo apt-get install kibana=9.5.*
sudo systemctl start kibana
```

## 5. Upgrade APM Server

```bash
sudo apt-get install apm-server=9.5.*
sudo systemctl start apm-server
```

## 6. Upgrade Filebeat (Central Node)

```bash
sudo apt-get install filebeat=9.5.*
sudo systemctl start filebeat
```

## 7. Ingest Agents on App Servers

8.19.x Filebeat/Metricbeat are fully compatible with a 9.5.1 cluster, upgrading them is optional. If you want them on 9.5.1 too, repeat the same stop, repo switch, install, start pattern used in Phase 1 step 9, pointed at the 9.x repo.

## 8. Validation

Same checks as Phase 1 step 10, confirming version 9.5.x everywhere.

---

# Phase 3: Ubuntu OS Upgrade (18.04.6 LTS)

Elastic 8.19.x and 9.5.1 both require Ubuntu 20.04 or newer. 18.04 must be upgraded before or alongside Phase 1. do-release-upgrade cannot skip LTS releases, so this is done in two hops: 18.04 to 20.04, then 20.04 to 22.04.

Do this on the AMI test clone first, and consider doing it before Phase 1 rather than after, since Elasticsearch 8.19 packages may refuse to install cleanly on 18.04's older glibc.

## 1. Stop ELK Services and Disable Auto-Start

```bash
sudo systemctl stop kibana apm-server filebeat elasticsearch
sudo systemctl disable kibana apm-server filebeat elasticsearch
```

## 2. Snapshot the AMI Again

Take a fresh AMI snapshot right before the OS upgrade so this step alone is reversible.

## 3. Hop 1: 18.04 to 20.04

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install update-manager-core -y
sudo do-release-upgrade
```

Reboot when prompted. Confirm:

```bash
lsb_release -a
```

## 4. Hop 2: 20.04 to 22.04

```bash
sudo apt update && sudo apt upgrade -y
sudo do-release-upgrade
```

Reboot when prompted. Confirm:

```bash
lsb_release -a
```

## 5. Re-enable and Start ELK Services

```bash
sudo systemctl enable kibana apm-server filebeat elasticsearch
sudo systemctl start elasticsearch
sudo systemctl start kibana apm-server filebeat
```

Check all services came back clean and the Elastic packages still match the versions from Phase 1 or 2, `do-release-upgrade` can sometimes touch package repos.

```bash
dpkg -l | grep -E 'elasticsearch|kibana|filebeat|apm-server'
```

Note: repeat Phase 3 on each app server too, since Filebeat/Metricbeat 8.19.x and 9.5.1 also require 20.04+.

---

# What May Break

- Kibana saved objects (old TSVB visualizations, Timelion) can fail to render after a major jump, check dashboards manually post-upgrade.
- Custom index templates or ILM policies referencing settings deprecated in 8.x or 9.x will need updates, Upgrade Assistant flags these.
- APM intake API and data stream naming changed between major versions, confirm APM agents in the app code still connect after the server upgrade.
- Filebeat and Metricbeat config syntax (especially module definitions and output blocks) changed meaningfully between 7.17 and 8.x, the two 7.17.29 servers need a manual config diff, not just a package swap.
- Security features (TLS, authentication) are enabled by default from 8.x onward, if any custom setup relied on security being off this will need reconciling.
- Custom Kibana plugins or dashboards built against old APIs may break, no custom plugins known here but worth confirming.
- glibc and OpenJDK versions bundled with Elastic packages assume a modern OS, this is why the Ubuntu hop matters before or alongside Phase 1.
- Disk space, both old and new package versions plus logs briefly coexist during `apt-get install`, confirm free space before each upgrade step.
- `do-release-upgrade` can silently disable or re-point third-party apt repos (including the Elastic repo), re-verify `/etc/apt/sources.list.d/elastic-*.list` after each OS hop.
