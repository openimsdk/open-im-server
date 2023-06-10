#!/usr/bin/env bash

# Get the external IP address
internet_ip=$(curl -s ifconfig.me)

# Print the IP address
echo "Internet IP address: ${internet_ip}"

# Read environment variables
source .env

# Check if MINIO_ENDPOINT variable is set correctly
if [[ -z "${MINIO_ENDPOINT}" || "${MINIO_ENDPOINT}" != "http://127.0.0.1:10005" ]]; then
  echo "Error: MINIO_ENDPOINT is not set or is not equal to http://127.0.0.1:10005."
  exit 1
fi

# Replace the IP address in the .env file
sed -i "s/127.0.0.1/${internet_ip}/" .env

# Enter the script directory
cd script || exit 1

# Add execute permission to scripts
chmod +x *.sh

# Initialize password
./init_pwd.sh

# Check environment variables
./env_check.sh

# Return to the parent directory
cd ..

# Start docker containers
docker-compose up -d

# Enter the script directory
cd script || exit 1

# Check if docker services are running properly
./docker_check_service.sh
