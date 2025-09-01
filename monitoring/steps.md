# Kubernetes Monitoring & Observability Setup

This guide walks through the setup of a containerized logging app, Kubernetes metrics server, Helm, NGINX ingress, Prometheus, Grafana, Loki, and Alloy for full observability.

---

## Build Docker Image & Push to Docker Hub


### Build the Docker image
```
docker build -t gocpulogger:latest .
```
### Login to Docker Hub
```
docker login -u meetgori1
```
### Tag and push the image
```
docker tag gocpulogger meetgori1/gocpulogger:latest
docker push meetgori1/gocpulogger:latest
```

## Install and Configure Metrics Server
### Apply metrics-server components
```
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml -n kube-system
```

### Patch the deployment for internal communication
```
kubectl patch deployment metrics-server -n kube-system --type='json' -p='[
  {
    "op": "add",
    "path": "/spec/template/spec/hostNetwork",
    "value": true
  },
  {
    "op": "replace",
    "path": "/spec/template/spec/containers/0/args",
    "value": [
      "--cert-dir=/tmp",
      "--secure-port=4443",
      "--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
      "--kubelet-use-node-status-port",
      "--metric-resolution=15s",
      "--kubelet-insecure-tls"
    ]
  },
  {
    "op": "replace",
    "path": "/spec/template/spec/containers/0/ports/0/containerPort",
    "value": 4443
  }
]'
```
### Test metrics:
```
kubectl top nodes -n kube-system
kubectl get all,ing -A
```

##  Install Helm
### Add Helm repo and key
```
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
```

### Install dependencies
```
sudo apt-get install apt-transport-https --yes
```

### Add Helm stable repo
```
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
```

### Update and install Helm
```
sudo apt update
sudo apt install helm
```

## Set Up Ingress Controller (NGINX)
### Create ingress namespace
```
kubectl create ns ingress
```

### Install NGINX Ingress CRDs
```
kubectl apply -f https://raw.githubusercontent.com/nginx/kubernetes-ingress/v5.0.0/deploy/crds.yaml -n ingress
```

### Install NGINX Ingress using Helm
```
helm install nginx-ingress oci://ghcr.io/nginx/charts/nginx-ingress \
  --version 2.1.0 \
  --set controller.debug.enable=false \
  --namespace ingress
```

### Apply Ingress Resource
```
kubectl apply -f ingress.yml -n ingress
```

### Enable Host Network for Ingress
### Check if hostNetwork is false (default)
```
helm get values -a nginx-ingress -n ingress | grep hostNetwork
```

### Upgrade ingress with hostNetwork enabled
```
helm upgrade nginx-ingress oci://ghcr.io/nginx/charts/nginx-ingress \
  --set controller.hostNetwork=true \
  --namespace ingress
```
### Network Testing
### Test connectivity
```
nc -vz techanek.work.gd 80
nslookup prometheus.techanek.work.gd
```
### Check app status
```
curl http://techanek.work.gd/api/status
```

## Install Prometheus and Grafana (via kube-prometheus-stack)
### Add Prometheus Helm repo
```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```
### Install the stack
```
helm install my-kube-prometheus-stack prometheus-community/kube-prometheus-stack -n monitoring
```
### Upgrade with custom values
```
helm upgrade my-kube-prometheus-stack prometheus-community/kube-prometheus-stack -n monitoring -f prom-values.yml
```

## Install Loki and Alloy
### Add Grafana repo
```
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```
### Install/upgrade Loki
```
helm install loki grafana/loki -n monitoring -f loki-values.yml
```
### Grafana Loki URL:
```
http://loki-gateway.monitoring.svc.cluster.local/
```

### Install Grafana Alloy 
```
helm install alloy grafana/alloy -n monitoring -f alloy-values.yml
```

### Install Grafana Tempo
```
helm install tempo grafana/tempo \
  -n monitoring \
  -f tempo-values.yml
```
### Summary
| Component      | Installed via | Namespace   |
| -------------- | ------------- | ----------- |
| app            | Deployment    | go-app      |
| Metrics Server | `kubectl`     | kube-system |
| Helm           | `apt`         | -           |
| Ingress NGINX  | Helm + YAML   | ingress     |
| Prometheus     | Helm Chart    | monitoring  |
| Grafana        | Helm Chart    | monitoring  |
| Loki           | Helm Chart    | monitoring  |
| Alloy          | Helm Chart    | monitoring  |
| Tempo          | Helm Chart    | monitoring  |

