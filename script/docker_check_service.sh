#!/usr/bin/env bash

source ./style_info.cfg

docker_compose_components=(
  etcd
  mongo
  mysql
  open-im-server
  redis
  kafka
  zookeeper
)

component_server_count=0

for ((i = 0; i < ${#docker_compose_components[*]}; i++)); do
  component_server="docker-compose ps|grep -w ${docker_compose_components[$i]}|grep Up"
  count="${component_server}|wc -l"

  if [ $(eval ${count}) -gt 0 ]; then
    echo -e "${SKY_BLUE_PREFIX}docker-compose ${docker_compose_components[$i]} is Up!${COLOR_SUFFIX}"
    let component_server_count+=1
  else
    echo -e "${RED_PREFIX} ${docker_compose_components[$i]} start failed!${COLOR_SUFFIX}"
  fi
done

if [ ${component_server_count} -eq 6 ]; then
  echo -e "${YELLOW_PREFIX}\ndocker-compose all services is Up!${COLOR_SUFFIX}"
else
  echo -e "${RED_PREFIX}\nsome docker-compose services start failed,please check red logs on console ${COLOR_SUFFIX}"
fi
