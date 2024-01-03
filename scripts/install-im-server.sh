#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
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


# Common utilities, variables and checks for all build scripts.
set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

trap 'openim::util::onCtrlC' INT

chmod +x "${OPENIM_ROOT}"/scripts/*.sh

openim::util::ensure_docker_daemon_connectivity

DOCKER_COMPOSE_COMMAND=
# Check if docker-compose command is available
openim::util::check_docker_and_compose_versions

if command -v docker compose &> /dev/null
then
    openim::log::info "docker compose command is available"
    DOCKER_COMPOSE_COMMAND="docker compose"
else
    DOCKER_COMPOSE_COMMAND="docker-compose"
fi

"${OPENIM_ROOT}"/scripts/init-config.sh
pushd "${OPENIM_ROOT}"
${DOCKER_COMPOSE_COMMAND} stop
curl https://gitee.com/openimsdk/openim-docker/raw/main/example/full-openim-server-and-chat.yml -o docker-compose.yml
${DOCKER_COMPOSE_COMMAND} up -d

# Wait for a short period to allow containers to initialize
sleep 10

# Check the status of the containers
if ! ${DOCKER_COMPOSE_COMMAND} ps | grep -q 'Up'; then
    echo "Error: One or more docker containers failed to start."
    ${DOCKER_COMPOSE_COMMAND} logs
    exit 1
fi

sleep 50  # Keep the original 60-second wait, adjusted for the 10-second check above
${DOCKER_COMPOSE_COMMAND} logs openim-server
${DOCKER_COMPOSE_COMMAND} ps

popd
