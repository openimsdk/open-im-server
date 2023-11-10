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

# this script is used to check whether the code is formatted by gofmt or not
#
# Usage: source scripts/lib/util.sh
################################################################################

# TODO Debug: Just for testing, please comment out
# OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
# source "${OPENIM_ROOT}/scripts/lib/logging.sh"

#1、将IP写在一个文件里，比如文件名为hosts_file，一行一个IP地址。
#2、修改ssh-mutual-trust.sh里面的用户名及密码，默认为root用户及密码123。
# hosts_file_path="path/to/your/hosts/file"
# openim:util::setup_ssh_key_copy "$hosts_file_path" "root" "123"
function openim:util::setup_ssh_key_copy() {
  local hosts_file="$1"
  local username="${2:-root}"
  local password="${3:-123}"

  local sshkey_file=~/.ssh/id_rsa.pub

  # check sshkey file 
  if [[ ! -e $sshkey_file ]]; then
    expect -c "
    spawn ssh-keygen -t rsa
    expect \"Enter*\" { send \"\n\"; exp_continue; }
    "
  fi

  # get hosts list
  local hosts=$(awk '/^[^#]/ {print $1}' "${hosts_file}")

  ssh_key_copy() {
    local target=$1

    # delete history
    sed -i "/$target/d" ~/.ssh/known_hosts

    # copy key 
    expect -c "
    set timeout 100
    spawn ssh-copy-id $username@$target
    expect {
      \"yes/no\" { send \"yes\n\"; exp_continue; }
      \"*assword\" { send \"$password\n\"; }
      \"already exist on the remote system\" { exit 1; }
    }
    expect eof
    "
  }

  # auto sshkey pair
  for host in $hosts; do
    if ! ping -i 0.2 -c 3 -W 1 "$host" > /dev/null 2>&1; then
      echo "[ERROR]: Can't connect $host"
      continue
    fi

    local host_entry=$(awk "/$host/"'{print $1, $2}' /etc/hosts)
    if [[ $host_entry ]]; then
      local hostaddr=$(echo "$host_entry" | awk '{print $1}')
      local hostname=$(echo "$host_entry" | awk '{print $2}')
      ssh_key_copy "$hostaddr"
      ssh_key_copy "$hostname"
    else
      ssh_key_copy "$host"
    fi
  done
}

function openim::util::sourced_variable {
  # Call this function to tell shellcheck that a variable is supposed to
  # be used from other calling context. This helps quiet an "unused
  # variable" warning from shellcheck and also document your code.
  true
}

openim::util::sortable_date() {
  date "+%Y%m%d-%H%M%S"
}

# arguments: target, item1, item2, item3, ...
# returns 0 if target is in the given items, 1 otherwise.
openim::util::array_contains() {
  local search="$1"
  local element
  shift
  for element; do
    if [[ "${element}" == "${search}" ]]; then
      return 0
     fi
  done
  return 1
}

openim::util::wait_for_url() {
  local url=$1
  local prefix=${2:-}
  local wait=${3:-1}
  local times=${4:-30}
  local maxtime=${5:-1}

  command -v curl >/dev/null || {
    openim::log::usage "curl must be installed"
    exit 1
  }

  local i
  for i in $(seq 1 "${times}"); do
    local out
    if out=$(curl --max-time "${maxtime}" -gkfs "${url}" 2>/dev/null); then
      openim::log::status "On try ${i}, ${prefix}: ${out}"
      return 0
    fi
    sleep "${wait}"
  done
  openim::log::error "Timed out waiting for ${prefix} to answer at ${url}; tried ${times} waiting ${wait} between each"
  return 1
}

# Example:  openim::util::wait_for_success 120 5 "openimctl get nodes|grep localhost"
# arguments: wait time, sleep time, shell command
# returns 0 if the shell command get output, 1 otherwise.
openim::util::wait_for_success(){
  local wait_time="$1"
  local sleep_time="$2"
  local cmd="$3"
  while [ "$wait_time" -gt 0 ]; do
    if eval "$cmd"; then
      return 0
    else
      sleep "$sleep_time"
      wait_time=$((wait_time-sleep_time))
    fi
  done
  return 1
}

# Example:  openim::util::trap_add 'echo "in trap DEBUG"' DEBUG
# See: http://stackoverflow.com/questions/3338030/multiple-bash-traps-for-the-same-signal
openim::util::trap_add() {
  local trap_add_cmd
  trap_add_cmd=$1
  shift

  for trap_add_name in "$@"; do
    local existing_cmd
    local new_cmd

    # Grab the currently defined trap commands for this trap
    existing_cmd=$(trap -p "${trap_add_name}" |  awk -F"'" '{print $2}')

    if [[ -z "${existing_cmd}" ]]; then
      new_cmd="${trap_add_cmd}"
    else
      new_cmd="${trap_add_cmd};${existing_cmd}"
    fi

    # Assign the test. Disable the shellcheck warning telling that trap
    # commands should be single quoted to avoid evaluating them at this
    # point instead evaluating them at run time. The logic of adding new
    # commands to a single trap requires them to be evaluated right away.
    # shellcheck disable=SC2064
    trap "${new_cmd}" "${trap_add_name}"
  done
}

# Opposite of openim::util::ensure-temp-dir()
openim::util::cleanup-temp-dir() {
  rm -rf "${OPENIM_TEMP}"
}

# Create a temp dir that'll be deleted at the end of this bash session.
#
# Vars set:
#   OPENIM_TEMP
openim::util::ensure-temp-dir() {
  if [[ -z ${OPENIM_TEMP-} ]]; then
    OPENIM_TEMP=$(mktemp -d 2>/dev/null || mktemp -d -t openimrnetes.XXXXXX)
    openim::util::trap_add openim::util::cleanup-temp-dir EXIT
  fi
}

openim::util::host_os() {
  local host_os
  case "$(uname -s)" in
    Darwin)
      host_os=darwin
      ;;
    Linux)
      host_os=linux
      ;;
    *)
      openim::log::error "Unsupported host OS.  Must be Linux or Mac OS X."
      exit 1
      ;;
  esac
  echo "${host_os}"
}

