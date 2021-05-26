#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

#begin path
begin_path=$PWD

for ((i = 0; i < ${#service_source_root[*]}; i++)); do
  cd $begin_path
  service_path=${service_source_root[$i]}
  cd $service_path && echo -e "${SKY_BLUE_PREFIX}Current directory: $PWD $COLOR_SUFFIX"
  make install && echo -e "${SKY_BLUE_PREFIX}build ${service_names[$i]} success,moving binary file to the bin directory${COLOR_SUFFIX}" &&
    echo -e "${SKY_BLUE_PREFIX}Successful moved ${service_names[$i]} to the bin directory${COLOR_SUFFIX}\n"
done
