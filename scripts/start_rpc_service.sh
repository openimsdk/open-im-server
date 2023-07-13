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

echo -e "${BACKGROUND_GREEN}${CYAN_PREFIX}=======>SCRIPTS_ROOT=$SCRIPTS_ROOT${COLOR_SUFFIX}"
echo -e "${BACKGROUND_GREEN}${CYAN_PREFIX}=======>OPENIM_ROOT=$OPENIM_ROOT${COLOR_SUFFIX}"
echo -e "${BACKGROUND_GREEN}${CYAN_PREFIX}=======>pwd=$PWD${COLOR_SUFFIX}"

bin_dir="$BIN_DIR"
logs_dir="$OPENIM_ROOT/logs"
sdk_db_dir="$OPENIM_ROOT/sdk/db/"

#service filename
service_filename=(
  #api
  openim-api
  #rpc
  openim-rpc-user
  openim-rpc-friend
  openim-rpc-group
  openim-rpc-auth
  ${msg_name}
  openim-rpc-conversation
  openim-rpc-third
)

#service config port name
service_port_name=(
  #api port name
  openImApiPort
  #rpc port name
  openImUserPort
  openImFriendPort
  openImGroupPort
  openImAuthPort
  openImMessagePort
  openImConversationPort
  openImThirdPort
)

service_prometheus_port_name=(
  #api port name
  openImApiPort
  #rpc port name
  userPrometheusPort
  friendPrometheusPort
  groupPrometheusPort
  authPrometheusPort
  messagePrometheusPort
  conversationPrometheusPort
  thirdPrometheusPort
)

for ((i = 0; i < ${#service_filename[*]}; i++)); do
  #Check whether the service exists
  service_name="ps -aux |grep -w ${service_filename[$i]} |grep -v grep"
  count="${service_name}| wc -l"

  if [ $(eval ${count}) -gt 0 ]; then
    pid="${service_name}| awk '{print \$2}'"
    echo  "${service_filename[$i]} service has been started,pid:$(eval $pid)"
    echo  "killing the service ${service_filename[$i]} pid:$(eval $pid)"
    #kill the service that existed
    kill -9 $(eval $pid)
    sleep 0.5
  fi
  cd $OPENIM_ROOT
  cd $BIN_DIR
  # Get the rpc port in the configuration file
  portList=$(cat $config_path | grep ${service_port_name[$i]} | awk -F '[:]' '{print $NF}')
  list_to_string ${portList}
  service_ports=($ports_array)

  portList2=$(cat $config_path | grep ${service_prometheus_port_name[$i]} | awk -F '[:]' '{print $NF}')
  list_to_string $portList2
  prome_ports=($ports_array)
  #Start related rpc services based on the number of ports
  for ((j = 0; j < ${#service_ports[*]}; j++)); do
    #Start the service in the background
    if [ -z "${prome_ports[$j]}" ]; then
      cmd="./${service_filename[$i]} --port ${service_ports[$j]} --config_folder_path ${configfile_path} "
    else
      cmd="./${service_filename[$i]} --port ${service_ports[$j]} --prometheus_port ${prome_ports[$j]}  --config_folder_path ${configfile_path} "
    fi
    if [ $i -eq 0 -o $i -eq 1 ]; then
      cmd="./${service_filename[$i]} --port ${service_ports[$j]}"
    fi
    echo $cmd
    echo "=====================start ${service_filename[$i]}======================">>$OPENIM_ROOT/logs/openIM.log
    nohup $cmd >>$OPENIM_ROOT/logs/openIM.log 2>&1 &
    sleep 1
    pid="netstat -ntlp|grep $j |awk '{printf \$7}'|cut -d/ -f1"
    echo -e "${GREEN_PREFIX}${service_filename[$i]} start success,port number:${service_ports[$j]} pid:$(eval $pid)$COLOR_SUFFIX"
  done
done
