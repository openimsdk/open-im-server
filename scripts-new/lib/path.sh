
#!/usr/bin/env bash

source "$(dirname "${BASH_SOURCE[0]}")/../lib/init.sh"




get_bin_full_path() {
  local bin_name="$1"
  local bin_full_path="${OPENIM_OUTPUT_HOSTBIN}/${bin_name}"
  echo ${bin_full_path}
}



get_tool_full_path() {
  local tool_name="$1"
  local tool_full_path="${OPENIM_OUTPUT_HOSTBIN_TOOLS}/${bin_name}"
  echo ${tool_full_path}
}
