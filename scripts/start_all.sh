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
source ./style_info.cfg
source ./path_info.cfg

# Print title
echo -e "${BOLD_PREFIX}${BLUE_PREFIX}OpenIM Server Start${COLOR_SUFFIX}"

# Get current time
time=$(date +"%Y-%m-%d %H:%M:%S")

# Print section separator
echo -e "${PURPLE_PREFIX}==========================================================${COLOR_SUFFIX}"

# Print server start time
echo -e "${BOLD_PREFIX}${CYAN_PREFIX}Server Start Time: ${time}${COLOR_SUFFIX}"

# Print section separator
echo -e "${PURPLE_PREFIX}==========================================================${COLOR_SUFFIX}"

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
  echo -e "=========> ${YELLOW_PREFIX}Executing ${i}...${COLOR_SUFFIX}"
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
echo "${PURPLE_PREFIX}==========================================================${COLOR_SUFFIX}"

# Print completion message
echo "${GREEN_PREFIX}${BOLD_PREFIX}OpenIM Server has been started successfully!${COLOR_SUFFIX}"