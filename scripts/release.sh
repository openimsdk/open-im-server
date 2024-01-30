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

# Description:
# This script automates the process of building and releasing OpenIM,
# including tasks like setting up the environment, verifying prerequisites,
# building commands, packaging tarballs, uploading tarballs, creating GitHub
# releases, and generating changelogs.
#
# Usage:
# ./scripts/release.sh [options]
# Options include:
#   -h, --help               : Show help message
#   -se, --setup-env         : Execute setup environment
#   -vp, --verify-prereqs    : Execute prerequisites verification
#   -bc, --build-command     : Execute build command
#   -bi, --build-image       : Execute build image (default: not executed)
#   -pt, --package-tarballs  : Execute package tarballs
#   -ut, --upload-tarballs   : Execute upload tarballs
#   -gr, --github-release    : Execute GitHub release
#   -gc, --generate-changelog: Execute generate changelog
#
# This script can also be executed via the 'make release' command as an alternative.
#
# Dependencies:
# This script depends on external scripts found in the 'scripts' directory and
# assumes the presence of necessary tools and permissions for building and
# releasing software.
#
# Note:
# The script uses standard bash script practices with error handling,
# and it defaults to executing all steps if no specific option is provided.
#
# Build a OpenIM release.  This will build the binaries, create the Docker
# images and other build artifacts.
# Build a OpenIM release. This script supports various flags for flexible execution control.

set -o errexit
set -o nounset
set -o pipefail
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/common.sh"
source "${OPENIM_ROOT}/scripts/lib/release.sh"

OPENIM_RELEASE_RUN_TESTS=${OPENIM_RELEASE_RUN_TESTS-y}

# Function to show help message
show_help() {
  echo "Usage: $(basename $0) [options]"
  echo "Options:"
  echo "  -h, --help               Show this help message"
  echo "  -se, --setup-env         Execute setup environment"
  echo "  -vp, --verify-prereqs    Execute prerequisites verification"
  echo "  -bc, --build-command     Execute build command"
  echo "  -bi, --build-image       Execute build image (default: not executed)"
  echo "  -pt, --package-tarballs  Execute package tarballs"
  echo "  -ut, --upload-tarballs   Execute upload tarballs"
  echo "  -gr, --github-release    Execute GitHub release"
  echo "  -gc, --generate-changelog Execute generate changelog"
}

# Initialize all actions to false
perform_setup_env=false
perform_verify_prereqs=false
perform_build_command=false
perform_build_image=false # New flag for build image
perform_package_tarballs=false
perform_upload_tarballs=false
perform_github_release=false
perform_generate_changelog=false

# Process command-line arguments
while getopts "hsevpbciptutgrgc-" opt; do
  case "${opt}" in
    h)  show_help; exit 0 ;;
    se) perform_setup_env=true ;;
    vp) perform_verify_prereqs=true ;;
    bc) perform_build_command=true ;;
    bi) perform_build_image=true ;; # Handling new option
    pt) perform_package_tarballs=true ;;
    ut) perform_upload_tarballs=true ;;
    gr) perform_github_release=true ;;
    gc) perform_generate_changelog=true ;;
    --) case "${OPTARG}" in
        help)               show_help; exit 0 ;;
        setup-env)          perform_setup_env=true ;;
        verify-prereqs)     perform_verify_prereqs=true ;;
        build-command)      perform_build_command=true ;;
        build-image)        perform_build_image=true ;; # Handling new long option
        package-tarballs)   perform_package_tarballs=true ;;
        upload-tarballs)    perform_upload_tarballs=true ;;
        github-release)     perform_github_release=true ;;
        generate-changelog) perform_generate_changelog=true ;;
        *) echo "Invalid option: --${OPTARG}"; show_help; exit 1 ;;
    esac ;;
    *) show_help; exit 1 ;;
  esac
done

# Enable all actions by default if no options are provided
if [ "$#" -eq 0 ]; then
  perform_setup_env=true
  perform_verify_prereqs=true
  perform_build_command=true
  perform_package_tarballs=true
  perform_upload_tarballs=true
  perform_github_release=true
  perform_generate_changelog=true
  # TODO: Not enabling build_image by default
  # perform_build_image=true
fi

# Function to perform actions
perform_action() {
  local flag=$1
  local message=$2
  local command=$3
  
  if [ "$flag" == true ]; then
    openim::log::info "## $message..."
    if ! eval "$command"; then
      openim::log::errexit "Error in $message"
    fi
  fi
}

echo "Starting script execution..."

perform_action $perform_setup_env "Setting up environment" "openim::golang::setup_env"
perform_action $perform_verify_prereqs "Verifying prerequisites" "openim::build::verify_prereqs && openim::release::verify_prereqs"
perform_action $perform_build_command "Executing build command" "openim::build::build_command"
perform_action $perform_build_image "Building image" "openim::build::build_image"
perform_action $perform_package_tarballs "Packaging tarballs" "openim::release::package_tarballs"
perform_action $perform_upload_tarballs "Uploading tarballs" "openim::release::upload_tarballs"
perform_action $perform_github_release "Creating GitHub release" "openim::release::github_release"
perform_action $perform_generate_changelog "Generating changelog" "openim::release::generate_changelog"

openim::log::success "OpenIM Relase Script Execution Completed."
