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


source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
echo -e  "\n"

echo -e "${BACKGROUND_BLUE}===============> Building all using make build binary files ${COLOR_SUFFIX}" 

echo -e  "\n"
echo -e "${BOLD_PREFIX}_____  _   _  _____  _____  _____  _____  _____  _____  _____  _____  _____ ${COLOR_SUFFIX}"


bin_dir="$BIN_DIR"
logs_dir="../logs"
sdk_db_dir="../db/sdk/"
# Automatically created when there is no bin, logs folder
if [ ! -d $bin_dir ]; then
  mkdir -p $bin_dir
fi
if [ ! -d $logs_dir ]; then
  mkdir -p $logs_dir
fi
if [ ! -d $sdk_db_dir ]; then
  mkdir -p $sdk_db_dir
fi

# OpenIM root path
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

#!/bin/bash

cd $OPENIM_ROOT

# Execute 'make build'
make build

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

# Determine if all scripts were successfully built
BUILD_SUCCESS=true
FAILED_SCRIPTS=()

for binary in $(find _output/bin/platforms/$REPO_DIR -type f); do
    if [[ ! -x $binary ]]; then
        FAILED_SCRIPTS+=("$binary")
        BUILD_SUCCESS=false
    fi
done

echo -e " \n"

echo -e "${BOLD_PREFIX}=====================>  Build Results <=====================${COLOR_SUFFIX}"

echo -e " \n"

if [[ "$BUILD_SUCCESS" == true ]]; then
    echo -e "${GREEN_PREFIX}All binaries built successfully.${COLOR_SUFFIX}"
else
    echo -e "${RED_PREFIX}Some binary builds failed. Please check the following binary files:${COLOR_SUFFIX}"
    for script in "${FAILED_SCRIPTS[@]}"; do
        echo -e "${RED_PREFIX}$script${COLOR_SUFFIX}"
    done
fi

echo -e " \n"

echo -e "${BOLD_PREFIX}============================================================${COLOR_SUFFIX}"

echo -e " \n"
