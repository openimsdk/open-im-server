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

set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/install/common.sh"

openim::log::info "\n# Begin to start all openim service scripts"

set +o errexit
openim::golang::check_openim_binaries
if [[ $? -ne 0 ]]; then
  openim::log::error "OpenIM binaries are not found. Please run 'make build' to build binaries."
  "${OPENIM_ROOT}"/scripts/build-all-service.sh
fi
set -o errexit

echo "You need to start the following scripts in order: ${OPENIM_SERVER_SCRIPTARIES[@]}"
openim::log::install_errexit

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
            "$script_path" "$arg"

            # Check if the script executed successfully.
            if [[ $? -eq 0 ]]; then
                openim::log::info "${script_path##*/} executed successfully."
            else
                openim::log::errexit "Error executing ${script_path##*/}."
            fi
        else
            openim::log::errexit "Script ${script_path##*/} is missing or not executable."
        fi
    done
    sleep 0.5
}


# TODO Prelaunch tools, simple for now, can abstract functions later
TOOLS_START_SCRIPTS_PATH=${START_SCRIPTS_PATH}/openim-tools.sh

openim::log::info "\n## Pre Starting OpenIM services"
${TOOLS_START_SCRIPTS_PATH} openim::tools::pre-start

openim::log::info "\n## Starting OpenIM services"
execute_scripts

openim::log::info "\n## Post Starting OpenIM services"
${TOOLS_START_SCRIPTS_PATH} openim::tools::post-start

openim::log::success "✨  All OpenIM services have been successfully started!"