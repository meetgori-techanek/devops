# Monitoring Stack — Deployment Steps (EKS)

> **Namespace:** `monitoring`  

---

## Prerequisites

### 1. Namespace

```bash
kubectl create namespace monitoring
```

## Helm Repos

**References:**
- https://www.linkedin.com/posts/jkroepke_helm-repository-migration-grafana-community-activity-741684…
- https://github.com/grafana/helm-charts/tree/main/charts

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo add grafana-community https://grafana-community.github.io/helm-charts
helm repo update
```

**Chart Mapping:**


| Component      | Chart to use                          | Github                                                                                                                                                                                  |
| -------------- | ------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| K8S Monitoring | `grafana/k8s-monitoring`              | [https://github.com/grafana/k8s-monitoring-helm](https://github.com/grafana/k8s-monitoring-helm)    
| Alloy          | `grafana/alloy`                       | [https://github.com/grafana/helm-charts/blob/main/charts/alloy/README.md](https://github.com/grafana/helm-charts/blob/main/charts/alloy/README.md)                                             |
| Mimir          | `grafana/mimir-distributed`           | [https://github.com/grafana/mimir/blob/main/operations/helm/charts/mimir-distributed/README.md](https://github.com/grafana/mimir/blob/main/operations/helm/charts/mimir-distributed/README.md) |                     |
| Loki           | `grafana-community/loki`              | [https://github.com/grafana-community/helm-charts/blob/main/charts/loki/README.md](https://github.com/grafana-community/helm-charts/blob/main/charts/loki/README.md)                           |
| Tempo          | `grafana-community/tempo-distributed` | [https://github.com/grafana-community/helm-charts/blob/main/charts/tempo-distributed/README.md](https://github.com/grafana-community/helm-charts/blob/main/charts/tempo-distributed/README.md) |
| Grafana        | `grafana-community/grafana`           | [https://github.com/grafana-community/helm-charts/blob/main/charts/grafana/README.md](https://github.com/grafana-community/helm-charts/blob/main/charts/grafana/README.md)
 

---

## Individual Commands

### k8s monitoring(contains node exporter, alloy, kube-state metrics)

```bash
helm upgrade --install k8s-monitoring grafana/k8s-monitoring -f k8s-values.yaml -n monitoring
helm uninstall k8s-monitoring -n monitoring
```
### OR 

### Alloy
```bash
helm upgrade --install alloy grafana/alloy -f alloy-values.yaml -n monitoring
helm uninstall alloy -n monitoring
```

### Mimir
```bash
helm upgrade --install mimir grafana/mimir-distributed -f mimir-values.yaml -n monitoring
helm uninstall mimir -n monitoring
```

### Loki
```bash
helm upgrade --install loki grafana-community/loki-distributed -f loki-values.yaml -n monitoring
helm uninstall loki -n monitoring 
```

### Tempo
```bash
helm upgrade --install tempo grafana-community/tempo-distributed -f tempo-values.yaml -n monitoring
helm uninstall tempo -n monitoring
```
### Grafana Admin Secret

```bash
kubectl create secret generic grafana-secret -n monitoring \
  --from-literal=admin-user=admin \
  --from-literal=admin-password=<YOUR_SECURE_PASSWORD>
```

### Grafana

```bash
helm upgrade --install grafana grafana-community/grafana -f grafana-values.yaml -n monitoring
helm uninstall grafana -n monitoring
```

## Buckets Required
Create these buckets in your MinIO instance before deploying:

| Bucket | Used by |
|--------|----------|
| `mimir-blocks` | Mimir blocks storage |
| `loki-chunks` | Loki chunks |
| `tempo-traces` | Tempo traces |
---


### 2. EKS StorageClass

Ensure the EBS CSI driver is installed and create the StorageClass if it doesn't exist:

```bash
kubectl apply -f gp3-storageclass.yaml
```

---

## Domains

| Component | URL |
|-----------|-----|
| Grafana | http://grafana.monitoring.labmeet.xyz |
| Mimir | https://mimir.monitoring.labmeet.xyz |
| Loki | https://loki.monitoring.labmeet.xyz |
| Tempo | not have ui |
| Alloy | https://alloy.monitoring.labmeet.xyz |