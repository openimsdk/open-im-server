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
#
# OpenIM Tools Control Script
#
# Description:
#  This script is responsible for managing the lifecycle of OpenIM tools, which include starting, stopping,
#  and handling pre and post operations. It's designed to be modular and extensible, ensuring that the 
#  individual operations can be managed separately, and integrated seamlessly with Linux systemd.
# 
# Features:
# 1. Robust error handling using Bash built-ins like 'errexit', 'nounset', and 'pipefail'.
# 2. The capability to source common utility functions and configurations to ensure uniform environmental settings.
# 3. Comprehensive logging functionalities, providing a detailed understanding of operational processes.
# 4. Provision for declaring and managing a set of OpenIM tools, each associated with its unique name and corresponding ports.
# 5. The ability to define and associate Prometheus ports for service monitoring purposes.
# 6. Functionalities to start each OpenIM tool, along with its designated ports, in a sequence.
#
# Usage:
# 1. Direct Script Execution:
#    This initiates all the OpenIM tools declared under the function openim::tools::service_name.
#    Example: ./openim-tools.sh  openim::tools::start
# 2. Controlling through Functions for systemctl operations:
#    Specific operations like installation, uninstallation, and status check can be executed by passing the respective function name as an argument to the script.
#    Example: ./openim-tools.sh openim::tools::install
#

set -o errexit
set +o nounset
set -o pipefail

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source "${OPENIM_ROOT}"/scripts/install/common.sh

SERVER_NAME="openim-tools"

openim::tools::start_name() {
  local targets=(
    imctl
  )
  echo "${targets[@]}"
}
IFS=" " read -ra OPENIM_TOOLS_NAME_TARGETS <<< "$(openim::tools::start_name)"
readonly OPENIM_TOOLS_NAME_TARGETS
readonly OPENIM_TOOLS_NAME_LISTARIES=("${OPENIM_TOOLS_NAME_TARGETS[@]##*/}")

openim::tools::pre_start_name() {
  local targets=(
    ncpu
    component
  )
  echo "${targets[@]}"
}
IFS=" " read -ra OPENIM_TOOLS_PRE_START_NAME_TARGETS <<< "$(openim::tools::pre_start_name)"
readonly OPENIM_TOOLS_PRE_START_NAME_TARGETS
readonly OPENIM_TOOLS_PRE_START_NAME_LISTARIES=("${OPENIM_TOOLS_PRE_START_NAME_TARGETS[@]##*/}")

openim::tools::post_start_name() {
  local targets=(
    infra
    versionchecker
  )
  echo "${targets[@]}"
}
IFS=" " read -ra OPENIM_TOOLS_POST_START_NAME_TARGETS <<< "$(openim::tools::post_start_name)"
readonly OPENIM_TOOLS_POST_START_NAME_TARGETS
readonly OPENIM_TOOLS_POST_START_NAME_LISTARIES=("${OPENIM_TOOLS_POST_START_NAME_TARGETS[@]##*/}")

function openim::tools::start_service() {
  local binary_name="$1"
  local config="$2"
  local service_port="$3"
  local prometheus_port="$4"

  local cmd="${OPENIM_OUTPUT_HOSTBIN_TOOLS}/${binary_name}"
  openim::log::info "Starting PATH: ${OPENIM_OUTPUT_HOSTBIN_TOOLS}/${binary_name}..."

  if [ -n "${config}" ]; then
    printf "Specifying config: %s\n" "${config}"
    cmd="${cmd} -c ${config}/config.yaml"
  fi

  if [ -n "${service_port}" ]; then
    printf "Specifying service port: %s\n" "${service_port}"
    cmd="${cmd} --port ${service_port}"
  fi

  if [ -n "${prometheus_port}" ]; then
    printf "Specifying prometheus port: %s\n" "${prometheus_port}"
    cmd="${cmd} --prometheus_port ${prometheus_port}"
  fi
  openim::log::info "Starting ${binary_name}..."
  ${cmd}
}

function openim::tools::start() {
    openim::log::info "Starting OpenIM Tools..."
    for tool in "${OPENIM_TOOLS_NAME_LISTARIES[@]}"; do
        openim::log::info "Starting ${tool}..."
        # openim::tools::start_service ${tool}
        sleep 0.2
    done
}


function openim::tools::pre-start() {
    openim::log::info "Preparing to start OpenIM Tools..."
    for tool in "${OPENIM_TOOLS_PRE_START_NAME_LISTARIES[@]}"; do
        openim::log::info "Starting ${tool}..."
        openim::tools::start_service ${tool} ${OPNEIM_CONFIG}
        sleep 0.2
    done
}

function openim::tools::post-start() {
    openim::log::info "Post-start actions for OpenIM Tools..."
    for tool in "${OPENIM_TOOLS_POST_START_NAME_LISTARIES[@]}"; do
        openim::log::info "Starting ${tool}..."
        openim::tools::start_service ${tool}
        sleep 0.2
    done
}

function openim::tools::stop() {
    openim::log::info "Stopping OpenIM Tools..."
    for tool in "${OPENIM_TOOLS_NAME_LISTARIES[@]}"; do
        openim::log::info "Stopping ${tool}..."
        # Similarly, place the actual command to stop the tool here.
        echo "Stopping service for ${tool}"
        sleep 0.2
    done
}

function openim::tools::pre-stop() {
    openim::log::info "Preparing to stop OpenIM Tools..."
    for tool in "${OPENIM_TOOLS_PRE_START_NAME_LISTARIES[@]}"; do
        openim::log::info "Setting up pre-stop for ${tool}..."
        echo "Pre-stop actions for ${tool}"
        sleep 0.2
    done
}

function openim::tools::post-stop() {
    openim::log::info "Post-stop actions for OpenIM Tools..."
    for tool in "${OPENIM_TOOLS_POST_START_NAME_LISTARIES[@]}"; do
        openim::log::info "Executing post-stop for ${tool}..."
        echo "Post-stop cleanup for ${tool}"
        sleep 0.2
    done
}

if [[ "$*" =~ openim::tools:: ]];then
  eval $*
fi
