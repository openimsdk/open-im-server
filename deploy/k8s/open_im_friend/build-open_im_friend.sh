#!/bin/bash

serviceName="open_im_friend"
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

cd /Open-IM-Server/cmd/rpc/${serviceName} && echo -e "${SKY_BLUE_PREFIX}Current directory: $PWD $COLOR_SUFFIX"
make install && echo -e "${SKY_BLUE_PREFIX}build ${serviceName} success,moving binary file to the bin directory${COLOR_SUFFIX}" &&
echo -e "${SKY_BLUE_PREFIX}Successful moved ${serviceName} to the bin directory${COLOR_SUFFIX}\n"
