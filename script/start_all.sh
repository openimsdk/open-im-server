#!/usr/bin/env bash
#fixme This script is the total startup script
#fixme The full name of the shell script that needs to be started is placed in the need_to_start_server_shell array

#fixme Put the shell script name here
dir_name=`dirname $0`
if [ "${dir_name:0:1}" = "/" ]; then
  cur_dir="`dirname $0`"
else
  cur_dir="`pwd`"/"`dirname $0`"
fi

need_to_start_server_shell=(
  $cur_dir/start_rpc_service.sh
  $cur_dir/push_start.sh
  $cur_dir/msg_transfer_start.sh
  $cur_dir/sdk_svr_start.sh
  $cur_dir/msg_gateway_start.sh
  $cur_dir/demo_svr_start.sh
#  start_cron.sh
)
time=`date +"%Y-%m-%d %H:%M:%S"`
echo "==========================================================">>$cur_dir/../logs/openIM.log 2>&1 &
echo "==========================================================">>$cur_dir/../logs/openIM.log 2>&1 &
echo "==========================================================">>$cur_dir/../logs/openIM.log 2>&1 &
echo "==========server start time:${time}===========">>$cur_dir/../logs/openIM.log 2>&1 &
echo "==========================================================">>$cur_dir/../logs/openIM.log 2>&1 &
echo "==========================================================">>$cur_dir/../logs/openIM.log 2>&1 &
echo "==========================================================">>$cur_dir/../logs/openIM.log 2>&1 &

for i in ${need_to_start_server_shell[*]}; do
  chmod +x $i
  $i
    if [ $? -ne 0 ]; then
        exit -1
  fi
done