openim::util::host_arch() {
  local host_arch
  case "$(uname -m)" in
    x86_64*)
      host_arch=amd64
      ;;
    i?86_64*)
      host_arch=amd64
      ;;
    amd64*)
      host_arch=amd64
      ;;
    aarch64*)
      host_arch=arm64
      ;;
    arm64*)
      host_arch=arm64
      ;;
    arm*)
      host_arch=arm
      ;;
    i?86*)
      host_arch=x86
      ;;
    s390x*)
      host_arch=s390x
      ;;
    ppc64le*)
      host_arch=ppc64le
      ;;
    *)
      openim::log::error "Unsupported host arch. Must be x86_64, 386, arm, arm64, s390x or ppc64le."
      exit 1
      ;;
  esac
  echo "${host_arch}"
}

# Define a bash function to check the versions of Docker and Docker Compose
openim::util::check_docker_and_compose_versions() {
    # Define the required versions of Docker and Docker Compose
    required_docker_version="20.10.0"
    required_compose_version="2.0"

    # Get the currently installed Docker version
    installed_docker_version=$(docker --version | awk '{print $3}' | sed 's/,//')

    # Check if the installed Docker version matches the required version
    if [[ "$installed_docker_version" < "$required_docker_version" ]]; then
        echo "Docker version mismatch. Installed: $installed_docker_version, Required: $required_docker_version"
        return 1
    fi

    # Check if the docker compose sub-command is available
    if ! docker compose version &> /dev/null; then
        echo "Docker does not support the docker compose sub-command"
        echo "You need to upgrade Docker to the right version"
        return 1
    fi

    # Get the currently installed Docker Compose version
    installed_compose_version=$(docker compose version --short)

    # Check if the installed Docker Compose version matches the required version
    if [[ "$installed_compose_version" < "$required_compose_version" ]]; then
        echo "Docker Compose version mismatch. Installed: $installed_compose_version, Required: $required_compose_version"
        return 1
    fi

}


