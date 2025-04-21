# ssl
## install certbot

1. install certbot on server:
```
sudo apt install certbot
sudo apt install python3-certbot-nginx
```

2. create wildcard certificate for domain
```
sudo certbot certonly --manual --preferred-challenges=dns --email meet@mailinator.com --server https://acme-v02.api.letsencrypt.org/directory --agree-tos -d *.excercise6.work.gd -d *.excercise6.work.gd
```
need to add txt records manually in domain providers and ser ttl 2 minutes and need to wait for 2 minutes 

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
