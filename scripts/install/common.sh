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


# Common utilities, variables and checks for all build scripts.
set -o errexit
set +o nounset
set -o pipefail

# Sourced flag
COMMON_SOURCED=true

# The root of the build/dist directory
OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)

source "${OPENIM_ROOT}/scripts/lib/init.sh"

# Make sure the environment is only called via common to avoid too much nesting
source "${OPENIM_ROOT}/scripts/install/environment.sh"

# Storing all the defined ports in an array for easy management and access.
# This array consolidates the port numbers for all the services defined above.
openim::common::service_port_name() {
  local targets=(
    $OPENIM_USER_PORT            # User service
    $OPENIM_FRIEND_PORT          # Friend service
    $OPENIM_MESSAGE_PORT         # Message service
    $OPENIM_MESSAGE_GATEWAY_PORT # Message gateway
    $OPENIM_GROUP_PORT           # Group service
    $OPENIM_AUTH_PORT            # Authorization service
    $OPENIM_PUSH_PORT            # Push service
    $OPENIM_CONVERSATION_PORT    # Conversation service
    $OPENIM_THIRD_PORT           # Third-party service
    $API_OPENIM_PORT             # API service
    $OPENIM_WS_PORT              # WebSocket service
  )
  echo "${targets[@]}"
}

IFS=" " read -ra OPENIM_SERVER_PORT_TARGETS <<< "$(openim::common::service_port_name)"
readonly OPENIM_SERVER_PORT_TARGETS
readonly OPENIM_SERVER_PORT_LISTARIES=("${OPENIM_SERVER_PORT_TARGETS[@]##*/}")

# Execute commands that require root permission without entering a password
function openim::common::sudo {
  echo ${LINUX_PASSWORD} | sudo -S $1
}
