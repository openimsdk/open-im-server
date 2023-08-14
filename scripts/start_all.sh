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
source "${OPENIM_ROOT}/scripts/lib/init.sh"

set +o errexit
openim::golang::check_openim_binaries
if [[ $? -ne 0 ]]; then
  openim::log::error "OpenIM binaries are not found. Please run 'make' to build binaries."
  ${OPENIM_ROOT}/scripts/build_all_service.sh
fi
set -o errexit

scripts_to_run=$(openim::golang::start_script_list)

for script in $scripts_to_run; do
    openim::log::info "Executing: $script"
    "$script"
    if [ $? -ne 0 ]; then
    # Print error message and exit
    openim::log::error "Error executing ${i}. Exiting..."
    exit -1
  fi
done

openim::log::success "OpenIM Server has been started successfully!"
