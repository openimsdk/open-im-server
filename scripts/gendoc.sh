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

#!/bin/bash

DEFAULT_DIRS=("pkg")
BASE_URL="github.com/openimsdk/open-im-server/v3"
REMOVE_DOC=false

usage() {
  echo "Usage: $0 [OPTIONS]"
  echo
  echo "This script iterates over directories. By default, it generates doc.go files for 'pkg' and 'internal/pkg'."
  echo
  echo "Options:"
  echo "  -d DIRS, --dirs DIRS    Specify directories to process, separated by commas (e.g., 'pkg,internal/pkg')."
  echo "  -u URL, --url URL       Set the base URL for the import path. Default is '$BASE_URL'."
  echo "  -r, --remove            Remove all doc.go files in the specified directories."
  echo "  -h, --help              Show this help message."
  echo
}

process_dir() {
  local dir="$1"
  local base_url="$2"
  local remove_doc="$3"

  find "$dir" -type d | while read -r d; do
    if [ "$remove_doc" = true ]; then
      if [ -f "$d/doc.go" ]; then
        echo "Removing $d/doc.go"
        rm -f "$d/doc.go"
      fi
    else
      if [ ! -f "$d/doc.go" ] && ls "$d/"*.go &>/dev/null; then
        echo "Creating $d/doc.go"
        echo "package $(basename "$d") // import \"$base_url/$(echo "$d" | sed "s|^\./||")\"" >"$d/doc.go"
      fi
    fi
  done
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -d|--dirs)
      IFS=',' read -ra DIRS <<< "$2"
      shift 2
      ;;
    -u|--url)
      BASE_URL="$2"
      shift 2
      ;;
    -r|--remove)
      REMOVE_DOC=true
      shift
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

DIRS=(${DIRS:-"${DEFAULT_DIRS[@]}"})

for dir in "${DIRS[@]}"; do
  process_dir "$dir" "$BASE_URL" "$REMOVE_DOC"
done
