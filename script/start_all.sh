#!/usr/bin/env bash
#fixme This script is the total startup script
#fixme The full name of the shell script that needs to be started is placed in the need_to_start_server_shell array

#fixme Put the shell script name here
need_to_start_server_shell=(
  start_rpc_service.sh
  push_start.sh
  msg_transfer_start.sh
  sdk_svr_start.sh
  msg_gateway_start.sh
  demo_svr_start.sh
#  start_cron.sh
)
time=`date +"%Y-%m-%d %H:%M:%S"`
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========server start time:${time}===========">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &

for i in ${need_to_start_server_shell[*]}; do
  chmod +x $i
  ./$i
    if [ $? -ne 0 ]; then
        exit -1
  fi
done
