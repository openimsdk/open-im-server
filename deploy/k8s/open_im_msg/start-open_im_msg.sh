#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

service_filename="open_im_msg"
service_port_name="openImOfflineMessagePort"
K8sServiceName="msg"

#Check whether the service exists
service_name="ps -aux |grep -w ${service_filename} |grep -v grep"
count="${service_name}| wc -l"

rm -rf /Open-IM-Server/config
mkdir /Open-IM-Server/config
cp /Open-IM-Server/config.tmp.yaml /Open-IM-Server/config/config.yaml
sed -i "s#openim-all#$POD_NAME.$K8sServiceName.$NAMESPACE.svc.cluster.local#g" /Open-IM-Server/config/config.yaml

if [ $(eval ${count}) -gt 0 ]; then
  pid="${service_name}| awk '{print \$2}'"
  echo -e "${SKY_BLUE_PREFIX}${service_filename} service has been started,pid:$(eval $pid)$COLOR_SUFFIX"
  echo -e "${SKY_BLUE_PREFIX}Killing the service ${service_filename} pid:$(eval $pid)${COLOR_SUFFIX}"
  #kill the service that existed
  kill -9 $(eval $pid)
  sleep 0.5
fi
cd ../bin && echo -e "${SKY_BLUE_PREFIX}${service_filename} service is starting${COLOR_SUFFIX}"
#Get the rpc port in the configuration file
portList=$(cat $config_path | grep ${service_port_name} | awk -F '[:]' '{print $NF}')
list_to_string ${portList}
#Start related rpc services based on the number of ports
for j in ${ports_array}; do
  echo -e "${SKY_BLUE_PREFIX}${POD_NAME} Service is starting,port number:$j $COLOR_SUFFIX"
  ./${service_filename} -port $j
done
