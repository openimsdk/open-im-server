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


# shellcheck disable=SC2034 # Variables sourced in other scripts.

# Common utilities, variables and checks for all build scripts.
set -o errexit
set +o nounset
set -o pipefail

# Unset CDPATH, having it set messes up with script import paths
unset CDPATH

USER_ID=$(id -u)
GROUP_ID=$(id -g)

DOCKER_OPTS=${DOCKER_OPTS:-""}
IFS=" " read -r -a DOCKER <<< "docker ${DOCKER_OPTS}"
DOCKER_HOST=${DOCKER_HOST:-""}
DOCKER_MACHINE_NAME=${DOCKER_MACHINE_NAME:-"openim-dev"}
readonly DOCKER_MACHINE_DRIVER=${DOCKER_MACHINE_DRIVER:-"virtualbox --virtualbox-cpu-count -1"}

# This will canonicalize the path
OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd -P)

# Please do not refer to lib after referring to common
. $(dirname ${BASH_SOURCE})/lib/init.sh

# Constants
readonly OPENIM_BUILD_IMAGE_REPO=openim-build
#readonly OPENIM_BUILD_IMAGE_CROSS_TAG="$(cat ""${OPENIM_ROOT}"/build/build-image/cross/VERSION")"

readonly OPENIM_DOCKER_REGISTRY="${OPENIM_DOCKER_REGISTRY:-k8s.gcr.io}"
readonly OPENIM_BASE_IMAGE_REGISTRY="${OPENIM_BASE_IMAGE_REGISTRY:-us.gcr.io/k8s-artifacts-prod/build-image}"

# This version number is used to cause everyone to rebuild their data containers
# and build image.  This is especially useful for automated build systems like
# Jenkins.
#
# Increment/change this number if you change the build image (anything under
# build/build-image) or change the set of volumes in the data container.
#readonly OPENIM_BUILD_IMAGE_VERSION_BASE="$(cat ""${OPENIM_ROOT}"/build/build-image/VERSION")"
#readonly OPENIM_BUILD_IMAGE_VERSION="${OPENIM_BUILD_IMAGE_VERSION_BASE}-${OPENIM_BUILD_IMAGE_CROSS_TAG}"

# Here we map the output directories across both the local and remote _output
# directories:
#
# *_OUTPUT_ROOT    - the base of all output in that environment.
# *_OUTPUT_SUBPATH - location where golang stuff is built/cached.  Also
#                    persisted across docker runs with a volume mount.
# *_OUTPUT_BINPATH - location where final binaries are placed.  If the remote
#                    is really remote, this is the stuff that has to be copied
#                    back.
# OUT_DIR can come in from the Makefile, so honor it.
readonly LOCAL_OUTPUT_ROOT=""${OPENIM_ROOT}"/${OUT_DIR:-_output}"
readonly LOCAL_OUTPUT_SUBPATH="${LOCAL_OUTPUT_ROOT}/platforms"
readonly LOCAL_OUTPUT_BINPATH="${LOCAL_OUTPUT_SUBPATH}"
readonly LOCAL_OUTPUT_GOPATH="${LOCAL_OUTPUT_SUBPATH}/go"
readonly LOCAL_OUTPUT_IMAGE_STAGING="${LOCAL_OUTPUT_ROOT}/images"

# This is the port on the workstation host to expose RSYNC on.  Set this if you
# are doing something fancy with ssh tunneling.
readonly OPENIM_RSYNC_PORT="${OPENIM_RSYNC_PORT:-}"

# This is the port that rsync is running on *inside* the container. This may be
# mapped to OPENIM_RSYNC_PORT via docker networking.
readonly OPENIM_CONTAINER_RSYNC_PORT=8730

