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

# This script automatically initializes the various configuration files
# Read: https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/init-config.md

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${OPENIM_ROOT}/scripts/lib/init.sh"

readonly ENV_FILE=${ENV_FILE:-"${OPENIM_ROOT}/scripts/install/environment.sh"}

declare -A TEMPLATES=(
  ["${OPENIM_ROOT}/deployments/templates/env-template.yaml"]="${OPENIM_ROOT}/.env"
  ["${OPENIM_ROOT}/deployments/templates/openim.yaml"]="${OPENIM_ROOT}/config/config.yaml"
  ["${OPENIM_ROOT}/deployments/templates/prometheus.yml"]="${OPENIM_ROOT}/config/prometheus.yml"
  ["${OPENIM_ROOT}/deployments/templates/alertmanager.yml"]="${OPENIM_ROOT}/config/alertmanager.yml"
)

openim::log::info "Read more configuration information: https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/environment.md"

# New variables for argument handling
FORCE_OVERWRITE=false
SKIP_EXISTING=false

# Function to display help
show_help() {
  echo "Usage: $(basename "$0") [options]"
  echo "Options:"
  echo "  -h, --help        Show this help message"
  echo "  --force           Overwrite existing files without prompt"
  echo "  --skip            Skip generation if file exists"
}

# Parse command-line options
while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      show_help
      exit 0
      ;;
    --force)
      FORCE_OVERWRITE=true
      shift
      ;;
    --skip)
      SKIP_EXISTING=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      show_help
      exit 1
      ;;
  esac
done

for template in "${!TEMPLATES[@]}"; do
  if [[ ! -f "${template}" ]]; then
    openim::log::error_exit "Template file ${template} does not exist..."
    exit 1
  fi
done

for template in "${!TEMPLATES[@]}"; do
  IFS=';' read -ra OUTPUT_FILES <<< "${TEMPLATES[$template]}"
  for output_file in "${OUTPUT_FILES[@]}"; do
    if [[ -f "${output_file}" ]]; then
      if [[ "${FORCE_OVERWRITE}" == true ]]; then
        openim::log::info "Force overwriting ${output_file}."
      elif [[ "${SKIP_EXISTING}" == true ]]; then
        openim::log::info "Skipping generation of ${output_file} as it already exists."
        continue
      else
        echo -n "File ${output_file} already exists. Overwrite? (Y/N): "
        read -r -n 1 REPLY
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
          openim::log::info "Skipping generation of ${output_file}."
          continue
        fi
      fi
    fi

    openim::log::info "⌚  Working with template file: ${template} to ${output_file}..."
    if [[ ! -f "${OPENIM_ROOT}/scripts/genconfig.sh" ]]; then
      openim::log::error "genconfig.sh script not found"
      exit 1
    fi
    "${OPENIM_ROOT}/scripts/genconfig.sh" "${ENV_FILE}" "${template}" > "${output_file}" || {
      openim::log::error "Error processing template file ${template}"
      exit 1
    }
    sleep 0.5
  done
done

openim::log::success "✨  All configuration files have been successfully generated!"
