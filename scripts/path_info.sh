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

# Determine the architecture and version
architecture=$(uname -m)
version=$(uname -s | tr '[:upper:]' '[:lower:]')

#Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/style_info.sh

cd $SCRIPTS_ROOT

# Define the supported architectures and corresponding bin directories
declare -A supported_architectures=(
    ["linux-amd64"]="_output/bin/platforms/linux/amd64"
    ["linux-arm64"]="_output/bin/platforms/linux/arm64"
    ["linux-mips64"]="_output/bin/platforms/linux/mips64"
    ["linux-mips64le"]="_output/bin/platforms/linux/mips64le"
    ["linux-ppc64le"]="_output/bin/platforms/linux/ppc64le"
    ["linux-s390x"]="_output/bin/platforms/linux/s390x"
    ["darwin-amd64"]="_output/bin/platforms/darwin/amd64"
    ["windows-amd64"]="_output/bin/platforms/windows/amd64"
    ["linux-x86_64"]="_output/bin/platforms/linux/amd64"  # Alias for linux-amd64
    ["darwin-x86_64"]="_output/bin/platforms/darwin/amd64"  # Alias for darwin-amd64
)

# Check if the architecture and version are supported
if [[ -z ${supported_architectures["$version-$architecture"]} ]]; then
    echo -e "${BLUE_PREFIX}================> Unsupported architecture: $architecture or version: $version${COLOR_SUFFIX}"
    exit 1
fi

echo -e "${BLUE_PREFIX}================> Architecture: $architecture${COLOR_SUFFIX}"

# Set the BIN_DIR based on the architecture and version
BIN_DIR=${supported_architectures["$version-$architecture"]}

echo -e "${BLUE_PREFIX}================> BIN_DIR: $OPENIM_ROOT/$BIN_DIR${COLOR_SUFFIX}"

# Don't put the space between "="
openim_msggateway="openim-msggateway"
msg_gateway_binary_root="$OPENIM_ROOT/$BIN_DIR"
msg_gateway_source_root="$OPENIM_ROOT/cmd/openim-msggateway/"

msg_name="openim-rpc-msg"
msg_binary_root="$OPENIM_ROOT/$BIN_DIR"
msg_source_root="$OPENIM_ROOT/cmd/openim-rpc/openim-rpc-msg/"

push_name="openim-push"
push_binary_root="$OPENIM_ROOT/$BIN_DIR"
push_source_root="$OPENIM_ROOT/cmd/openim-push/"

openim_msgtransfer="openim-msgtransfer"
msg_transfer_binary_root="$OPENIM_ROOT/$BIN_DIR"
msg_transfer_source_root="$OPENIM_ROOT/cmd/openim-msgtransfer/"
msg_transfer_service_num=4

cron_task_name="openim-crontask"
cron_task_binary_root="$OPENIM_ROOT/$BIN_DIR"
cron_task_source_root="$OPENIM_ROOT/cmd/openim-crontask/"

cmd_utils_name="openim-cmdutils"
cmd_utils_binary_root="$OPENIM_ROOT/$BIN_DIR"
cmd_utils_source_root="$OPENIM_ROOT/cmd/openim-cmdutils/"

# Global configuration file default dir
config_path="$OPENIM_ROOT/config/config.yaml"
configfile_path="$OPENIM_ROOT/config"
log_path="$OPENIM_ROOT/log"

# servicefile dir path
service_source_root=(
  # api service file
  "$OPENIM_ROOT/cmd/api/"
  # rpc service file
  "$OPENIM_ROOT/cmd/openim-rpc/openim-rpc-user/"
  "$OPENIM_ROOT/cmd/openim-rpc/openim-rpc-friend/"
  "$OPENIM_ROOT/cmd/openim-rpc/openim-rpc-group/"
  "$OPENIM_ROOT/cmd/openim-rpc/openim-rpc-auth/"
  "$OPENIM_ROOT/cmd/openim-rpc/openim-rpc-conversation/"
  "$OPENIM_ROOT/cmd/openim-rpc/openim-rpc-third/"
  "$OPENIM_ROOT/cmd/openim-crontask"
  "${msg_gateway_source_root}"
  "${msg_transfer_source_root}"
  "${msg_source_root}"
  "${push_source_root}"
  # "${sdk_server_source_root}"
)

# service filename
service_names=(
  # api service filename
  "openim-api"
  # rpc service filename
  "openim-rpc-user"
  "openim-rpc-friend"
  "openim-rpc-group"
  "openim-rpc-auth"
  "openim-rpc-conversation"
  "openim-rpc-third"
  "openim-crontask"
  "${openim_msggateway}"
  "${openim_msgtransfer}"
  "${msg_name}"
  "${push_name}"
  # "${sdk_server_name}"
)
