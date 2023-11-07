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
#
# OpenIM RPC Service Control Script
# 
# Description:
# This script provides a control interface for the OpenIM RPC service within a Linux environment. It offers functionalities to start multiple RPC services, each denoted by their respective names under openim::rpc::service_name.
# 
# Features:
# 1. Robust error handling using Bash built-ins like 'errexit', 'nounset', and 'pipefail'.
# 2. The capability to source common utility functions and configurations to ensure uniform environmental settings.
# 3. Comprehensive logging functionalities, providing a detailed understanding of operational processes.
# 4. Provision for declaring and managing a set of RPC services, each associated with its unique name and corresponding ports.
# 5. The ability to define and associate Prometheus ports for service monitoring purposes.
# 6. Functionalities to start each RPC service, along with its designated ports, in a sequence.
#
# Usage:
# 1. Direct Script Execution:
#    This initiates all the RPC services declared under the function openim::rpc::service_name.
#    Example: ./openim-rpc-{rpc-name}.sh  openim::rpc::start
# 2. Controlling through Functions for systemctl operations:
#    Specific operations like installation, uninstallation, and status check can be executed by passing the respective function name as an argument to the script.
#    Example: ./openim-rpc-{rpc-name}.sh openim::rpc::install
#
# Note: Before executing this script, ensure that the necessary permissions are granted and relevant environmental variables are set.
#

set -o errexit
set +o nounset
set -o pipefail

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source "${OPENIM_ROOT}"/scripts/install/common.sh

SERVER_NAME="openim-rpc"
readonly OPENIM_RPC_CONFIG="${OPENIM_ROOT}"/config

openim::rpc::service_name() {
  local targets=(
    openim-rpc-user
    openim-rpc-friend
    openim-rpc-msg
    openim-rpc-group
    openim-rpc-auth
    openim-rpc-conversation
    openim-rpc-third
  )
  echo "${targets[@]}"
}
IFS=" " read -ra OPENIM_RPC_SERVICE_TARGETS <<< "$(openim::rpc::service_name)"
readonly OPENIM_RPC_SERVICE_TARGETS
readonly OPENIM_RPC_SERVICE_LISTARIES=("${OPENIM_RPC_SERVICE_TARGETS[@]##*/}")

# Make sure the environment is only called via common to avoid too much nesting
openim::rpc::service_port() {
  local targets=(
    ${OPENIM_USER_PORT}            # User service 10110
    ${OPENIM_FRIEND_PORT}          # Friend service 10120
    ${OPENIM_MESSAGE_PORT}         # Message service 10130
    # ${OPENIM_MESSAGE_GATEWAY_PORT} # Message gateway 10140
    ${OPENIM_GROUP_PORT}           # Group service 10150
    ${OPENIM_AUTH_PORT}            # Authorization service 10160
    # ${OPENIM_PUSH_PORT}            # Push service 10170
    ${OPENIM_CONVERSATION_PORT}    # Conversation service 10180
    ${OPENIM_THIRD_PORT}           # Third-party service 10190
  )
  echo "${targets[@]}"
}
IFS=" " read -ra OPENIM_RPC_PORT_TARGETS <<< "$(openim::rpc::service_port)"
readonly OPENIM_RPC_PORT_TARGETS
readonly OPENIM_RPC_PORT_LISTARIES=("${OPENIM_RPC_PORT_TARGETS[@]##*/}")

openim::rpc::prometheus_port() {
  # Declare an array to hold all the Prometheus ports for different services
  local targets=(
    ${USER_PROM_PORT}               # Prometheus port for user service
    ${FRIEND_PROM_PORT}             # Prometheus port for friend service
    ${MESSAGE_PROM_PORT}            # Prometheus port for message service
    ${GROUP_PROM_PORT}              # Prometheus port for group service
    ${AUTH_PROM_PORT}               # Prometheus port for authentication service
    ${CONVERSATION_PROM_PORT}       # Prometheus port for conversation service
    ${THIRD_PROM_PORT}              # Prometheus port for third-party integrations service
  )
  # Print the list of ports
  echo "${targets[@]}"
}
IFS=" " read -ra OPENIM_RPC_PROM_PORT_TARGETS <<< "$(openim::rpc::prometheus_port)"
readonly OPENIM_RPC_PROM_PORT_TARGETS
readonly OPENIM_RPC_PROM_PORT_LISTARIES=("${OPENIM_RPC_PROM_PORT_TARGETS[@]##*/}")

