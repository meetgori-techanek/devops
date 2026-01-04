# kubectl – Kubernetes CLI (Ubuntu Installation Guide)

## 1. Update System
```
sudo apt update
```

---

## 2. Download kubectl (Official – Stable)
```
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
```

---

## 3. Install kubectl Binary
```
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```

---

## 4. Verify Installation
```
kubectl version --client
```

Expected output:
```text
Client Version: v1.xx.x
```

---

## 5. Enable kubectl  Autocomplete (Optional)
```
sudo apt install -y -completion
kubectl completion  | sudo tee /etc/_completion.d/kubectl
source ~/.rc
```

---

## 6. Configure kubectl for EKS

Update kubeconfig using AWS CLI:
```
aws eks update-kubeconfig \
  --region ap-south-1 \
  --name <EKS_CLUSTER_NAME>
```

Config file location:
```text
~/.kube/config
```

---

## 7. Verify Cluster Access
```
kubectl get nodes
```

Expected output:
```text
NAME           STATUS   ROLES    AGE   VERSION
ip-xxx-xxx     Ready    <none>   ...   v1.xx.x
```

---

## Result

kubectl is installed and configured to communicate with the EKS cluster for cluster management and operations.

---

# eksctl – Amazon EKS Cluster Management Tool (Ubuntu Installation Guide)

## 1. Update System
```
sudo apt update
```

---

## 2. Download eksctl (Official)
```
curl -sL "https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_Linux_amd64.tar.gz" \
| tar xz -C /tmp
```

---

## 3. Install eksctl
```
sudo mv /tmp/eksctl /usr/local/bin
```

---

## 4. Verify Installation
```
eksctl version
```

Expected output:
```text
0.xx.x
```

---

## 5. Verify AWS CLI Access (Required)

eksctl uses AWS CLI credentials or an attached IAM role.
```
aws sts get-caller-identity
```

---

## 6. Create an EKS Cluster (Example)
```
eksctl create cluster \
  --name demo-cluster \
  --region ap-south-1 \
  --nodegroup-name linux-nodes \
  --node-type t3.medium \
  --nodes 2 \
  --managed
```

---

## 7. List Existing Clusters
```
eksctl get clusters --region ap-south-1
```

---

## 8. Delete Cluster (Cleanup)
```
eksctl delete cluster \
  --name demo-cluster \
  --region ap-south-1
```

---

## 9. Required IAM Permissions

Attach one of the following:

### Recommended (Managed Policy)
```text
AdministratorAccess
```

### OR Minimum Required Policies
```text
AmazonEKSClusterPolicy
AmazonEKSWorkerNodePolicy
IAMFullAccess
AmazonEC2FullAccess
CloudFormationFullAccess
```

eksctl relies heavily on CloudFormation, EC2, IAM, and EKS APIs.
