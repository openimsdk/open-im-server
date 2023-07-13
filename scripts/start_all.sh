#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
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

#FIXME This script is the startup script for multiple servers.
#FIXME The full names of the shell scripts that need to be started are placed in the `need_to_start_server_shell` array.

#Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${SCRIPTS_ROOT}")/..

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/style_info.sh
source $SCRIPTS_ROOT/path_info.sh
source $SCRIPTS_ROOT/function.sh

cd $SCRIPTS_ROOT

echo -e "${YELLOW_PREFIX}=======>SCRIPTS_ROOT=$SCRIPTS_ROOT${COLOR_SUFFIX}"
echo -e "${YELLOW_PREFIX}=======>OPENIM_ROOT=$OPENIM_ROOT${COLOR_SUFFIX}"
echo -e "${YELLOW_PREFIX}=======>pwd=$PWD${COLOR_SUFFIX}"

bin_dir="$BIN_DIR"
logs_dir="$OPENIM_ROOT/logs"
sdk_db_dir="$OPENIM_ROOT/sdk/db/"

if [ ! -d "$bin_dir" ]; then
    mkdir -p "$bin_dir"
fi

if [ ! -d "$logs_dir" ]; then
    mkdir -p "$logs_dir"
fi

if [ ! -d "$sdk_db_dir" ]; then
    mkdir -p "$sdk_db_dir"
fi

# Print title
echo -e "${BOLD_PREFIX}${BLUE_PREFIX}================> OpenIM Server Start${COLOR_SUFFIX}"

# Get current time
time=$(date +"%Y-%m-%d %H:%M:%S")

# Print section separator
echo -e "${PURPLE_PREFIX}==========================================================${COLOR_SUFFIX}"

# Print server start time
echo -e "${BOLD_PREFIX}${CYAN_PREFIX}Server Start Time: ${time}${COLOR_SUFFIX}"

# Print section separator
echo -e "${PURPLE_PREFIX}==========================================================${COLOR_SUFFIX}"

cd  $SCRIPTS_ROOT

# FIXME Put the shell script names here
need_to_start_server_shell=(
  start_rpc_service.sh
  push_start.sh
  msg_transfer_start.sh
  msg_gateway_start.sh
  start_cron.sh
)

# Loop through the script names and execute them
for i in ${need_to_start_server_shell[*]}; do
  chmod +x $i

  echo -e ""
  # Print script execution message
  echo -e "=========> ${BACKGROUND_GREEN}Executing ${i}...${COLOR_SUFFIX}"
  echo -e ""

  ./$i

  # Check if the script executed successfully
  if [ $? -ne 0 ]; then
    # Print error message and exit
    echo "${BOLD_PREFIX}${RED_PREFIX}Error executing ${i}. Exiting...${COLOR_SUFFIX}"
    exit -1
  fi
done

# Print section separator
echo -e "${PURPLE_PREFIX}==========================================================${COLOR_SUFFIX}"

# Print completion message
echo -e "${GREEN_PREFIX}${BOLD_PREFIX}OpenIM Server has been started successfully!${COLOR_SUFFIX}"
