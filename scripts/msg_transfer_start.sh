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

chmod +x build_all_service.sh
./build_all_service.sh

echo -e "${YELLOW_PREFIX}=======>SCRIPTS_ROOT=$SCRIPTS_ROOT${COLOR_SUFFIX}"
echo -e "${YELLOW_PREFIX}=======>OPENIM_ROOT=$OPENIM_ROOT${COLOR_SUFFIX}"
echo -e "${YELLOW_PREFIX}=======>pwd=$PWD${COLOR_SUFFIX}"

bin_dir="$BIN_DIR"
logs_dir="$OPENIM_ROOT/logs"
sdk_db_dir="$OPENIM_ROOT/sdk/db/"

cd $OPENIM_ROOT

list1=$(cat $config_path | grep messageTransferPrometheusPort | awk -F '[:]' '{print $NF}')
list_to_string $list1
prome_ports=($ports_array)

#Check if the service exists
#If it is exists,kill this process
check=`ps aux | grep -w ./${openim_msgtransfer} | grep -v grep| wc -l`
if [ $check -ge 1 ]
then
oldPid=`ps aux | grep -w ./${openim_msgtransfer} | grep -v grep|awk '{print $2}'`
 kill -9 $oldPid
fi
#Waiting port recycling
sleep 1

cd ${msg_transfer_binary_root}
for ((i = 0; i < ${msg_transfer_service_num}; i++)); do
      prome_port=${prome_ports[$i]}
      cmd="nohup ./${openim_msgtransfer}  --config_folder_path ${configfile_path} "
      if [ $prome_port != "" ]; then
        cmd="$cmd --prometheus_port $prome_port  --config_folder_path ${configfile_path} "
      fi
      echo "==========================start msg_transfer server===========================">>$OPENIM_ROOT/logs/openIM.log
      $cmd >>$OPENIM_ROOT/logs/openIM.log 2>&1 &
done

#Check launched service process
check=`ps aux | grep -w ./${openim_msgtransfer} | grep -v grep| wc -l`
if [ $check -ge 1 ]
then
newPid=`ps aux | grep -w ./${openim_msgtransfer} | grep -v grep|awk '{print $2}'`
allPorts=""
    echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS "${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${openim_msgtransfer}${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${newPid}${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${allPorts}${COLOR_SUFFIX}
else
    echo -e ${BACKGROUND_GREEN}${openim_msgtransfer}${COLOR_SUFFIX}${RED_PREFIX}"\n SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
