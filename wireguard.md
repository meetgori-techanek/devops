# WireGuard Server Setup on AWS EC2 (Ubuntu)
#### This guide walks through the steps to install and configure a WireGuard VPN server on an AWS EC2 instance running Ubuntu.

## Prerequisites
- An AWS EC2 instance running Ubuntu (e.g., Ubuntu 20.04 or 22.04).
- SSH access to the instance (you'll likely be using the ubuntu user).
- The EC2 Security Group must have an Inbound Rule to allow UDP traffic on port 51820 from your client IPs.

1. Update System and Install WireGuard
First, update your package lists and install the WireGuard package.
### Update package lists
```
sudo apt update
```

### Install WireGuard
```
sudo apt install wireguard -y
```
2. Prepare Configuration Directory and Keys\
Create the necessary configuration directory and generate the server's private and public keys.\

Create Folder and Set Permissions

### Create the WireGuard configuration folder
```
sudo mkdir -p /etc/wireguard
```

### Set secure permissions (only root can read/write)
```
sudo chmod 700 /etc/wireguard
```
## Generate Server Keys
### Generate private key and save it securely
```
sudo wg genkey | sudo tee /etc/wireguard/server_private.key > /dev/null
```
### Generate the public key from the private key
```
sudo cat /etc/wireguard/server_private.key | sudo wg pubkey | sudo tee /etc/wireguard/server_public.key > /dev/null
```

### Verify Key Files
### Check that both keys have been created in the correct location.
```
sudo ls -l /etc/wireguard/
```

3. Create WireGuard Server Configuration (wg0.conf)
The server configuration file (wg0.conf) defines the VPN tunnel interface, IP, port, private key, and crucial firewall rules for Network Address Translation (NAT) and IP forwarding.\

Bash

sudo bash -c 'cat > /etc/wireguard/wg0.conf <<EOF
[Interface]
# Internal VPN IP for the server
Address = 10.8.0.1/24
# Listening port (must match Security Group rule)
ListenPort = 51820
# Server's private key (fetched automatically)
PrivateKey = $(sudo cat /etc/wireguard/server_private.key)
# Save client configurations automatically upon addition
SaveConfig = true

# PostUp/PostDown rules for NAT (Masquerading) and traffic forwarding
# Assuming your public-facing interface is 'eth0' on EC2.
PostUp = iptables -A FORWARD -i %i -j ACCEPT
PostUp = iptables -A FORWARD -o %i -j ACCEPT
PostUp = iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT
PostDown = iptables -D FORWARD -o %i -j ACCEPT
PostDown = iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
EOF'
4. Enable IP Forwarding
For the EC2 instance to route traffic between the VPN tunnel and the internet, you must enable IPv4 forwarding.

Bash

# Set the forwarding parameter in sysctl.conf
echo "net.ipv4.ip_forward=1" | sudo tee -a /etc/sysctl.conf

# Apply the new settings immediately
sudo sysctl -p
5. Configure the Firewall (UFW)
You need to open the WireGuard port and ensure SSH access is maintained.

Bash

# Install UFW if it's not present
sudo apt install ufw -y

# Allow WireGuard UDP port
sudo ufw allow 51820/udp

# Allow standard SSH access (port 22)
sudo ufw allow OpenSSH

# Enable the firewall (confirm with 'y' if prompted)
sudo ufw enable

# Check the firewall status
sudo ufw status
You should see rules for 51820/udp and OpenSSH.

6. Enable and Start WireGuard Service
Start the WireGuard interface (wg0) and ensure it starts automatically on boot.

Bash

# Enable the service to start on boot
sudo systemctl enable wg-quick@wg0

# Start the WireGuard service now
sudo systemctl start wg-quick@wg0

# Check the service status (should show 'active (exited)')
sudo systemctl status wg-quick@wg0
7. Verification
Verify that the wg0 interface is up and running.

Bash

# Show WireGuard interface details
sudo wg show
The output should display your public key, the private key hash, and the listening port (51820).


### dump
#### Check server’s allowed IPs and keys
ensure your server /etc/wireguard/wg0.conf has the correct peer public key and allowed IPs.
For example:\
Server:
```
[Interface]
Address = 10.8.0.1/24
ListenPort = 51820
PrivateKey = <server_private_key>

[Peer]
PublicKey = <client_public_key>
AllowedIPs = 10.8.0.2/32
```

Client:
```
[Interface]
Address = 10.8.0.2/24
PrivateKey = <client_private_key>
DNS = 8.8.8.8

[Peer]
PublicKey = <server_public_key>
AllowedIPs = 0.0.0.0/0
Endpoint = 54.91.198.90:51820
PersistentKeepalive = 25
```
