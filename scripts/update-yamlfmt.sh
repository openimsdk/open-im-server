#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

cd "${OPENIM_ROOT}"

find_files() {
  find . -not \( \
      \( \
        -wholename './output' \
        -o -wholename './.git' \
        -o -wholename './_output' \
        -o -wholename './_gopath' \
        -o -wholename './release' \
        -o -wholename './target' \
        -o -wholename '*/vendor/*' \
      \) -prune \
    \) -name 'OWNERS*'
}

export GO111MODULE=on
find_files | xargs go run tools/yamlfmt/main.go