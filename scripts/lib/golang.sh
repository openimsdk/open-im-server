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


# shellcheck disable=SC2034 # Variables sourced in other scripts.

# The server platform we are building on.
readonly OPENIM_SUPPORTED_SERVER_PLATFORMS=(
  linux/amd64
  linux/arm64
)

# If we update this we should also update the set of platforms whose standard
# library is precompiled for in build/build-image/cross/Dockerfile
readonly OPENIM_SUPPORTED_CLIENT_PLATFORMS=(
  linux/amd64
  linux/arm64
)

# The set of server targets that we are only building for Linux
# If you update this list, please also update build/BUILD.
openim::golang::server_targets() {
  local targets=(
    openim-api
    openim-cmdutils
    openim-cmdutils
    openim-crontask
    openim-msggateway
    openim-msgtransfer
    openim-push
    openim-rpc-auth
    openim-rpc-conversation
    openim-rpc-friend
    openim-rpc-group
    openim-rpc-msg
    openim-rpc-third
    openim-rpc-user
  )
  echo "${targets[@]}"
}

IFS=" " read -ra OPENIM_SERVER_TARGETS <<< "$(openim::golang::server_targets)"
readonly OPENIM_SERVER_TARGETS
readonly OPENIM_SERVER_BINARIES=("${OPENIM_SERVER_TARGETS[@]##*/}")

# The set of server targets we build docker images for
openim::golang::server_image_targets() {
  # NOTE: this contains cmd targets for openim::build::get_docker_wrapped_binaries
  local targets=(
    cmd/openim-apiserver
    cmd/openim-authz-server
    cmd/openim-pump
    cmd/openim-watcher
  )
  echo "${targets[@]}"
}

IFS=" " read -ra OPENIM_SERVER_IMAGE_TARGETS <<< "$(openim::golang::server_image_targets)"
readonly OPENIM_SERVER_IMAGE_TARGETS
readonly OPENIM_SERVER_IMAGE_BINARIES=("${OPENIM_SERVER_IMAGE_TARGETS[@]##*/}")

# ------------
# NOTE: All functions that return lists should use newlines.
# bash functions can't return arrays, and spaces are tricky, so newline
# separators are the preferred pattern.
# To transform a string of newline-separated items to an array, use openim::util::read-array:
# openim::util::read-array FOO < <(openim::golang::dups a b c a)
#
# ALWAYS remember to quote your subshells. Not doing so will break in
# bash 4.3, and potentially cause other issues.
# ------------

# Returns a sorted newline-separated list containing only duplicated items.
openim::golang::dups() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort | uniq -d
}

# Returns a sorted newline-separated list with duplicated items removed.
openim::golang::dedup() {
  # We use printf to insert newlines, which are required by sort.
  printf "%s\n" "$@" | sort -u
}

