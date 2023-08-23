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


# This script lints each shell script by `shellcheck`.
# Usage: `scripts/verify-shellcheck.sh`.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

# allow overriding docker cli, which should work fine for this script
DOCKER="${DOCKER:-docker}"

# required version for this script, if not installed on the host we will
# use the official docker image instead. keep this in sync with SHELLCHECK_IMAGE
SHELLCHECK_VERSION="0.8.0"
SHELLCHECK_IMAGE="docker.io/koalaman/shellcheck-alpine:v0.8.0@sha256:f42fde76d2d14a645a848826e54a4d650150e151d9c81057c898da89a82c8a56"

# disabled lints
disabled=(
  # this lint disallows non-constant source, which we use extensively without
  # any known bugs
  1090
  # this lint warns when shellcheck cannot find a sourced file
  # this wouldn't be a bad idea to warn on, but it fails on lots of path
  # dependent sourcing, so just disable enforcing it
  1091
  # this lint prefers command -v to which, they are not the same
  2230
)
# comma separate for passing to shellcheck
join_by() {
  local IFS="$1";
  shift;
  echo "$*";
}
SHELLCHECK_DISABLED="$(join_by , "${disabled[@]}")"
readonly SHELLCHECK_DISABLED

# ensure we're linting the k8s source tree
cd "${OPENIM_ROOT}"

# Find all shell scripts excluding:
# - Anything git-ignored - No need to lint untracked files.
# - ./_* - No need to lint output directories.
# - ./.git/* - Ignore anything in the git object store.
# - ./vendor* - Vendored code should be fixed upstream instead.
# - ./third_party/*, but re-include ./third_party/forked/*  - only code we
#    forked should be linted and fixed.
all_shell_scripts=()
while IFS=$'\n' read -r script;
  do git check-ignore -q "$script" || all_shell_scripts+=("$script");
done < <(find . -name "*.sh" \
  -not \( \
    -path ./_\*      -o \
    -path ./.git\*   -o \
    -path ./vendor\* -o \
    \( -path ./third_party\* -a -not -path ./third_party/forked\* \) \
  \))

# detect if the host machine has the required shellcheck version installed
# if so, we will use that instead.
HAVE_SHELLCHECK=false
if which shellcheck &>/dev/null; then
  detected_version="$(shellcheck --version | grep 'version: .*')"
  if [[ "${detected_version}" = "version: ${SHELLCHECK_VERSION}" ]]; then
    HAVE_SHELLCHECK=true
  fi
fi

# if KUBE_JUNIT_REPORT_DIR is set, disable colorized output.
# Colorized output causes malformed XML in the JUNIT report.
SHELLCHECK_COLORIZED_OUTPUT="auto"
if [[ -n "${KUBE_JUNIT_REPORT_DIR:-}" ]]; then
  SHELLCHECK_COLORIZED_OUTPUT="never"
fi

# common arguments we'll pass to shellcheck
SHELLCHECK_OPTIONS=(
  # allow following sourced files that are not specified in the command,
  # we need this because we specify one file at a time in order to trivially
  # detect which files are failing
  "--external-sources"
  # include our disabled lints
  "--exclude=${SHELLCHECK_DISABLED}"
  # set colorized output
  "--color=${SHELLCHECK_COLORIZED_OUTPUT}"
)

# tell the user which we've selected and lint all scripts
# The shellcheck errors are printed to stdout by default, hence they need to be redirected
# to stderr in order to be well parsed for Junit representation by juLog function
res=0
if ${HAVE_SHELLCHECK}; then
  openim::log::info "Using host shellcheck ${SHELLCHECK_VERSION} binary."
  shellcheck "${SHELLCHECK_OPTIONS[@]}" "${all_shell_scripts[@]}" >&2 || res=$?
else
  openim::log::info "Using shellcheck ${SHELLCHECK_VERSION} docker image."
  "${DOCKER}" run \
    --rm -v "${OPENIM_ROOT}:${OPENIM_ROOT}" -w "${OPENIM_ROOT}" \
    "${SHELLCHECK_IMAGE}" \
  shellcheck "${SHELLCHECK_OPTIONS[@]}" "${all_shell_scripts[@]}" >&2 || res=$?
fi

# print a message based on the result
if [ $res -eq 0 ]; then
  echo 'Congratulations! All shell files are passing lint :-)'
else
  {
    echo
    echo 'Please review the above warnings. You can test via "./scripts/verify-shellcheck.sh"'
    echo 'If the above warnings do not make sense, you can exempt this warning with a comment'
    echo ' (if your reviewer is okay with it).'
    echo 'In general please prefer to fix the error, we have already disabled specific lints'
    echo ' that the project chooses to ignore.'
    echo 'See: https://github.com/koalaman/shellcheck/wiki/Ignore#ignoring-one-specific-instance-in-a-file'
    echo
  } >&2
  exit 1
fi

# preserve the result
exit $res
