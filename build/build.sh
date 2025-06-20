#!/usr/bin/env bash

CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NO_COLOR='\033[0m'

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
IMAGE_DIR="$BASE_DIR/images"

if [[ ! -d "$IMAGE_DIR" ]]; then
  echo -e "${RED}dir $IMAGE_DIR not exist!${NO_COLOR}"
  exit 1
fi

for dir in "$IMAGE_DIR"/*/; do
  [[ -d "$dir" ]] || continue

  name="$(basename "$dir")"
  dockerfile="$dir/Dockerfile"

  if [[ -f "$dockerfile" ]]; then
    echo -e "${CYAN}Building ${name}:test...${NO_COLOR}"
    build_context="${dir}../../../"
    if docker build -t "${name}:test" -f "$dockerfile" "$build_context"; then
      echo -e "${GREEN}Successfully built ${name}:test${NO_COLOR}"
    else
      echo -e "${RED}Failed to build ${name}:test${NO_COLOR}"
    fi
  else
    echo -e "${YELLOW}Skipping ${name}: Dockerfile not found${NO_COLOR}"
  fi
done