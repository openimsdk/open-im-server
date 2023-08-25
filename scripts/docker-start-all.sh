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

#fixme This scripts is the total startup scripts
#fixme The full name of the shell scripts that needs to be started is placed in the need_to_start_server_shell array

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/install/common.sh"

trap 'openim::util::onCtrlC' INT

nohup "${OPENIM_ROOT}"/scripts/start-all.sh >> ${LOG_FILE} 2>&1 &

sleep 15

nohup "${OPENIM_ROOT}"/scripts/check-all.sh >> ${LOG_FILE} 2>&1 &

tail -f ${LOG_FILE}