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

if ! command -v pv &> /dev/null
then
    echo "pv not found, installing..."
    if [ -e /etc/debian_version ]; then
        sudo apt-get update
        sudo apt-get install -y pv
    elif [ -e /etc/redhat-release ]; then
        sudo yum install -y pv
    else
        echo "Unsupported OS, please install pv manually."
        exit 1
    fi
fi

readonly t_reset=$(tput sgr0)
readonly  green=$(tput bold; tput setaf 2)
readonly yellow=$(tput bold; tput setaf 3)
readonly   blue=$(tput bold; tput setaf 6)
readonly timeout=$(if [ "$(uname)" == "Darwin" ]; then echo "1"; else echo "0.1"; fi)
readonly ipv6regex='(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))'

clear
. $(dirname ${BASH_SOURCE})/lib/util.sh

trap 'openim::util::onCtrlC' INT

function openim::util::onCtrlC() {
    echo -e "\n${t_reset}Ctrl+C Press it. It's exiting openim make init..."
    exit 0
}

openim::util::desc "========> Welcome to the OpenIM Demo"
openim::util::desc "========> We'll help you get started with OpenIM quickly"
openim::util::desc "========> Press Enter to continue...."
openim::util::run "make advertise"
clear

openim::util::desc "========> Initialize the project and generate configuration files"
openim::util::run "make init"
clear

# openim::util::desc "========> You can look git diff"
# openim::util::run "git diff"
# clear

openim::util::desc "You can learn a lot about automation using make help"
openim::util::run "make help"
clear

openim::util::desc "You can learn a lot about automation using make help-all"
openim::util::run "make help-all"
clear

openim::util::desc "First, let's verify and install some necessary tools"
openim::util::run "make tools"
clear

openim::util::desc "========> Start the basic openim docker components"
openim::util::desc "========> You can use docker-compose ps to check the status of the container"
openim::util::run "docker compose up -d"
clear

openim::util::desc "========> Use make init-githooks Initialize git hooks "
openim::util::run "make init-githooks"
clear

openim::util::desc "The specification is pretty high, you need to be bound on your branch name, as well as commit messages"
openim::util::run "git commit -a -s -m 'feta: commit demo against specification'"
openim::util::run "# git commit -a -s -m 'feat: commit demo against specification' --amend"
clear

openim::util::desc "How did we teach you how to build OpenIM"
openim::util::desc "A full build startup check"
openim::util::run "# make all"
openim::util::desc "Build one OpenIM binary"
openim::util::desc "BINS: openim-api openim-cmdutils openim-crontask openim-msggateway openim-msgtransfer openim-push openim-rpc changelog infra ncpu yamlfmt"
openim::util::run "make build BINS=openim-api"
openim::util::run "make build"

openim::util::desc "Build binaries for all platforms"
openim::util::run "make multiarch -j BINS=openim-crontask PLATFORMS='linux_arm64 linux_amd64' "

openim::util::desc "If you wish to use dlv for debugging, either binary or process"
openim::util::desc "You need to enable debug mode"
openim::util::run "make build BINS=openim-cmdutils DEBUG=1"
clear

openim::util::desc "Next, let's learn how to start the OpenIM service. For starting, we have two ways"
openim::util::desc "The first is Background startup"
openim::util::run "make start"
openim::util::desc "The second way is through the Linux system way"
openim::util::run "./scripts/install/install.sh --help"
clear

openim::util::desc "Next, let's learn how to check the OpenIM service. For checking, we have two ways"
openim::util::run "make check"
clear

openim::util::desc "Next, let's learn how to stop the OpenIM service. For stopping, we have two ways"
openim::util::run "make stop"
clear

openim::util::desc "Run tidy to format and fix imports"
openim::util::run "make tidy"
clear

openim::util::desc "Vendor go.mod dependencies"
openim::util::run "# make vendor"
clear

openim::util::desc "Run unit tests"
openim::util::run "# make test"
clear

openim::util::desc "Run unit tests and get test coverage"
openim::util::run "# make cover"
clear

openim::util::desc "Check for updates to go.mod dependencies"
openim::util::run "# make updates"
clear

openim::util::desc "You can learn a lot about automation using make clean, remove all files that are created by building"
openim::util::run "make clean"
clear

openim::util::desc "Generate all necessary files"
openim::util::run "make gen"
clear

openim::util::desc "Verify the license headers for all files"
openim::util::run "make verify-copyright"
clear

openim::util::desc "Add copyright"
openim::util::run "make add-copyright"
clear

exit 0
