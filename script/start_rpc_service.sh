#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

#service filename
service_filename=(
  #api
  open_im_api
  open_im_cms_api
  #rpc
  open_im_user
  open_im_friend
  open_im_group
  open_im_auth
  open_im_admin_cms
  ${msg_name}
  open_im_office
  open_im_organization
  open_im_conversation
  open_im_cache
)

#service config port name
service_port_name=(
  #api port name
  openImApiPort
  openImCmsApiPort
  #rpc port name
  openImUserPort
  openImFriendPort
  openImGroupPort
  openImAuthPort
  openImAdminCmsPort
  openImMessagePort
  openImOfficePort
  openImOrganizationPort
  openImConversationPort
  openImCachePort
)

service_prometheus_port_name=(
  #api port name
  openImApiPort
  openImCmsApiPort
  #rpc port name
  userPrometheusPort
  friendPrometheusPort
  groupPrometheusPort
  authPrometheusPort
  adminCmsPrometheusPort
  messagePrometheusPort
  officePrometheusPort
  organizationPrometheusPort
  conversationPrometheusPort
  cachePrometheusPort
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
  cd ../bin
  #Get the rpc port in the configuration file
  portList=$(cat $config_path | grep ${service_port_name[$i]} | awk -F '[:]' '{print $NF}')
  list_to_string ${portList}
  service_ports=($ports_array)

  portList2=$(cat $config_path | grep ${service_prometheus_port_name[$i]} | awk -F '[:]' '{print $NF}')
  list_to_string $portList2
  prome_ports=($ports_array)
  #Start related rpc services based on the number of ports
  for ((j = 0; j < ${#service_ports[*]}; j++)); do
    #Start the service in the background
    cmd="./${service_filename[$i]} -port ${service_ports[$j]} -prometheus_port ${prome_ports[$j]}"
    if [ $i -eq 0 -o $i -eq 1 ]; then
      cmd="./${service_filename[$i]} -port ${service_ports[$j]}"
    fi
    echo $cmd
    nohup $cmd >>../logs/openIM.log 2>&1 &
    sleep 1
    pid="netstat -ntlp|grep $j |awk '{printf \$7}'|cut -d/ -f1"
    echo -e "${GREEN_PREFIX}${service_filename[$i]} start success,port number:${service_ports[$j]} pid:$(eval $pid)$COLOR_SUFFIX"
  done
done
