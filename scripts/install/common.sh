#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# Common utilities, variables and checks for all build scripts.
set -o errexit
set +o nounset
set -o pipefail

# Sourced flag
COMMON_SOURCED=true

# The root of the build/dist directory
OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)

source "${OPENIM_ROOT}/scripts/lib/init.sh"
source "${OPENIM_ROOT}/scripts/install/environment.sh"

# Execute commands that require root permission without entering a password
function openim::common::sudo {
  echo ${LINUX_PASSWORD} | sudo -S $1
}
