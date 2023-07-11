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

# Unset CDPATH so that path interpolation can work correctly
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
OPENIM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${OPENIM_ROOT}/scripts/lib/util.sh"
source "${OPENIM_ROOT}/scripts/lib/logging.sh"
source "${OPENIM_ROOT}/scripts/lib/color.sh"

openim::log::install_errexit

source "${OPENIM_ROOT}/scripts/lib/version.sh"
source "${OPENIM_ROOT}/scripts/lib/golang.sh"
