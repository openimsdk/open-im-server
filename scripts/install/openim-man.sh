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

# openim-man.sh Script to manage man pages for openim
#
# Description:
#   This script manages the man pages for the OpenIM software suite.
#   It provides facilities to install, uninstall, and verify the 
#   installation status of the man pages related to OpenIM components.
#
# Usage:
#   ./openim-man.sh openim::man::install      - Install man pages
#   ./openim-man.sh openim::man::uninstall    - Uninstall man pages
#   ./openim-man.sh openim::man::status       - Check installation status
#
# Dependencies:
#   - Assumes there's a common.sh in ""${OPENIM_ROOT}"/scripts/install/" 
#     containing shared functions and variables.
#   - Relies on the script ""${OPENIM_ROOT}"/scripts/update-generated-docs.sh" 
#     to generate the man pages.
#
# Notes:
#   - This script must be run with appropriate permissions to modify the 
#     system man directories.
#   - Always ensure you're in the script's directory or provide the correct 
#     path when executing.
################################################################################

# Define the root of the build/dist directory
OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

# Ensure the common script is sourced
[[ -z ${COMMON_SOURCED} ]] && source "${OPENIM_ROOT}/scripts/install/common.sh"

# Print usage information after installation
function openim::man::info() {
cat <<- EOF
Usage:
  man openim-server  # Display the man page for openim-server
EOF
}

# Install the man pages for openim
function openim::man::install() {
    # Navigate to the openim root directory
    pushd "${OPENIM_ROOT}" > /dev/null

    # Generate man pages for each component
    ""${OPENIM_ROOT}"/scripts/update-generated-docs.sh"
    openim::common::sudo "cp docs/man/man1/* /usr/share/man/man1/"
    
    # Verify installation status
    if openim::man::status; then
        openim::log::info "Installed openim-server man page successfully"
        openim::man::info
    fi

    # Return to the original directory
    popd > /dev/null
}

# Uninstall the man pages for openim
function openim::man::uninstall() {
    # Turn off exit-on-error temporarily to handle non-existing files gracefully
    set +o errexit
    openim::common::sudo "rm -f /usr/share/man/man1/openim-*"
    set -o errexit
    
    openim::log::info "Uninstalled openim man pages successfully"
}

# Check the installation status of the man pages
function openim::man::status() {
    if ! ls /usr/share/man/man1/openim-* &> /dev/null; then
        openim::log::error "OpenIM man files not found. Perhaps they were not installed correctly."
        return 1
    fi
    return 0
}

# Execute the appropriate function based on the given arguments
if [[ "$*" =~ openim::man:: ]]; then
    eval "$*"
fi
