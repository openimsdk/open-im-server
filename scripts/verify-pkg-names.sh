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


OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

openim::golang::verify_go_version

cd "${OPENIM_ROOT}"
if git --no-pager grep -E $'^(import |\t)[a-z]+[A-Z_][a-zA-Z]* "[^"]+"$' -- '**/*.go' ':(exclude)vendor/*' ':(exclude)**/*.pb.go'; then
  openim::log::error "Some package aliases break go conventions."
  echo "To fix these errors, do not use capitalized or underlined characters"
  echo "in pkg aliases. Refer to https://blog.golang.org/package-names for more info."
  exit 1
fi
