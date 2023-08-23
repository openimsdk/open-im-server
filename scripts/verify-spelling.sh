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

# This script checks commonly misspelled English words in all files in the
# working directory by client9/misspell package.
# Usage: `scripts/verify-spelling.sh`.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
export OPENIM_ROOT
source "${OPENIM_ROOT}/scripts/lib/init.sh"

# Ensure that we find the binaries we build before anything else.
export GOBIN="${KUBE_OUTPUT_BINPATH}"
PATH="${GOBIN}:${PATH}"

# Install tools we need
pushd ""${OPENIM_ROOT}"/tools" >/dev/null
  GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
popd >/dev/null

# Spell checking
# All the skipping files are defined in scripts/.spelling_failures
skipping_file=""${OPENIM_ROOT}"/scripts/.spelling_failures"
failing_packages=$(sed "s| | -e |g" "${skipping_file}")
git ls-files | grep -v -e "${failing_packages}" | xargs misspell -i "Creater,creater,ect" -error -o stderr
