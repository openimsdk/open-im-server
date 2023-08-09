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


# This script checks whether the OWNERS files need to be formatted or not by
# `yamlfmt`. Run `scripts/update-yamlfmt.sh` to actually format sources.
#
# Usage: `scripts/verify-yamlfmt.sh`.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

openim::util::ensure_clean_working_dir
# This sets up the environment, like GOCACHE, which keeps the worktree cleaner.
openim::golang::setup_env

_tmpdir="$(openim::realpath "$(mktemp -d -t "$(basename "$0").XXXXXX")")"
git worktree add -f -q "${_tmpdir}" HEAD
openim::util::trap_add "git worktree remove -f ${_tmpdir}" EXIT
cd "${_tmpdir}"

# Format YAML files
hack/update-yamlfmt.sh

# Test for diffs
diffs=$(git status --porcelain | wc -l)
if [[ ${diffs} -gt 0 ]]; then
  echo "YAML files need to be formatted" >&2
  git diff
  echo "Please run 'hack/update-yamlfmt.sh'" >&2
  exit 1
fi