# Get the set of master binaries that run in Docker (on Linux)
# Entry format is "<name-of-binary>,<base-image>".
# Binaries are placed in /usr/local/bin inside the image.
#
# $1 - server architecture
openim::build::get_docker_wrapped_binaries() {
  local arch=$1
  local debian_base_version=v2.1.0
  local debian_iptables_version=v12.1.0
  ### If you change any of these lists, please also update DOCKERIZED_BINARIES
  ### in build/BUILD. And openim::golang::server_image_targets

  local targets=(
    "openim-api,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-cmdutils,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-crontask,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-msggateway,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-msgtransfer,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-push,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-rpc-auth,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-rpc-conversation,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-rpc-friend,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-rpc-group,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-rpc-msg,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-rpc-third,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "openim-rpc-user,${OPENIM_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
  )
  echo "${targets[@]}"
}

# ---------------------------------------------------------------------------
# Basic setup functions

# Verify that the right utilities and such are installed for building openim. Set
# up some dynamic constants.
# Args:
#   $1 - boolean of whether to require functioning docker (default true)
#
# Vars set:
#   OPENIM_ROOT_HASH
#   OPENIM_BUILD_IMAGE_TAG_BASE
#   OPENIM_BUILD_IMAGE_TAG
#   OPENIM_BUILD_IMAGE
#   OPENIM_BUILD_CONTAINER_NAME_BASE
#   OPENIM_BUILD_CONTAINER_NAME
#   OPENIM_DATA_CONTAINER_NAME_BASE
#   OPENIM_DATA_CONTAINER_NAME
#   OPENIM_RSYNC_CONTAINER_NAME_BASE
#   OPENIM_RSYNC_CONTAINER_NAME
#   DOCKER_MOUNT_ARGS
#   LOCAL_OUTPUT_BUILD_CONTEXT
function openim::build::verify_prereqs() {
  local -r require_docker=${1:-true}
  openim::log::status "Verifying Prerequisites...."
  openim::build::ensure_tar || return 1
  openim::build::ensure_rsync || return 1
  if ${require_docker}; then
    openim::build::ensure_docker_in_path || return 1
    openim::util::ensure_docker_daemon_connectivity || return 1

    if (( OPENIM_VERBOSE > 6 )); then
      openim::log::status "Docker Version:"
      "${DOCKER[@]}" version | openim::log::info_from_stdin
    fi
  fi

  OPENIM_GIT_BRANCH=$(git symbolic-ref --short -q HEAD 2>/dev/null || true)
  OPENIM_ROOT_HASH=$(openim::build::short_hash "${HOSTNAME:-}:"${OPENIM_ROOT}":${OPENIM_GIT_BRANCH}")
  OPENIM_BUILD_IMAGE_TAG_BASE="build-${OPENIM_ROOT_HASH}"
  #OPENIM_BUILD_IMAGE_TAG="${OPENIM_BUILD_IMAGE_TAG_BASE}-${OPENIM_BUILD_IMAGE_VERSION}"
  #OPENIM_BUILD_IMAGE="${OPENIM_BUILD_IMAGE_REPO}:${OPENIM_BUILD_IMAGE_TAG}"
  OPENIM_BUILD_CONTAINER_NAME_BASE="openim-build-${OPENIM_ROOT_HASH}"
  #OPENIM_BUILD_CONTAINER_NAME="${OPENIM_BUILD_CONTAINER_NAME_BASE}-${OPENIM_BUILD_IMAGE_VERSION}"
  OPENIM_RSYNC_CONTAINER_NAME_BASE="openim-rsync-${OPENIM_ROOT_HASH}"
  #OPENIM_RSYNC_CONTAINER_NAME="${OPENIM_RSYNC_CONTAINER_NAME_BASE}-${OPENIM_BUILD_IMAGE_VERSION}"
  OPENIM_DATA_CONTAINER_NAME_BASE="openim-build-data-${OPENIM_ROOT_HASH}"
  #OPENIM_DATA_CONTAINER_NAME="${OPENIM_DATA_CONTAINER_NAME_BASE}-${OPENIM_BUILD_IMAGE_VERSION}"
  #DOCKER_MOUNT_ARGS=(--volumes-from "${OPENIM_DATA_CONTAINER_NAME}")
  #LOCAL_OUTPUT_BUILD_CONTEXT="${LOCAL_OUTPUT_IMAGE_STAGING}/${OPENIM_BUILD_IMAGE}"

  openim::version::get_version_vars
  #openim::version::save_version_vars ""${OPENIM_ROOT}"/.dockerized-openim-version-defs"
}

# ---------------------------------------------------------------------------
# Utility functions

function openim::build::docker_available_on_osx() {
  if [[ -z "${DOCKER_HOST}" ]]; then
    if [[ -S "/var/run/docker.sock" ]]; then
      openim::log::status "Using Docker for MacOS"
      return 0
    fi

    openim::log::status "No docker host is set. Checking options for setting one..."
    if [[ -z "$(which docker-machine)" ]]; then
      openim::log::status "It looks like you're running Mac OS X, yet neither Docker for Mac nor docker-machine can be found."
      openim::log::status "See: https://docs.docker.com/engine/installation/mac/ for installation instructions."
      return 1
    elif [[ -n "$(which docker-machine)" ]]; then
      openim::build::prepare_docker_machine
    fi
  fi
}

function openim::build::prepare_docker_machine() {
  openim::log::status "docker-machine was found."

  local available_memory_bytes
  available_memory_bytes=$(sysctl -n hw.memsize 2>/dev/null)

  local bytes_in_mb=1048576

  # Give virtualbox 1/2 the system memory. Its necessary to divide by 2, instead
  # of multiple by .5, because bash can only multiply by ints.
  local memory_divisor=2

  local virtualbox_memory_mb=$(( available_memory_bytes / (bytes_in_mb * memory_divisor) ))

  docker-machine inspect "${DOCKER_MACHINE_NAME}" &> /dev/null || {
    openim::log::status "Creating a machine to build OPENIM"
    docker-machine create --driver "${DOCKER_MACHINE_DRIVER}" \
      --virtualbox-memory "${virtualbox_memory_mb}" \
      --engine-env HTTP_PROXY="${OPENIMRNETES_HTTP_PROXY:-}" \
      --engine-env HTTPS_PROXY="${OPENIMRNETES_HTTPS_PROXY:-}" \
      --engine-env NO_PROXY="${OPENIMRNETES_NO_PROXY:-127.0.0.1}" \
      "${DOCKER_MACHINE_NAME}" > /dev/null || {
      openim::log::error "Something went wrong creating a machine."
      openim::log::error "Try the following: "
      openim::log::error "docker-machine create -d ${DOCKER_MACHINE_DRIVER} --virtualbox-memory ${virtualbox_memory_mb} ${DOCKER_MACHINE_NAME}"
      return 1
    }
  }
  docker-machine start "${DOCKER_MACHINE_NAME}" &> /dev/null
  # it takes `docker-machine env` a few seconds to work if the machine was just started
  local docker_machine_out
  while ! docker_machine_out=$(docker-machine env "${DOCKER_MACHINE_NAME}" 2>&1); do
    if [[ ${docker_machine_out} =~ "Error checking TLS connection" ]]; then
      echo "${docker_machine_out}"
      docker-machine regenerate-certs "${DOCKER_MACHINE_NAME}"
    else
      sleep 1
    fi
  done
  eval "$(docker-machine env "${DOCKER_MACHINE_NAME}")"
  openim::log::status "A Docker host using docker-machine named '${DOCKER_MACHINE_NAME}' is ready to go!"
  return 0
}

function openim::build::is_gnu_sed() {
  [[ $(sed --version 2>&1) == *GNU* ]]
}

function openim::build::ensure_rsync() {
  if [[ -z "$(which rsync)" ]]; then
    openim::log::error "Can't find 'rsync' in PATH, please fix and retry."
    return 1
  fi
}

function openim::build::update_dockerfile() {
  if openim::build::is_gnu_sed; then
    sed_opts=(-i)
  else
    sed_opts=(-i '')
  fi
  sed "${sed_opts[@]}" "s/OPENIM_BUILD_IMAGE_CROSS_TAG/${OPENIM_BUILD_IMAGE_CROSS_TAG}/" "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
}

function  openim::build::set_proxy() {
  if [[ -n "${OPENIMRNETES_HTTPS_PROXY:-}" ]]; then
    echo "ENV https_proxy $OPENIMRNETES_HTTPS_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
  if [[ -n "${OPENIMRNETES_HTTP_PROXY:-}" ]]; then
    echo "ENV http_proxy $OPENIMRNETES_HTTP_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
  if [[ -n "${OPENIMRNETES_NO_PROXY:-}" ]]; then
    echo "ENV no_proxy $OPENIMRNETES_NO_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
}

function openim::build::ensure_docker_in_path() {
  if [[ -z "$(which docker)" ]]; then
    openim::log::error "Can't find 'docker' in PATH, please fix and retry."
    openim::log::error "See https://docs.docker.com/installation/#installation for installation instructions."
    return 1
  fi
}

function openim::build::ensure_tar() {
  if [[ -n "${TAR:-}" ]]; then
    return
  fi

  # Find gnu tar if it is available, bomb out if not.
  TAR=tar
  if which gtar &>/dev/null; then
      TAR=gtar
  else
      if which gnutar &>/dev/null; then
	  TAR=gnutar
      fi
  fi
  if ! "${TAR}" --version | grep -q GNU; then
    echo "  !!! Cannot find GNU tar. Build on Linux or install GNU tar"
    echo "      on Mac OS X (brew install gnu-tar)."
    return 1
  fi
}

function openim::build::has_docker() {
  which docker &> /dev/null
}

function openim::build::has_ip() {
  which ip &> /dev/null && ip -Version | grep 'iproute2' &> /dev/null
}

# Detect if a specific image exists
#
# $1 - image repo name
# $2 - image tag
function openim::build::docker_image_exists() {
  [[ -n $1 && -n $2 ]] || {
    openim::log::error "Internal error. Image not specified in docker_image_exists."
    exit 2
  }

  [[ $("${DOCKER[@]}" images -q "${1}:${2}") ]]
}

# Delete all images that match a tag prefix except for the "current" version
#
# $1: The image repo/name
# $2: The tag base. We consider any image that matches $2*
# $3: The current image not to delete if provided
function openim::build::docker_delete_old_images() {
  # In Docker 1.12, we can replace this with
  #    docker images "$1" --format "{{.Tag}}"
  for tag in $("${DOCKER[@]}" images "${1}" | tail -n +2 | awk '{print $2}') ; do
    if [[ "${tag}" != "${2}"* ]] ; then
      V=3 openim::log::status "Keeping image ${1}:${tag}"
      continue
    fi

    if [[ -z "${3:-}" || "${tag}" != "${3}" ]] ; then
      V=2 openim::log::status "Deleting image ${1}:${tag}"
      "${DOCKER[@]}" rmi "${1}:${tag}" >/dev/null
    else
      V=3 openim::log::status "Keeping image ${1}:${tag}"
    fi
  done
}

# Stop and delete all containers that match a pattern
#
# $1: The base container prefix
# $2: The current container to keep, if provided
function openim::build::docker_delete_old_containers() {
  # In Docker 1.12 we can replace this line with
  #   docker ps -a --format="{{.Names}}"
  for container in $("${DOCKER[@]}" ps -a | tail -n +2 | awk '{print $NF}') ; do
    if [[ "${container}" != "${1}"* ]] ; then
      V=3 openim::log::status "Keeping container ${container}"
      continue
    fi
    if [[ -z "${2:-}" || "${container}" != "${2}" ]] ; then
      V=2 openim::log::status "Deleting container ${container}"
      openim::build::destroy_container "${container}"
    else
      V=3 openim::log::status "Keeping container ${container}"
    fi
  done
}

# Takes $1 and computes a short has for it. Useful for unique tag generation
function openim::build::short_hash() {
  [[ $# -eq 1 ]] || {
    openim::log::error "Internal error.  No data based to short_hash."
    exit 2
  }

  local short_hash
  if which md5 >/dev/null 2>&1; then
    short_hash=$(md5 -q -s "$1")
  else
    short_hash=$(echo -n "$1" | md5sum)
  fi
  echo "${short_hash:0:10}"
}

# Pedantically kill, wait-on and remove a container. The -f -v options
# to rm don't actually seem to get the job done, so force kill the
# container, wait to ensure it's stopped, then try the remove. This is
# a workaround for bug https://github.com/docker/docker/issues/3968.
function openim::build::destroy_container() {
  "${DOCKER[@]}" kill "$1" >/dev/null 2>&1 || true
  if [[ $("${DOCKER[@]}" version --format '{{.Server.Version}}') = 17.06.0* ]]; then
    # Workaround https://github.com/moby/moby/issues/33948.
    # TODO: remove when 17.06.0 is not relevant anymore
    DOCKER_API_VERSION=v1.29 "${DOCKER[@]}" wait "$1" >/dev/null 2>&1 || true
  else
    "${DOCKER[@]}" wait "$1" >/dev/null 2>&1 || true
  fi
  "${DOCKER[@]}" rm -f -v "$1" >/dev/null 2>&1 || true
}

# ---------------------------------------------------------------------------
# Building


function openim::build::clean() {
  if openim::build::has_docker ; then
    openim::build::docker_delete_old_containers "${OPENIM_BUILD_CONTAINER_NAME_BASE}"
    openim::build::docker_delete_old_containers "${OPENIM_RSYNC_CONTAINER_NAME_BASE}"
    openim::build::docker_delete_old_containers "${OPENIM_DATA_CONTAINER_NAME_BASE}"
    openim::build::docker_delete_old_images "${OPENIM_BUILD_IMAGE_REPO}" "${OPENIM_BUILD_IMAGE_TAG_BASE}"

    V=2 openim::log::status "Cleaning all untagged docker images"
    "${DOCKER[@]}" rmi "$("${DOCKER[@]}" images -q --filter 'dangling=true')" 2> /dev/null || true
  fi

  if [[ -d "${LOCAL_OUTPUT_ROOT}" ]]; then
    openim::log::status "Removing _output directory"
    rm -rf "${LOCAL_OUTPUT_ROOT}"
  fi
}

# Set up the context directory for the openim-build image and build it.
function openim::build::build_image() {
  mkdir -p "${LOCAL_OUTPUT_BUILD_CONTEXT}"
  # Make sure the context directory owned by the right user for syncing sources to container.
  chown -R "${USER_ID}":"${GROUP_ID}" "${LOCAL_OUTPUT_BUILD_CONTEXT}"

  cp /etc/localtime "${LOCAL_OUTPUT_BUILD_CONTEXT}/"

  cp ""${OPENIM_ROOT}"/build/build-image/Dockerfile" "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  cp ""${OPENIM_ROOT}"/build/build-image/rsyncd.sh" "${LOCAL_OUTPUT_BUILD_CONTEXT}/"
  dd if=/dev/urandom bs=512 count=1 2>/dev/null | LC_ALL=C tr -dc 'A-Za-z0-9' | dd bs=32 count=1 2>/dev/null > "${LOCAL_OUTPUT_BUILD_CONTEXT}/rsyncd.password"
  chmod go= "${LOCAL_OUTPUT_BUILD_CONTEXT}/rsyncd.password"

  openim::build::update_dockerfile
  openim::build::set_proxy
  openim::build::docker_build "${OPENIM_BUILD_IMAGE}" "${LOCAL_OUTPUT_BUILD_CONTEXT}" 'false'

  # Clean up old versions of everything
  openim::build::docker_delete_old_containers "${OPENIM_BUILD_CONTAINER_NAME_BASE}" "${OPENIM_BUILD_CONTAINER_NAME}"
  openim::build::docker_delete_old_containers "${OPENIM_RSYNC_CONTAINER_NAME_BASE}" "${OPENIM_RSYNC_CONTAINER_NAME}"
  openim::build::docker_delete_old_containers "${OPENIM_DATA_CONTAINER_NAME_BASE}" "${OPENIM_DATA_CONTAINER_NAME}"
  openim::build::docker_delete_old_images "${OPENIM_BUILD_IMAGE_REPO}" "${OPENIM_BUILD_IMAGE_TAG_BASE}" "${OPENIM_BUILD_IMAGE_TAG}"

  openim::build::ensure_data_container
  openim::build::sync_to_container
}

# Build a docker image from a Dockerfile.
# $1 is the name of the image to build
# $2 is the location of the "context" directory, with the Dockerfile at the root.
# $3 is the value to set the --pull flag for docker build; true by default
function openim::build::docker_build() {
  local -r image=$1
  local -r context_dir=$2
  local -r pull="${3:-true}"
  local -ra build_cmd=("${DOCKER[@]}" build -t "${image}" "--pull=${pull}" "${context_dir}")

  openim::log::status "Building Docker image ${image}"
  local docker_output
  docker_output=$("${build_cmd[@]}" 2>&1) || {
    cat <<EOF >&2
+++ Docker build command failed for ${image}

${docker_output}

To retry manually, run:

${build_cmd[*]}

EOF
    return 1
  }
}

function openim::build::ensure_data_container() {
  # If the data container exists AND exited successfully, we can use it.
  # Otherwise nuke it and start over.
  local ret=0
  local code=0

  code=$(docker inspect \
      -f '{{.State.ExitCode}}' \
      "${OPENIM_DATA_CONTAINER_NAME}" 2>/dev/null) || ret=$?
  if [[ "${ret}" == 0 && "${code}" != 0 ]]; then
    openim::build::destroy_container "${OPENIM_DATA_CONTAINER_NAME}"
    ret=1
  fi
  if [[ "${ret}" != 0 ]]; then
    openim::log::status "Creating data container ${OPENIM_DATA_CONTAINER_NAME}"
    # We have to ensure the directory exists, or else the docker run will
    # create it as root.
    mkdir -p "${LOCAL_OUTPUT_GOPATH}"
    # We want this to run as root to be able to chown, so non-root users can
    # later use the result as a data container.  This run both creates the data
    # container and chowns the GOPATH.
    #
    # The data container creates volumes for all of the directories that store
    # intermediates for the Go build. This enables incremental builds across
    # Docker sessions. The *_cgo paths are re-compiled versions of the go std
    # libraries for true static building.
    local -ra docker_cmd=(
      "${DOCKER[@]}" run
      --volume "${REMOTE_ROOT}"   # white-out the whole output dir
      --volume /usr/local/go/pkg/linux_386_cgo
      --volume /usr/local/go/pkg/linux_amd64_cgo
      --volume /usr/local/go/pkg/linux_arm_cgo
      --volume /usr/local/go/pkg/linux_arm64_cgo
      --volume /usr/local/go/pkg/linux_ppc64le_cgo
      --volume /usr/local/go/pkg/darwin_amd64_cgo
      --volume /usr/local/go/pkg/darwin_386_cgo
      --volume /usr/local/go/pkg/windows_amd64_cgo
      --volume /usr/local/go/pkg/windows_386_cgo
      --name "${OPENIM_DATA_CONTAINER_NAME}"
      --hostname "${HOSTNAME}"
      "${OPENIM_BUILD_IMAGE}"
      chown -R "${USER_ID}":"${GROUP_ID}"
        "${REMOTE_ROOT}"
        /usr/local/go/pkg/
    )
    "${docker_cmd[@]}"
  fi
}

# Build all openim commands.
function openim::build::build_command() {
  openim::log::status "Running build command..."
  make -C "${OPENIM_ROOT}" multiarch
}
