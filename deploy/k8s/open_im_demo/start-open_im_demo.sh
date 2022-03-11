#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
ulimit -n 200000

service_filename="open_im_demo"
K8sServiceName="demo"
binary_root="/Open-IM-Server/bin"

rm -rf /Open-IM-Server/config
mkdir /Open-IM-Server/config
cp /Open-IM-Server/config.tmp.yaml /Open-IM-Server/config/config.yaml
sed -i "s#openim-all#$POD_NAME.$K8sServiceName.$NAMESPACE.svc.cluster.local#g" /Open-IM-Server/config/config.yaml

switch=$(cat $config_path | grep demoswitch |awk -F '[:]' '{print $NF}')
if [ ${switch} != "true" ]; then
      echo -e ${YELLOW_PREFIX}" demo service switch is false not start demo "${COLOR_SUFFIX}
      exit 0
fi
list1=$(cat $config_path | grep openImDemoPort | awk -F '[:]' '{print $NF}')
list_to_string $list1
api_ports=($ports_array)

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
./${service_filename} -port ${api_ports[$i]}
