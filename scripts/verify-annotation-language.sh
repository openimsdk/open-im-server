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


# This script verifies whether codes follow golang convention.
# Usage: `scripts/verify-pkg-names.sh`.

set -o errexit

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

openim::golang::verify_go_version

openim::golang::verify_go_version

OPENIM_OUTPUT_HOSTBIN_TOOLS="${OPENIM_ROOT}/_output/bin/tools/linux/amd64"
CODESCAN_BINARY="${OPENIM_OUTPUT_HOSTBIN_TOOLS}/codescan"

if [[ ! -f "${CODESCAN_BINARY}" ]]; then
    echo "codescan binary not found, building..."
    pushd "${OPENIM_ROOT}" >/dev/null
    make build BINS="codescan"
    popd >/dev/null
fi

if [[ ! -f "${CODESCAN_BINARY}" ]]; then
    echo "Failed to build codescan binary."
    exit 1
fi

CONFIG_PATH="${OPENIM_ROOT}/tools/codescan/config.yaml"

"${CODESCAN_BINARY}" -config "${CONFIG_PATH}"