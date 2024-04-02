
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/
source "${OPENIM_ROOT}/lib/util.sh"
source "${OPENIM_ROOT}/define/binaries.sh"
source "${OPENIM_ROOT}/lib/path.sh"
source "${OPENIM_ROOT}/lib/logging.sh"


#停止所有的二进制对应的进程
stop_binaries() {
  for binary in "${!binaries[@]}"; do
    full_path=$(get_bin_full_path "$binary")
    openim::util::kill_exist_binary "$full_path"
  done
}

LOG_FILE=log.file
ERR_LOG_FILE=err.log.file

#启动所有的二进制
start_binaries() {
  local project_dir="$OPENIM_ROOT"  # You should adjust this path as necessary
  # Iterate over binaries defined in binary_path.sh
  for binary in "${!binaries[@]}"; do
    local count=${binaries[$binary]}
    local bin_full_path=$(get_bin_full_path "$binary")
    # Loop to start binary the specified number of times
    for ((i=0; i<count; i++)); do
      echo "Starting $binary instance $i: $bin_full_path -i $i -c $OPENIM_OUTPUT_CONFIG"
      cmd=("$bin_full_path" -i "$i" -c "$OPENIM_OUTPUT_CONFIG")
      nohup "${cmd[@]}" >> "${LOG_FILE}" 2> >(tee -a "$ERR_LOG_FILE" | while read line; do echo -e "\e[31m${line}\e[0m"; done >&2) &
      done
  done
}


#kill二进制全路径对应的进程
kill_exist_binaries(){
  for binary in "${!binaries[@]}"; do
    full_path=$(get_bin_full_path "$binary")
    result=$(openim::util::kill_exist_binary "$full_path" | tail -n1)
  done
}


#检查所有的二进制是否退出
check_binaries_stop() {
  local running_binaries=0

  for binary in "${!binaries[@]}"; do
    full_path=$(get_bin_full_path "$binary")

    result=$(openim::util::check_process_names_exist "$full_path")
    if [ "$result" -ne 0 ]; then
      echo "Process for $binary is still running."
      running_binaries=$((running_binaries + 1))
    fi
  done

  if [ "$running_binaries" -ne 0 ]; then
    echo "There are $running_binaries binaries still running. Aborting..."
    return 1
  else
    echo "All processes have been stopped."
    return 0
  fi
}



#检查所有的二进制是否运行
check_binaries_running(){
  local no_running_binaries=0
  for binary in "${!binaries[@]}"; do
    expected_count=${binaries[$binary]}
    full_path=$(get_bin_full_path "$binary")

    result=$(openim::util::check_process_names "$full_path" "$expected_count")
    ret_val=$?
    if [ "$ret_val" -eq 0 ]; then
      echo "$binary is running normally."
    else
      no_running_binaries=$((no_running_binaries + 1))
      echo $result
    fi
  done

  if [ "$no_running_binaries" -ne 0 ]; then
      return 1
    else
      return 0
    fi
}




#打印所有的二进制对应的进程所所监听的端口
print_listened_ports_by_binaries(){
  for binary in "${!binaries[@]}"; do
    expected_count=${binaries[$binary]}
    base_path=$(get_bin_full_path "$binary")
    for ((i=0; i<expected_count; i++)); do
      full_path="${base_path} -i ${i} -c $OPENIM_OUTPUT_CONFIG"
      openim::util::print_binary_ports "$full_path"
    done
  done
}

