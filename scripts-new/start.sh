
#!/usr/bin/env bash


OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
OPENIM_SCRIPTS=$(dirname "${BASH_SOURCE[0]}")/
source "$OPENIM_SCRIPTS/lib/path.sh"
source "$OPENIM_SCRIPTS/define/binaries.sh"


# Assuming 'openim::util::host_platform' is defined in one of the sourced scripts or elsewhere.
# If not, you'll need to define it to return the appropriate platform directory name.

# Main function to start binaries
start_binaries() {
  local project_dir="$OPENIM_ROOT"  # You should adjust this path as necessary
  # Iterate over binaries defined in binary_path.sh
  for binary in "${!binaries[@]}"; do
    local count=${binaries[$binary]}
    local bin_full_path=$(get_bin_full_path "$binary")
    # Loop to start binary the specified number of times
    for ((i=0; i<count; i++)); do
      echo "Starting $binary instance $i: $bin_full_path -i $i -c $OPENIM_OUTPUT_CONFIG"
      nohup "$bin_full_path" -i "$i" -c "$OPENIM_OUTPUT_CONFIG" > "test.log" 2>&1 &

      done
  done
}

kill_exist_binaries(){
  for binary in "${!binaries[@]}"; do
    full_path=$(get_bin_full_path "$binary")
    result=$(openim::util::kill_exist_binary "$full_path" | tail -n1)
    echo "result     $result"
   if [ "$result" -eq 0 ]; then
       echo "$full_path no exist"
     else
       echo "$full_path running. waiting stop"
     fi
  done
}


check_all_stop() {
  for binary in "${!binaries[@]}"; do
    expected_count=${binaries[$binary]}
    full_path=$(get_bin_full_path "$binary")

    result=$(openim::util::check_process_names_exist "$full_path")
    if [ "$result" -ne 0 ]; then
          echo "Process for $binary is still running. Aborting..."
          exit 1
        fi
      done
      echo "All processes have been stopped."
}



kill_exist_binaries

check_all_stop

# Call the main function
start_binaries



