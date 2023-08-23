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


OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"
source "${OPENIM_ROOT}/scripts/install/environment.sh"

openim::util::onCtrlC

docker_compose_file_name="docker-compose.yaml"

# Load environment variables from .env file
load_env() {
    source "${OPENIM_ROOT}"/.env
}

# Replace local IP with public IP in .env
replace_ip() {
    if [ "$API_URL" == "http://127.0.0.1:10002/object/" ]; then
        sed -i "s/127.0.0.1/${internet_ip}/" "${OPENIM_ROOT}"/.env
    fi

    if [ "$MINIO_ENDPOINT" == "http://127.0.0.1:10005" ]; then
        sed -i "s/127.0.0.1/${internet_ip}/" "${OPENIM_ROOT}"/.env
    fi 

    openim::log::info "Your minio endpoint is ${MINIO_ENDPOINT}"
}

# Execute necessary scripts
execute_scripts() {
    chmod +x "${OPENIM_ROOT}"/scripts/*.sh
    openim::log::info "Executing init_pwd.sh"
    "${OPENIM_ROOT}"/scripts/init_pwd.sh

    openim::log::info "Executing env_check.sh"
    "${OPENIM_ROOT}"/scripts/env_check.sh
}

# Start docker compose
start_docker_compose() {
    openim::log::info "Checking if docker-compose command is available"
    if command -v docker-compose &> /dev/null; then
        docker-compose up -d
    else
        docker compose up -d
    fi

    "${OPENIM_ROOT}"/scripts/docker-check-service.sh
}

main() {
    load_env
    openim::util::get_server_ip
    replace_ip
    execute_scripts
    start_docker_compose
    openim::log::success "Script executed successfully"
}

# Run the main function
main