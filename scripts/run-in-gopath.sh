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


# This script sets up a temporary openim GOPATH and runs an arbitrary
# command under it. Go tooling requires that the current directory be under
# GOPATH or else it fails to find some things, such as the vendor directory for
# the project.
# Usage: `scripts/run-in-gopath.sh <command>`.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

# This sets up a clean GOPATH and makes sure we are currently in it.
openim::golang::setup_env

# Run the user-provided command.
"${@}"
