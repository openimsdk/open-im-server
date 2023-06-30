#!/usr/bin/env bash
#fixme This scripts is the total startup scripts
#fixme The full name of the shell scripts that needs to be started is placed in the need_to_start_server_shell array

#fixme Put the shell scripts name here
need_to_start_server_shell=(
  start_rpc_service.sh
  msg_gateway_start.sh
  push_start.sh
  msg_transfer_start.sh
  sdk_svr_start.sh
  start_cron.sh
)

#fixme The 10 second delay to start the project is for the docker-compose one-click to start openIM when the infrastructure dependencies are not started

sleep 10
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
done

sleep 15

#fixme prevents the openIM service exit after execution in the docker container
tail -f /dev/null
