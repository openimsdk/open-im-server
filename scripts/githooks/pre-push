#!/usr/bin/env bash
# Copyright © 2023 OpenIMSDK.
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
# ==============================================================================
#

YELLOW="\e[93m"
GREEN="\e[32m"
RED="\e[31m"
ENDCOLOR="\e[0m"

local_branch="$(git rev-parse --abbrev-ref HEAD)"
valid_branch_regex="^(main|master|develop|release(-[a-zA-Z0-9._-]+)?)$|(feature|feat|openim|hotfix|test|bug|ci|cicd|style|)\/[a-z0-9._-]+$|^HEAD$"

printMessage() {
   printf "${YELLOW}OpenIM : $1${ENDCOLOR}\n"
}

printSuccess() {
   printf "${GREEN}OpenIM : $1${ENDCOLOR}\n"
}

printError() {
   printf "${RED}OpenIM : $1${ENDCOLOR}\n"
}

printMessage "Running local OpenIM pre-push hook."

if [[ `git status --porcelain` ]]; then
  printError "This scripts needs to run against committed code only. Please commit or stash you changes."
  exit 1
fi

COLOR_SUFFIX="\033[0m"

BLACK_PREFIX="\033[30m"
RED_PREFIX="\033[31m"
GREEN_PREFIX="\033[32m"
BACKGROUND_GREEN="\033[33m"
BLUE_PREFIX="\033[34m"
PURPLE_PREFIX="\033[35m"
SKY_BLUE_PREFIX="\033[36m"
WHITE_PREFIX="\033[37m"
BOLD_PREFIX="\033[1m"
UNDERLINE_PREFIX="\033[4m"
ITALIC_PREFIX="\033[3m"

# Function to print colored text
print_color() {
  local text=$1
  local color=$2
  echo -e "${color}${text}${COLOR_SUFFIX}"
}

# Function to print section separator
print_separator() {
  print_color "==========================================================" ${PURPLE_PREFIX}
}

# Get current time
time=$(date +"%Y-%m-%d %H:%M:%S")

# Print section separator
print_separator

# Print time of submission
print_color "PTIME: ${time}" "${BOLD_PREFIX}${CYAN_PREFIX}"
echo ""
author=$(git config user.name)
repository=$(basename -s .git $(git config --get remote.origin.url))

# Print additional information if needed
print_color "Repository: ${repository}" "${BLUE_PREFIX}"
echo ""

print_color "Author: ${author}" "${PURPLE_PREFIX}"

# Print section separator
print_separator

file_list=$(git diff --name-status HEAD @{u})
added_files=$(grep -c '^A' <<< "$file_list")
modified_files=$(grep -c '^M' <<< "$file_list")
deleted_files=$(grep -c '^D' <<< "$file_list")

print_color "Added Files: ${added_files}" "${BACKGROUND_GREEN}"
print_color "Modified Files: ${modified_files}" "${BACKGROUND_GREEN}"
print_color "Deleted Files: ${deleted_files}" "${BACKGROUND_GREEN}"

if [[ ! $local_branch =~ $valid_branch_regex ]]
then
    printError "There is something wrong with your branch name. Branch names in this project must adhere to this contract: $valid_branch_regex. 
Your commit will be rejected. You should rename your branch to a valid name(feat/name OR bug/name) and try again."
    printError "For more on this, read on: https://gist.github.com/cubxxw/126b72104ac0b0ca484c9db09c3e5694"
    exit 1
fi

#
#printMessage "Running the Flutter analyzer"
#flutter analyze
#
#if [ $? -ne 0 ]; then
#  printError "Flutter analyzer error"
#  exit 1
#fi
#
#printMessage "Finished running the Flutter analyzer"
