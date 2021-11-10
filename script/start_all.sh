#!/usr/bin/env bash
#fixme This script is the total startup script
#fixme The full name of the shell script that needs to be started is placed in the need_to_start_server_shell array

#fixme Put the shell script name here
need_to_start_server_shell=(
  start_rpc_service.sh
  msg_gateway_start.sh
  push_start.sh
  msg_transfer_start.sh
  sdk_svr_start.sh
  timer_start.sh
)

for i in ${need_to_start_server_shell[*]}; do
  chmod +x $i
  ./$i
done

