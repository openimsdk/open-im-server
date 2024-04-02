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






OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/
source "${OPENIM_ROOT}/lib/util.sh"
source "${OPENIM_ROOT}/define/binaries.sh"
source "${OPENIM_ROOT}/lib/path.sh"


for binary in "${!binaries[@]}"; do
  expected_count=${binaries[$binary]}
  full_path=$(get_bin_full_path "$binary")

  result=$(openim::util::check_process_names "$full_path" "$expected_count")

 if [ "$result" -eq 0 ]; then
     echo "Startup successful for $binary"
   else
     echo "Startup failed for $binary, $result processes missing."
   fi
done

for binary in "${!binaries[@]}"; do
  expected_count=${binaries[$binary]}
  base_path=$(get_bin_full_path "$binary")

  for ((i=0; i<expected_count; i++)); do
    full_path="${base_path} -i ${i} -c $OPENIM_OUTPUT_CONFIG"
    check_binary_ports "$full_path"
  done
done





