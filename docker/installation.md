## Docker installation
### Install Docker using apt
```
sudo apt install docker.io 
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
