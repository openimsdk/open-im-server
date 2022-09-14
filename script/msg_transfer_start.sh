#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

list1=$(cat $config_path | grep messageTransferPrometheusPort | awk -F '[:]' '{print $NF}')
list_to_string $list1
prome_ports=($ports_array)


#Check if the service exists
#If it is exists,kill this process
check=`ps aux | grep -w ./${msg_transfer_name} | grep -v grep| wc -l`
if [ $check -ge 1 ]
then
oldPid=`ps aux | grep -w ./${msg_transfer_name} | grep -v grep|awk '{print $2}'`
 kill -9 $oldPid
fi
#Waiting port recycling
sleep 1

cd ${msg_transfer_binary_root}
for ((i = 0; i < ${msg_transfer_service_num}; i++)); do
      prome_port=${prome_ports[$i]}
      cmd="nohup ./${msg_transfer_name}"
      if [ $prome_port != "" ]; then
        cmd="$cmd -prometheus_port $prome_port"
      fi
      $cmd >>../logs/openIM.log 2>&1 &
done

#Check launched service process
check=`ps aux | grep -w ./${msg_transfer_name} | grep -v grep| wc -l`
if [ $check -ge 1 ]
then
newPid=`ps aux | grep -w ./${msg_transfer_name} | grep -v grep|awk '{print $2}'`
allPorts=""
    echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS "${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${msg_transfer_name}${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${newPid}${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
    echo -e ${YELLOW_PREFIX}${msg_transfer_name}${COLOR_SUFFIX}${RED_PREFIX}"SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
