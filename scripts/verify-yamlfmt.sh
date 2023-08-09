#!/usr/bin/env bash

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