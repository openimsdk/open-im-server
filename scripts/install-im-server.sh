#!/usr/bin/env bash
#
# OpenIM Docker Deployment Script
#
# This script automates the process of building the OpenIM server image
# and deploying it using Docker Compose.
#
# Variables:
#   - SERVER_IMAGE_VERSION: Version of the server image (default: test)
#   - IMAGE_REGISTRY: Docker image registry (default: openim)
#   - DOCKER_COMPOSE_FILE_URL: URL to the docker-compose.yml file
#
# Usage:
#   SERVER_IMAGE_VERSION=latest IMAGE_REGISTRY=myregistry ./this_script.sh

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

trap 'openim::util::onCtrlC' INT

chmod +x "${OPENIM_ROOT}"/scripts/*.sh

openim::util::ensure_docker_daemon_connectivity

# Default values for variables
: ${SERVER_IMAGE_VERSION:=test}
: ${IMAGE_REGISTRY:=openim}
: ${DOCKER_COMPOSE_FILE_URL:="https://raw.githubusercontent.com/openimsdk/openim-docker/main/docker-compose.yaml"}

DOCKER_COMPOSE_COMMAND=
# Check if docker-compose command is available
openim::util::check_docker_and_compose_versions
if command -v docker compose &> /dev/null; then
    openim::log::info "docker compose command is available"
    DOCKER_COMPOSE_COMMAND="docker compose"
else
    DOCKER_COMPOSE_COMMAND="docker-compose"
fi

export SERVER_IMAGE_VERSION
export IMAGE_REGISTRY
"${OPENIM_ROOT}"/scripts/init-config.sh

pushd "${OPENIM_ROOT}"
docker build -t "${IMAGE_REGISTRY}/openim-server:${SERVER_IMAGE_VERSION}" .
${DOCKER_COMPOSE_COMMAND} stop
curl "${DOCKER_COMPOSE_FILE_URL}" -o docker-compose.yml
${DOCKER_COMPOSE_COMMAND} up -d

# Function to check container status
check_containers() {
    if ! ${DOCKER_COMPOSE_COMMAND} ps | grep -q 'Up'; then
        echo "Error: One or more docker containers failed to start."
        ${DOCKER_COMPOSE_COMMAND} logs
        return 1
    fi
    return 0
}

# Wait for a short period to allow containers to initialize
sleep 30
check_containers

${DOCKER_COMPOSE_COMMAND} logs openim-server
${DOCKER_COMPOSE_COMMAND} ps

popd