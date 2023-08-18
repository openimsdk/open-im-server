#!/usr/bin/env bash

# Copyright © 2023 OpenIM open source community. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/usr/bin/env bash


# Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "$SCRIPTS_ROOT/style_info.sh"
source "$SCRIPTS_ROOT/path_info.sh"
source "$SCRIPTS_ROOT/function.sh"

printf "${YELLOW_PREFIX}=======>SCRIPTS_ROOT=%s${COLOR_SUFFIX}\n" "$SCRIPTS_ROOT"
printf "${YELLOW_PREFIX}=======>OPENIM_ROOT=%s${COLOR_SUFFIX}\n" "$OPENIM_ROOT"
printf "${YELLOW_PREFIX}=======>pwd=%s${COLOR_SUFFIX}\n" "$PWD"

bin_dir="$BIN_DIR"
logs_dir="$SCRIPTS_ROOT/../logs"

printf "${YELLOW_PREFIX}=======>bin_dir=%s${COLOR_SUFFIX}\n" "$bin_dir"
printf "${YELLOW_PREFIX}=======>logs_dir=%s${COLOR_SUFFIX}\n" "$logs_dir"
printf "${YELLOW_PREFIX}=======>sdk_db_dir=%s${COLOR_SUFFIX}\n" "$sdk_db_dir"

# Service filenames
service_filenames=(
  chat-api
  admin-api
  #rpc
  admin-rpc
  chat-rpc
)

# Service config port names
service_port_names=(
  openImChatApiPort
  openImAdminApiPort
  #api port name
  openImAdminPort
  openImChatPort
)

service_prometheus_port_names=()

cd "$SCRIPTS_ROOT"

# Function to kill a service
kill_service() {
  local service_name=$1
  local pid=$(pgrep -f "$service_name")
  if [ -n "$pid" ]; then
    echo "$service_name service has been started, pid: $pid"
    echo "Killing the service $service_name, pid: $pid"
    killall "$service_name"
    sleep 0.5
  fi
}

for ((i = 0; i < ${#service_filenames[*]}; i++)); do
  service_name="${service_filenames[$i]}"
  kill_service "$service_name"
  cd "$SCRIPTS_ROOT"

  # Get the rpc ports from the configuration file
  readarray -t portList < "$config_path"
  service_ports=()
  for line in "${portList[@]}"; do
    if [[ $line == *"${service_port_names[$i]}"* ]]; then
      port=$(echo "$line" | awk -F ':' '{print $NF}')
      service_ports+=("$port")
    fi
  done

  # Start related rpc services based on the number of ports
  for port in "${service_ports[@]}"; do
    # Start the service in the background
    cmd="$bin_dir/$service_name -port $port --config_folder_path $config_path"
    if [[ $i -eq 0 || $i -eq 1 ]]; then
      cmd="$bin_dir/$service_name -port $port --config_folder_path $config_path"
    fi
    echo "$cmd"
    nohup "$cmd" >> "${logs_dir}/openIM.log" 2>&1 &
    sleep 1
  done
done

sleep infinity
