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
set +o nounset
set -o pipefail

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source ${OPENIM_ROOT}/scripts/install/common.sh

readonly SERVER_NAME="openim-crontask"
readonly SERVER_PATH="${OPENIM_OUTPUT_HOSTBIN}/openim-crontask"

openim::log::status "Start OpenIM Cron, binary root: ${SERVER_NAME}"
openim::log::info "Start OpenIM Cron, path: ${SERVER_PATH}"

openim::util::stop_services_with_name ${SERVER_NAME}

sleep 1

openim::log::status "start cron_task process"

nohup ./${cron_task_name}  >>$OPENIM_ROOT/logs/openIM.log 2>&1 &
#done

#Check launched service process
check=`ps  | grep -w ./${cron_task_name} | grep -v grep| wc -l`
if [ $check -ge 1 ]
then
newPid=`ps  | grep -w ./${cron_task_name} | grep -v grep|awk '{print $2}'`
allPorts=""
    echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS "${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${cron_task_name}${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${newPid}${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${BACKGROUND_GREEN}${allPorts}${COLOR_SUFFIX}
else
    echo -e ${BACKGROUND_GREEN}${cron_task_name}${COLOR_SUFFIX}${RED_PREFIX}"\n SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
