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

# This script is stop all openim service
#
# Usage: `scripts/stop.sh`.
# Encapsulated as: `make stop`.





OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${OPENIM_ROOT}/scripts/install/common.sh"

openim::log::info "\n# Begin to stop all openim service"

echo -e "\n++ Stop all processes in the path ${OPENIM_OUTPUT_HOSTBIN}"

openim::util::stop_services_with_name "${OPENIM_OUTPUT_HOSTBIN}"
# todo OPENIM_ALL_SERVICE_LIBRARIES





result=$(openim::util::check_process_names_for_stop)
if [[ $? -ne 0 ]]; then
  echo "+++ cat openim log file >>> ${LOG_FILE}"
  openim::log::error "stop process failed.\n" "${result}"
  return 1
fi

openim::log::success "✨  All processes to be stopped"