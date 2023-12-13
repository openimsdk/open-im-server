#!/usr/bin/env bash

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

# Function to generate configuration files
generate_config_files() {
  for template in "${!TEMPLATES[@]}"; do
    local output_file="${TEMPLATES[$template]}"
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
    else
      if [[ "${SKIP_EXISTING}" == true ]]; then
        openim::log::info "Generating ${output_file} as it does not exist."
      fi
    fi

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
}

# Function to generate example files
generate_example_files() {
  for template in "${!EXAMPLES[@]}"; do
    local example_file="${EXAMPLES[$template]}"
    if [[ -f "${example_file}" ]]; then
      if [[ "${FORCE_OVERWRITE}" == true ]]; then
        openim::log::info "Force overwriting example file: ${example_file}."
      elif [[ "${SKIP_EXISTING}" == true ]]; then
        openim::log::info "Skipping generation of example file: ${example_file} as it already exists."
        continue
      else
        echo -n "Example file ${example_file} already exists. Overwrite? (Y/N): "
        read -r -n 1 REPLY
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
          openim::log::info "Skipping generation of example file: ${example_file}."
          continue
        fi
      fi
    elif [[ "${SKIP_EXISTING}" == true ]]; then
      openim::log::info "Generating example file: ${example_file} as it does not exist."
    fi

    openim::log::info "⌚  Working with template file: ${template} to generate example file: ${example_file}..."
    if [[ ! -f "${OPENIM_ROOT}/scripts/genconfig.sh" ]]; then
      openim::log::error "genconfig.sh script not found"
      exit 1
    fi
    "${OPENIM_ROOT}/scripts/genconfig.sh" "${ENV_FILE}" "${template}" > "${example_file}" || {
      openim::log::error "Error processing template file ${template}"
      exit 1
    }
    sleep 0.5
  done
}


# Function to clean configuration files
clean_config_files() {
  for output_file in "${TEMPLATES[@]}"; do
    if [[ -f "${output_file}" ]]; then
      rm -f "${output_file}"
      openim::log::info "Removed configuration file: ${output_file}"
    fi
  done
}

# Function to clean example files
clean_example_files() {
  for example_file in "${EXAMPLES[@]}"; do
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
