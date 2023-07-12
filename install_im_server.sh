#!/usr/bin/env bash

# Get the public internet IP address
internet_ip=$(curl ifconfig.me -s)
echo $internet_ip

# Load environment variables from .env file
source .env
echo $MINIO_ENDPOINT

# Replace local IP address with the public IP address in .env file
if [ $MINIO_ENDPOINT == "http://127.0.0.1:10005" ]; then
    sed -i "s/127.0.0.1/${internet_ip}/" .env
fi

# Change directory to scripts folder
cd scripts
chmod +x *.sh

# Execute necessary scripts
./init_pwd.sh
./env_check.sh

# Go back to the previous directory
cd ..

# Check if docker-compose command is available
if command -v docker-compose &> /dev/null
then
    docker-compose up -d
else
    docker compose up -d
fi

# Change directory to scripts folder again
cd scripts

# Check docker services
./docker_check_service.sh