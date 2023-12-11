#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script automatically initializes various configuration files and can generate example files.

set -o errexit
set -o nounset
set -o pipefail

# Root directory of the OpenIM project
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# Source initialization script
source "${OPENIM_ROOT}/scripts/lib/init.sh"

# Default environment file
readonly ENV_FILE=${ENV_FILE:-"${OPENIM_ROOT}/scripts/install/environment.sh"}

# Templates for configuration files
declare -A TEMPLATES=(
  ["${OPENIM_ROOT}/deployments/templates/env-template.yaml"]="${OPENIM_ROOT}/.env"
  ["${OPENIM_ROOT}/deployments/templates/openim.yaml"]="${OPENIM_ROOT}/config/config.yaml"
  ["${OPENIM_ROOT}/deployments/templates/prometheus.yml"]="${OPENIM_ROOT}/config/prometheus.yml"
  ["${OPENIM_ROOT}/deployments/templates/alertmanager.yml"]="${OPENIM_ROOT}/config/alertmanager.yml"
)

# Templates for example files
declare -A EXAMPLES=(
  ["${OPENIM_ROOT}/deployments/templates/env-template.yaml"]="${OPENIM_ROOT}/config/templates/env.template"
  ["${OPENIM_ROOT}/deployments/templates/openim.yaml"]="${OPENIM_ROOT}/config/templates/config.yaml.template"
  ["${OPENIM_ROOT}/deployments/templates/prometheus.yml"]="${OPENIM_ROOT}/config/templates/prometheus.yml.template"
  ["${OPENIM_ROOT}/deployments/templates/alertmanager.yml"]="${OPENIM_ROOT}/config/templates/alertmanager.yml.template"
)

# Command-line options
FORCE_OVERWRITE=false
SKIP_EXISTING=false
GENERATE_EXAMPLES=false
CLEAN_ENV_EXAMPLES=false

# Function to display help information
show_help() {
  echo "Usage: $(basename "$0") [options]"
  echo "Options:"
  echo "  -h, --help             Show this help message"
  echo "  --force                Overwrite existing files without prompt"
  echo "  --skip                 Skip generation if file exists"
  echo "  --examples             Generate example files"
  echo "  --clean-env-examples   Generate example files in a clean environment"
}

# Function to generate configuration files
generate_config_files() {
  # Loop through each template in TEMPLATES
  for template in "${!TEMPLATES[@]}"; do
    # Read the corresponding output files for the template
    IFS=';' read -ra OUTPUT_FILES <<< "${TEMPLATES[$template]}"
    for output_file in "${OUTPUT_FILES[@]}"; do
      # Check if the output file already exists
      if [[ -f "${output_file}" ]]; then
        # Handle existing file based on command-line options
        if [[ "${FORCE_OVERWRITE}" == true ]]; then
          openim::log::info "Force overwriting ${output_file}."
        elif [[ "${SKIP_EXISTING}" == true ]]; then
          openim::log::info "Skipping generation of ${output_file} as it already exists."
          continue
        else
          # Ask user for confirmation to overwrite
          echo -n "File ${output_file} already exists. Overwrite? (Y/N): "
          read -r -n 1 REPLY
          echo
          if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            openim::log::info "Skipping generation of ${output_file}."
            continue
          fi
        fi
      fi

      # Process the template file to generate the output file
      openim::log::info "⌚  Working with template file: ${template} to generate ${output_file}..."
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
}

# Function to generate example files
generate_example_files() {
  for template in "${!EXAMPLES[@]}"; do
    local example_file="${EXAMPLES[$template]}"
    if [[ ! -f "${example_file}" ]]; then
      openim::log::info "Generating example file: ${example_file} from ${template}..."
      cp "${template}" "${example_file}"
    fi
  done
}

declare -A env_vars=(
    ["OPENIM_IP"]="172.28.0.1"
    ["DATA_DIR"]="./"
    ["LOG_STORAGE_LOCATION"]="../logs/"
)

generate_clean_environment_examples() {
  env_cmd="env -i"
  for var in "${!env_vars[@]}"; do
      env_cmd+=" $var='${env_vars[$var]}'"
  done

  for template in "${!EXAMPLES[@]}"; do
    local example_file="${EXAMPLES[$template]}"
    openim::log::info "Generating example file: ${example_file} from ${template}..."

    eval "$env_cmd ${OPENIM_ROOT}/scripts/genconfig.sh '${ENV_FILE}' '${template}' > '${example_file}'" || {
      openim::log::error "Error processing template file ${template}"
      exit 1
    }
  done
}

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
    --examples)
      GENERATE_EXAMPLES=true
      shift
      ;;
    --clean-env-examples)
      CLEAN_ENV_EXAMPLES=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      show_help
      exit 1
      ;;
  esac
done

# Generate configuration files if requested
if [[ "${FORCE_OVERWRITE}" == true || "${SKIP_EXISTING}" == false ]]; then
  generate_config_files
fi

# Generate example files if --examples option is provided
if [[ "${GENERATE_EXAMPLES}" == true ]]; then
  generate_example_files
fi

# Generate example files in a clean environment if --clean-env-examples option is provided
if [[ "${CLEAN_ENV_EXAMPLES}" == true ]]; then
  generate_clean_environment_examples
fi

openim::log::success "Configuration and example files generation complete!"