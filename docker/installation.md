## Docker installation
### Update package lists:
```
sudo apt update
```
### Install dependencies: Install the packages needed to use apt over HTTPS.
```
sudo apt install ca-certificates curl gnupg
```
### Add Docker's official GPG key: Add the GPG key to verify the downloaded packages.
```
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
```
### Add the Docker repository
```
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```
### Update package lists again
```
sudo apt update
```
### Install Docker using apt
```
sudo apt install docker.io 
```

## Allow non-root users to run Docker
```
sudo usermod -aG docker $USER
newgrp docker
```

## Install Docker Compose 
1. Define the Docker configuration directory environment variable, defaulting to ~/.docker if not set:
```
DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
```

2. Create the directory for Docker CLI plugins:
```
mkdir -p $DOCKER_CONFIG/cli-plugins 
```

3. Download the Docker Compose plugin binary (version 2.34.0) directly from the official GitHub releases into the CLI plugins directory:
```
curl -SL https://github.com/docker/compose/releases/download/v2.34.0/docker-compose-linux-x86_64 -o $DOCKER_CONFIG/cli-plugins/docker-compose 
```

4. Make the Docker Compose binary executable:
```
chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose 
```
5. Verify the installation by checking the version:
```
docker compose version 
```
