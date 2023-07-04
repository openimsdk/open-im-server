#!/usr/bin/env bash

# Build a IAM release.  This will build the binaries, create the Docker
# images and other build artifacts.

set -o errexit
set -o nounset
set -o pipefail

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${IAM_ROOT}/scripts/common.sh"
source "${IAM_ROOT}/scripts/lib/release.sh"

IAM_RELEASE_RUN_TESTS=${IAM_RELEASE_RUN_TESTS-y}

iam::golang::setup_env
iam::build::verify_prereqs
iam::release::verify_prereqs
#iam::build::build_image
iam::build::build_command
iam::release::package_tarballs
iam::release::updload_tarballs
git push origin ${VERSION}
#iam::release::github_release
#iam::release::generate_changelog
