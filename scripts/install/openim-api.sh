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

set -o errexit
set +o nounset
set -o pipefail

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source "${OPENIM_ROOT}"/scripts/install/common.sh

SERVER_NAME="openim-api"

readonly OPENIM_API_PORT_TARGETS=(
  ${API_OPENIM_PORT}
)
readonly OPENIM_API_PORT_LISTARIES=("${OPENIM_API_PORT_TARGETS[@]##*/}")

readonly OPENIM_API_SERVICE_TARGETS=(
  openim-api
)
readonly OPENIM_API_SERVICE_LISTARIES=("${OPENIM_API_SERVICE_TARGETS[@]##*/}")

function openim::api::start() {
    echo "++ OPENIM_API_SERVICE_LISTARIES: ${OPENIM_API_SERVICE_LISTARIES[@]}"
    echo "++ OPENIM_API_PORT_LISTARIES: ${OPENIM_API_PORT_LISTARIES[@]}"
    echo "++ OpenIM API config path: ${OPENIM_API_CONFIG}"

    openim::log::info "Starting ${SERVER_NAME} ..."

    printf "+------------------------+--------------+\n"
    printf "| Service Name           | Port         |\n"
    printf "+------------------------+--------------+\n"

    length=${#OPENIM_API_SERVICE_LISTARIES[@]}

    for ((i=0; i<$length; i++)); do
    printf "| %-22s | %6s       |\n" "${OPENIM_API_SERVICE_LISTARIES[$i]}" "${OPENIM_API_PORT_LISTARIES[$i]}"
    printf "+------------------------+--------------+\n"
    done
    # start all api services
    for ((i = 0; i < ${#OPENIM_API_SERVICE_LISTARIES[*]}; i++)); do
    openim::util::stop_services_on_ports ${OPENIM_API_PORT_LISTARIES[$i]}
    openim::log::info "OpenIM ${OPENIM_API_SERVICE_LISTARIES[$i]} config path: ${OPENIM_API_CONFIG}"

    # Get the service and Prometheus ports.
    OPENIM_API_SERVICE_PORTS=( $(openim::util::list-to-string ${OPENIM_API_PORT_LISTARIES[$i]}) )

    # TODO Only one port is supported. An error occurs on multiple ports
    if [ ${#OPENIM_API_SERVICE_PORTS[@]} -ne 1 ]; then
        openim::log::error_exit "Set only one port for ${OPENIM_API_SERVICE_LISTARIES[$i]} service."
    fi

    for ((j = 0; j < ${#OPENIM_API_SERVICE_PORTS[@]}; j++)); do
        openim::log::info "Starting ${OPENIM_API_SERVICE_LISTARIES[$i]} service, port: ${OPENIM_API_SERVICE_PORTS[j]}, binary root: ${OPENIM_OUTPUT_HOSTBIN}/${OPENIM_API_SERVICE_LISTARIES[$i]}"
        openim::api::start_service "${OPENIM_API_SERVICE_LISTARIES[$i]}" "${OPENIM_API_PORT_LISTARIES[j]}"
        sleep 1
      done
    done

    OPENIM_API_PORT_STRINGARIES=( $(openim::util::list-to-string ${OPENIM_API_PORT_LISTARIES[@]}) )
    openim::util::check_ports ${OPENIM_API_PORT_STRINGARIES[@]}
}

function openim::api::start_service() {
  local binary_name="$1"
  local service_port="$2"
  local prometheus_port="$3"

  local cmd="${OPENIM_OUTPUT_HOSTBIN}/${binary_name} --port ${service_port} -c ${OPENIM_API_CONFIG}"

  nohup ${cmd} >> "${LOG_FILE}" 2>&1 &

  if [ $? -ne 0 ]; then
    openim::log::error_exit "Failed to start ${binary_name} on port ${service_port}."
  fi
}

###################################### Linux Systemd ######################################
SYSTEM_FILE_PATH="/etc/systemd/system/${SERVER_NAME}.service"

# Print the necessary information after installation
function openim::api::info() {
cat << EOF
openim-api listen on: ${OPENIM_API_HOST}:${API_OPENIM_PORT}
EOF
}

# install openim-api
function openim::api::install() {
    openim::log::info "Installing ${SERVER_NAME} ..."

    pushd "${OPENIM_ROOT}"

    # 1. Build openim-api
    make build BINS=${SERVER_NAME}
    openim::common::sudo "cp -r ${OPENIM_OUTPUT_HOSTBIN}/${SERVER_NAME} ${OPENIM_INSTALL_DIR}/${SERVER_NAME}"
    openim::log::status "${SERVER_NAME} binary: ${OPENIM_INSTALL_DIR}/${SERVER_NAME}/${SERVER_NAME}"

    # 2. Generate and install the openim-api configuration file (config)
    openim::log::status "${SERVER_NAME} config file: ${OPENIM_CONFIG_DIR}/config.yaml"

    # 3. Create and install the ${SERVER_NAME} systemd unit file
    echo ${LINUX_PASSWORD} | sudo -S bash -c \
        "SERVER_NAME=${SERVER_NAME} ./scripts/genconfig.sh ${ENV_FILE} deployments/templates/openim.service > ${SYSTEM_FILE_PATH}"
    openim::log::status "${SERVER_NAME} systemd file: ${SYSTEM_FILE_PATH}"

    # 4. Start the openim-api service
    openim::common::sudo "systemctl daemon-reload"
    openim::common::sudo "systemctl restart ${SERVER_NAME}"
    openim::common::sudo "systemctl enable ${SERVER_NAME}"
    openim::api::status || return 1
    openim::api::info

    openim::log::info "install ${SERVER_NAME} successfully"
    popd
}

# Unload
function openim::api::uninstall() {
    openim::log::info "Uninstalling ${SERVER_NAME} ..."

    set +o errexit
    openim::common::sudo "systemctl stop ${SERVER_NAME}"
    openim::common::sudo "systemctl disable ${SERVER_NAME}"
    openim::common::sudo "rm -f ${OPENIM_INSTALL_DIR}/${SERVER_NAME}"
    openim::common::sudo "rm -f ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"
    openim::common::sudo "rm -f /etc/systemd/system/${SERVER_NAME}.service"
    set -o errexit
    openim::log::info "uninstall ${SERVER_NAME} successfully"
}

# Status Check
function openim::api::status() {
    openim::log::info "Checking ${SERVER_NAME} status ..."

    # Check the running status of the ${SERVER_NAME}. If active (running) is displayed, the ${SERVER_NAME} is started successfully.
    systemctl status ${SERVER_NAME}|grep -q 'active' || {
        openim::log::error "${SERVER_NAME} failed to start, maybe not installed properly"
        return 1
    }

    openim::util::check_ports ${OPENIM_API_PORT_LISTARIES[@]}
}

if [[ "$*" =~ openim::api:: ]];then
    eval $*
fi
