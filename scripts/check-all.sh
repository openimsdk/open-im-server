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

# This script is check openim service is running normally
# 
# Usage: `scripts/check-all.sh`.
# Encapsulated as: `make check`.
# READ: https://github.com/openimsdk/open-im-server/tree/main/scripts/install/environment.sh 

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/install/common.sh"

OPENIM_VERBOSE=4

openim::log::info "\n# Begin to check all openim service"

# OpenIM status
# Elegant printing function
print_services_and_ports() {
    # 获取数组
    declare -g service_names=("${!1}")
    declare -g service_ports=("${!2}")

    echo "+-------------------------+----------+"
    echo "| Service Name            | Port     |"
    echo "+-------------------------+----------+"

    for index in "${!service_names[@]}"; do
        printf "| %-23s | %-8s |\n" "${service_names[$index]}" "${service_ports[$index]}"
    done

    echo "+-------------------------+----------+"
}


# Print out services and their ports
print_services_and_ports OPENIM_SERVER_NAME_TARGETS OPENIM_SERVER_PORT_TARGETS

# Print out dependencies and their ports
print_services_and_ports OPENIM_DEPENDENCY_TARGETS OPENIM_DEPENDENCY_PORT_TARGETS


# OpenIM check
echo "++ The port being checked: ${OPENIM_SERVER_PORT_LISTARIES[@]}"
openim::log::info "\n## Check all dependent service ports"
echo "+++ The port being checked: ${OPENIM_DEPENDENCY_PORT_LISTARIES[@]}"

set +e

# Later, after discarding Docker, the Docker keyword is unreliable, and Kubepods is used
if grep -qE 'docker|kubepods' /proc/1/cgroup || [ -f /.dockerenv ]; then
    openim::color::echo ${COLOR_CYAN} "Environment in the interior of the container"
else
    openim::color::echo ${COLOR_CYAN} "The environment is outside the container"
    openim::util::check_ports ${OPENIM_DEPENDENCY_PORT_LISTARIES[@]} || return 0
fi

if [[ $? -ne 0 ]]; then
  openim::log::error_exit "The service does not start properly, please check the port, query variable definition!"
  echo "+++ https://github.com/openimsdk/open-im-server/tree/main/scripts/install/environment.sh +++"
else
  echo "++++ Check all dependent service ports successfully !"
fi

openim::log::info "\n## Check OpenIM service name"
. $(dirname ${BASH_SOURCE})/install/openim-msgtransfer.sh openim::msgtransfer::check

openim::log::info "\n## Check all OpenIM service ports"
echo "+++ The port being checked: ${OPENIM_SERVER_PORT_LISTARIES[@]}"
openim::util::check_ports ${OPENIM_SERVER_PORT_LISTARIES[@]}
if [[ $? -ne 0 ]]; then
  echo "+++ cat openim log file >>> ${LOG_FILE}"
  openim::log::error_exit "The service does not start properly, please check the port, query variable definition!"
else
  echo "++++ Check all openim service ports successfully !"
fi

set -e