# The `openim::util::check_ports` function analyzes the state of processes based on given ports.
# It accepts multiple ports as arguments and prints:
# 1. The state of the process (whether it's running or not).
# 2. The start time of the process if it's running.
# User:
# openim::util::check_ports 8080 8081 8082
# The function returns a status of 1 if any of the processes is not running.
openim::util::check_ports() {
    # An array to collect ports of processes that are not running.
    local not_started=()

    # An array to collect information about processes that are running.
    local started=()

    openim::log::info "Checking ports: $*"
    # Iterate over each given port.
    for port in "$@"; do
        # Use the `ss` command to find process information related to the given port.
        local info=$(ss -ltnp | grep ":$port" || true)
        
        # If there's no process information, it means the process associated with the port is not running.
        if [[ -z $info ]]; then
            not_started+=($port)
        else
            # Extract relevant details: Process Name, PID, and FD.
            local details=$(echo $info | sed -n 's/.*users:(("\([^"]*\)",pid=\([^,]*\),fd=\([^)]*\))).*/\1 \2 \3/p')
            local command=$(echo $details | awk '{print $1}')
            local pid=$(echo $details | awk '{print $2}')
            local fd=$(echo $details | awk '{print $3}')
            
            # Get the start time of the process using the PID
            if [[ -z $pid ]]; then
                local start_time="N/A"
            else
                # Get the start time of the process using the PID
                local start_time=$(ps -p $pid -o lstart=)
            fi
            
            started+=("Port $port - Command: $command, PID: $pid, FD: $fd, Started: $start_time")
        fi
    done

    # Print information about ports whose processes are not running.
    if [[ ${#not_started[@]} -ne 0 ]]; then
        openim::log::info "\n### Not started ports:"
        for port in "${not_started[@]}"; do
            openim::log::error "Port $port is not started."
        done
    fi

    # Print information about ports whose processes are running.
    if [[ ${#started[@]} -ne 0 ]]; then
        openim::log::info "\n### Started ports:"
        for info in "${started[@]}"; do
            openim::log::info "$info"
        done
    fi

    # If any of the processes is not running, return a status of 1.
    if [[ ${#not_started[@]} -ne 0 ]]; then
        echo "++++ OpenIM Log >> cat ${LOG_FILE}"
        return 1
    else
        openim::log::success "All specified processes are running."
        return 0
    fi
}
# set +o errexit
# Sample call for testing:
# openim::util::check_ports 10002 1004 12345 13306
# set -o errexit

# The `openim::util::check_process_names` function analyzes the state of processes based on given names.
# It accepts multiple process names as arguments and prints:
# 1. The state of the process (whether it's running or not).
# 2. The start time of the process if it's running.
# User:
# openim::util::check_process_names nginx mysql redis
# The function returns a status of 1 if any of the processes is not running.
openim::util::check_process_names() {
    # Arrays to collect details of processes
    local not_started=()
    local started=()

    openim::log::info "Checking processes: $*"
    # Iterate over each given process name
    for process_name in "$@"; do
        # Use `pgrep` to find process IDs related to the given process name
        local pids=($(pgrep -f $process_name))
        
        # Check if any process IDs were found
        if [[ ${#pids[@]} -eq 0 ]]; then
            not_started+=($process_name)
        else
            # If there are PIDs, loop through each one
            for pid in "${pids[@]}"; do
                local command=$(ps -p $pid -o cmd=)
                local start_time=$(ps -p $pid -o lstart=)
                local port=$(ss -ltnp 2>/dev/null | grep $pid | awk '{print $4}' | cut -d ':' -f2)

                # Check if port information was found for the PID
                if [[ -z $port ]]; then
                    port="N/A"
                fi

                started+=("Process $process_name - Command: $command, PID: $pid, Port: $port, Start time: $start_time")
            done
        fi
    done

    # Print information
    if [[ ${#not_started[@]} -ne 0 ]]; then
        openim::log::info "Not started processes:"
        for process_name in "${not_started[@]}"; do
            openim::log::error "Process $process_name is not started."
        done
    fi

    if [[ ${#started[@]} -ne 0 ]]; then
        echo
        openim::log::info "Started processes:"
        for info in "${started[@]}"; do
            openim::log::info "$info"
        done
    fi

    # Return status
    if [[ ${#not_started[@]} -ne 0 ]]; then
        echo "++++ OpenIM Log >> cat ${LOG_FILE}"
        return 1
    else
        openim::log::success "All processes are running."
        return 0
    fi
}
# openim::util::check_process_names docker-pr

# The `openim::util::stop_services_on_ports` function stops services running on specified ports.
# It accepts multiple ports as arguments and performs the following:
# 1. Attempts to stop any services running on the specified ports.
# 2. Prints details of services successfully stopped and those that failed to stop.
# Usage:
# openim::util::stop_services_on_ports 8080 8081 8082
# The function returns a status of 1 if any service couldn't be stopped.
openim::util::stop_services_on_ports() {
    # An array to collect ports of processes that couldn't be stopped.
    local not_stopped=()

    # An array to collect information about processes that were stopped.
    local stopped=()

    openim::log::info "Stopping services on ports: $*"
    # Iterate over each given port.
    for port in "$@"; do
        # Use the `lsof` command to find process information related to the given port.
        info=$(lsof -i :$port -n -P | grep LISTEN || true)
                
        # If there's process information, it means the process associated with the port is running.
        if [[ -n $info ]]; then
            # Extract the Process ID.
            while read -r line; do
                local pid=$(echo $line | awk '{print $2}')
                    
                # Try to stop the service by killing its process.
                if kill -TERM $pid; then
                    stopped+=($port)
                else
                    not_stopped+=($port)
                fi
            done <<< "$info"
        fi
    done

    # Print information about ports whose processes couldn't be stopped.
    if [[ ${#not_stopped[@]} -ne 0 ]]; then
        openim::log::info "Ports that couldn't be stopped:"
        for port in "${not_stopped[@]}"; do
            openim::log::status "Failed to stop service on port $port."
        done
    fi

    # Print information about ports whose processes were successfully stopped.
    if [[ ${#stopped[@]} -ne 0 ]]; then
        echo
        openim::log::info "Stopped services on ports:"
        for port in "${stopped[@]}"; do
            openim::log::info "Successfully stopped service on port $port."
        done
    fi

    # If any of the processes couldn't be stopped, return a status of 1.
    if [[ ${#not_stopped[@]} -ne 0 ]]; then
        return 1
    else
        openim::log::success "All specified services were stopped."
        return 0
    fi
}
# nc -l -p 12345
# nc -l -p 123456
# ps -ef | grep "nc -l"
# openim::util::stop_services_on_ports 1234 12345 


# The `openim::util::stop_services_with_name` function stops services with specified names.
# It accepts multiple service names as arguments and performs the following:
# 1. Attempts to stop any services with the specified names.
# 2. Prints details of services successfully stopped and those that failed to stop.
# Usage:
# openim::util::stop_services_with_name nginx apache
# The function returns a status of 1 if any service couldn't be stopped.
openim::util::stop_services_with_name() {
    # An array to collect names of processes that couldn't be stopped.
    local not_stopped=()

    # An array to collect information about processes that were stopped.
    local stopped=()

    openim::log::info "Stopping services with names: $*"
    # Iterate over each given service name.
    for server_name in "$@"; do
        # Use the `pgrep` command to find process IDs related to the given service name.
        local pids=$(pgrep -f "$server_name")

        # If no process was found with the name, add it to the not_stopped list
        if [[ -z $pids ]]; then
            not_stopped+=("$server_name")
            continue
        fi
        local stopped_this_time=false
        for pid in $pids; do

            # Exclude the PID of the current script
            if [[ "$pid" == "$$" ]]; then
                continue
            fi

            # If there's a Process ID, it means the service with the name is running.
            if [[ -n $pid ]]; then
                # Try to stop the service by killing its process.
                if kill -TERM $pid 2>/dev/null; then
                    stopped_this_time=true
                fi
            fi
        done

        if $stopped_this_time; then
            stopped+=("$server_name")
        else
            not_stopped+=("$server_name")
        fi
    done

    # Print information about services whose processes couldn't be stopped.
    if [[ ${#not_stopped[@]} -ne 0 ]]; then
        openim::log::info "Services that couldn't be stopped:"
        for name in "${not_stopped[@]}"; do
            openim::log::status "Failed to stop the $name service."
        done
    fi

    # Print information about services whose processes were successfully stopped.
    if [[ ${#stopped[@]} -ne 0 ]]; then
        echo
        openim::log::info "Stopped services:"
        for name in "${stopped[@]}"; do
            openim::log::info "Successfully stopped the $name service."
        done
    fi

    openim::log::success "All specified services were stopped."
}
# sleep 333333&
# sleep 444444&
# ps -ef | grep "sleep"
# openim::util::stop_services_with_name "sleep 333333" "sleep 444444"

# This figures out the host platform without relying on golang.  We need this as
# we don't want a golang install to be a prerequisite to building yet we need
# this info to figure out where the final binaries are placed.
openim::util::host_platform() {
  echo "$(openim::util::host_os)/$(openim::util::host_arch)"
}

# looks for $1 in well-known output locations for the platform ($2)
# $OPENIM_ROOT must be set
openim::util::find-binary-for-platform() {
  local -r lookfor="$1"
  local -r platform="$2"
  local locations=(
    ""${OPENIM_ROOT}"/_output/bin/${lookfor}"
    ""${OPENIM_ROOT}"/_output/${platform}/${lookfor}"
    ""${OPENIM_ROOT}"/_output/local/bin/${platform}/${lookfor}"
    ""${OPENIM_ROOT}"/_output/platforms/${platform}/${lookfor}"
    ""${OPENIM_ROOT}"/_output/platforms/bin/${platform}/${lookfor}"
  )

  # List most recently-updated location.
  local -r bin=$( (ls -t "${locations[@]}" 2>/dev/null || true) | head -1 )
  echo -n "${bin}"
}

# looks for $1 in well-known output locations for the host platform
# $OPENIM_ROOT must be set
openim::util::find-binary() {
  openim::util::find-binary-for-platform "$1" "$(openim::util::host_platform)"
}

# Run all known doc generators (today gendocs and genman for openimctl)
# $1 is the directory to put those generated documents
openim::util::gen-docs() {
  local dest="$1"

  # Find binary
  gendocs=$(openim::util::find-binary "gendocs")
  genopenimdocs=$(openim::util::find-binary "genopenimdocs")
  genman=$(openim::util::find-binary "genman")
  genyaml=$(openim::util::find-binary "genyaml")
  genfeddocs=$(openim::util::find-binary "genfeddocs")

  # TODO: If ${genfeddocs} is not used from anywhere (it isn't used at
  # least from k/k tree), remove it completely.
  openim::util::sourced_variable "${genfeddocs}"

  mkdir -p "${dest}/docs/guide/en-US/cmd/openimctl/"
  "${gendocs}" "${dest}/docs/guide/en-US/cmd/openimctl/"

  mkdir -p "${dest}/docs/guide/en-US/cmd/"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-api"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-cmdutils"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-crontask"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-msggateway"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-msgtransfer"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-push"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-rpc-auth"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-rpc-conversation"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-rpc-friend"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-rpc-group"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-rpc-msg"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-rpc-third"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/" "openim-rpc-user"
  "${genopenimdocs}" "${dest}/docs/guide/en-US/cmd/openimctl" "openimctl"

  mkdir -p "${dest}/docs/man/man1/"
"${genman}" "${dest}/docs/man/man1/" "openim-api"
"${genman}" "${dest}/docs/man/man1/" "openim-cmdutils"
"${genman}" "${dest}/docs/man/man1/" "openim-crontask"
"${genman}" "${dest}/docs/man/man1/" "openim-msggateway"
"${genman}" "${dest}/docs/man/man1/" "openim-msgtransfer"
"${genman}" "${dest}/docs/man/man1/" "openim-push"
"${genman}" "${dest}/docs/man/man1/" "openim-rpc-auth"
"${genman}" "${dest}/docs/man/man1/" "openim-rpc-conversation"
"${genman}" "${dest}/docs/man/man1/" "openim-rpc-friend"
"${genman}" "${dest}/docs/man/man1/" "openim-rpc-group"
"${genman}" "${dest}/docs/man/man1/" "openim-rpc-msg"
"${genman}" "${dest}/docs/man/man1/" "openim-rpc-third"
"${genman}" "${dest}/docs/man/man1/" "openim-rpc-user"

  mkdir -p "${dest}/docs/guide/en-US/yaml/openimctl/"
  "${genyaml}" "${dest}/docs/guide/en-US/yaml/openimctl/"

  # create the list of generated files
  pushd "${dest}" > /dev/null || return 1
  touch docs/.generated_docs
  find . -type f | cut -sd / -f 2- | LC_ALL=C sort > docs/.generated_docs
  popd > /dev/null || return 1
}

# Removes previously generated docs-- we don't want to check them in. $OPENIM_ROOT
# must be set.
openim::util::remove-gen-docs() {
  if [ -e ""${OPENIM_ROOT}"/docs/.generated_docs" ]; then
    # remove all of the old docs; we don't want to check them in.
    while read -r file; do
      rm ""${OPENIM_ROOT}"/${file}" 2>/dev/null || true
    done <""${OPENIM_ROOT}"/docs/.generated_docs"
    # The docs/.generated_docs file lists itself, so we don't need to explicitly
    # delete it.
  fi
}

# Returns the name of the upstream remote repository name for the local git
# repo, e.g. "upstream" or "origin".
openim::util::git_upstream_remote_name() {
  git remote -v | grep fetch |\
    grep -E 'github.com[/:]openimsdk/open-im-server|openim.cc/server' |\
    head -n 1 | awk '{print $1}'
}

# Exits script if working directory is dirty. If it's run interactively in the terminal
# the user can commit changes in a second terminal. This script will wait.
openim::util::ensure_clean_working_dir() {
  while ! git diff HEAD --exit-code &>/dev/null; do
    echo -e "\nUnexpected dirty working directory:\n"
    if tty -s; then
        git status -s
    else
        git diff -a # be more verbose in log files without tty
        exit 1
    fi | sed 's/^/  /'
    echo -e "\nCommit your changes in another terminal and then continue here by pressing enter."
    read -r
  done 1>&2
}

# Find the base commit using:
# $PULL_BASE_SHA if set (from Prow)
# current ref from the remote upstream branch
openim::util::base_ref() {
  local -r git_branch=$1

  if [[ -n ${PULL_BASE_SHA:-} ]]; then
    echo "${PULL_BASE_SHA}"
    return
  fi

  full_branch="$(openim::util::git_upstream_remote_name)/${git_branch}"

  # make sure the branch is valid, otherwise the check will pass erroneously.
  if ! git describe "${full_branch}" >/dev/null; then
    # abort!
    exit 1
  fi

  echo "${full_branch}"
}

# Checks whether there are any files matching pattern $2 changed between the
# current branch and upstream branch named by $1.
# Returns 1 (false) if there are no changes
#         0 (true) if there are changes detected.
openim::util::has_changes() {
  local -r git_branch=$1
  local -r pattern=$2
  local -r not_pattern=${3:-totallyimpossiblepattern}

  local base_ref
  base_ref=$(openim::util::base_ref "${git_branch}")
  echo "Checking for '${pattern}' changes against '${base_ref}'"

  # notice this uses ... to find the first shared ancestor
  if git diff --name-only "${base_ref}...HEAD" | grep -v -E "${not_pattern}" | grep "${pattern}" > /dev/null; then
    return 0
  fi
  # also check for pending changes
  if git status --porcelain | grep -v -E "${not_pattern}" | grep "${pattern}" > /dev/null; then
    echo "Detected '${pattern}' uncommitted changes."
    return 0
  fi
  echo "No '${pattern}' changes detected."
  return 1
}

openim::util::download_file() {
  local -r url=$1
  local -r destination_file=$2

  rm "${destination_file}" 2&> /dev/null || true

  for i in $(seq 5)
  do
    if ! curl -fsSL --retry 3 --keepalive-time 2 "${url}" -o "${destination_file}"; then
      echo "Downloading ${url} failed. $((5-i)) retries left."
      sleep 1
    else
      echo "Downloading ${url} succeed"
      return 0
    fi
  done
  return 1
}

# Test whether openssl is installed.
# Sets:
#  OPENSSL_BIN: The path to the openssl binary to use
function openim::util::test_openssl_installed {
    if ! openssl version >& /dev/null; then
      echo "Failed to run openssl. Please ensure openssl is installed"
      exit 1
    fi

    OPENSSL_BIN=$(command -v openssl)
}

# creates a client CA, args are sudo, dest-dir, ca-id, purpose
# purpose is dropped in after "key encipherment", you usually want
# '"client auth"'
# '"server auth"'
# '"client auth","server auth"'
function openim::util::create_signing_certkey {
    local sudo=$1
    local dest_dir=$2
    local id=$3
    local purpose=$4
    # Create client ca
    ${sudo} /usr/bin/env bash -e <<EOF
    rm -f "${dest_dir}/${id}-ca.crt" "${dest_dir}/${id}-ca.key"
    ${OPENSSL_BIN} req -x509 -sha256 -new -nodes -days 365 -newkey rsa:2048 -keyout "${dest_dir}/${id}-ca.key" -out "${dest_dir}/${id}-ca.crt" -subj "/C=xx/ST=x/L=x/O=x/OU=x/CN=ca/emailAddress=x/"
    echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment",${purpose}]}}}' > "${dest_dir}/${id}-ca-config.json"
EOF
}

# signs a client certificate: args are sudo, dest-dir, CA, filename (roughly), username, groups...
function openim::util::create_client_certkey {
    local sudo=$1
    local dest_dir=$2
    local ca=$3
    local id=$4
    local cn=${5:-$4}
    local groups=""
    local SEP=""
    shift 5
    while [ -n "${1:-}" ]; do
        groups+="${SEP}{\"O\":\"$1\"}"
        SEP=","
        shift 1
    done
    ${sudo} /usr/bin/env bash -e <<EOF
    cd ${dest_dir}
    echo '{"CN":"${cn}","names":[${groups}],"hosts":[""],"key":{"algo":"rsa","size":2048}}' | ${CFSSL_BIN} gencert -ca=${ca}.crt -ca-key=${ca}.key -config=${ca}-config.json - | ${CFSSLJSON_BIN} -bare client-${id}
    mv "client-${id}-key.pem" "client-${id}.key"
    mv "client-${id}.pem" "client-${id}.crt"
    rm -f "client-${id}.csr"
EOF
}

# signs a serving certificate: args are sudo, dest-dir, ca, filename (roughly), subject, hosts...
function openim::util::create_serving_certkey {
    local sudo=$1
    local dest_dir=$2
    local ca=$3
    local id=$4
    local cn=${5:-$4}
    local hosts=""
    local SEP=""
    shift 5
    while [ -n "${1:-}" ]; do
        hosts+="${SEP}\"$1\""
        SEP=","
        shift 1
    done
    ${sudo} /usr/bin/env bash -e <<EOF
    cd ${dest_dir}
    echo '{"CN":"${cn}","hosts":[${hosts}],"key":{"algo":"rsa","size":2048}}' | ${CFSSL_BIN} gencert -ca=${ca}.crt -ca-key=${ca}.key -config=${ca}-config.json - | ${CFSSLJSON_BIN} -bare serving-${id}
    mv "serving-${id}-key.pem" "serving-${id}.key"
    mv "serving-${id}.pem" "serving-${id}.crt"
    rm -f "serving-${id}.csr"
EOF
}

# creates a self-contained openimconfig: args are sudo, dest-dir, ca file, host, port, client id, token(optional)
function openim::util::write_client_openimconfig {
    local sudo=$1
    local dest_dir=$2
    local ca_file=$3
    local api_host=$4
    local api_port=$5
    local client_id=$6
    local token=${7:-}
    cat <<EOF | ${sudo} tee "${dest_dir}"/"${client_id}".openimconfig > /dev/null
apiVersion: v1
kind: Config
clusters:
  - cluster:
      certificate-authority: ${ca_file}
      server: https://${api_host}:${api_port}/
    name: local-up-cluster
users:
  - user:
      token: ${token}
      client-certificate: ${dest_dir}/client-${client_id}.crt
      client-key: ${dest_dir}/client-${client_id}.key
    name: local-up-cluster
contexts:
  - context:
      cluster: local-up-cluster
      user: local-up-cluster
    name: local-up-cluster
current-context: local-up-cluster
EOF

    # flatten the openimconfig files to make them self contained
    username=$(whoami)
    ${sudo} /usr/bin/env bash -e <<EOF
    $(openim::util::find-binary openimctl) --openimconfig="${dest_dir}/${client_id}.openimconfig" config view --minify --flatten > "/tmp/${client_id}.openimconfig"
    mv -f "/tmp/${client_id}.openimconfig" "${dest_dir}/${client_id}.openimconfig"
    chown ${username} "${dest_dir}/${client_id}.openimconfig"
EOF
}

# Determines if docker can be run, failures may simply require that the user be added to the docker group.
function openim::util::ensure_docker_daemon_connectivity {
  IFS=" " read -ra DOCKER <<< "${DOCKER_OPTS}"
  # Expand ${DOCKER[@]} only if it's not unset. This is to work around
  # Bash 3 issue with unbound variable.
  DOCKER=(docker ${DOCKER[@]:+"${DOCKER[@]}"})
  if ! "${DOCKER[@]}" info > /dev/null 2>&1 ; then
    cat <<'EOF' >&2
Can't connect to 'docker' daemon.  please fix and retry.

Possible causes:
  - Docker Daemon not started
    - Linux: confirm via your init system
    - macOS w/ docker-machine: run `docker-machine ls` and `docker-machine start <name>`
    - macOS w/ Docker for Mac: Check the menu bar and start the Docker application
  - DOCKER_HOST hasn't been set or is set incorrectly
    - Linux: domain socket is used, DOCKER_* should be unset. In Bash run `unset ${!DOCKER_*}`
    - macOS w/ docker-machine: run `eval "$(docker-machine env <name>)"`
    - macOS w/ Docker for Mac: domain socket is used, DOCKER_* should be unset. In Bash run `unset ${!DOCKER_*}`
  - Other things to check:
    - Linux: User isn't in 'docker' group.  Add and relogin.
      - Something like 'sudo usermod -a -G docker ${USER}'
      - RHEL7 bug and workaround: https://bugzilla.redhat.com/show_bug.cgi?id=1119282#c8
EOF
    return 1
  fi
}

# Wait for background jobs to finish. Return with
# an error status if any of the jobs failed.
openim::util::wait-for-jobs() {
  local fail=0
  local job
  for job in $(jobs -p); do
    wait "${job}" || fail=$((fail + 1))
  done
  return ${fail}
}

# openim::util::join <delim> <list...>
# Concatenates the list elements with the delimiter passed as first parameter
#
# Ex: openim::util::join , a b c
#  -> a,b,c
function openim::util::join {
  local IFS="$1"
  shift
  echo "$*"
}

# Function: openim::util::list-to-string <list...>
# Description: Converts a list to a string, removing spaces, brackets, and commas.
# Example input: [1002 3 ,  2 32 3 ,  3 434 ,]
# Example output: 10023 2323 3434
# Example usage:
# result=$(openim::util::list-to-string "[10023, 2323, 3434]")
# echo $result
function openim::util::list-to-string() {
    # Capture all arguments into a single string
    ports_list="$*"

    # Use sed for transformations:
    # 1. Remove spaces
    # 2. Replace commas with spaces
    # 3. Remove opening and closing brackets
    ports_array=$(echo "$ports_list" | sed 's/ //g; s/,/ /g; s/^\[\(.*\)\]$/\1/')
    # For external use, we might want to echo the result so that it can be captured by callers
    echo "$ports_array"
}
# MSG_GATEWAY_PROM_PORTS=$(openim::util::list-to-string "10023, 2323, 34 34")
# read -a MSG_GATEWAY_PROM_PORTS <<< $(openim::util::list-to-string "10023, 2323, 34 34")
# echo ${MSG_GATEWAY_PROM_PORTS}
# echo "${#MSG_GATEWAY_PROM_PORTS[@]}"
# Downloads cfssl/cfssljson/cfssl-certinfo into $1 directory if they do not already exist in PATH
#
# Assumed vars:
#   $1 (cfssl directory) (optional)
#
# Sets:
#  CFSSL_BIN: The path of the installed cfssl binary
#  CFSSLJSON_BIN: The path of the installed cfssljson binary
#  CFSSLCERTINFO_BIN: The path of the installed cfssl-certinfo binary
#
function openim::util::ensure-cfssl {
  if command -v cfssl &>/dev/null && command -v cfssljson &>/dev/null && command -v cfssl-certinfo &>/dev/null; then
    CFSSL_BIN=$(command -v cfssl)
    CFSSLJSON_BIN=$(command -v cfssljson)
    CFSSLCERTINFO_BIN=$(command -v cfssl-certinfo)
    return 0
  fi

  host_arch=$(openim::util::host_arch)

  if [[ "${host_arch}" != "amd64" ]]; then
    echo "Cannot download cfssl on non-amd64 hosts and cfssl does not appear to be installed."
    echo "Please install cfssl, cfssljson and cfssl-certinfo and verify they are in \$PATH."
    echo "Hint: export PATH=\$PATH:\$GOPATH/bin; go get -u github.com/cloudflare/cfssl/cmd/..."
    exit 1
  fi

  # Create a temp dir for cfssl if no directory was given
  local cfssldir=${1:-}
  if [[ -z "${cfssldir}" ]]; then
    cfssldir="$HOME/bin"
  fi

  mkdir -p "${cfssldir}"
  pushd "${cfssldir}" > /dev/null || return 1

  echo "Unable to successfully run 'cfssl' from ${PATH}; downloading instead..."
  kernel=$(uname -s)
  case "${kernel}" in
    Linux)
      curl --retry 10 -L -o cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
      curl --retry 10 -L -o cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
      curl --retry 10 -L -o cfssl-certinfo https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64
      ;;
    Darwin)
      curl --retry 10 -L -o cfssl https://pkg.cfssl.org/R1.2/cfssl_darwin-amd64
      curl --retry 10 -L -o cfssljson https://pkg.cfssl.org/R1.2/cfssljson_darwin-amd64
      curl --retry 10 -L -o cfssl-certinfo https://pkg.cfssl.org/R1.2/cfssl-certinfo_darwin-amd64
      ;;
    *)
      echo "Unknown, unsupported platform: ${kernel}." >&2
      echo "Supported platforms: Linux, Darwin." >&2
      exit 2
  esac

  chmod +x cfssl || true
  chmod +x cfssljson || true
  chmod +x cfssl-certinfo || true

  CFSSL_BIN="${cfssldir}/cfssl"
  CFSSLJSON_BIN="${cfssldir}/cfssljson"
  CFSSLCERTINFO_BIN="${cfssldir}/cfssl-certinfo"
  if [[ ! -x ${CFSSL_BIN} || ! -x ${CFSSLJSON_BIN} || ! -x ${CFSSLCERTINFO_BIN} ]]; then
    echo "Failed to download 'cfssl'."
    echo "Please install cfssl, cfssljson and cfssl-certinfo and verify they are in \$PATH."
    echo "Hint: export PATH=\$PATH:\$GOPATH/bin; go get -u github.com/cloudflare/cfssl/cmd/..."
    exit 1
  fi
  popd > /dev/null || return 1
}

function openim::util::ensure-docker-buildx {
  # podman returns 0 on `docker buildx version`, docker on `docker buildx`. One of them must succeed.
  if docker buildx version >/dev/null 2>&1 || docker buildx >/dev/null 2>&1; then
    return 0
  else
    echo "ERROR: docker buildx not available. Docker 19.03 or higher is required with experimental features enabled"
    exit 1
  fi
}

# openim::util::ensure-bash-version
# Check if we are using a supported bash version
#
function openim::util::ensure-bash-version {
  # shellcheck disable=SC2004
  if ((${BASH_VERSINFO[0]}<4)) || ( ((${BASH_VERSINFO[0]}==4)) && ((${BASH_VERSINFO[1]}<2)) ); then
    echo "ERROR: This script requires a minimum bash version of 4.2, but got version of ${BASH_VERSINFO[0]}.${BASH_VERSINFO[1]}"
    if [ "$(uname)" = 'Darwin' ]; then
      echo "On macOS with homebrew 'brew install bash' is sufficient."
    fi
    exit 1
  fi
}

# openim::util::ensure-install-nginx
# Check if nginx is installed
#
function openim::util::ensure-install-nginx {
  if ! command -v nginx &>/dev/null; then
    echo "ERROR: nginx not found. Please install nginx."
    exit 1
  fi

  for port in 80
  do
    if echo |telnet 127.0.0.1 $port 2>&1|grep refused &>/dev/null;then
      exit 1
    fi
  done
}

# openim::util::ensure-gnu-sed
# Determines which sed binary is gnu-sed on linux/darwin
#
# Sets:
#  SED: The name of the gnu-sed binary
#
function openim::util::ensure-gnu-sed {
  # NOTE: the echo below is a workaround to ensure sed is executed before the grep.
  # see: https://github.com/openimrnetes/openimrnetes/issues/87251
  sed_help="$(LANG=C sed --help 2>&1 || true)"
  if echo "${sed_help}" | grep -q "GNU\|BusyBox"; then
    SED="sed"
  elif command -v gsed &>/dev/null; then
    SED="gsed"
  else
    openim::log::error "Failed to find GNU sed as sed or gsed. If you are on Mac: brew install gnu-sed." >&2
    return 1
  fi
  openim::util::sourced_variable "${SED}"
}

# openim::util::ensure-gnu-date
# Determines which date binary is gnu-date on linux/darwin
#
# Sets:
#  DATE: The name of the gnu-date binary
#
function openim::util::ensure-gnu-date {
  # NOTE: the echo below is a workaround to ensure date is executed before the grep.
  date_help="$(LANG=C date --help 2>&1 || true)"
  if echo "${date_help}" | grep -q "GNU\|BusyBox"; then
    DATE="date"
  elif command -v gdate &>/dev/null; then
    DATE="gdate"
  else
    openim::log::error "Failed to find GNU date as date or gdate. If you are on Mac: brew install coreutils." >&2
    return 1
  fi
  openim::util::sourced_variable "${DATE}"
}

# openim::util::check-file-in-alphabetical-order <file>
# Check that the file is in alphabetical order
#
function openim::util::check-file-in-alphabetical-order {
  local failure_file="$1"
  if ! diff -u "${failure_file}" <(LC_ALL=C sort "${failure_file}"); then
    {
      echo
      echo "${failure_file} is not in alphabetical order. Please sort it:"
      echo
      echo "  LC_ALL=C sort -o ${failure_file} ${failure_file}"
      echo
    } >&2
  false
  fi
}

# openim::util::require-jq
# Checks whether jq is installed.
function openim::util::require-jq {
  if ! command -v jq &>/dev/null; then
    echo "jq not found. Please install." 1>&2
    return 1
  fi
}

# openim::util::require-dig
# Checks whether dig is installed and provides installation instructions if it is not.
function openim::util::require-dig {
  if ! command -v dig &>/dev/null; then
    echo "dig command not found."
    echo "Please install 'dig' to use this feature."
    echo "Installation instructions:"
    echo "  For Ubuntu/Debian: sudo apt-get install dnsutils"
    echo "  For CentOS/RedHat: sudo yum install bind-utils"
    echo "  For macOS: 'dig' should be preinstalled. If missing, try: brew install bind"
    echo "  For Windows: Install BIND9 tools from https://www.isc.org/download/"
    return 1
  fi
}

# outputs md5 hash of $1, works on macOS and Linux
function openim::util::md5() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

# openim::util::read-array
# Reads in stdin and adds it line by line to the array provided. This can be
# used instead of "mapfile -t", and is bash 3 compatible.
#
# Assumed vars:
#   $1 (name of array to create/modify)
#
# Example usage:
# openim::util::read-array files < <(ls -1)
#
function openim::util::read-array {
  local i=0
  unset -v "$1"
  while IFS= read -r "$1[i++]"; do :; done
  eval "[[ \${$1[--i]} ]]" || unset "$1[i]" # ensures last element isn't empty
}

# Some useful colors.
if [[ -z "${color_start-}" ]]; then
  declare -r color_start="\033["
  declare -r color_red="${color_start}0;31m"
  declare -r color_yellow="${color_start}0;33m"
  declare -r color_green="${color_start}0;32m"
  declare -r color_blue="${color_start}1;34m"
  declare -r color_cyan="${color_start}1;36m"
  declare -r color_norm="${color_start}0m"

  openim::util::sourced_variable "${color_start}"
  openim::util::sourced_variable "${color_red}"
  openim::util::sourced_variable "${color_yellow}"
  openim::util::sourced_variable "${color_green}"
  openim::util::sourced_variable "${color_blue}"
  openim::util::sourced_variable "${color_cyan}"
  openim::util::sourced_variable "${color_norm}"
fi

# ex: ts=2 sw=2 et filetype=sh

function openim::util::desc() {
    openim::util:run::maybe_first_prompt
    rate=25
    if [ -n "$DEMO_RUN_FAST" ]; then
      rate=1000
    fi
    echo "$blue# $@$reset" | pv -qL $rate
    openim::util:run::prompt
}

function openim::util:run::prompt() {
    echo -n "$yellow\$ $reset"
}

started=""
function openim::util:run::maybe_first_prompt() {
    if [ -z "$started" ]; then
        openim::util:run::prompt
        started=true
    fi
}

# After a `run` this variable will hold the stdout of the command that was run.
# If the command was interactive, this will likely be garbage.
DEMO_RUN_STDOUT=""

function openim::util::run() {
    openim::util:run::maybe_first_prompt
    rate=25
    if [ -n "$DEMO_RUN_FAST" ]; then
      rate=1000
    fi
    echo "$green$1$reset" | pv -qL $rate
    if [ -n "$DEMO_RUN_FAST" ]; then
      sleep 0.5
    fi
    OFILE="$(mktemp -t $(basename $0).XXXXXX)"
    if [ "$(uname)" == "Darwin" ]; then
       script -q "$OFILE" $1
    else
       script -eq -c "$1" -f "$OFILE"
    fi
    r=$?
    read -d '' -t "${timeout}" -n 10000 # clear stdin
    openim::util:run::prompt
    if [ -z "$DEMO_AUTO_RUN" ]; then
      read -s
    fi
    DEMO_RUN_STDOUT="$(tail -n +2 $OFILE | sed 's/\r//g')"
    return $r
}

function openim::util::run::relative() {
    for arg; do
        echo "$(realpath $(dirname $(which $0)))/$arg" | sed "s|$(realpath $(pwd))|.|"
    done
}

# This function retrieves the IP address of the current server.
# It primarily uses the `curl` command to fetch the public IP address from ifconfig.me.
# If curl or the service is not available, it falls back 
# to the internal IP address provided by the hostname command.
# TODO: If a delay is found, the delay needs to be addressed
function openim::util::get_server_ip() {
    # Check if the 'curl' command is available
    if command -v curl &> /dev/null; then
        # Try to retrieve the public IP address using curl and ifconfig.me
        IP=$(dig TXT +short o-o.myaddr.l.google.com @ns1.google.com | sed 's/"//g' | tr -d '\n')
        
        # Check if IP retrieval was successful
        if [[ -z "$IP" ]]; then
            # If not, get the internal IP address
            IP=$(ip addr show | grep 'inet ' | grep -v 127.0.0.1 | awk '{print $2}' | cut -d'/' -f1 | head -n 1)
        fi
    else
        # If curl is not available, get the internal IP address
        IP=$(ip addr show | grep 'inet ' | grep -v 127.0.0.1 | awk '{print $2}' | cut -d'/' -f1 | head -n 1)
    fi
    
    # Return the fetched IP address
    echo "$IP"
}

function openim::util::onCtrlC() {
    echo -e "\n${t_reset}Ctrl+C Press it. It's exiting openim make init..."
    exit 1
}

# Function Function: Remove Spaces in the string
function openim::util::remove_space() {
    value=$*  # 获取传入的参数
    result=$(echo $value | sed 's/ //g')  # 去除空格
}

function openim::util::gencpu() {
    # Check the system type
    system_type=$(uname)

    if [[ "$system_type" == "Darwin" ]]; then
        # macOS (using sysctl)
        cpu_count=$(sysctl -n hw.ncpu)
    elif [[ "$system_type" == "Linux" ]]; then
        # Linux (using lscpu)
        cpu_count=$(lscpu --parse | grep -E '^([^#].*,){3}[^#]' | sort -u | wc -l)
    else
        echo "Unsupported operating system: $system_type"
        exit 1
    fi
    echo $cpu_count
}

function openim::util::set_max_fd() {
    local desired_fd=$1
    local max_fd_limit

    # Check if we're not on cygwin or darwin.
    if [ "$(uname -s | tr '[:upper:]' '[:lower:]')" != "cygwin" ] && [ "$(uname -s | tr '[:upper:]' '[:lower:]')" != "darwin" ]; then
        # Try to get the hard limit.
        max_fd_limit=$(ulimit -H -n)
        if [ $? -eq 0 ]; then
            # If desired_fd is 'maximum' or 'max', set it to the hard limit.
            if [ "$desired_fd" = "maximum" ] || [ "$desired_fd" = "max" ]; then
                desired_fd="$max_fd_limit"
            fi
            
            # Check if desired_fd is less than or equal to max_fd_limit.
            if [ "$desired_fd" -le "$max_fd_limit" ]; then
                ulimit -n "$desired_fd"
                if [ $? -ne 0 ]; then
                    echo "Warning: Could not set maximum file descriptor limit to $desired_fd"
                fi
            else
                echo "Warning: Desired file descriptor limit ($desired_fd) is greater than the hard limit ($max_fd_limit)"
            fi
        else
            echo "Warning: Could not query the maximum file descriptor hard limit."
        fi
    else
        echo "Warning: Not attempting to modify file descriptor limit on Cygwin or Darwin."
    fi
}


function openim::util::gen_os_arch() {
    # Get the current operating system and architecture
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    # Select the repository home directory based on the operating system and architecture
    if [[ "$OS" == "darwin" ]]; then
        if [[ "$ARCH" == "x86_64" ]]; then
            REPO_DIR="darwin/amd64"
        else
            REPO_DIR="darwin/386"
        fi
    elif [[ "$OS" == "linux" ]]; then
        if [[ "$ARCH" == "x86_64" ]]; then
            REPO_DIR="linux/amd64"
        elif [[ "$ARCH" == "arm64" ]]; then
            REPO_DIR="linux/arm64"
        elif [[ "$ARCH" == "mips64" ]]; then
            REPO_DIR="linux/mips64"
        elif [[ "$ARCH" == "mips64le" ]]; then
            REPO_DIR="linux/mips64le"
        elif [[ "$ARCH" == "ppc64le" ]]; then
            REPO_DIR="linux/ppc64le"
        elif [[ "$ARCH" == "s390x" ]]; then
            REPO_DIR="linux/s390x"
        else
            REPO_DIR="linux/386"
        fi
    elif [[ "$OS" == "windows" ]]; then
        if [[ "$ARCH" == "x86_64" ]]; then
            REPO_DIR="windows/amd64"
        else
            REPO_DIR="windows/386"
        fi
    else
        echo -e "${RED_PREFIX}Unsupported OS: $OS${COLOR_SUFFIX}"
        exit 1
    fi
}

if [[ "$*" =~ openim::util:: ]];then
  eval $*
fi