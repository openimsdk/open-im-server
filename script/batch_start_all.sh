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
  demo_svr_start.sh
)
time=`date +"%Y-%m-%d %H:%M:%S"`
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========server start time:${time}===========">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &
echo "==========================================================">>../logs/openIM.log 2>&1 &

build_pid_array=()
idx=0
for i in ${need_to_start_server_shell[*]}; do
  chmod +x $i
  ./$i &
  build_pid=$!
  echo "build_pid " $build_pid
  build_pid_array[idx]=$build_pid
  let idx=idx+1
done

echo "wait all start finish....."

exit 0

success_num=0
for ((i = 0; i < ${#need_to_start_server_shell[*]}; i++)); do
  echo "wait pid: " ${build_pid_array[i]} ${need_to_start_server_shell[$i]}
  wait ${build_pid_array[i]}
  stat=$?
  echo ${build_pid_array[i]}  " " $stat
 if [ $stat == 0 ]
 then
     # echo -e "${GREEN_PREFIX}${need_to_start_server_shell[$i]} successfully be built ${COLOR_SUFFIX}\n"
      let success_num=$success_num+1

 else
      #echo -e "${RED_PREFIX}${need_to_start_server_shell[$i]} build failed ${COLOR_SUFFIX}\n"
      exit -1
 fi
done

echo "success_num" $success_num  "service num:" ${#need_to_start_server_shell[*]}
if [ $success_num == ${#need_to_start_server_shell[*]} ]
then
  echo -e ${YELLOW_PREFIX}"all services build success"${COLOR_SUFFIX}
fi


