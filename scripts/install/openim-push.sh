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
# OpenIM Push Control Script
# 
# Description:
# This script provides a control interface for the OpenIM Push service within a Linux environment. It supports two installation methods: installation via function calls to systemctl, and direct installation through background processes.
# 
# Features:
# 1. Robust error handling leveraging Bash built-ins such as 'errexit', 'nounset', and 'pipefail'.
# 2. Capability to source common utility functions and configurations, ensuring environmental consistency.
# 3. Comprehensive logging tools, offering clear operational insights.
# 4. Support for creating, managing, and interacting with Linux systemd services.
# 5. Mechanisms to verify the successful running of the service.
#
# Usage:
# 1. Direct Script Execution:
#    This will start the OpenIM push directly through a background process.
#    Example: ./openim-push.sh
# 
# 2. Controlling through Functions for systemctl operations:
#    Specific operations like installation, uninstallation, and status check can be executed by passing the respective function name as an argument to the script.
#    Example: ./openim-push.sh openim::push::install
#
# ENVIRONMENT VARIABLES:
# export OPENIM_PUSH_BINARY="8080 8081 8082"
# export OPENIM_PUSH_PORT="9090 9091 9092"
#
# Note: Ensure that the appropriate permissions and environmental variables are set prior to script execution.
# 
set -o errexit
set +o nounset
set -o pipefail

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source "${OPENIM_ROOT}"/scripts/install/common.sh

SERVER_NAME="openim-push"

function openim::push::start() {
    openim::log::status "Start OpenIM Push, binary root: ${SERVER_NAME}"
    openim::log::info "Start OpenIM Push, path: ${OPENIM_PUSH_BINARY}"

    openim::log::status "prepare start push process, path: ${OPENIM_PUSH_BINARY}"
    openim::log::status "prepare start push process, port: ${OPENIM_PUSH_PORT}, prometheus port: ${PUSH_PROM_PORT}"

    OPENIM_PUSH_PORTS_ARRAY=$(openim::util::list-to-string ${OPENIM_PUSH_PORT} )
    PUSH_PROM_PORTS_ARRAY=$(openim::util::list-to-string ${PUSH_PROM_PORT} )

    openim::util::stop_services_with_name ${SERVER_NAME}

    openim::log::status "push port list: ${OPENIM_PUSH_PORTS_ARRAY[@]}"
    openim::log::status "prometheus port list: ${PUSH_PROM_PORTS_ARRAY[@]}"

    if [ ${#OPENIM_PUSH_PORTS_ARRAY[@]} -ne ${#PUSH_PROM_PORTS_ARRAY[@]} ]; then
        openim::log::error_exit "The length of the two port lists is different!"
    fi

    for (( i=0; i<${#OPENIM_PUSH_PORTS_ARRAY[@]}; i++ )); do
        openim::log::info "start push process, port: ${OPENIM_PUSH_PORTS_ARRAY[$i]}, prometheus port: ${PUSH_PROM_PORTS_ARRAY[$i]}"
        nohup ${OPENIM_PUSH_BINARY} --port ${OPENIM_PUSH_PORTS_ARRAY[$i]} -c ${OPENIM_PUSH_CONFIG} --prometheus_port ${PUSH_PROM_PORTS_ARRAY[$i]} >> ${LOG_FILE} 2>&1 &
    done

    openim::util::check_process_names ${SERVER_NAME}
}

###################################### Linux Systemd ######################################
SYSTEM_FILE_PATH="/etc/systemd/system/${SERVER_NAME}.service"

# Print the necessary information after installation
function openim::push::info() {
cat << EOF
openim-push listen on: ${OPENIM_PUSH_HOST}:${OPENIM_PUSH_PORT}
EOF
}

# install openim-push
function openim::push::install() {
  pushd "${OPENIM_ROOT}"

  # 1. Build openim-push
  make build BINS=${SERVER_NAME}
  openim::common::sudo "cp -r ${OPENIM_OUTPUT_HOSTBIN}/${SERVER_NAME} ${OPENIM_INSTALL_DIR}/${SERVER_NAME}"
  openim::log::status "${SERVER_NAME} binary: ${OPENIM_INSTALL_DIR}/${SERVER_NAME}/${SERVER_NAME}"

  # 2. Generate and install the openim-push configuration file (config)
  openim::log::status "${SERVER_NAME} config file: ${OPENIM_CONFIG_DIR}/config.yaml"

  # 3. Create and install the ${SERVER_NAME} systemd unit file
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "SERVER_NAME=${SERVER_NAME} ./scripts/genconfig.sh ${ENV_FILE} deployments/templates/openim.service > ${SYSTEM_FILE_PATH}"
  openim::log::status "${SERVER_NAME} systemd file: ${SYSTEM_FILE_PATH}"

  # 4. Start the openim-push service
  openim::common::sudo "systemctl daemon-reload"
  openim::common::sudo "systemctl restart ${SERVER_NAME}"
  openim::common::sudo "systemctl enable ${SERVER_NAME}"
  openim::push::status || return 1
  openim::push::info

  openim::log::info "install ${SERVER_NAME} successfully"
  popd
}

# Unload
function openim::push::uninstall() {
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
function openim::push::status() {
  # Check the running status of the ${SERVER_NAME}. If active (running) is displayed, the ${SERVER_NAME} is started successfully.
  systemctl status ${SERVER_NAME}|grep -q 'active' || {
    openim::log::error "${SERVER_NAME} failed to start, maybe not installed properly"
    return 1
  }

  # The listening port is hardcode in the configuration file
  if echo | telnet ${OPENIM_MSGGATEWAY_HOST} ${OPENIM_PUSH_PORT} 2>&1|grep refused &>/dev/null;then  # Assuming a different port for push
    openim::log::error "cannot access health check port, ${SERVER_NAME} maybe not startup"
    return 1
  fi
}

if [[ "$*" =~ openim::push:: ]];then
  eval $*
fi
