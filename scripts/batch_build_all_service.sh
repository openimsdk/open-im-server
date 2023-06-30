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


source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

bin_dir="../bin"
logs_dir="../logs"
sdk_db_dir="../db/sdk/"
#Automatically created when there is no bin, logs folder
if [ ! -d $bin_dir ]; then
  mkdir -p $bin_dir
fi
if [ ! -d $logs_dir ]; then
  mkdir -p $logs_dir
fi
if [ ! -d $sdk_db_dir ]; then
  mkdir -p $sdk_db_dir
fi

#begin path
begin_path=$PWD


build_pid_array=()

for ((i = 0; i < ${#service_source_root[*]}; i++)); do
  cd $begin_path
  service_path=${service_source_root[$i]}
  cd $service_path
  make install > /dev/null &
  build_pid=$!
  build_pid_array[i]=$build_pid
done


echo "wait all build finish....."

success_num=0
for ((i = 0; i < ${#service_source_root[*]}; i++)); do
  echo "wait pid: " ${build_pid_array[i]} ${service_names[$i]}
  wait ${build_pid_array[i]}
  stat=$?
  echo ${service_names[$i]} "pid: " ${build_pid_array[i]}  "stat: " $stat
 if [ $stat == 0 ]
 then
      echo -e "${GREEN_PREFIX}${service_names[$i]} successfully be built ${COLOR_SUFFIX}\n"
      let success_num=$success_num+1

 else
      echo -e "${RED_PREFIX}${service_names[$i]} build failed ${COLOR_SUFFIX}\n"
      exit -1
 fi
done

echo "success_num" $success_num  "service num:" ${#service_source_root[*]}
if [ $success_num == ${#service_source_root[*]} ]
then
  echo -e ${YELLOW_PREFIX}"all services build success"${COLOR_SUFFIX}
fi
