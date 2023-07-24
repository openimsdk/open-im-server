#!/usr/bin/env bash

mark=''
for ((ratio=0;${ratio}<=100;ratio+=5))
do
        sleep 0.2
        printf "progress:[%-40s]%d%%\r" "${mark}" "${ratio}"
        mark="##${mark}"
done
echo

set -e

# Change directory to the 'scripts' folder
cd scripts

# Grant execute permissions to all shell scripts in the 'scripts' folder
chmod +x *.sh

# Run the 'env_check.sh' script for environment checks
./env_check.sh

# Move back to the parent directory
cd ..

# Check if Docker is installed
if ! command -v docker >/dev/null 2>&1; then
  echo "Error: Docker is not installed. Please install Docker before running this script."
  exit 1
fi

# Start Docker services using docker-compose
if command -v docker-compose &> /dev/null
then
    docker-compose up -d
else
    docker compose up -d
fi

# Move back to the 'scripts' folder
cd scripts

# Run the 'docker_check_service.sh' script for Docker service checks
./docker_check_service.sh
