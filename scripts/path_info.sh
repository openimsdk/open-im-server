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

# Don't put the space between "="
openim_msggateway="openim-msggateway"
msg_gateway_binary_root="$OPENIM_ROOT/$BIN_DIR"

msg_name="openim-rpc-msg"
msg_binary_root="$OPENIM_ROOT/$BIN_DIR"

push_name="openim-push"
push_binary_root="$OPENIM_ROOT/$BIN_DIR"
push_source_root="$OPENIM_ROOT/cmd/openim-push/"

openim_msgtransfer="openim-msgtransfer"
msg_transfer_binary_root="$OPENIM_ROOT/$BIN_DIR"
msg_transfer_service_num=4

cron_task_name="openim-crontask"
cron_task_binary_root="$OPENIM_ROOT/$BIN_DIR"

cmd_utils_name="openim-cmdutils"
cmd_utils_binary_root="$OPENIM_ROOT/$BIN_DIR"

# Global configuration file default dir
config_path="$OPENIM_ROOT/config/config.yaml"
configfile_path="$OPENIM_ROOT/config"

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