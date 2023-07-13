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

#Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/style_info.sh
source $SCRIPTS_ROOT/path_info.sh
source $SCRIPTS_ROOT/function.sh

cd $SCRIPTS_ROOT

echo -e "${YELLOW_PREFIX}=======>SCRIPTS_ROOT=$SCRIPTS_ROOT${COLOR_SUFFIX}"
echo -e "${YELLOW_PREFIX}=======>OPENIM_ROOT=$OPENIM_ROOT${COLOR_SUFFIX}"
echo -e "${YELLOW_PREFIX}=======>pwd=$PWD${COLOR_SUFFIX}"

bin_dir="$BIN_DIR"
logs_dir="$OPENIM_ROOT/logs"
sdk_db_dir="$OPENIM_ROOT/sdk/db/"

ulimit -n 200000

list1=$(cat $config_path | grep openImMessageGatewayPort | awk -F '[:]' '{print $NF}')
list2=$(cat $config_path | grep openImWsPort | awk -F '[:]' '{print $NF}')
list3=$(cat $config_path | grep messageGatewayPrometheusPort | awk -F '[:]' '{print $NF}')
list_to_string $list1
rpc_ports=($ports_array)
list_to_string $list2
ws_ports=($ports_array)
list_to_string $list3
prome_ports=($ports_array)
if [ ${#rpc_ports[@]} -ne ${#ws_ports[@]} ]; then

  echo -e ${RED_PREFIX}"ws_ports does not match push_rpc_ports or prome_ports in quantity!!!"${COLOR_SUFFIX}
  exit -1

fi
#Check if the service exists
#If it is exists,kill this process
check=$(ps aux | grep -w ./${openim_msggateway} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps aux | grep -w ./${openim_msggateway} | grep -v grep | awk '{print $2}')
    kill -9 ${oldPid}
fi
#Waiting port recycling
sleep 1
cd ${msg_gateway_binary_root}
for ((i = 0; i < ${#ws_ports[@]}; i++)); do
  echo "==========================start msg_gateway server===========================">>$OPENIM_ROOT/logs/openIM.log
  nohup ./${openim_msggateway} --port ${rpc_ports[$i]} --ws_port ${ws_ports[$i]} --prometheus_port ${prome_ports[$i]} --config_folder_path ${configfile_path}  >>$OPENIM_ROOT/logs/openIM.log 2>&1 &
done

#Check launched service process
sleep 3
check=$(ps aux | grep -w ./${openim_msggateway} | grep -v grep | wc -l)
allPorts=""
if [ $check -ge 1 ]; then
  allNewPid=$(ps aux | grep -w ./${openim_msggateway} | grep -v grep | awk '{print $2}')
  for i in $allNewPid; do
    ports=$(netstat -netulp | grep -w ${i} | awk '{print $4}' | awk -F '[:]' '{print $NF}')
      allPorts=${allPorts}"$ports "
  done
  echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS"${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${openim_msggateway}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${allNewPid}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${allPorts}${COLOR_SUFFIX}
else
  echo -e ${BACKGROUND_GREEN}${openim_msggateway}${COLOR_SUFFIX}${RED_PREFIX}"\n SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
