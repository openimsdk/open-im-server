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


#!/bin/bash





OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/install/common.sh"



# Function to execute the scripts.
function execute_start_scripts() {
  for script_path in "${OPENIM_SERVER_SCRIPT_START_LIST[@]}"; do
    # Extract the script name without extension for argument generation.
    script_name_with_prefix=$(basename "$script_path" .sh)

    # Remove the "openim-" prefix.
    script_name=${script_name_with_prefix#openim-}

    # Construct the argument based on the script name.
    arg="openim::${script_name}::start"

    # Check if the script file exists and is executable.
    if [[ -x "$script_path" ]]; then
       openim::log::print_blue "Starting script: ${script_path##*/}"     # Log the script name.

      # Execute the script with the constructed argument.
      result=$("$script_path" "$arg")
     if [[ $? -ne 0 ]]; then
        openim::log::error "Start script: ${script_path##*/} failed"
        openim::log::error "$result"
        return 1
      fi

    else
      openim::log::errexit "Script ${script_path##*/} is missing or not executable."
      return 1
    fi
  done
}




if openim::util::is_running_in_container; then
  exec > ${DOCKER_LOG_FILE} 2>&1
fi



openim::golang::check_openim_binaries
if [[ $? -ne 0 ]]; then
  openim::log::error "OpenIM binaries are not found. Please run 'make build' to build binaries."
  "${OPENIM_ROOT}"/scripts/build-all-service.sh
fi


"${OPENIM_ROOT}"/scripts/init-config.sh --skip

#openim::log::print_blue "Execute the following script in sequence: ${OPENIM_SERVER_SCRIPTARIES[@]}"


# TODO Prelaunch tools, simple for now, can abstract functions later
TOOLS_START_SCRIPTS_PATH=${START_SCRIPTS_PATH}/openim-tools.sh

openim::log::print_blue "\n## Pre Starting OpenIM services"



if ! ${TOOLS_START_SCRIPTS_PATH} openim::tools::pre-start; then
  openim::log::error "Pre Starting OpenIM services failed, aborting..."
  exit 1
fi


openim::log::print_blue "Pre Starting OpenIM services processed successfully"

result=$("${OPENIM_ROOT}"/scripts/stop-all.sh)
if [[ $? -ne 0 ]]; then
  openim::log::error "View the error logs from this startup. ${LOG_FILE} \n"
  openim::log::error "Some programs have not exited; the start process is aborted .\n $result"
  exit 1
fi



openim::log::status "\n## Starting openim scripts: "
execute_start_scripts

sleep 2

result=$(. $(dirname ${BASH_SOURCE})/install/openim-msgtransfer.sh openim::msgtransfer::check)
if [[ $? -ne 0 ]]; then
  openim::log::error "The program may fail to start.\n $result"
  exit 1
fi


result=$(openim::util::check_process_names ${OPENIM_ALL_SERVICE_LIBRARIES_NO_TRANSFER[@]})
if [[ $? -ne 0 ]]; then
  openim::log::error "The program may fail to start.\n $result"
  exit 1
fi


openim::log::info "\n## Post Starting openim services"
${TOOLS_START_SCRIPTS_PATH} openim::tools::post-start

openim::log::success "All openim services have been successfully started!"