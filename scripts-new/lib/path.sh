
#!/usr/bin/env bash

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/lib/util.sh"


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
  local bin_full_path="${project_path}/get_bin_path/${host_platform}/${bin_name}"
  echo "${bin_full_path}" 1111111111111
}



