#!/usr/bin/env bash
#fixme This script is to stop the service
dir_name=`dirname $0`
if [ "${dir_name:0:1}" = "/" ]; then
  cur_dir="`dirname $0`"
else
  cur_dir="`pwd`"/"`dirname $0`"
fi

source "$cur_dir/style_info.cfg"
source "$cur_dir/path_info.cfg"


for i in ${service_names[*]}; do
  #Check whether the service exists
  name="ps -aux |grep -w $i |grep -v grep"
  count="${name}| wc -l"
  if [ $(eval ${count}) -gt 0 ]; then
    pid="${name}| awk '{print \$2}'"
    echo -e "${SKY_BLUE_PREFIX}Killing service:$i pid:$(eval $pid)${COLOR_SUFFIX}"
    #kill the service that existed
    kill -9 $(eval $pid)
    echo -e "${SKY_BLUE_PREFIX}service:$i was killed ${COLOR_SUFFIX}"
  fi
done
