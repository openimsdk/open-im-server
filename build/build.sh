#!/usr/bin/env bash

set -euo pipefail

CYAN='\033[0;36m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NO_COLOR='\033[0m'

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="$BASE_DIR/images/openim-server/docker-compose.build.yml"
RELEASE="${RELEASE:-false}"
PUSH="${PUSH:-false}"
DRY_RUN="${DRY_RUN:-false}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"
IMAGE_TAGS="${IMAGE_TAGS:-}"
IMAGE_REGISTRIES="${IMAGE_REGISTRIES:-}"

if [[ ! -f "$COMPOSE_FILE" ]]; then
  echo -e "${RED}docker-compose.build.yml not found: $COMPOSE_FILE${NO_COLOR}"
  exit 1
fi

cd "$BASE_DIR/.." || exit 1

split_values() {
  echo "$1" | grep -o '[^, ]\+'
}

run_or_print() {
  if [[ "$DRY_RUN" == "true" ]]; then
    printf '%q ' "$@"
    printf '\n'
  else
    "$@"
  fi
}

build_local() {
  echo -e "${CYAN}Building all services...${NO_COLOR}"
  while IFS= read -r service; do
    RELEASE="$RELEASE" docker compose -f "$COMPOSE_FILE" build "$service"
  done < <(docker compose -f "$COMPOSE_FILE" config --services)

  echo -e "${CYAN}Tagging compatibility images for Kubernetes...${NO_COLOR}"
  while IFS= read -r built_image; do
    [[ -n "$built_image" ]] || continue
    compatibility_tag="${built_image##*/}"

    if [[ "$built_image" != "$compatibility_tag" ]]; then
      docker tag "$built_image" "$compatibility_tag"
    fi
  done < <(docker compose -f "$COMPOSE_FILE" config --images)

  echo -e "${GREEN}Successfully built all services${NO_COLOR}"
}

build_push() {
  if ! command -v jq >/dev/null 2>&1; then
    echo -e "${RED}jq is required for PUSH=true${NO_COLOR}"
    exit 1
  fi

  if [[ -z "$IMAGE_TAGS" || -z "$IMAGE_REGISTRIES" ]]; then
    echo -e "${RED}IMAGE_TAGS and IMAGE_REGISTRIES are required for PUSH=true${NO_COLOR}"
    exit 1
  fi

  image_tags=()
  while IFS= read -r tag; do
    image_tags+=("$tag")
  done < <(split_values "$IMAGE_TAGS")

  image_registries=()
  while IFS= read -r registry; do
    image_registries+=("$registry")
  done < <(split_values "$IMAGE_REGISTRIES")

  if [[ ${#image_tags[@]} -eq 0 || ${#image_registries[@]} -eq 0 ]]; then
    echo -e "${RED}IMAGE_TAGS and IMAGE_REGISTRIES must contain at least one value${NO_COLOR}"
    exit 1
  fi

  compose_config=$(docker compose -f "$COMPOSE_FILE" config --format json)

  echo -e "${CYAN}Building and pushing service images...${NO_COLOR}"
  while IFS= read -r service; do
    context=$(jq -r --arg service "$service" '.services[$service].build.context // empty' <<< "$compose_config")
    dockerfile=$(jq -r --arg service "$service" '.services[$service].build.dockerfile // empty' <<< "$compose_config")
    cmd_path=$(jq -r --arg service "$service" '.services[$service].build.args.CMD_PATH // empty' <<< "$compose_config")
    binary_name=$(jq -r --arg service "$service" '.services[$service].build.args.BINARY_NAME // empty' <<< "$compose_config")

    if [[ -z "$context" || -z "$dockerfile" || -z "$cmd_path" || -z "$binary_name" ]]; then
      echo -e "${RED}Invalid build config for $service${NO_COLOR}"
      exit 1
    fi

    if [[ ! -d "$cmd_path" && ! -f "$cmd_path/main.go" ]]; then
      echo -e "${CYAN}Skipping $service because $cmd_path does not exist${NO_COLOR}"
      continue
    fi

    tag_args=()
    for registry in "${image_registries[@]}"; do
      for tag in "${image_tags[@]}"; do
        tag_args+=(--tag "$registry/$binary_name:$tag")
      done
    done

    echo -e "${CYAN}Building $binary_name for $PLATFORMS...${NO_COLOR}"
    run_or_print docker buildx build \
      --platform "$PLATFORMS" \
      --file "$dockerfile" \
      --build-arg "CMD_PATH=$cmd_path" \
      --build-arg "BINARY_NAME=$binary_name" \
      --build-arg "RELEASE=$RELEASE" \
      "${tag_args[@]}" \
      --push \
      "$context"
  done < <(jq -r '.services | keys[]' <<< "$compose_config" | sort)

  echo -e "${GREEN}Successfully pushed service images${NO_COLOR}"
}

if [[ "$PUSH" == "true" ]]; then
  build_push
else
  build_local
fi
