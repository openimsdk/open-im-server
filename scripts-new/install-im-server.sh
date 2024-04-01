#!/usr/bin/env bash
# Copyright © 2024 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
    ${DOCKER_COMPOSE_COMMAND} logs openim-server
    ${DOCKER_COMPOSE_COMMAND} logs openim-chat
    return 1
  fi
  return 0
}

# Wait for a short period to allow containers to initialize
sleep 100

${DOCKER_COMPOSE_COMMAND} ps

check_containers

popd