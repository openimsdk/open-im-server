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


# Build a OpenIM release.  This will build the binaries, create the Docker
# images and other build artifacts.

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/common.sh"
source "${OPENIM_ROOT}/scripts/lib/release.sh"

OPENIM_RELEASE_RUN_TESTS=${OPENIM_RELEASE_RUN_TESTS-y}

openim::golang::setup_env
openim::build::verify_prereqs
openim::release::verify_prereqs
#openim::build::build_image
openim::build::build_command
openim::release::package_tarballs
openim::release::updload_tarballs
git push origin ${VERSION}
#openim::release::github_release
#openim::release::generate_changelog
