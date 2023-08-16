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

SERVER_NAME="openim-rpc"
readonly OPENIM_RPC_CONFIG=${OPENIM_ROOT}/config

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

readonly OPENIM_API_SERVICE_TARGETS=(
  openim-api
)
readonly OPENIM_API_SERVICE_LISTARIES=("${OPENIM_API_SERVICE_TARGETS[@]##*/}")

readonly OPENIM_API_PORT_TARGETS=(
  ${API_OPENIM_PORT}
)
readonly OPENIM_API_PORT_LISTARIES=("${OPENIM_API_PORT_TARGETS[@]##*/}")

readonly OPENIM_RPC_ALL_NAME_TARGETS=(
  "${OPENIM_API_SERVICE_TARGETS[@]}"
  "${OPENIM_RPC_SERVICE_TARGETS[@]}"
)
readonly OPENIM_RPC_ALL_NAME_LISTARIES=("${OPENIM_RPC_ALL_NAME_TARGETS[@]##*/}")

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

readonly OPENIM_RPC_ALL_PORT_TARGETS=(
  "${OPENIM_API_PORT_TARGETS[@]}"
  "${OPENIM_RPC_PORT_TARGETS[@]}"
)
readonly OPENIM_RPC_ALL_PORT_LISTARIES=("${OPENIM_RPC_ALL_PORT_TARGETS[@]##*/}")

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

readonly OPENIM_RPC_ALL_PROM_PORT_TARGETS=(
  "${OPENIM_API_PORT_TARGETS[@]}"
  "${OPENIM_RPC_PROM_PORT_TARGETS[@]}"
)
readonly OPENIM_RPC_ALL_PROM_PORT_LISTARIES=("${OPENIM_RPC_ALL_PROM_PORT_TARGETS[@]##*/}")

echo "OPENIM_RPC_ALL_NAME_TARGETS: ${OPENIM_RPC_ALL_NAME_TARGETS[@]}"
echo "OPENIM_RPC_ALL_PROM_PORT_TARGETS: ${OPENIM_RPC_ALL_PROM_PORT_TARGETS[@]}"
echo "OPENIM_RPC_ALL_PORT_TARGETS: ${OPENIM_RPC_ALL_PORT_TARGETS[@]}"

openim::log::info "Starting ${SERVER_NAME} ..."

printf "+------------------------+-------+-----------------+\n"
printf "| Service Name           | Port  | Prometheus Port |\n"
printf "+------------------------+-------+-----------------+\n"

length=${#OPENIM_RPC_ALL_NAME_LISTARIES[@]}

for ((i=0; i<$length; i++)); do
  printf "| %-22s | %-5s | %-15s |\n" "${OPENIM_RPC_ALL_NAME_LISTARIES[$i]}" "${OPENIM_RPC_ALL_PORT_LISTARIES[$i]}" "${OPENIM_RPC_ALL_PROM_PORT_LISTARIES[$i]}"
  printf "+------------------------+-------+-----------------+\n"
done

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

# start all rpc services
for ((i = 0; i < ${#OPENIM_RPC_ALL_NAME_TARGETS[*]}; i++)); do
  openim::util::stop_services_with_name ${OPENIM_RPC_ALL_NAME_TARGETS[$i]}
  openim::log::info "OpenIM ${OPENIM_RPC_ALL_NAME_TARGETS[$i]} config path: ${OPENIM_RPC_CONFIG}"

  # Get the service and Prometheus ports.
  OPENIM_RPC_SERVICE_PORTS=( $(openim::util::list-to-string ${OPENIM_RPC_ALL_PORT_LISTARIES[$i]}) )
  OPENIM_RPC_PROM_PORTS=( $(openim::util::list-to-string ${OPENIM_RPC_ALL_PROM_PORT_LISTARIES[$i]}) )

  for ((j = 0; j < ${#OPENIM_RPC_SERVICE_PORTS[@]}; j++)); do
    openim::log::info "Starting ${OPENIM_RPC_ALL_NAME_TARGETS[$i]} service, port: ${OPENIM_RPC_SERVICE_PORTS[j]}, prometheus port: ${OPENIM_RPC_PROM_PORTS[j]}, binary root: ${OPENIM_OUTPUT_HOSTBIN}/${OPENIM_RPC_ALL_NAME_TARGETS[$i]}"
    openim::rpc::start_service "${OPENIM_RPC_ALL_NAME_TARGETS[$i]}" "${OPENIM_RPC_SERVICE_PORTS[j]}" "${OPENIM_RPC_PROM_PORTS[j]}"
  done
done


openim::util::check_ports ${OPENIM_RPC_ALL_PORT_TARGETS[@]}
openim::util::check_ports ${OPENIM_RPC_ALL_PROM_PORT_TARGETS[@]}

# ---
# for ((i = 0; i < ${#service_filename[*]}; i++)); do
#   #Check whether the service exists
#   service_name="ps |grep -w ${service_filename[$i]} |grep -v grep"
#   count="${service_name}| wc -l"

#   if [ $(eval ${count}) -gt 0 ]; then
#     pid="${service_name}| awk '{print \$2}'"
#     echo  "${service_filename[$i]} service has been started,pid:$(eval $pid)"
#     echo  "killing the service ${service_filename[$i]} pid:$(eval $pid)"
#     #kill the service that existed
#     kill -9 $(eval $pid)
#     sleep 0.5
#   fi
#   cd $OPENIM_ROOT
#   cd $BIN_DIR
#   # Get the rpc port in the configuration file
#   portList=$(cat $config_path | grep ${service_port_name[$i]} | awk -F '[:]' '{print $NF}')
#   openim::util::list-to-string ${portList}
#   service_ports=($ports_array)

#   portList2=$(cat $config_path | grep ${service_prometheus_port_name[$i]} | awk -F '[:]' '{print $NF}')
#   openim::util::list-to-string $portList2
#   prome_ports=($ports_array)
#   #Start related rpc services based on the number of ports
#   for ((j = 0; j < ${#service_ports[*]}; j++)); do
#     #Start the service in the background
#     if [ -z "${prome_ports[$j]}" ]; then
#       cmd="./${service_filename[$i]} --port ${service_ports[$j]} -c ${configfile_path} "
#     else
#       cmd="./${service_filename[$i]} --port ${service_ports[$j]} --prometheus_port ${prome_ports[$j]}  -c ${configfile_path} "
#     fi
#     if [ $i -eq 0 -o $i -eq 1 ]; then
#       cmd="./${service_filename[$i]} --port ${service_ports[$j]}"
#     fi
#     echo "=====================start ${service_filename[$i]}======================">>$OPENIM_ROOT/logs/openIM.log
#     nohup $cmd >>$OPENIM_ROOT/logs/openIM.log 2>&1 &
#     sleep 1
#     pid="netstat -ntlp|grep $j |awk '{printf \$7}'|cut -d/ -f1"
#     echo -e "${GREEN_PREFIX}${service_filename[$i]} start success,port number:${service_ports[$j]} pid:$(eval $pid)$COLOR_SUFFIX"
#   done
# done
