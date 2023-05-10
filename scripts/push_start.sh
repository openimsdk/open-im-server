#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh



list1=$(cat $config_path | grep openImPushPort | awk -F '[:]' '{print $NF}')
list2=$(cat $config_path | grep pushPrometheusPort | awk -F '[:]' '{print $NF}')
list_to_string $list1
rpc_ports=($ports_array)
list_to_string $list2
prome_ports=($ports_array)

#Check if the service exists
#If it is exists,kill this process
check=$(ps aux | grep -w ./${push_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps aux | grep -w ./${push_name} | grep -v grep | awk '{print $2}')
  kill -9 $oldPid
fi
#Waiting port recycling
sleep 1
cd ${push_binary_root}

for ((i = 0; i < ${#rpc_ports[@]}; i++)); do
  nohup ./${push_name} -port ${rpc_ports[$i]} -prometheus_port ${prome_ports[$i]} >>../logs/openIM.log 2>&1 &
done

sleep 3
#Check launched service process
check=$(ps aux | grep -w ./${push_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  newPid=$(ps aux | grep -w ./${push_name} | grep -v grep | awk '{print $2}')
  ports=$(netstat -netulp | grep -w ${newPid} | awk '{print $4}' | awk -F '[:]' '{print $NF}')
  allPorts=""

  for i in $ports; do
    allPorts=${allPorts}"$i "
  done
  echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS "${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${push_name}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${newPid}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
  echo -e ${YELLOW_PREFIX}${push_name}${COLOR_SUFFIX}${RED_PREFIX}"SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
