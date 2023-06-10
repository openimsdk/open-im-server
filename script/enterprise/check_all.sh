#!/usr/bin/env bash

source ./style_info.cfg
source ./enterprise/path_info.cfg
source ./enterprise/function.sh
service_port_name=(
 openImChatApiPort
 openImAdminApiPort
   #api port name
   openImAdminPort
   openImChatPort
)
switch=$(cat $config_path | grep demoswitch |awk -F '[:]' '{print $NF}')
for i in ${service_port_name[*]}; do
  list=$(cat $config_path | grep -w ${i} | awk -F '[:]' '{print $NF}')
  list_to_string $list
  for j in ${ports_array}; do
    port=$(ss -tunlp| grep open_im | awk '{print $5}' | grep -w ${j} | awk -F '[:]' '{print $NF}')
    if [[ ${port} -ne ${j} ]]; then
      echo -e ${YELLOW_PREFIX}${i}${COLOR_SUFFIX}${RED_PREFIX}" service does not start normally,not initiated port is "${COLOR_SUFFIX}${YELLOW_PREFIX}${j}${COLOR_SUFFIX}
      echo -e ${RED_PREFIX}"please check ../logs/openIM.log "${COLOR_SUFFIX}
      exit -1
    else
      echo -e ${j}${GREEN_PREFIX}" port has been listening,belongs service is "${i}${COLOR_SUFFIX}
    fi
  done
done

