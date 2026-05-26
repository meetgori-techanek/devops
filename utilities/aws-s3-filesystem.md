# Amazon S3 Files – POC Setup Guide

---

## 1. Create S3 Bucket and Enable Versioning

Go to **S3 Console → Buckets → Create bucket**

- Block all public access: **ON**
- Bucket Versioning: **Enable**
- Default encryption: **SSE-S3 (AES-256)**

---

## 2. Create S3 File System

Go to **S3 Console → General purpose buckets → your bucket → File systems tab → Create file system**

- Select your bucket
- VPC: **Default VPC** (or your custom VPC)
- Console auto-creates mount targets in every AZ

---

## 3. Launch EC2 Instance

Go to **EC2 Console → Launch Instance**

- AMI: **Ubuntu 22.04 or Amazon Linux 2023**
- Instance type: `t3.medium` or higher
- VPC: same VPC used in Step 2
- Security group: note the SG ID (used in Step 5)

---

## 4. Create IAM Role and Attach to EC2

Go to **IAM Console → Roles → Create role**

- Trusted entity: **EC2**
- Attach policy: `AmazonS3FilesClientFullAccess`
- Role name: `ec2-s3files-role`

Then attach to your EC2 instance:

**EC2 Console → Instance → Actions → Security → Modify IAM Role → Select `ec2-s3files-role`**

---

## 5. Configure Security Group for NFS (Port 2049)

Go to **EC2 Console → Security Groups → Select the SG attached to your mount target**

Add inbound rule:

| Type | Protocol | Port | Source |
|------|----------|------|--------|
| Custom TCP | TCP | 2049 | EC2 Security Group ID |

> Source must be the security group ID of your EC2 instance, not `0.0.0.0/0`

---

## 6. Install Dependencies and Mount

SSH into your EC2 instance and run:

```bash
curl https://amazon-efs-utils.aws.com/efs-utils-installer.sh | sudo sh -s -- \
  --install-launch-wizard \
  --mount-s3files fs-0fe3ff6ef4ca0ae1b /mnt/s3/fs1
```

Replace `fs-0fe3ff6ef4ca0ae1b` with your actual File System ID from Step 2.

Verify mount:

```bash
df -h /mnt/s3/fs1
ls /mnt/s3/fs1
```

---

## 7. Persist Mount on Reboot

```bash
echo "fs-0fe3ff6ef4ca0ae1b /mnt/s3/fs1 s3files defaults,iam,_netdev 0 0" | sudo tee -a /etc/fstab
```

---

## 8. Write Heavy Data and Verify

**Write 5 x 1GB files:**

```bash
for i in {1..5}; do
  wget -O "/mnt/s3/fs1/1GB_$i.bin" https://ash-speed.hetzner.com/1GB.bin
done
```

**Verify files on mount:**

```bash
ls -lh /mnt/s3/fs1/
```

**Verify files synced to S3 (wait ~60 seconds):**

```bash
aws s3 ls s3://$BUCKET --recursive --human-readable --summarize
```

---

## 9. Lifecycle Rule (Keep Versioning Cost Low)

```bash
aws s3api put-bucket-lifecycle-configuration \
  --bucket $BUCKET \
  --lifecycle-configuration '{
    "Rules": [{
      "ID": "delete-old-versions",
      "Status": "Enabled",
      "Filter": { "Prefix": "" },
      "NoncurrentVersionExpiration": { "NoncurrentDays": 1 },
      "ExpiredObjectDeleteMarker": true,
      "AbortIncompleteMultipartUpload": { "DaysAfterInitiation": 1 }
    }]
  }'
```

---

## Quick Reference

| What | Value |
|------|-------|
| File System ID | `fs-0fe3ff6ef4ca0ae1b` |
| Mount Path | `/mnt/s3/fs1` |
| NFS Port | `2049` |
| Write sync delay to S3 | ~60 seconds |
| Versioning | Required (mandatory for S3 Files) |
