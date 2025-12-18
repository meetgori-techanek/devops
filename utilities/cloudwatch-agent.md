# Amazon CloudWatch Agent – Ubuntu Installation Guide

## 1. Update System
```bash
sudo apt update
```

---

## 2. Download CloudWatch Agent (Official)
```bash
cd /tmp
wget https://s3.amazonaws.com/amazoncloudwatch-agent/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb
```

---

## 3. Install the Agent
```bash
sudo dpkg -i amazon-cloudwatch-agent.deb
```

### Fix Dependencies (if prompted)
```bash
sudo apt -f install -y
```

---

## 4. Verify Installation
```bash
ls /opt/aws/amazon-cloudwatch-agent
```

Expected output:
```
bin  etc  logs  var
```

---

## 5. Create / Update Config File
```bash
sudo nano /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.d/file_config.json
```

Paste the following configuration:
```json
{
  "metrics": {
    "metrics_collected": {
      "cpu": {
        "measurement": ["cpu_usage_idle"],
        "metrics_collection_interval": 60,
        "totalcpu": true
      },
      "mem": {
        "measurement": [
          "mem_used_percent",
          "mem_available",
          "mem_total"
        ],
        "metrics_collection_interval": 60
      },
      "disk": {
        "measurement": [
          "disk_used_percent",
          "disk_free"
        ],
        "metrics_collection_interval": 60,
        "resources": ["/"]
      }
    },
    "append_dimensions": {
      "InstanceId": "${aws:InstanceId}"
    }
  }
}
```

Save and exit.

---

## 6. Start CloudWatch Agent
```bash
sudo systemctl enable amazon-cloudwatch-agent
sudo systemctl restart amazon-cloudwatch-agent
```

---

## 7. Verify Agent Status
```bash
sudo systemctl status amazon-cloudwatch-agent
```

---

## 8. Check Agent Logs
```bash
sudo tail -f /opt/aws/amazon-cloudwatch-agent/logs/amazon-cloudwatch-agent.log
```

---

## 9. Verify in AWS Console

Navigate to:
```
CloudWatch → Metrics → CWAgent → InstanceId
```

Expected metrics:
- cpu_usage_idle  
- mem_used_percent  
- disk_used_percent  

---

## 10. Required IAM Permissions (Attach to EC2)
```json
{
  "Effect": "Allow",
  "Action": [
    "cloudwatch:PutMetricData",
    "ec2:DescribeTags",
    "logs:CreateLogGroup",
    "logs:CreateLogStream",
    "logs:PutLogEvents"
  ],
  "Resource": "*"
}
```
OR\
Attach
```
CloudWatchAgentServerPolicy
```
---

## Result

CloudWatch Agent is correctly installed on Ubuntu and collecting metrics suitable for cost analysis and future capacity planning.
