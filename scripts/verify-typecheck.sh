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
# Usage: `scripts/verify-typecheck.sh`.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

openim::golang::verify_go_version

cd "${OPENIM_ROOT}"

# As of June, 2020 the typecheck tool is written in terms of go/packages, but
# that library doesn't work well with multiple modules.  Until that is done,
# force this tooling to run in a fake GOPATH.
ret=0
TYPECHECK_SERIAL="${TYPECHECK_SERIAL:-false}"
scripts/run-in-gopath.sh \
    go run test/typecheck/typecheck.go "$@" "--serial=$TYPECHECK_SERIAL" || ret=$?
if [[ $ret -ne 0 ]]; then
  openim::log::error "Type Check has failed. This may cause cross platform build failures." >&2
  openim::log::error "Please see https://github.com/openimsdk/open-im-server/tree/main/test/typecheck for more information." >&2
  exit 1
fi
