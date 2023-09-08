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

SERVER_NAME="openim-msggateway"

function openim::msggateway::start()
{
    openim::log::info "Start OpenIM Msggateway, binary root: ${SERVER_NAME}"
    openim::log::status "Start OpenIM Msggateway, path: ${OPENIM_MSGGATEWAY_BINARY}"

    openim::util::stop_services_with_name ${OPENIM_MSGGATEWAY_BINARY}

    # OpenIM message gateway service port
    OPENIM_MESSAGE_GATEWAY_PORTS=$(openim::util::list-to-string ${OPENIM_MESSAGE_GATEWAY_PORT} )
    read -a OPENIM_MSGGATEWAY_PORTS_ARRAY <<< ${OPENIM_MESSAGE_GATEWAY_PORTS}
    openim::util::stop_services_on_ports ${OPENIM_MSGGATEWAY_PORTS_ARRAY[*]}
    # OpenIM WS port
    OPENIM_WS_PORTS=$(openim::util::list-to-string ${OPENIM_WS_PORT} )
    read -a OPENIM_WS_PORTS_ARRAY <<< ${OPENIM_WS_PORTS}
    
    # Message Gateway Prometheus port of the service
    MSG_GATEWAY_PROM_PORTS=$(openim::util::list-to-string ${MSG_GATEWAY_PROM_PORT} )
    read -a MSG_GATEWAY_PROM_PORTS_ARRAY <<< ${MSG_GATEWAY_PROM_PORTS}

    openim::log::status "OpenIM Mssage Getway ports: ${OPENIM_MESSAGE_GATEWAY_PORTS[*]}"
    openim::log::status "OpenIM WS ports: ${OPENIM_WS_PORTS[*]}"
    openim::log::status "OpenIM Prometheus ports: ${MSG_GATEWAY_PROM_PORTS[*]}"

    openim::log::status "OpenIM Msggateway config path: ${OPENIM_MSGGATEWAY_CONFIG}"

    if [ ${#OPENIM_MSGGATEWAY_PORTS_ARRAY[@]} -ne ${#OPENIM_WS_PORTS_ARRAY[@]} ]; then
        openim::log::error_exit "ws_ports does not match push_rpc_ports or prome_ports in quantity!!!"
    fi

    for ((i = 0; i < ${#OPENIM_WS_PORTS_ARRAY[@]}; i++)); do
        openim::log::info "start push process, port: ${OPENIM_MSGGATEWAY_PORTS_ARRAY[$i]}, prometheus port: ${MSG_GATEWAY_PROM_PORTS_ARRAY[$i]}"

        PROMETHEUS_PORT_OPTION=""
        if [[ -n "${MSG_GATEWAY_PROM_PORTS_ARRAY[$i]}" ]]; then
            PROMETHEUS_PORT_OPTION="--prometheus_port ${MSG_GATEWAY_PROM_PORTS_ARRAY[$i]}"
        fi

        nohup ${OPENIM_MSGGATEWAY_BINARY} --port ${OPENIM_MSGGATEWAY_PORTS_ARRAY[$i]} --ws_port ${OPENIM_WS_PORTS_ARRAY[$i]} $PROMETHEUS_PORT_OPTION -c ${OPENIM_MSGGATEWAY_CONFIG} >> ${LOG_FILE} 2>&1 &
    done

    openim::util::check_process_names ${SERVER_NAME}
}

###################################### Linux Systemd ######################################
SYSTEM_FILE_PATH="/etc/systemd/system/${SERVER_NAME}.service"

# Print the necessary information after installation
function openim::msggateway::info() {
cat << EOF
openim-msggateway listen on: ${OPENIM_MSGGATEWAY_HOST}
EOF
}

# install openim-msggateway
function openim::msggateway::install()
{
  pushd "${OPENIM_ROOT}"

  # 1. Build openim-msggateway
  make build BINS=${SERVER_NAME}
  openim::common::sudo "cp ${OPENIM_OUTPUT_HOSTBIN}/${SERVER_NAME} ${OPENIM_INSTALL_DIR}/bin"

  openim::log::status "${SERVER_NAME} binary: ${OPENIM_INSTALL_DIR}/bin/${SERVER_NAME}"

  # 2. Generate and install the openim-msggateway configuration file (openim-msggateway.yaml)
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} deployments/templates/${SERVER_NAME}.yaml > ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"
  openim::log::status "${SERVER_NAME} config file: ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"

  # 3. Create and install the ${SERVER_NAME} systemd unit file
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} deployments/templates/init/${SERVER_NAME}.service > ${SYSTEM_FILE_PATH}"
  openim::log::status "${SERVER_NAME} systemd file: ${SYSTEM_FILE_PATH}"

  # 4. Start the openim-msggateway service
  openim::common::sudo "systemctl daemon-reload"
  openim::common::sudo "systemctl restart ${SERVER_NAME}"
  openim::common::sudo "systemctl enable ${SERVER_NAME}"
  openim::msggateway::status || return 1
  openim::msggateway::info

  openim::log::info "install ${SERVER_NAME} successfully"
  popd
}


# Unload
function openim::msggateway::uninstall()
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
function openim::msggateway::status()
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

if [[ "$*" =~ openim::msggateway:: ]];then
  eval $*
fi
