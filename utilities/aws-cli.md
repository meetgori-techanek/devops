# AWS CLI – Ubuntu Installation Guide

## 1. Update System
```
sudo apt update
```

---

## 2. Install Prerequisites

Ensure required packages are installed:
```
sudo apt install -y curl unzip
```

---

## 3. Download AWS CLI v2 (Official)
```
cd /tmp
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
```

---

## 4. Extract and Install
```
unzip awscliv2.zip
sudo ./aws/install
```

---

## 5. Verify Installation
```
aws --version
```

---

## 6. Verify AWS CLI Access

Check caller identity:
```
aws sts get-caller-identity
```
