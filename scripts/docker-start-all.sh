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

set -o errexit
set -o nounset
set -o pipefail

#fixme This scripts is the total startup scripts
#fixme The full name of the shell scripts that needs to be started is placed in the need_to_start_server_shell array

# Fixed ports inside the docker startup container
export OPENIM_WS_PORT=10001
export API_OPENIM_PORT=10002
export API_PROM_PORT=20100
export USER_PROM_PORT=20110
export FRIEND_PROM_PORT=20120
export MESSAGE_PROM_PORT=20130
export MSG_GATEWAY_PROM_PORT=20140
export GROUP_PROM_PORT=20150
export AUTH_PROM_PORT=20160
export PUSH_PROM_PORT=20170
export CONVERSATION_PROM_PORT=20230
export RTC_PROM_PORT=21300
export THIRD_PROM_PORT=21301
export MSG_TRANSFER_PROM_PORT=21400
export MSG_TRANSFER_PROM_PORT=21401
export MSG_TRANSFER_PROM_PORT=21402
export MSG_TRANSFER_PROM_PORT=21403

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/install/common.sh"

openim::log::info "\n# Use Docker to start all openim service"

trap 'openim::util::onCtrlC' INT

"${OPENIM_ROOT}"/scripts/init-config.sh --skip

"${OPENIM_ROOT}"/scripts/start-all.sh

sleep 5

"${OPENIM_ROOT}"/scripts/check-all.sh

tail -f ${LOG_FILE}