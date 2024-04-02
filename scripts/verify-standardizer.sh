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

# This script does a fast type check of script srnetes code for all platforms.
# Usage: `scripts/verify-standardizer.sh`.

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

openim::golang::verify_go_version

cd "${OPENIM_ROOT}"
ret=0
scripts/run-in-gopath.sh \
make tools.verify.standardizer
${OPENIM_ROOT}/_output/tools/standardizer || ret=$?
if [[ $ret -ne 0 ]]; then
  openim::log::error "Failed to check the directory name or file name. Your name may not meet the specification. Please check the configuration file and the directory or file name." >&2
  openim::log::error "Please see https://github.com/kubecub/standardizer for more information." >&2
  exit 1
fi
