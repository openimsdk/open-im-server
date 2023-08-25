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

set -o errexit
set +o nounset
set -o pipefail

# Short-circuit if init.sh has already been sourced
[[ $(type -t openim::init::loaded) == function ]] && return 0

# Unset CDPATH so that path interpolation can work correctly
unset CDPATH

# Until all GOPATH references are removed from all build scripts as well,
# explicitly disable module mode to avoid picking up user-set GO111MODULE preferences.
# As individual scripts (like hack/update-vendor.sh) make use of go modules,
# they can explicitly set GO111MODULE=on
export GO111MODULE=on

# The root of the build/dist directory
OPENIM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

OPENIM_OUTPUT_SUBPATH="${OPENIM_OUTPUT_SUBPATH:-_output}"
OPENIM_OUTPUT=""${OPENIM_ROOT}"/${OPENIM_OUTPUT_SUBPATH}"

OPENIM_OUTPUT_BINPATH="${OPENIM_OUTPUT}/bin/platforms"
OPENIM_OUTPUT_BINTOOLPATH="${OPENIM_OUTPUT}/bin/tools"
OPENIM_OUTPUT_TOOLS="${OPENIM_OUTPUT}/tools"
OPENIM_OUTPUT_TMP="${OPENIM_OUTPUT}/tmp"
OPENIM_OUTPUT_LOGS="${OPENIM_OUTPUT}/logs"

# This controls rsync compression. Set to a value > 0 to enable rsync
# compression for build container
OPENIM_RSYNC_COMPRESS="${KUBE_RSYNC_COMPRESS:-0}"

# Set no_proxy for localhost if behind a proxy, otherwise,
# the connections to localhost in scripts will time out
export no_proxy="127.0.0.1,localhost${no_proxy:+,${no_proxy}}"

# This is a symlink to binaries for "this platform", e.g. build tools.
export THIS_PLATFORM_BIN=""${OPENIM_ROOT}"/_output/bin/platforms"
export THIS_PLATFORM_BIN_TOOLS=""${OPENIM_ROOT}"/_output/bin/tools"

. $(dirname ${BASH_SOURCE})/color.sh
. $(dirname ${BASH_SOURCE})/util.sh
. $(dirname ${BASH_SOURCE})/logging.sh

openim::log::install_errexit
openim::util::ensure-bash-version

. $(dirname ${BASH_SOURCE})/version.sh
. $(dirname ${BASH_SOURCE})/golang.sh
. $(dirname ${BASH_SOURCE})/release.sh
. $(dirname ${BASH_SOURCE})/chat.sh

OPENIM_OUTPUT_HOSTBIN="${OPENIM_OUTPUT_BINPATH}/$(openim::util::host_platform)"
export OPENIM_OUTPUT_HOSTBIN
OPENIM_OUTPUT_HOSTBIN_TOOLS="${OPENIM_OUTPUT_BINTOOLPATH}/$(openim::util::host_platform)"
export OPENIM_OUTPUT_HOSTBIN_TOOLS

export OPENIM_NONSERVER_GROUP_VERSIONS


# This emulates "readlink -f" which is not available on MacOS X.
# Test:
# T=/tmp/$$.$RANDOM
# mkdir $T
# touch $T/file
# mkdir $T/dir
# ln -s $T/file $T/linkfile
# ln -s $T/dir $T/linkdir
# function testone() {
#   X=$(readlink -f $1 2>&1)
#   Y=$(kube::readlinkdashf $1 2>&1)
#   if [ "$X" != "$Y" ]; then
#     echo readlinkdashf $1: expected "$X", got "$Y"
#   fi
# }
# testone /
# testone /tmp
# testone $T
# testone $T/file
# testone $T/dir
# testone $T/linkfile
# testone $T/linkdir
# testone $T/nonexistant
# testone $T/linkdir/file
# testone $T/linkdir/dir
# testone $T/linkdir/linkfile
# testone $T/linkdir/linkdir
function openim::readlinkdashf {
  # run in a subshell for simpler 'cd'
  (
    if [[ -d "${1}" ]]; then # This also catch symlinks to dirs.
      cd "${1}"
      pwd -P
    else
      cd "$(dirname "${1}")"
      local f
      f=$(basename "${1}")
      if [[ -L "${f}" ]]; then
        readlink "${f}"
      else
        echo "$(pwd -P)/${f}"
      fi
    fi
  )
}

# This emulates "readlink -f" which is not available on MacOS X.
# Test:
# T=/tmp/$$.$RANDOM
# mkdir $T
# touch $T/file
# mkdir $T/dir
# ln -s $T/file $T/linkfile
# ln -s $T/dir $T/linkdir
# function testone() {
#   X=$(readlink -f $1 2>&1)
#   Y=$(kube::readlinkdashf $1 2>&1)
#   if [ "$X" != "$Y" ]; then
#     echo readlinkdashf $1: expected "$X", got "$Y"
#   fi
# }
# testone /
# testone /tmp
# testone $T
# testone $T/file
# testone $T/dir
# testone $T/linkfile
# testone $T/linkdir
# testone $T/nonexistant
# testone $T/linkdir/file
# testone $T/linkdir/dir
# testone $T/linkdir/linkfile
# testone $T/linkdir/linkdir
function openim::readlinkdashf {
  # run in a subshell for simpler 'cd'
  (
    if [[ -d "${1}" ]]; then # This also catch symlinks to dirs.
      cd "${1}"
      pwd -P
    else
      cd "$(dirname "${1}")"
      local f
      f=$(basename "${1}")
      if [[ -L "${f}" ]]; then
        readlink "${f}"
      else
        echo "$(pwd -P)/${f}"
      fi
    fi
  )
}

# This emulates "realpath" which is not available on MacOS X
# Test:
# T=/tmp/$$.$RANDOM
# mkdir $T
# touch $T/file
# mkdir $T/dir
# ln -s $T/file $T/linkfile
# ln -s $T/dir $T/linkdir
# function testone() {
#   X=$(realpath $1 2>&1)
#   Y=$(kube::realpath $1 2>&1)
#   if [ "$X" != "$Y" ]; then
#     echo realpath $1: expected "$X", got "$Y"
#   fi
# }
# testone /
# testone /tmp
# testone $T
# testone $T/file
# testone $T/dir
# testone $T/linkfile
# testone $T/linkdir
# testone $T/nonexistant
# testone $T/linkdir/file
# testone $T/linkdir/dir
# testone $T/linkdir/linkfile
# testone $T/linkdir/linkdir
openim::realpath() {
  if [[ ! -e "${1}" ]]; then
    echo "${1}: No such file or directory" >&2
    return 1
  fi
  openim::readlinkdashf "${1}"
}

# Marker function to indicate init.sh has been fully sourced
openim::init::loaded() {
  return 0
}