function openim::rpc::start() {
    echo "OPENIM_RPC_SERVICE_LISTARIES: ${OPENIM_RPC_SERVICE_LISTARIES[@]}"
    echo "OPENIM_RPC_PROM_PORT_LISTARIES: ${OPENIM_RPC_PROM_PORT_LISTARIES[@]}"
    echo "OPENIM_RPC_PORT_LISTARIES: ${OPENIM_RPC_PORT_LISTARIES[@]}"

    openim::log::info "Starting ${SERVER_NAME} ..."

    printf "+------------------------+-------+-----------------+\n"
    printf "| Service Name           | Port  | Prometheus Port |\n"
    printf "+------------------------+-------+-----------------+\n"

    length=${#OPENIM_RPC_SERVICE_LISTARIES[@]}

    for ((i=0; i<$length; i++)); do
    printf "| %-22s | %-5s | %-15s |\n" "${OPENIM_RPC_SERVICE_LISTARIES[$i]}" "${OPENIM_RPC_PORT_LISTARIES[$i]}" "${OPENIM_RPC_PROM_PORT_LISTARIES[$i]}"
    printf "+------------------------+-------+-----------------+\n"
    done

    # start all rpc services
    for ((i = 0; i < ${#OPENIM_RPC_SERVICE_LISTARIES[*]}; i++)); do
        # openim::util::stop_services_with_name ${OPENIM_RPC_SERVICE_LISTARIES
        openim::util::stop_services_on_ports ${OPENIM_RPC_PORT_LISTARIES[$i]}
        openim::log::info "OpenIM ${OPENIM_RPC_SERVICE_LISTARIES[$i]} config path: ${OPENIM_RPC_CONFIG}"
    
        # Get the service and Prometheus ports.
        OPENIM_RPC_SERVICE_PORTS=( $(openim::util::list-to-string ${OPENIM_RPC_PORT_LISTARIES[$i]}) )
        read -a OPENIM_RPC_SERVICE_PORTS_ARRAY <<< ${OPENIM_RPC_SERVICE_PORTS}
        
        OPENIM_RPC_PROM_PORTS=( $(openim::util::list-to-string ${OPENIM_RPC_PROM_PORT_LISTARIES[$i]}) )
        read -a OPENIM_RPC_PROM_PORTS_ARRAY <<< ${OPENIM_RPC_PROM_PORTS}

        for ((j = 0; j < ${#OPENIM_RPC_SERVICE_PORTS_ARRAY[@]}; j++)); do
            openim::log::info "Starting ${OPENIM_RPC_SERVICE_LISTARIES[$i]} service, port: ${OPENIM_RPC_SERVICE_PORTS[j]}, prometheus port: ${OPENIM_RPC_PROM_PORTS[j]}, binary root: ${OPENIM_OUTPUT_HOSTBIN}/${OPENIM_RPC_SERVICE_LISTARIES[$i]}"
            openim::rpc::start_service "${OPENIM_RPC_SERVICE_LISTARIES[$i]}" "${OPENIM_RPC_SERVICE_PORTS[j]}" "${OPENIM_RPC_PROM_PORTS[j]}"
        done
    done

    sleep 0.5

    openim::util::check_ports ${OPENIM_RPC_PORT_TARGETS[@]}
    # openim::util::check_ports ${OPENIM_RPC_PROM_PORT_TARGETS[@]}

}

function openim::rpc::start_service() {
  local binary_name="$1"
  local service_port="$2"
  local prometheus_port="$3"

  local cmd="${OPENIM_OUTPUT_HOSTBIN}/${binary_name} --port ${service_port} -c ${OPENIM_RPC_CONFIG}"

  if [ -n "${prometheus_port}" ]; then
    printf "Specifying prometheus port: %s\n" "${prometheus_port}"
    cmd="${cmd} --prometheus_port ${prometheus_port}"
  fi
  nohup ${cmd} >> "${LOG_FILE}" 2>&1 &
}

###################################### Linux Systemd ######################################
declare -A SYSTEM_FILE_PATHS
for service in "${OPENIM_RPC_SERVICE_LISTARIES[@]}"; do
    SYSTEM_FILE_PATHS["$service"]="/etc/systemd/system/${service}.service"
done

# Print the necessary information after installation
function openim::rpc::info() {
    for service in "${OPENIM_RPC_SERVICE_LISTARIES[@]}"; do
        echo "${service} listen on: ${OPENIM_RPC_PORT_LISTARIES[@]}"
    done
}

# install openim-rpc
function openim::rpc::install() {
    pushd "${OPENIM_ROOT}"

    # 1. Build openim-rpc
    for service in "${OPENIM_RPC_SERVICE_LISTARIES[@]}"; do
        make build BINS=${service}
        openim::common::sudo "cp -r ${OPENIM_OUTPUT_HOSTBIN}/${service} ${OPENIM_INSTALL_DIR}/${service}"
        openim::log::status "${service} binary: ${OPENIM_INSTALL_DIR}/${service}/${service}"
    done

    # 2. Generate and install the openim-rpc configuration file (config)
    openim::log::status "openim-rpc config file: ${OPENIM_CONFIG_DIR}/config.yaml"

    # 3. Create and install the systemd unit files
    for service in "${OPENIM_RPC_SERVICE_LISTARIES[@]}"; do
        echo ${LINUX_PASSWORD} | sudo -S bash -c \
            "SERVER_NAME=${service} ./scripts/genconfig.sh ${ENV_FILE} deployments/templates/openim.service > ${SYSTEM_FILE_PATHS[$service]}"
        openim::log::status "${service} systemd file: ${SYSTEM_FILE_PATHS[$service]}"
    done

    # 4. Start the openim-rpc services
    openim::common::sudo "systemctl daemon-reload"
    for service in "${OPENIM_RPC_SERVICE_LISTARIES[@]}"; do
        openim::common::sudo "systemctl restart ${service}"
        openim::common::sudo "systemctl enable ${service}"
    done
    openim::rpc::status || return 1
    openim::rpc::info

    openim::log::info "install openim-rpc successfully"
    popd
}

# Unload
function openim::rpc::uninstall() {
    set +o errexit
    for service in "${OPENIM_RPC_SERVICE_LISTARIES[@]}"; do
        openim::common::sudo "systemctl stop ${service}"
        openim::common::sudo "systemctl disable ${service}"
        openim::common::sudo "rm -f ${OPENIM_INSTALL_DIR}/${service}"
        openim::common::sudo "rm -f ${OPENIM_CONFIG_DIR}/${service}.yaml"
        openim::common::sudo "rm -f ${SYSTEM_FILE_PATHS[$service]}"
    done
    set -o errexit
    openim::log::info "uninstall openim-rpc successfully"
}

# Status Check
function openim::rpc::status() {
    for service in "${OPENIM_RPC_SERVICE_LISTARIES[@]}"; do
        # Check the running status of the ${service}. If active (running) is displayed, the ${service} is started successfully.
        systemctl status ${service}|grep -q 'active' || {
            openim::log::error "${service} failed to start, maybe not installed properly"
            return 1
        }

        # The listening port is hardcoded in the configuration file
        if echo | telnet ${OPENIM_MSGGATEWAY_HOST} ${OPENIM_RPC_PORT_LISTARIES[@]} 2>&1|grep refused &>/dev/null;then
            openim::log::error "cannot access health check port, ${service} maybe not startup"
            return 1
        fi
    done
}

if [[ "$*" =~ openim::rpc:: ]];then
    eval $*
fi
