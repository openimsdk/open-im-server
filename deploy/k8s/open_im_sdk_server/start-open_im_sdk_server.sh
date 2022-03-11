#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
ulimit -n 200000


service_filename="open_im_sdk_server"
binary_root="/Open-IM-Server/bin"
K8sServiceName="sdk-server"
K8sServiceNameForApi="api"
rm -rf /Open-IM-Server/config
mkdir /Open-IM-Server/config
cp /Open-IM-Server/config.tmp.yaml /Open-IM-Server/config/config.yaml
sed -i "s#openim-all#$K8sServiceNameForApi.$NAMESPACE#g" /Open-IM-Server/config/config.yaml

list1=$(cat $config_path | grep openImApiPort | awk -F '[:]' '{print $NF}')
list2=$(cat $config_path | grep openImWsPort | awk -F '[:]' '{print $NF}')
list3=$(cat $config_path | grep openImSdkWsPort | awk -F '[:]' '{print $NF}')
logLevel=$(cat $config_path | grep remainLogLevel | awk -F '[:]' '{print $NF}')
list_to_string $list1
api_ports=($ports_array)
list_to_string $list2
ws_ports=($ports_array)
list_to_string $list3
sdk_ws_ports=($ports_array)
list_to_string $list4


check=$(ps aux | grep -w ./${service_filename} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps aux | grep -w ./${service_filename} | grep -v grep | awk '{print $2}')
    kill -9 ${oldPid}
fi
#Check if the service exists
#If it is exists,kill this process
check=$(ps aux | grep -w ./${service_filename} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps aux | grep -w ./${service_filename} | grep -v grep | awk '{print $2}')
    kill -9 ${oldPid}
fi
#Waiting port recycling
sleep 1
cd ${binary_root}
./${service_filename} -openIM_api_port ${api_ports[0]} -openIM_ws_port ${ws_ports[0]} -sdk_ws_port ${sdk_ws_ports[0]}
