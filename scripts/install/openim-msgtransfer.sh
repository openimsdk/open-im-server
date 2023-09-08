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
set +o nounset
set -o pipefail

ulimit -n 200000

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source "${OPENIM_ROOT}"/scripts/install/common.sh

SERVER_NAME="openim-msgtransfer"

function openim::msgtransfer::start()
{
    openim::log::info "Start OpenIM Msggateway, binary root: ${SERVER_NAME}"
    openim::log::status "Start OpenIM Msggateway, path: ${OPENIM_MSGTRANSFER_BINARY}"

    openim::util::stop_services_with_name ${OPENIM_MSGTRANSFER_BINARY}

    # Message Transfer Prometheus port list
    MSG_TRANSFER_PROM_PORTS=(openim::util::list-to-string ${MSG_TRANSFER_PROM_PORT} )

    openim::log::status "OpenIM Prometheus ports: ${MSG_TRANSFER_PROM_PORTS[*]}"

    openim::log::status "OpenIM Msggateway config path: ${OPENIM_MSGTRANSFER_CONFIG}"

    openim::log::info "openim maggateway num: ${OPENIM_MSGGATEWAY_NUM}"

    if [ "${OPENIM_MSGGATEWAY_NUM}" -lt 1 ]; then
    opeim::log::error_exit "OPENIM_MSGGATEWAY_NUM must be greater than 0"
    fi

    if [ ${OPENIM_MSGGATEWAY_NUM} -ne $((${#MSG_TRANSFER_PROM_PORTS[@]} - 1)) ]; then
    openim::log::error_exit "OPENIM_MSGGATEWAY_NUM must be equal to the number of MSG_TRANSFER_PROM_PORTS"
    fi

    for (( i=1; i<=$OPENIM_MSGGATEWAY_NUM; i++ )) do
      openim::log::info "prometheus port: ${MSG_TRANSFER_PROM_PORTS[$i]}"
      PROMETHEUS_PORT_OPTION=""
      if [[ -n "${OPENIM_PROMETHEUS_PORTS[$i]}" ]]; then
        PROMETHEUS_PORT_OPTION="--prometheus_port ${OPENIM_PROMETHEUS_PORTS[$i]}"
      fi
      nohup ${OPENIM_MSGTRANSFER_BINARY} ${PROMETHEUS_PORT_OPTION} -c ${OPENIM_MSGTRANSFER_CONFIG} >> ${LOG_FILE} 2>&1 &
    done

    openim::util::check_process_names  "${OPENIM_OUTPUT_HOSTBIN}/${SERVER_NAME}"
}

function openim::msgtransfer::check()
{
    PIDS=$(pgrep -f "${OPENIM_OUTPUT_HOSTBIN}/openim-msgtransfer")

    NUM_PROCESSES=$(echo "$PIDS" | wc -l)
    # NUM_PROCESSES=$(($NUM_PROCESSES - 1))

    if [ "$NUM_PROCESSES" -eq "$OPENIM_MSGGATEWAY_NUM" ]; then
      openim::log::info "Found $OPENIM_MSGGATEWAY_NUM processes named $OPENIM_OUTPUT_HOSTBIN"
      for PID in $PIDS; do
        ps -p $PID -o pid,cmd
      done
    else
      openim::log::error_exit "Expected $OPENIM_MSGGATEWAY_NUM openim msgtransfer processes, but found $NUM_PROCESSES msgtransfer processes."
    fi
}

###################################### Linux Systemd ######################################
SYSTEM_FILE_PATH="/etc/systemd/system/${SERVER_NAME}.service"

# Print the necessary information after installation
function openim::msgtransfer::info() {
cat << EOF
openim-msgtransfer listen on: ${OPENIM_MSGTRANSFER_HOST}
EOF
}

# install openim-msgtransfer
function openim::msgtransfer::install()
{
  pushd "${OPENIM_ROOT}"

  # 1. Build openim-msgtransfer
  make build BINS=${SERVER_NAME}
  openim::common::sudo "cp ${OPENIM_OUTPUT_HOSTBIN}/${SERVER_NAME} ${OPENIM_INSTALL_DIR}/bin"

  openim::log::status "${SERVER_NAME} binary: ${OPENIM_INSTALL_DIR}/bin/${SERVER_NAME}"

  # 2. Generate and install the openim-msgtransfer configuration file (openim-msgtransfer.yaml)
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} deployments/templates/${SERVER_NAME}.yaml > ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"
  openim::log::status "${SERVER_NAME} config file: ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"

  # 3. Create and install the ${SERVER_NAME} systemd unit file
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} deployments/templates/init/${SERVER_NAME}.service > ${SYSTEM_FILE_PATH}"
  openim::log::status "${SERVER_NAME} systemd file: ${SYSTEM_FILE_PATH}"

  # 4. Start the openim-msgtransfer service
  openim::common::sudo "systemctl daemon-reload"
  openim::common::sudo "systemctl restart ${SERVER_NAME}"
  openim::common::sudo "systemctl enable ${SERVER_NAME}"
  openim::msgtransfer::status || return 1
  openim::msgtransfer::info

  openim::log::info "install ${SERVER_NAME} successfully"
  popd
}


# Unload
function openim::msgtransfer::uninstall()
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
function openim::msgtransfer::status()
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

if [[ "$*" =~ openim::msgtransfer:: ]];then
  eval $*
fi
