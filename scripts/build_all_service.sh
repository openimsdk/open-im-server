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

#Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/path_info.sh
source $SCRIPTS_ROOT/lib/init.sh

bin_dir="$BIN_DIR"
logs_dir="$OPENIM_ROOT/logs"
sdk_db_dir="$OPENIM_ROOT/db/sdk/"

echo "==> bin_dir=$bin_dir"
echo "==> logs_dir=$logs_dir"
echo "==> sdk_db_dir=$sdk_db_dir"

# Automatically created when there is no bin, logs folder
if [ ! -d $logs_dir ]; then
  mkdir -p $logs_dir
fi
if [ ! -d $sdk_db_dir ]; then
  mkdir -p $sdk_db_dir
fi

cd $OPENIM_ROOT

# CPU core number
# Check the system type
system_type=$(uname)

if [[ "$system_type" == "Darwin" ]]; then
    # macOS (using sysctl)
    cpu_count=$(sysctl -n hw.ncpu)
elif [[ "$system_type" == "Linux" ]]; then
    # Linux (using lscpu)
    cpu_count=$(lscpu --parse | grep -E '^([^#].*,){3}[^#]' | sort -u | wc -l)
else
    echo "Unsupported operating system: $system_type"
    exit 1
fi
echo -e "${GREEN_PREFIX}======> cpu_count=$cpu_count${COLOR_SUFFIX}"

# Count the number of concurrent compilations (half the number of cpus)
compile_count=$((cpu_count / 2))

# Execute 'make build' run the make command for concurrent compilation
make -j$compile_count build

if [ $? -ne 0 ]; then
  echo "make build Error, script exits"
  exit 1
fi

openim::util::gen_os_arch

# Determine if all scripts were successfully built
BUILD_SUCCESS=true
FAILED_SCRIPTS=()

for binary in $(find _output/bin/platforms/$REPO_DIR -type f); do
    if [[ ! -x $binary ]]; then
        FAILED_SCRIPTS+=("$binary")
        BUILD_SUCCESS=false
    fi
done

echo -e " "

echo -e "${BOLD_PREFIX}=====================>  Build Results <=====================${COLOR_SUFFIX}"

echo -e " "

if [[ "$BUILD_SUCCESS" == true ]]; then
    echo -e "${GREEN_PREFIX}All binaries built successfully.${COLOR_SUFFIX}"
else
    echo -e "${RED_PREFIX}Some binary builds failed. Please check the following binary files:${COLOR_SUFFIX}"
    for script in "${FAILED_SCRIPTS[@]}"; do
        echo -e "${RED_PREFIX}$script${COLOR_SUFFIX}"
    done
fi

echo -e " "

echo -e "${BOLD_PREFIX}============================================================${COLOR_SUFFIX}"

echo -e " "
