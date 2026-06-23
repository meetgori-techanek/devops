# Runbook: Jenkins Upgrade

**Date:** 2026-06-23
**Scope:** Ubuntu 24.04.4 LTS, Java 25 LTS, Jenkins 2.555.3, Plugin upgrades, Version pinning

---

> ⚠️ **Before you start:** Take a full server snapshot (EC2 AMI or equivalent). Back up `/var/lib/jenkins` and `/etc/default/jenkins`.

```bash
# Backup Jenkins home
sudo tar -czvf /tmp/jenkins-home-backup-$(date +%Y%m%d).tar.gz /var/lib/jenkins
```

---

## Step 1: Stop Jenkins

```bash
sudo systemctl stop jenkins
sudo systemctl status jenkins   # confirm: inactive (dead)
sudo systemctl disable jenkins
```

---

## Step 2: Upgrade Ubuntu to 24.04.4 LTS

### If you are currently on Ubuntu 22.04

```bash
sudo apt update
sudo apt dist-upgrade -y
sudo reboot
```

After reboot:

```bash
sudo do-release-upgrade
```

**Prompts to expect:**

| Prompt | Action |
|---|---|
| Restart services without asking? | Yes |
| Pink box - config file maintainer version | Choose "install the package maintainer's version" |
| sshd_config modified - what to do? | Choose "install the package maintainer's version" |
| Additional SSH daemon on port 1022? | `y` (continue) |

After the upgrade completes and reboots, verify:

```bash
lsb_release -a
# Expected: Ubuntu 24.04.x LTS
```

### If you are already on Ubuntu 24.04.x (getting to 24.04.4)

```bash
sudo apt update
sudo apt dist-upgrade -y
sudo reboot
```

After reboot, verify the point release:

```bash
lsb_release -a
# Expected: Ubuntu 24.04.4 LTS
```

---

## Step 3: Install Java 25 LTS

Ubuntu 24.04 does not ship OpenJDK 25 in its default repos. Use the `openjdk-r` PPA.

```bash
sudo apt update
sudo apt install -y software-properties-common
sudo add-apt-repository ppa:openjdk-r/ppa -y
sudo apt update
sudo apt install -y openjdk-25-jdk
```

Verify installation:

```bash
java -version
# Expected: openjdk version "25" ...

javac -version
# Expected: javac 25
```

Set Java 25 as the system default (if multiple Java versions are installed):

```bash
sudo update-alternatives --config java
# Select the entry for java-25-openjdk-amd64
```

Set `JAVA_HOME` for Jenkins:

```bash
# Find the Java 25 path
update-java-alternatives -l | grep java-25

# Edit Jenkins defaults
sudo nano /etc/default/jenkins
```

Add or update this line:

```bash
JAVA_HOME=/usr/lib/jvm/java-25-openjdk-amd64
```

---

## Step 4: Upgrade Jenkins to 2.555.3

Add the Jenkins stable apt repo and install the specific version.

```bash
sudo apt update

# Add Jenkins signing key
wget -q -O - https://pkg.jenkins.io/debian-stable/jenkins.io-2023.key \
  | sudo tee /usr/share/keyrings/jenkins-keyring.asc > /dev/null

# Add Jenkins stable repo
echo "deb [signed-by=/usr/share/keyrings/jenkins-keyring.asc] \
  https://pkg.jenkins.io/debian-stable binary/" \
  | sudo tee /etc/apt/sources.list.d/jenkins.list > /dev/null

sudo apt update

# Install the specific version
sudo apt install jenkins=2.555.3
```

> ⚠️ If `2.555.3` is not available in the stable repo, install it from the weekly repo (`https://pkg.jenkins.io/debian`) using the same key steps but the weekly repo URL.

Start Jenkins and check logs:

```bash
sudo systemctl daemon-reload
sudo systemctl start jenkins
sudo systemctl disable jenkins
sudo systemctl status jenkins

# Watch logs for startup errors
sudo journalctl -u jenkins -f
```

Verify version in browser: `http://<server-ip>:8080` or via CLI:

```bash
java -jar /usr/share/jenkins/jenkins.war --version 2>/dev/null || \
  curl -s http://localhost:8080/api/json?tree=version | python3 -m json.tool
```

---

## Step 5: Upgrade Plugins

**Option A: Via Jenkins UI (recommended)**

1. Go to `Manage Jenkins` > `Plugins` > `Updates`
2. Click `Select All`
3. Click `Download now and install after restart`
4. Check `Restart Jenkins when installation is complete`

**Option B: Via Jenkins CLI**

```bash
# Download Jenkins CLI jar
wget http://localhost:8080/jnlpJars/jenkins-cli.jar

# List plugins with updates
java -jar jenkins-cli.jar -s http://localhost:8080/ \
  -auth admin:<api-token> list-plugins | grep -E '\(.*\)' | awk '{print $1}'

# Install all updates
java -jar jenkins-cli.jar -s http://localhost:8080/ \
  -auth admin:<api-token> install-plugin \
  $(java -jar jenkins-cli.jar -s http://localhost:8080/ \
    -auth admin:<api-token> list-plugins | grep -E '\(.*\)' | awk '{print $1}') \
  -restart
```

After the restart, verify all plugins loaded without errors:

`Manage Jenkins` > `Manage Old Data` and `System Log` - no red errors.

---

## Step 6: Pin Versions (Prevent Auto-Upgrades)

### Pin Java 25

```bash
sudo apt-mark hold openjdk-25-jdk openjdk-25-jre openjdk-25-jre-headless
```

Verify hold:

```bash
apt-mark showhold
# Should list all three openjdk-25 packages
```

### Pin Jenkins

```bash
sudo apt-mark hold jenkins
```

### Prevent Ubuntu Unattended Upgrades

This stops the OS from auto-applying security/kernel patches without your control.

```bash
sudo nano /etc/apt/apt.conf.d/20auto-upgrades
```

Set both values to `"0"`:

```
APT::Periodic::Update-Package-Lists "0";
APT::Periodic::Unattended-Upgrade "0";
```

> ⚠️ Disabling unattended upgrades means you are responsible for applying security patches manually on a schedule. Recommended: run `sudo apt update && sudo apt upgrade` during your maintenance windows.

Disable the timer completely if preferred:

```bash
sudo systemctl disable apt-daily.timer
sudo systemctl disable apt-daily-upgrade.timer
sudo systemctl stop apt-daily.timer
sudo systemctl stop apt-daily-upgrade.timer
```

---

## Step 7: Post-Upgrade Verification

```bash
# OS
lsb_release -a

# Java
java -version

# Jenkins service
sudo systemctl status jenkins

# Jenkins version
curl -s http://localhost:8080/api/json?tree=version

# Held packages
apt-mark showhold

# Unattended upgrade status
systemctl is-enabled apt-daily.timer
systemctl is-enabled apt-daily-upgrade.timer
```

---

## Rollback Plan

If Jenkins fails to start after the upgrade:

```bash
# Stop Jenkins
sudo systemctl stop jenkins

# Restore Jenkins home from backup
sudo tar -xzvf /tmp/jenkins-home-backup-<date>.tar.gz -C /

# Downgrade Jenkins (example to previous version)
sudo apt-get install jenkins=<previous-version>

# Restore previous Java (if needed)
sudo update-alternatives --config java

sudo systemctl start jenkins
```

---

## References

- Jenkins 2.555.3 changelog: https://www.jenkins.io/changelog/2.555.3/
- Jenkins Java support matrix: https://www.jenkins.io/doc/administration/requirements/java/
- Ubuntu release upgrade guide: https://ubuntu.com/server/docs/upgrade-introduction
