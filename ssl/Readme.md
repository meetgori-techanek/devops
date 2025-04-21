# ssl

1. install certbot on server:
```
sudo apt install certbot
sudo apt install python3-certbot-nginx
```

2. create ssl certificate for domain
- wildcard certificate mannualy 
```
sudo certbot certonly --manual --preferred-challenges=dns --email <email address> --server https://acme-v02.api.letsencrypt.org/directory --agree-tos -d *.<domain name>
```
once certificate issue then need to add this in nginx file 
this will redirect http traffic to https 
```
server{
        listen 80;
        server_name <domain-name>;
        return 301 https://$host$request_uri;
        #root /var/www/html;
        }
```
this will listen port 443 for https traffic 
need to remove port 80 from config 
```
  listen 443 ssl default_server;
         listen [::]:443 ssl default_server;
         ssl_certificate  /etc/letsencrypt/live/<domain-name>/fullchain.pem;
         ssl_certificate_key /etc/letsencrypt/live/<domain-name>/privkey.pem;
```

once this added then vertify and restart nginx 


- wildcard certificate with automaticalluy with nginx plugin
```
sudo certbot certonly --manual --preferred-challenges=dns  --installer nginx --email <email address> --server https://acme-v02.api.letsencrypt.org/directory --agree-tos -d *.<domain name>
```

note: need to add txt records manually in domain providers and ser ttl 2 minutes and need to wait for 2 minutes 

note: it will not work with free domain try use paid domain

3. need to configure nginx file and verify it then restart
```
sudo nginx -t
sudo systemctl restart nginx
```

4. to check ssl certificates
```
sudo certbot certificates
```
```
openssl x509 -in /etc/letsencrypt/live/<domain-name>/fullchain.pem -text -noout
```

5. delete ssl certificates
```
sudo certbot delete
```

command refrances : https://eff-certbot.readthedocs.io/en/latest/using.html#certbot-commands
