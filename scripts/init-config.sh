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

# This script automatically initializes various configuration files and can generate example files.





# Root directory of the OpenIM project
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# Source initialization script
source "${OPENIM_ROOT}/scripts/lib/init.sh"

# Default environment file
readonly ENV_FILE=${ENV_FILE:-"${OPENIM_ROOT}/scripts/install/environment.sh"}

# Templates for configuration files
declare -A TEMPLATES=(
  ["${OPENIM_ROOT}/deployments/templates/env-template.yaml"]="${OPENIM_ROOT}/.env"
  ["${OPENIM_ROOT}/deployments/templates/config.yaml"]="${OPENIM_ROOT}/config/config.yaml"
  ["${OPENIM_ROOT}/deployments/templates/prometheus.yml"]="${OPENIM_ROOT}/config/prometheus.yml"
  ["${OPENIM_ROOT}/deployments/templates/alertmanager.yml"]="${OPENIM_ROOT}/config/alertmanager.yml"
)

# Templates for example files
declare -A EXAMPLES=(
  ["${OPENIM_ROOT}/deployments/templates/env-template.yaml"]="${OPENIM_ROOT}/config/templates/env.template"
  ["${OPENIM_ROOT}/deployments/templates/config.yaml"]="${OPENIM_ROOT}/config/templates/config.yaml.template"
  ["${OPENIM_ROOT}/deployments/templates/prometheus.yml"]="${OPENIM_ROOT}/config/templates/prometheus.yml.template"
  ["${OPENIM_ROOT}/deployments/templates/alertmanager.yml"]="${OPENIM_ROOT}/config/templates/alertmanager.yml.template"
)

# Templates for config Copy file
declare -A COPY_TEMPLATES=(
  ["${OPENIM_ROOT}/deployments/templates/email.tmpl"]="${OPENIM_ROOT}/config/email.tmpl"
  ["${OPENIM_ROOT}/deployments/templates/instance-down-rules.yml"]="${OPENIM_ROOT}/config/instance-down-rules.yml"
  ["${OPENIM_ROOT}/deployments/templates/notification.yaml"]="${OPENIM_ROOT}/config/notification.yaml"
)

# Templates for config Copy file
declare -A COPY_EXAMPLES=(
  ["${OPENIM_ROOT}/deployments/templates/email.tmpl"]="${OPENIM_ROOT}/config/templates/email.tmpl.template"
  ["${OPENIM_ROOT}/deployments/templates/instance-down-rules.yml"]="${OPENIM_ROOT}/config/templates/instance-down-rules.yml.template"
  ["${OPENIM_ROOT}/deployments/templates/notification.yaml"]="${OPENIM_ROOT}/config/templates/notification.yaml.template"
)

# Command-line options
FORCE_OVERWRITE=false
SKIP_EXISTING=false
GENERATE_EXAMPLES=false
CLEAN_CONFIG=false
CLEAN_EXAMPLES=false

# Function to display help information
show_help() {
  echo "Usage: $(basename "$0") [options]"
  echo "Options:"
  echo "  -h, --help             Show this help message"
  echo "  --force                Overwrite existing files without prompt"
  echo "  --skip                 Skip generation if file exists"
  echo "  --examples             Generate example files"
  echo "  --clean-config         Clean all configuration files"
  echo "  --clean-examples       Clean all example files"
}

# Function to generate and copy configuration files
generate_config_files() {
  # Handle TEMPLATES array
  for template in "${!TEMPLATES[@]}"; do
    local output_file="${TEMPLATES[$template]}"
    process_file "$template" "$output_file" true
  done
  
  # Handle COPY_TEMPLATES array
  for template in "${!COPY_TEMPLATES[@]}"; do
    local output_file="${COPY_TEMPLATES[$template]}"
    process_file "$template" "$output_file" false
  done
}

# Function to generate example files
generate_example_files() {
  env_cmd="env -i"
  
  env_vars["OPENIM_IP"]="127.0.0.1"
  env_vars["LOG_STORAGE_LOCATION"]="../../"
  
  for var in "${!env_vars[@]}"; do
    env_cmd+=" $var='${env_vars[$var]}'"
  done
  
  # Processing EXAMPLES array
  for template in "${!EXAMPLES[@]}"; do
    local example_file="${EXAMPLES[$template]}"
    process_file "$template" "$example_file" true
  done
  
  # Processing COPY_EXAMPLES array
  for template in "${!COPY_EXAMPLES[@]}"; do
    local example_file="${COPY_EXAMPLES[$template]}"
    process_file "$template" "$example_file" false
  done
}

