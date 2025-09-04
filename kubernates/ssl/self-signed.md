# Self-Signed SSL with OpenSSL and Kubernetes

1. Generate Self-Signed Certificate\
Use openssl to generate a 2048-bit RSA key and a self-signed certificate valid for 1 year (365 days):
```
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
   -keyout tls.key \
   -out tls.crt \
   -subj "/CN=techanek.work.gd/O=test-app"
```

2. Create Kubernetes TLS Secret\
Create a Kubernetes TLS secret from the generated certificate and key: 
```
   kubectl create secret tls self-signed-tls \
     --cert=tls.crt \
     --key=tls.key \
     -n <app-namespace>
```

4. Update Ingress Resource
```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-ingress
  namespace: go-app
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - secretName: self-signed-tls
    hosts: 
      - "techanek.work.gd"
  rules:
  - host: techanek.work.gd
   ....
```
`nginx.ingress.kubernetes.io/ssl-redirect: "true"` will redirect all http trafic to https automatically

5. verify
```
curl -v https://techanek.work.gd -k
```