# Depends on values of user-facing OPENIM_BUILD_PLATFORMS, OPENIM_FASTBUILD,
# and OPENIM_BUILDER_OS.
# Configures OPENIM_SERVER_PLATFORMS and OPENIM_CLIENT_PLATFORMS, then sets them
# to readonly.
# The configured vars will only contain platforms allowed by the
# OPENIM_SUPPORTED* vars at the top of this file.
declare -a OPENIM_SERVER_PLATFORMS
declare -a OPENIM_CLIENT_PLATFORMS
openim::golang::setup_platforms() {
  if [[ -n "${OPENIM_BUILD_PLATFORMS:-}" ]]; then
    # OPENIM_BUILD_PLATFORMS needs to be read into an array before the next
    # step, or quoting treats it all as one element.
    local -a platforms
    IFS=" " read -ra platforms <<< "${OPENIM_BUILD_PLATFORMS}"

    # Deduplicate to ensure the intersection trick with openim::golang::dups
    # is not defeated by duplicates in user input.
    openim::util::read-array platforms < <(openim::golang::dedup "${platforms[@]}")

    # Use openim::golang::dups to restrict the builds to the platforms in
    # OPENIM_SUPPORTED_*_PLATFORMS. Items should only appear at most once in each
    # set, so if they appear twice after the merge they are in the intersection.
    openim::util::read-array OPENIM_SERVER_PLATFORMS < <(openim::golang::dups \
        "${platforms[@]}" \
        "${OPENIM_SUPPORTED_SERVER_PLATFORMS[@]}" \
      )
    readonly OPENIM_SERVER_PLATFORMS

    openim::util::read-array OPENIM_CLIENT_PLATFORMS < <(openim::golang::dups \
        "${platforms[@]}" \
        "${OPENIM_SUPPORTED_CLIENT_PLATFORMS[@]}" \
      )
    readonly OPENIM_CLIENT_PLATFORMS

  elif [[ "${OPENIM_FASTBUILD:-}" == "true" ]]; then
    OPENIM_SERVER_PLATFORMS=(linux/amd64)
    readonly OPENIM_SERVER_PLATFORMS
    OPENIM_CLIENT_PLATFORMS=(linux/amd64)
    readonly OPENIM_CLIENT_PLATFORMS
  else
    OPENIM_SERVER_PLATFORMS=("${OPENIM_SUPPORTED_SERVER_PLATFORMS[@]}")
    readonly OPENIM_SERVER_PLATFORMS

    OPENIM_CLIENT_PLATFORMS=("${OPENIM_SUPPORTED_CLIENT_PLATFORMS[@]}")
    readonly OPENIM_CLIENT_PLATFORMS
  fi
}

openim::golang::setup_platforms

# The set of client targets that we are building for all platforms
# If you update this list, please also update build/BUILD.
readonly OPENIM_CLIENT_TARGETS=(
  imctl
)
readonly OPENIM_CLIENT_BINARIES=("${OPENIM_CLIENT_TARGETS[@]##*/}")

readonly OPENIM_ALL_TARGETS=(
  "${OPENIM_SERVER_TARGETS[@]}"
  "${OPENIM_CLIENT_TARGETS[@]}"
)
readonly OPENIM_ALL_BINARIES=("${OPENIM_ALL_TARGETS[@]##*/}")

# Asks golang what it thinks the host platform is. The go tool chain does some
# slightly different things when the target platform matches the host platform.
openim::golang::host_platform() {
  echo "$(go env GOHOSTOS)/$(go env GOHOSTARCH)"
}

# Ensure the go tool exists and is a viable version.
openim::golang::verify_go_version() {
  if [[ -z "$(command -v go)" ]]; then
    openim::log::usage_from_stdin <<EOF
Can't find 'go' in PATH, please fix and retry.
See http://golang.org/doc/install for installation instructions.
EOF
    return 2
  fi

  local go_version
  IFS=" " read -ra go_version <<< "$(go version)"
  local minimum_go_version
  minimum_go_version=go1.13.4
  if [[ "${minimum_go_version}" != $(echo -e "${minimum_go_version}\n${go_version[2]}" | sort -s -t. -k 1,1 -k 2,2n -k 3,3n | head -n1) && "${go_version[2]}" != "devel" ]]; then
    openim::log::usage_from_stdin <<EOF
Detected go version: ${go_version[*]}.
OpenIM requires ${minimum_go_version} or greater.
Please install ${minimum_go_version} or later.
EOF
    return 2
  fi
}

# openim::golang::setup_env will check that the `go` commands is available in
# ${PATH}. It will also check that the Go version is good enough for the
# OpenIM build.
#
# Outputs:
#   env-var GOBIN is unset (we want binaries in a predictable place)
#   env-var GO15VENDOREXPERIMENT=1
#   env-var GO111MODULE=on
openim::golang::setup_env() {
  openim::golang::verify_go_version

  # Unset GOBIN in case it already exists in the current session.
  unset GOBIN

  # This seems to matter to some tools
  export GO15VENDOREXPERIMENT=1

  # Open go module feature
  export GO111MODULE=on

  # This is for sanity.  Without it, user umasks leak through into release
  # artifacts.
  umask 0022
}
