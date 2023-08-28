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

# Controls verbosity of the script output and logging.
OPENIM_VERBOSE="${OPENIM_VERBOSE:-5}"

# Enable logging by default. Set to false to disable.
ENABLE_LOGGING=true

# If OPENIM_OUTPUT is not set, set it to the default value
if [[ ! -v OPENIM_OUTPUT ]]; then
    OPENIM_OUTPUT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../_output" && pwd -P)"
fi

# Set the log file path
LOG_FILE="${OPENIM_OUTPUT}/logs/openim_$(date '+%Y%m%d').log"

if [[ ! -d "${OPENIM_OUTPUT}/logs" ]]; then
    mkdir -p "${OPENIM_OUTPUT}/logs"
    touch "$LOG_FILE"
fi

# Define the logging function
function echo_log() {
    if $ENABLE_LOGGING; then
        echo -e "$@" | tee -a "${LOG_FILE}"
    else
        echo -e "$@"
    fi
}

# MAX_LOG_SIZE=10485760 # 10MB

# Clear logs from 5 days ago
# find $OPENIM_OUTPUT_LOGS -type f -name "*.log" -mtime +5 -exec rm -f {} \;

# Handler for when we exit automatically on an error.
# Borrowed from https://gist.github.com/ahendrix/7030300
openim::log::errexit() {
  local err="${PIPESTATUS[*]}"

  # If the shell we are in doesn't have errexit set (common in subshells) then
  # don't dump stacks.
  set +o | grep -qe "-o errexit" || return

  set +o xtrace
  local code="${1:-1}"
  # Print out the stack trace described by $function_stack
  if [ ${#FUNCNAME[@]} -gt 2 ]
  then
    openim::log::error "Call tree:"
    for ((i=1;i<${#FUNCNAME[@]}-1;i++))
    do
      openim::log::error " ${i}: ${BASH_SOURCE[${i}+1]}:${BASH_LINENO[${i}]} ${FUNCNAME[${i}]}(...)"
    done
  fi
  openim::log::error_exit "Error in ${BASH_SOURCE[1]}:${BASH_LINENO[0]}. '${BASH_COMMAND}' exited with status ${err}" "${1:-1}" 1
}

openim::log::install_errexit() {
  # trap ERR to provide an error handler whenever a command exits nonzero  this
  # is a more verbose version of set -o errexit
  trap 'openim::log::errexit' ERR

  # setting errtrace allows our ERR trap handler to be propagated to functions,
  # expansions and subshells
  set -o errtrace
}

# Print out the stack trace
#
# Args:
#   $1 The number of stack frames to skip when printing.
openim::log::stack() {
  local stack_skip=${1:-0}
  stack_skip=$((stack_skip + 1))
  if [[ ${#FUNCNAME[@]} -gt ${stack_skip} ]]; then
    echo_log "Call stack:" >&2
    local i
    for ((i=1 ; i <= ${#FUNCNAME[@]} - stack_skip ; i++))
    do
      local frame_no=$((i - 1 + stack_skip))
      local source_file=${BASH_SOURCE[${frame_no}]}
      local source_lineno=${BASH_LINENO[$((frame_no - 1))]}
      local funcname=${FUNCNAME[${frame_no}]}
      echo_log "  ${i}: ${source_file}:${source_lineno} ${funcname}(...)" >&2
    done
  fi
}

# Log an error and exit.
# Args:
#   $1 Message to log with the error
#   $2 The error code to return
#   $3 The number of stack frames to skip when printing.
openim::log::error_exit() {
  local message="${1:-}"
  local code="${2:-1}"
  local stack_skip="${3:-0}"
  stack_skip=$((stack_skip + 1))

  if [[ ${OPENIM_VERBOSE} -ge 4 ]]; then
    local source_file=${BASH_SOURCE[${stack_skip}]}
    local source_line=${BASH_LINENO[$((stack_skip - 1))]}
    echo_log -e "${COLOR_RED}!!! Error in ${source_file}:${source_line} ${COLOR_SUFFIX}" >&2
    [[ -z ${1-} ]] || {
      echo_log "  ${1}" >&2
    }

    openim::log::stack ${stack_skip}

    echo_log "Exiting with status ${code}" >&2
  fi

  exit "${code}"
}

# Log an error but keep going.  Don't dump the stack or exit.
openim::log::error() {
  timestamp=$(date +"[%m%d %H:%M:%S]")
  echo_log "!!! ${timestamp} ${1-}" >&2
  shift
  for message; do
    echo_log "    ${message}" >&2
  done
}

# Print an usage message to stderr.  The arguments are printed directly.
openim::log::usage() {
  echo_log >&2
  local message
  for message; do
    echo_log "${message}" >&2
  done
  echo_log >&2
}

openim::log::usage_from_stdin() {
  local messages=()
  while read -r line; do
    messages+=("${line}")
  done

  openim::log::usage "${messages[@]}"
}

# Print out some info that isn't a top level status line
openim::log::info() {
  local V="${V:-0}"
  if [[ ${OPENIM_VERBOSE} < ${V} ]]; then
    return
  fi

  for message; do
    echo_log "${message}"
  done
}

# Just like openim::log::info, but no \n, so you can make a progress bar
openim::log::progress() {
  for message; do
    echo_log -e -n "${message}"
  done
}

# Print out some info that isn't a top level status line
openim::log::info_from_stdin() {
  local messages=()
  while read -r line; do
    messages+=("${line}")
  done

  openim::log::info "${messages[@]}"
}

# Print a status line.  Formatted to show up in a stream of output.
openim::log::status() {
  local V="${V:-0}"
  if [[ ${OPENIM_VERBOSE} < ${V} ]]; then
    return
  fi

  timestamp=$(date +"[%m%d %H:%M:%S]")
  echo_log "+++ ${timestamp} ${1}"
  shift
  for message; do
    echo_log "    ${message}"
  done
}

openim::log::success()
{
  local V="${V:-0}"
  if [[ ${OPENIM_VERBOSE} < ${V} ]]; then
      return
  fi
  timestamp=$(date +"%m%d %H:%M:%S")
  echo_log -e "${COLOR_GREEN}[success ${timestamp}] ${COLOR_SUFFIX}==> " "$@"
}

function openim::log::test_log() {
    echo_log "test log"
    openim::log::info "openim::log::info"
    openim::log::progress "openim::log::progress"
    openim::log::status "openim::log::status"
    openim::log::success "openim::log::success"
    openim::log::error "openim::log::error"
    openim::log::error_exit "openim::log::error_exit"
}

# openim::log::test_log