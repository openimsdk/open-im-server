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
#
# -----------------------------------------------------------------------------
# init-githooks.sh
#
# This script assists in managing Git hooks for the OpenIM project.
# When executed:
# 1. It prompts the user to enable git hooks.
# 2. If the user accepts, it copies predefined hook scripts to the appropriate
#    Git directory, making them executable.
# 3. If requested, it can delete the added hooks.
#
# This script equal runs `make init-githooks` command.
# Usage:
# ./init-githooks.sh              Prompt to enable git hooks.
# ./init-githooks.sh --delete     Delete previously added git hooks.
# ./init-githooks.sh --help       Show the help message.
#
# Example: `scripts/build-go.sh --help`.
# Documentation & related context can be found at:
# https://gist.github.com/cubxxw/126b72104ac0b0ca484c9db09c3e5694
#
# -----------------------------------------------------------------------------

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
HOOKS_DIR="${OPENIM_ROOT}/.git/hooks"

help_info() {
    echo "Usage: $0 [options]"
    echo
    echo "This script helps to manage git hooks."
    echo
    echo "Options:"
    echo "  -h, --help       Show this help message and exit."
    echo "  -d, --delete     Delete the hooks that have been added."
    echo "  By default, it will prompt to enable git hooks."
}

delete_hooks() {
    for file in ${OPENIM_ROOT}/scripts/githooks/*.sh; do
        hook_name=$(basename "$file" .sh)  # This removes the .sh extension
        rm -f "$HOOKS_DIR/$hook_name"
    done
    echo "Git hooks have been deleted."
}

enable_hooks() {
    echo "Would you like to:"
    echo "1) Enable git hooks mode"
    echo "2) Delete existing git hooks"
    echo "Please select a number (or any other key to exit):"
    read -r choice

    case "$choice" in
        1)
            for file in ${OPENIM_ROOT}/scripts/githooks/*.sh; do
                hook_name=$(basename "$file" .sh)  # This removes the .sh extension
                cp -f "$file" "$HOOKS_DIR/$hook_name"
            done

            chmod +x $HOOKS_DIR/*
            
            echo "Git hooks mode has been enabled."
            echo "With git hooks enabled, every time you perform a git action (e.g. git commit), the corresponding hooks script will be triggered automatically."
            echo "This means that if the size of the file you're committing exceeds the set limit (e.g. 42MB), the commit will be rejected."
            ;;
        2)
            delete_hooks
            ;;
        *)
            echo "Exiting without making changes."
            ;;
    esac
}


case "$1" in
    -h|--help)
        help_info
        ;;
    -d|--delete)
        delete_hooks
        ;;
    *)
        enable_hooks
        ;;
esac
