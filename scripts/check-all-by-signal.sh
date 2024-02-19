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

handle_error() {
  echo "An error occurred. Printing ${STDERR_LOG_FILE} contents:"
  cat "${STDERR_LOG_FILE}"
  exit 1
}

trap handle_error ERR

. $(dirname ${BASH_SOURCE})/install/openim-msgtransfer.sh openim::msgtransfer::check_by_signal

# Assuming openim::util::check_ports_by_signal function sets a proper exit status
# based on whether services are running or not.
if openim::util::check_ports_by_signal ${OPENIM_SERVER_PORT_LISTARIES[@]}; then
  echo "+++ cat openim log file >>> ${LOG_FILE}"
    openim::log::error_exit "The service does not stop properly, there are still processes running, please check!"
else
  echo "++++ All openim service ports stop successfully !"
fi

set -e

trap - ERR