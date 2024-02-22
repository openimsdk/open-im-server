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




OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/install/common.sh"


# Function to execute the scripts.
function execute_scripts() {
  for script_path in "${OPENIM_SERVER_SCRIPT_START_LIST[@]}"; do
    # Extract the script name without extension for argument generation.
    script_name_with_prefix=$(basename "$script_path" .sh)

    # Remove the "openim-" prefix.
    script_name=${script_name_with_prefix#openim-}

    # Construct the argument based on the script name.
    arg="openim::${script_name}::start"

    # Check if the script file exists and is executable.
    if [[ -x "$script_path" ]]; then
      openim::log::status "Starting script: ${script_path##*/}"     # Log the script name.

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




openim::log::info "\n# Begin to start all openim service scripts"


openim::golang::check_openim_binaries
if [[ $? -ne 0 ]]; then
  openim::log::error "OpenIM binaries are not found. Please run 'make build' to build binaries."
  "${OPENIM_ROOT}"/scripts/build-all-service.sh
fi


"${OPENIM_ROOT}"/scripts/init-config.sh --skip

echo "You need to start the following scripts in order: ${OPENIM_SERVER_SCRIPTARIES[@]}"


# TODO Prelaunch tools, simple for now, can abstract functions later
TOOLS_START_SCRIPTS_PATH=${START_SCRIPTS_PATH}/openim-tools.sh

openim::log::info "\n## Pre Starting OpenIM services"
${TOOLS_START_SCRIPTS_PATH} openim::tools::pre-start


"${OPENIM_ROOT}"/scripts/stop-all.sh

sleep 30


openim::log::info "\n## Starting OpenIM services"
execute_scripts

sleep 2

openim::log::info "\n## Check OpenIM service name"
. $(dirname ${BASH_SOURCE})/install/openim-msgtransfer.sh openim::msgtransfer::check

echo "+++ The process being checked:"
for item in "${OPENIM_ALL_SERVICE_LIBRARIES_NO_TRANSFER[@]}"; do
    echo "$item"
done

openim::util::check_process_names ${OPENIM_ALL_SERVICE_LIBRARIES_NO_TRANSFER[@]}

openim::log::info "\n## Post Starting OpenIM services"
${TOOLS_START_SCRIPTS_PATH} openim::tools::post-start

openim::color::echo $COLOR_BLUE "✨  All OpenIM services have been successfully started!"