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
#

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source ${OPENIM_ROOT}/scripts/install/common.sh

source ${OPENIM_ROOT}/scripts/install/dependency.sh
source ${OPENIM_ROOT}/scripts/install/openim-msggateway.sh
source ${OPENIM_ROOT}/scripts/install/openim-msgtransfer.sh
source ${OPENIM_ROOT}/scripts/install/openim-push.sh
source ${OPENIM_ROOT}/scripts/install/openim-rpc.sh
source ${OPENIM_ROOT}/scripts/install/openim-crontask.sh
source ${OPENIM_ROOT}/scripts/install/openim-api.sh
source ${OPENIM_ROOT}/scripts/install/test.sh
source ${OPENIM_ROOT}/scripts/install/man.sh

# Detailed help function
function show_help() {
    echo "OpenIM Installer"
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  -i, --install        Install all OpenIM components."
    echo "  -u, --uninstall      Remove all OpenIM components."
    echo "  -s, --status         Check the current status of OpenIM components."
    echo "  -h, --help           Show this help menu."
    echo ""
    echo "Example: "
    echo "  $0 -i     Will install all OpenIM components."
    echo "  $0 --install  Same as above."
}

function openim::install::install_openim()
{
    openim::log::info "check openim dependency"
    openim::util::check_ports ${OPENIM_DEPENDENCY_PORT_LISTARIES[@]}

    openim::msggateway::install || return 1
    openim::msgtransfer::install || return 1
    openim::push::install || return 1
    openim::rpc::install || return 1
    openim::crontask::install || return 1
    openim::api::install || return 1

    openim::log::success "openim install success"
}

function openim::uninstall::uninstall_openim()
{
    openim::log::info "uninstall openim"

    openim::msggateway::uninstall || return 1
    openim::msgtransfer::uninstall || return 1
    openim::push::uninstall || return 1
    openim::rpc::uninstall || return 1
    openim::crontask::uninstall || return 1
    openim::api::uninstall || return 1

    openim::log::success "openim uninstall success"
}

function openim::install::status()
{
    openim::log::info "check openim status"

    openim::msggateway::status || return 1
    openim::msgtransfer::status || return 1
    openim::push::status || return 1
    openim::rpc::status || return 1
    openim::crontask::status || return 1
    openim::api::status || return 1

    openim::log::success "openim status success"
}

# Argument parsing to call functions based on user input
while (( "$#" )); do
    case "$1" in
        -i|--install)
            openim::install::install_openim
            shift
            ;;
        -u|--uninstall)
            openim::uninstall::uninstall_openim
            shift
            ;;
        -s|--status)
            openim::install::status
            shift
            ;;
        -h|--help|*)
            show_help
            exit 0
            ;;
    esac
done

# If no arguments are provided, show help
if [[ $# -eq 0 ]]; then
    show_help
    exit 0
fi