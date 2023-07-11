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

service_port_name=(
  openImWsPort
  openImApiPort
  openImUserPort
  openImFriendPort
  openImMessagePort
  openImMessageGatewayPort
  openImGroupPort
  openImAuthPort
  openImPushPort
  openImConversationPort
  openImThirdPort
)
for i in ${service_port_name[*]}; do
  list=$(cat $config_path | grep -w ${i} | awk -F '[:]' '{print $NF}')
  list_to_string $list
  for j in ${ports_array}; do
    port=$(ss -tunlp| grep openim | awk '{print $5}' | grep -w ${j} | awk -F '[:]' '{print $NF}')
    if [[ ${port} -ne ${j} ]]; then
      echo -e ${BACKGROUND_GREEN}${i}${COLOR_SUFFIX}${RED_PREFIX}" service does not start normally,not initiated port is "${COLOR_SUFFIX}${BACKGROUND_GREEN}${j}${COLOR_SUFFIX}
      echo -e ${RED_PREFIX}"please check $OPENIM_ROOT/logs/openIM.log "${COLOR_SUFFIX}
      exit -1
    else
      echo -e ${j}${GREEN_PREFIX}" port has been listening,belongs service is "${i}${COLOR_SUFFIX}
    fi
  done
done

#Check launched service process
check=$(ps aux | grep -w ./${openim_msgtransfer} | grep -v grep | wc -l)
if [ $check -eq ${msg_transfer_service_num} ]; then
  echo -e ${GREEN_PREFIX}"none  port has been listening,belongs service is openImMsgTransfer"${COLOR_SUFFIX}
else
  echo	 $check ${msg_transfer_service_num}
  echo -e ${RED_PREFIX}"openImMsgTransfer service does not start normally, num err"${COLOR_SUFFIX}
        echo -e ${RED_PREFIX}"please check $OPENIM_ROOT/logs/openIM.log "${COLOR_SUFFIX}
      exit -1
fi


check=$(ps aux | grep -w ./${cron_task_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  echo -e ${GREEN_PREFIX}"none  port has been listening,belongs service is openImCronTask"${COLOR_SUFFIX}
else
  echo -e ${RED_PREFIX}"cron_task_name service does not start normally"${COLOR_SUFFIX}
        echo -e ${RED_PREFIX}"please check $OPENIM_ROOT/logs/openIM.log "${COLOR_SUFFIX}
      exit -1
fi

echo -e ${BACKGROUND_GREEN}"all services launch success"${COLOR_SUFFIX}
