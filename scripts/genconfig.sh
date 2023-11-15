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

# Function of this script: Generate the OPENIM component YAML configuration file according to the scripts/environment.sh configuration.
# eg：./scripts/genconfig.sh scripts/install/environment.sh scripts/template/config.yaml
# Read: https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/init-config.md

env_file="$1"
template_file="$2"

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${OPENIM_ROOT}/scripts/lib/init.sh"

if [ $# -ne 2 ];then
    openim::log::error "Usage: scripts/genconfig.sh scripts/environment.sh configs/openim-api.yaml"
    exit 1
fi

openim::util::require-dig
result=$?
if [ $result -ne 0 ]; then
    openim::log::info "Please install 'dig' to use this feature."
    openim::log::info "Installation instructions:"
    openim::log::info "  For Ubuntu/Debian: sudo apt-get install dnsutils"
    openim::log::info "  For CentOS/RedHat: sudo yum install bind-utils"
    openim::log::info "  For macOS: 'dig' should be preinstalled. If missing, try: brew install bind"
    openim::log::info "  For Windows: Install BIND9 tools from https://www.isc.org/download/"
    openim::log::error_exit "Error: 'dig' command is required but not installed."
fi

source "${env_file}"

declare -A envs

set +u
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file})
do
    if [ -z "$(eval echo \$${env})" ];then
        openim::log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ "${missing}" ];then
    openim::log::error 'You may run `source scripts/environment.sh` to set these environment'
    exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"
