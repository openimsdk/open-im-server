
#!/usr/bin/env bash

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..


start_binaries() {
  for bin in "${!binaries[@]}"; do
    local count=${binaries[$bin]}
    local bin_path=$(get_bin_full_path "$OPENIM_ROOT" "$bin")
    local conf_dir=$(get_conf_dir)

    for ((i=0; i<count; i++)); do
      echo "Starting $bin instance $i: $bin_path -i $i -c $conf_dir"
      "$bin_path" -i "$i" -c "$conf_dir"
    done
  done
}

start_binaries

#!/bin/bash

# Assuming the current script is located in the 'scripts' directory,
# adjust these paths according to your actual directory structure.
source "$(dirname "$0")/define/binaries.sh"
source "$(dirname "$0")/lib/path.sh"

# Assuming 'openim::util::host_platform' is defined in one of the sourced scripts or elsewhere.
# If not, you'll need to define it to return the appropriate platform directory name.

# Main function to start binaries
start_binaries() {
  local project_dir="$OPENIM_ROOT"  # You should adjust this path as necessary

  # Iterate over binaries defined in binary_path.sh
  for binary in "${!binaries[@]}"; do
    local count=${binaries[$binary]}
    local bin_full_path=$(get_bin_full_path "$project_dir" "$binary")
    local conf_dir=$(get_conf_dir "$project_dir")

    # Loop to start binary the specified number of times
    for ((i=0; i<count; i++)); do
      echo "Starting $binary instance $i: $bin_full_path -i $i -f $conf_dir"
      "$bin_full_path" -i "$i" -f "$conf_dir"
    done
  done
}

# Call the main function
start_binaries