# Function to process a single file, either by generating or copying
process_file() {
  local template=$1
  local output_file=$2
  local use_genconfig=$3
  
  if [[ -f "${output_file}" ]]; then
    if [[ "${FORCE_OVERWRITE}" == true ]]; then
      openim::log::info "Force overwriting ${output_file}."
      elif [[ "${SKIP_EXISTING}" == true ]]; then
      openim::log::info "Skipping generation of ${output_file} as it already exists."
      return
    else
      echo -n "File ${output_file} already exists. Overwrite? (Y/N): "
      read -r -n 1 REPLY
      echo
      if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        openim::log::info "Skipping generation of ${output_file}."
        return
      fi
    fi
  else
    if [[ "${SKIP_EXISTING}" == true ]]; then
      openim::log::info "Generating ${output_file} as it does not exist."
    fi
  fi
  
  if [[ "$use_genconfig" == true ]]; then
    openim::log::info "⌚  Working with template file: ${template} to generate ${output_file}..."
    if [[ ! -f "${OPENIM_ROOT}/scripts/genconfig.sh" ]]; then
      openim::log::error "genconfig.sh script not found"
      exit 1
    fi
    if [[ -n "${env_cmd}" ]]; then

    {
        printf "debugggggggggggggggggggg file: %s template: %s\n" "${ENV_FILE}" "${template}"
    } | tee /tmp/debug.log


      eval "$env_cmd ${OPENIM_ROOT}/scripts/genconfig.sh '${ENV_FILE}' '${template}' > '${output_file}'" || {
        openim::log::error "Error processing template file ${template}"
        exit 1
      }
    else
      "${OPENIM_ROOT}/scripts/genconfig.sh" "${ENV_FILE}" "${template}" > "${output_file}" || {
        openim::log::error "Error processing template file ${template}"
        exit 1
      }
    fi
  else
    openim::log::info "📋 Copying ${template} to ${output_file}..."
    cp "${template}" "${output_file}" || {
      openim::log::error "Error copying template file ${template}"
      exit 1
    }
  fi
  
  sleep 0.5
}

clean_config_files() {
  local all_templates=("${TEMPLATES[@]}" "${COPY_TEMPLATES[@]}")
  
  for output_file in "${all_templates[@]}"; do
    if [[ -f "${output_file}" ]]; then
      rm -f "${output_file}"
      openim::log::info "Removed configuration file: ${output_file}"
    fi
  done
}

# Function to clean example files
clean_example_files() {
  local all_examples=("${EXAMPLES[@]}" "${COPY_EXAMPLES[@]}")
  
  for example_file in "${all_examples[@]}"; do
    if [[ -f "${example_file}" ]]; then
      rm -f "${example_file}"
      openim::log::info "Removed example file: ${example_file}"
    fi
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
    --clean-config)
      CLEAN_CONFIG=true
      shift
    ;;
    --clean-examples)
      CLEAN_EXAMPLES=true
      shift
    ;;
    *)
      echo "Unknown option: $1"
      show_help
      exit 1
    ;;
  esac
done

# Clean configuration files if --clean-config option is provided
if [[ "${CLEAN_CONFIG}" == true ]]; then
  clean_config_files
fi

# Clean example files if --clean-examples option is provided
if [[ "${CLEAN_EXAMPLES}" == true ]]; then
  clean_example_files
fi

# Generate configuration files if requested
if [[ "${FORCE_OVERWRITE}" == true || "${SKIP_EXISTING}" == false ]] && [[ "${CLEAN_CONFIG}" == false ]]; then
  generate_config_files
fi

# Generate configuration files if requested
if [[ "${SKIP_EXISTING}" == true ]]; then
  generate_config_files
fi

# Generate example files if --examples option is provided
if [[ "${GENERATE_EXAMPLES}" == true ]] && [[ "${CLEAN_EXAMPLES}" == false ]]; then
  generate_example_files
fi

openim::log::success "Configuration and example files operation complete!"
