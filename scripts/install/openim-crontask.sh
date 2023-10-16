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
# OpenIM CronTask Control Script
# 
# Description:
# This script provides a control interface for the OpenIM CronTask service within a Linux environment. It supports two installation methods: installation via function calls to systemctl, and direct installation through background processes.
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
#    This will start the OpenIM CronTask directly through a background process.
#    Example: ./openim-crontask.sh openim::crontask::start
# 
# 2. Controlling through Functions for systemctl operations:
#    Specific operations like installation, uninstallation, and status check can be executed by passing the respective function name as an argument to the script.
#    Example: ./openim-crontask.sh openim::crontask::install
# 
# Note: Ensure that the appropriate permissions and environmental variables are set prior to script execution.
# 

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source "${OPENIM_ROOT}"/scripts/install/common.sh

SERVER_NAME="openim-crontask"

function openim::crontask::start()
{
    openim::log::info "Start OpenIM Cron, binary root: ${SERVER_NAME}"
    openim::log::status "Start OpenIM Cron, path: ${OPENIM_CRONTASK_BINARY}"

    openim::util::stop_services_with_name ${OPENIM_CRONTASK_BINARY}

    openim::log::status "start cron_task process, path: ${OPENIM_CRONTASK_BINARY}"
    nohup ${OPENIM_CRONTASK_BINARY} >> ${LOG_FILE} 2>&1 &
    openim::util::check_process_names ${SERVER_NAME}
}

###################################### Linux Systemd ######################################
SYSTEM_FILE_PATH="/etc/systemd/system/${SERVER_NAME}.service"

# Print the necessary information after installation
function openim::crontask::info() {
cat << EOF
openim-crontask listen on: ${OPENIM_CRONTASK_HOST}
EOF
}

# install openim-crontask
function openim::crontask::install()
{
  pushd "${OPENIM_ROOT}"

  # 1. Build openim-crontask
  make build BINS=${SERVER_NAME}
  openim::common::sudo "cp ${OPENIM_OUTPUT_HOSTBIN}/${SERVER_NAME} ${OPENIM_INSTALL_DIR}/bin"

  openim::log::status "${SERVER_NAME} binary: ${OPENIM_INSTALL_DIR}/bin/${SERVER_NAME}"

  # 2. Generate and install the openim-crontask configuration file (openim-crontask.yaml)
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} deployments/templates/${SERVER_NAME}.yaml > ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"
  openim::log::status "${SERVER_NAME} config file: ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"

  # 3. Create and install the ${SERVER_NAME} systemd unit file
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} deployments/templates/init/${SERVER_NAME}.service > ${SYSTEM_FILE_PATH}"
  openim::log::status "${SERVER_NAME} systemd file: ${SYSTEM_FILE_PATH}"

  # 4. Start the openim-crontask service
  openim::common::sudo "systemctl daemon-reload"
  openim::common::sudo "systemctl restart ${SERVER_NAME}"
  openim::common::sudo "systemctl enable ${SERVER_NAME}"
  openim::crontask::status || return 1
  openim::crontask::info

  openim::log::info "install ${SERVER_NAME} successfully"
  popd
}


# Unload
function openim::crontask::uninstall()
{
  set +o errexit
  openim::common::sudo "systemctl stop ${SERVER_NAME}"
  openim::common::sudo "systemctl disable ${SERVER_NAME}"
  openim::common::sudo "rm -f ${OPENIM_INSTALL_DIR}/bin/${SERVER_NAME}"
  openim::common::sudo "rm -f ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"
  openim::common::sudo "rm -f /etc/systemd/system/${SERVER_NAME}.service"
  set -o errexit
  openim::log::info "uninstall ${SERVER_NAME} successfully"
}

# Status Check
function openim::crontask::status()
{
  # Check the running status of the ${SERVER_NAME}. If active (running) is displayed, the ${SERVER_NAME} is started successfully.
  systemctl status ${SERVER_NAME}|grep -q 'active' || {
    openim::log::error "${SERVER_NAME} failed to start, maybe not installed properly"
    return 1
  }

  # The listening port is hardcode in the configuration file
  if echo | telnet 127.0.0.1 7070 2>&1|grep refused &>/dev/null;then
    openim::log::error "cannot access health check port, ${SERVER_NAME} maybe not startup"
    return 1
  fi
}

if [[ "$*" =~ openim::crontask:: ]];then
  eval $*
fi
