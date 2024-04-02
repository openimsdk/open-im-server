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






OPENIM_SCRIPTS=$(dirname "${BASH_SOURCE[0]}")/
source "$OPENIM_SCRIPTS/bricks.sh"



result=$(check_binaries_running)
ret_val=$?
if [ $ret_val -eq 0 ]; then
    echo "All binaries are running."
else
    echo "$result"
    echo "abort..."
    exit 1
fi


print_listened_ports_by_binaries


