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

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd -P)
source "${OPENIM_ROOT}/scripts/install/common.sh"

openim::log::info "\n# Begin Install OpenIM Config"

for file in "${OPENIM_SERVER_TARGETS[@]}"; do
    VARNAME="$(echo $file | tr '[:lower:]' '[:upper:]' | tr '.' '_' | tr '-' '_')"
    VARVALUE="$OPENIM_OUTPUT_HOSTBIN/$file"
    # /etc/profile.d/openim-env.sh
    echo "export $VARNAME=$VARVALUE" > /etc/profile.d/openim-env.sh
    source /etc/profile.d/openim-env.sh
done
