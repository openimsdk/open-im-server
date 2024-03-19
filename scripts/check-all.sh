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





OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/install/common.sh"

if openim::util::is_running_in_container; then
  exec >> ${DOCKER_LOG_FILE} 2>&1
fi


OPENIM_VERBOSE=4

openim::log::info "\n# Begin to check all openim service"

openim::log::status "Check all dependent service ports"
# Elegant printing function
# Elegant printing function
print_services_and_ports() {
  local service_names=("$@")
  local half_length=$((${#service_names[@]} / 2))
  local service_ports=("${service_names[@]:half_length}")
  
  echo "+-------------------------+----------+"
  echo "| Service Name            | Port     |"
  echo "+-------------------------+----------+"
  
  for ((index=0; index < half_length; index++)); do
    printf "| %-23s | %-8s |\n" "${service_names[$index]}" "${service_ports[$index]}"
  done
  
  echo "+-------------------------+----------+"
}


# Assuming OPENIM_SERVER_NAME_TARGETS and OPENIM_SERVER_PORT_TARGETS are defined
# Similarly for OPENIM_DEPENDENCY_TARGETS and OPENIM_DEPENDENCY_PORT_TARGETS

# Print out services and their ports
print_services_and_ports "${OPENIM_SERVER_NAME_TARGETS[@]}" "${OPENIM_SERVER_PORT_TARGETS[@]}"

# Print out dependencies and their ports
print_services_and_ports "${OPENIM_DEPENDENCY_TARGETS[@]}" "${OPENIM_DEPENDENCY_PORT_TARGETS[@]}"

# OpenIM check
echo "++ The port being checked: ${OPENIM_SERVER_PORT_LISTARIES[@]}"
openim::log::info "\n## Check all dependent service ports"
echo "++ The port being checked: ${OPENIM_DEPENDENCY_PORT_LISTARIES[@]}"


# Later, after discarding Docker, the Docker keyword is unreliable, and Kubepods is used
if grep -qE 'docker|kubepods' /proc/1/cgroup || [ -f /.dockerenv ]; then
  openim::color::echo ${COLOR_CYAN} "Environment in the interior of the container"
else
  openim::color::echo ${COLOR_CYAN} "The environment is outside the container"
  openim::util::check_ports ${OPENIM_DEPENDENCY_PORT_LISTARIES[@]}
fi

if [[ $? -ne 0 ]]; then
  openim::log::error_exit "The service does not start properly, please check the port, query variable definition!"
  echo "+++ https://github.com/openimsdk/open-im-server/tree/main/scripts/install/environment.sh +++"
else
  openim::log::success "All components depended on by openim are running normally! "
fi


openim::log::info "\n## Check openim service name:\n${OPENIM_OUTPUT_HOSTBIN}/openim-msgtransfer"
result=$(. $(dirname ${BASH_SOURCE})/install/openim-msgtransfer.sh openim::msgtransfer::check)
if [[ $? -ne 0 ]]; then
  echo "+++ cat openim log file >>> ${LOG_FILE}"
  openim::log::error "check process failed.\n $result"
fi


echo "Check openim service name:"
for item in "${OPENIM_ALL_SERVICE_LIBRARIES_NO_TRANSFER[@]}"; do
    echo "$item"
done

echo "11111111111111111111111:"
openim::util::check_process_names ${OPENIM_ALL_SERVICE_LIBRARIES_NO_TRANSFER[@]}
echo "22222222222222222222222222:"
result=$(openim::util::check_process_names ${OPENIM_ALL_SERVICE_LIBRARIES_NO_TRANSFER[@]})
if [[ $? -ne 0 ]]; then
  echo "+++ cat openim log file >>> ${LOG_FILE}"
  openim::log::error "check process failed.\n "
  echo "$result"
  exit 1
else
  openim::log::success "All openim services are running normally! "
fi

