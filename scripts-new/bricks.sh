
source "$(dirname "${BASH_SOURCE[0]}")/lib/util.sh"
source "$(dirname "${BASH_SOURCE[0]}")/define/binaries.sh"
source "$(dirname "${BASH_SOURCE[0]}")/lib/path.sh"
source "$(dirname "${BASH_SOURCE[0]}")/lib/logging.sh"


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
  # Iterate over binaries defined in binary_path.sh
  for binary in "${!binaries[@]}"; do
    local count=${binaries[$binary]}
    local bin_full_path=$(get_bin_full_path "$binary")
    # Loop to start binary the specified number of times
    for ((i=0; i<count; i++)); do
      echo "Starting $bin_full_path -i $i -c $OPENIM_OUTPUT_CONFIG"
      cmd=("$bin_full_path" -i "$i" -c "$OPENIM_OUTPUT_CONFIG")
      nohup "${cmd[@]}" >> "${LOG_FILE}" 2> >(tee -a "$ERR_LOG_FILE" | while read line; do echo -e "\e[31m${line}\e[0m"; done >&2) &
      done
  done
}


start_tools() {
  # Assume tool_binaries=("ncpu" "infra")
  for binary in "${tool_binaries[@]}"; do
    local bin_full_path=$(get_tool_full_path "$binary")
    cmd=("$bin_full_path" -c "$OPENIM_OUTPUT_CONFIG")
    echo "Starting ${cmd[@]}"
    result=$( "${cmd[@]}" 2>&1 )
    ret_val=$?
    if [ $ret_val -eq 0 ]; then
        echo "Started $bin_full_path successfully." $result
    else
        echo "Failed to start $bin_full_path with exit code $ret_val." $result
        return 1
    fi
  done

  return 0
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
    return 1
  else
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
    if [ "$ret_val" -ne 0 ]; then
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

