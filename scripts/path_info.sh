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