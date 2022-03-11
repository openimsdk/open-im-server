#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
ulimit -n 200000

service_filename="open_im_push"
K8sServiceName="push"
binary_root="/Open-IM-Server/bin"

rm -rf /Open-IM-Server/config
mkdir /Open-IM-Server/config
cp /Open-IM-Server/config.tmp.yaml /Open-IM-Server/config/config.yaml
sed -i "s#openim-all#$POD_NAME.$K8sServiceName.$NAMESPACE.svc.cluster.local#g" /Open-IM-Server/config/config.yaml


list1=$(cat $config_path | grep openImPushPort | awk -F '[:]' '{print $NF}')
list_to_string $list1
rpc_ports=($ports_array)

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
./${service_filename} -port ${rpc_ports[$i]}
