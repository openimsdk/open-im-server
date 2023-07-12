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

#fixme This scripts is the total startup scripts
#fixme The full name of the shell scripts that needs to be started is placed in the need_to_start_server_shell array

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

need_to_start_server_shell=(
  "start_rpc_service.sh"
  "msg_gateway_start.sh"
  "push_start.sh"
  "msg_transfer_start.sh"
)

time=$(date +"%Y-%m-%d %H:%M:%S")

for _ in {1..3}; do
  echo "==========================================================" >> $OPENIM_ROOT/logs/openIM.log 2>&1
done

echo "==========server start time:${time}===========" >> $OPENIM_ROOT/logs/openIM.log 2>&1

for _ in {1..3}; do
  echo "==========================================================" >> $OPENIM_ROOT/logs/openIM.log 2>&1
done

build_pid_array=()
idx=0
for i in ${need_to_start_server_shell[*]}; do
  chmod +x $i
  ./$i &
  build_pid=$!
  echo "build_pid " $build_pid
  build_pid_array[idx]=$build_pid
  let idx=idx+1
done

echo "wait all start finish....."

exit 0

success_num=0
for ((i = 0; i < ${#need_to_start_server_shell[*]}; i++)); do
  echo "wait pid: " ${build_pid_array[i]} ${need_to_start_server_shell[$i]}
  wait ${build_pid_array[i]}
  stat=$?
  echo ${build_pid_array[i]}  " " $stat
 if [ $stat == 0 ]
 then
     # echo -e "${GREEN_PREFIX}${need_to_start_server_shell[$i]} successfully be built ${COLOR_SUFFIX}\n"
      let success_num=$success_num+1

 else
      #echo -e "${RED_PREFIX}${need_to_start_server_shell[$i]} build failed ${COLOR_SUFFIX}\n"
      exit -1
 fi
done

echo "success_num" $success_num  "service num:" ${#need_to_start_server_shell[*]}
if [ $success_num == ${#need_to_start_server_shell[*]} ]
then
  echo -e ${BACKGROUND_GREEN}"all services build success"${COLOR_SUFFIX}
fi
