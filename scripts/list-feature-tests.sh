#!/usr/bin/env bash

# This script lists all of the [Feature:.+] tests in our e2e suite.
#
# Usage: `scripts/list-feature-tests.sh`.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
grep "\[Feature:\w+\]" "${OPENIM_ROOT}"/test/e2e/**/*.go -Eoh | LC_ALL=C sort -u