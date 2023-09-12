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

DEFAULT_DIRS=(
    "pkg"
    "internal/pkg"
)
BASE_URL="github.com/openimsdk/open-im-server"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "This script iterates over directories and generates doc.go if necessary."
    echo "By default, it processes 'pkg' and 'internal/pkg' directories."
    echo
    echo "Options:"
    echo "  -d DIRS, --dirs DIRS    Specify the directories to be processed, separated by commas. E.g., 'pkg,internal/pkg'."
    echo "  -u URL,  --url URL     Set the base URL for the import path. Default is '$BASE_URL'."
    echo "  -h,      --help        Show this help message."
    echo
}

process_dir() {
    local dir=$1
    local base_url=$2

    for d in $(find $dir -type d); do
        if [ ! -f $d/doc.go ]; then
            if ls $d/*.go > /dev/null 2>&1; then
                echo $d/doc.go
                echo "package $(basename $d) // import \"$base_url/$d\"" > $d/doc.go
            fi
        fi
    done
}

while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
        -d|--dirs)
            IFS=',' read -ra DIRS <<< "$2"
            shift # shift past argument
            shift # shift past value
            ;;
        -u|--url)
            BASE_URL="$2"
            shift # shift past argument
            shift # shift past value
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            usage
            exit 1
            ;;
    esac
done

DIRS=${DIRS:-${DEFAULT_DIRS[@]}}

for dir in "${DIRS[@]}"; do
    process_dir $dir $BASE_URL
done
