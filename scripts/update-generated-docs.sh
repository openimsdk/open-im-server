#!/usr/bin/env bash

# This file is not intended to be run automatically. It is meant to be run
# immediately before exporting docs. We do not want to check these documents in
# by default.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/hack/lib/init.sh"

openim::golang::setup_env

BINS=(
	cmd/gendocs
	cmd/genopenimdocs
	cmd/genman
	cmd/genyaml
)
make -C "${OPENIM_ROOT}" WHAT="${BINS[*]}"

openim::util::ensure-temp-dir

openim::util::gen-docs "${OPENIM_TEMP}"

# remove all of the old docs
openim::util::remove-gen-docs

# Copy fresh docs into the repo.
# the shopt is so that we get docs/.generated_docs from the glob.
shopt -s dotglob
cp -af "${OPENIM_TEMP}"/* "${OPENIM_ROOT}"
shopt -u dotglob
