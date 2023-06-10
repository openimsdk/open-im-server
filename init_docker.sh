#!/usr/bin/env bash

# This script is used to check the environment and start the docker containers

# Define the directory path
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Define the functions
function check_command {
    if ! command -v $1 &> /dev/null; then
        echo "$1 command not found. Please install it first."
        exit 1
    fi
}

function check_docker {
    if ! docker ps &> /dev/null; then
        echo "Docker is not running. Please start it first."
        exit 1
    fi
}

# Check if the necessary commands are installed
check_command docker
check_command docker-compose

# Check if Docker is running
check_docker

# Change to the script directory
cd $SCRIPT_DIR

# Set permissions for the scripts
chmod +x *.sh

# Check the environment
./env_check.sh

# Start the docker containers
docker-compose up -d

# Check the docker services
./docker_check_service.sh