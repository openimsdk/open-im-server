
#!/usr/bin/env bash

source "$(dirname "${BASH_SOURCE[0]}")/../lib/util.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/init.sh"


get_conf_dir() {
  local project_path="$1"
  echo "${project_path}/conf/"
}

get_log_dir() {
  local project_path="$1"
  echo "${project_path}/_output/logs/"
}

get_bin_dir() {
   local project_path="$1"
    echo "${project_path}/_output/bin/"
}


get_bin_full_path() {
  local project_path="$1"
  local bin_name="$2"

  local host_platform=$(openim::util::host_platform)

  local bin_dir=$(get_bin_dir "$project_path")

  local bin_full_path="${OPENIM_OUTPUT_BINPATH}/${bin_name}"
  echo ${bin_full_path}
}



