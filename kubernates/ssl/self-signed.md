# Self Signed ssl with openssl

1. d
   ```
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
   -keyout tls.key \
   -out tls.crt \
   -subj "/CN=techanek.work.gd/O=test-app"
   ```

2. d
   ```
   kubectl create secret tls self-signed-tls \
     --cert=tls.crt \
     --key=tls.key \
     -n <app-namespace>
   ```

3. update ingress
   ```
   ....
   spec:
    ingressClassName: nginx
     tls:
     - secretName: self-signed-tls
       hosts: 
        - "techanek.work.gd"
     rules:
   ....
   ```

4. verify
```
curl -v https://techanek.work.gd -k
```
