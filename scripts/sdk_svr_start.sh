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
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
ulimit -n 200000

ws_address=$(cat $config_path | grep openImWsAddress | awk -F '[ ]' '{print $NF}')
api_address=$(cat $config_path | grep openImApiAddress | awk -F '[ ]' '{print $NF}')
list3=$(cat $config_path | grep openImSdkWsPort | awk -F '[:]' '{print $NF}')
logLevel=$(cat $config_path | grep remainLogLevel | awk -F '[:]' '{print $NF}')
list_to_string $list3
sdkws_ports=($ports_array)



#Check if the service exists
#If it is exists,kill this process
check=$(ps aux | grep -w ./${sdk_server_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps aux | grep -w ./${sdk_server_name} | grep -v grep | awk '{print $2}')
    kill -9 ${oldPid}
fi
#Waiting port recycling
sleep 1
cd ${sdk_server_binary_root}
  echo "==========================start js sdk server===========================">>../logs/openIM.log
  nohup ./${sdk_server_name}  -openIM_ws_address ${ws_address}  -sdk_ws_port ${sdkws_ports[0]} -openIM_api_address ${api_address} -openIM_log_level ${logLevel} >>../logs/openIM.log 2>&1 &

#Check launched service process
sleep 3
check=$(ps aux | grep -w ./${sdk_server_name} | grep -v grep | wc -l)
allPorts=""
if [ $check -ge 1 ]; then
  allNewPid=$(ps aux | grep -w ./${sdk_server_name} | grep -v grep | awk '{print $2}')
  for i in $allNewPid; do
    ports=$(netstat -netulp | grep -w ${i} | awk '{print $4}' | awk -F '[:]' '{print $NF}')
      allPorts=${allPorts}"$ports "
  done
  echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS "${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${sdk_server_name}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allNewPid}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
  echo -e ${YELLOW_PREFIX}${sdk_server_name}${COLOR_SUFFIX}${RED_PREFIX}"SERVICE START ERROR PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
