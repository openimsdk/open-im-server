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
[[ -z ${COMMON_SOURCED} ]] && source ${OPENIM_ROOT}/scripts/install/common.sh

SERVER_NAME="openim-push"

# openim::log::status "Start OpenIM Push, binary root: ${SERVER_NAME}"
# openim::log::info "Start OpenIM Push, path: ${OPENIM_PUSH_BINARY}"

# openim::util::stop_services_with_name ${SERVER_NAME}

# openim::log::status "start push process, path: ${OPENIM_PUSH_BINARY}"
# nohup ${OPENIM_PUSH_BINARY} >>${LOG_FILE} 2>&1 &
# openim::util::check_process_names ${SERVER_NAME}

###################################### Linux Systemd ######################################
SYSTEM_FILE_PATH="/etc/systemd/system/${SERVER_NAME}.service"

# Print the necessary information after installation
function openim::push::info() {
cat << EOF
openim-push listen on: ${OPENIM_PUSH_HOST}
EOF
}

# install openim-push
function openim::push::install()
{
  pushd ${OPENIM_ROOT}

  # 1. Build openim-push
  make build BINS=${SERVER_NAME}
  openim::common::sudo "cp ${OPENIM_OUTPUT_HOSTBIN}/${SERVER_NAME} ${OPENIM_INSTALL_DIR}/bin"

  openim::log::status "${SERVER_NAME} binary: ${OPENIM_INSTALL_DIR}/bin/${SERVER_NAME}"

  # 2. Generate and install the openim-push configuration file (openim-push.yaml)
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} deployments/templates/${SERVER_NAME}.yaml > ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"
  openim::log::status "${SERVER_NAME} config file: ${OPENIM_CONFIG_DIR}/${SERVER_NAME}.yaml"

  # 3. Create and install the ${SERVER_NAME} systemd unit file
  echo ${LINUX_PASSWORD} | sudo -S bash -c \
    "./scripts/genconfig.sh ${ENV_FILE} deployments/templates/init/${SERVER_NAME}.service > ${SYSTEM_FILE_PATH}"
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
function openim::push::uninstall()
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
function openim::push::status()
{
  # Check the running status of the ${SERVER_NAME}. If active (running) is displayed, the ${SERVER_NAME} is started successfully.
  systemctl status ${SERVER_NAME}|grep -q 'active' || {
    openim::log::error "${SERVER_NAME} failed to start, maybe not installed properly"
    return 1
  }

  # The listening port is hardcode in the configuration file
  if echo | telnet 127.0.0.1 7071 2>&1|grep refused &>/dev/null;then  # Assuming a different port for push
    openim::log::error "cannot access health check port, ${SERVER_NAME} maybe not startup"
    return 1
  fi
}

if [[ "$*" =~ ${SERVER_NAME}:: ]];then
  eval $*
fi



cd $SCRIPTS_ROOT

bin_dir="$BIN_DIR"
logs_dir="$OPENIM_ROOT/logs"

cd "$OPENIM_ROOT/scripts/"

list1=$(cat $config_path | grep openImPushPort | awk -F '[:]' '{print $NF}')
list2=$(cat $config_path | grep pushPrometheusPort | awk -F '[:]' '{print $NF}')
openim::util::list-to-string $list1
rpc_ports=($ports_array)
openim::util::list-to-string $list2
prome_ports=($ports_array)

#Check if the service exists
#If it is exists,kill this process
check=$(ps -aux | grep -w ./${push_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps -aux | grep -w ./${push_name} | grep -v grep | awk '{print $2}')
  kill -9 $oldPid
fi
#Waiting port recycling
sleep 1
cd ${push_binary_root}

for ((i = 0; i < ${#rpc_ports[@]}; i++)); do
  echo "==========================start push server===========================">>$OPENIM_ROOT/logs/openIM.log
  nohup ./${push_name} --port ${rpc_ports[$i]} --prometheus_port ${prome_ports[$i]} >>$OPENIM_ROOT/logs/openIM.log 2>&1 &
done

sleep 3
#Check launched service process
check=$(ps -aux | grep -w ./${push_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  newPid=$(ps -aux | grep -w ./${push_name} | grep -v grep | awk '{print $2}')
  ports=$(netstat -netulp | grep -w ${newPid} | awk '{print $4}' | awk -F '[:]' '{print $NF}')
  allPorts=""

  for i in $ports; do
    allPorts=${allPorts}"$i "
  done
  echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS "${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${push_name}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${newPid}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${allPorts}${COLOR_SUFFIX}
else
  echo -e ${BACKGROUND_GREEN}${push_name}${COLOR_SUFFIX}${RED_PREFIX}"\n SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
