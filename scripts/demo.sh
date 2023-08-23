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

clear
. $(dirname ${BASH_SOURCE})/lib/util.sh

trap 'openim::util::onCtrlC' INT

openim::util::desc "========> Welcome to the OpenIM Demo"
openim::util::desc "========> We'll help you get started with OpenIM quickly"
openim::util::desc "========> Press Enter to continue...."
openim::util::run "make advertise"
clear

openim::util::desc "========> Initialize the project and generate configuration files"
openim::util::run "make init"

openim::util::desc "========> You can look git diff"
openim::util::run "git diff"
clear

openim::util::desc "You can learn a lot about automation using make help"
openim::util::run "make help"
clear

openim::util::desc "You can learn a lot about automation using make help-all"
openim::util::run "make help-all"
clear

openim::util::desc "How did we teach you how to build OpenIM"
openim::util::desc "A full build startup check"
openim::util::run "make all"

openim::util::desc "Build one OpenIM binary"
openim::util::desc "BINS: openim-api openim-cmdutils openim-crontask openim-msggateway openim-msgtransfer openim-push openim-rpc changelog infra ncpu yamlfmt"
openim::util::run "make build BINS=openim-api"

openim::util::desc "Build binaries for all platforms"
openim::util::run "make multiarch -j BINS=openim-api PLATFORMS='linux_arm64 linux_amd64' "

openim::util::desc "If you wish to use dlv for debugging, either binary or process"
openim::util::desc "You need to enable debug mode"
openim::util::run "make build BINS=openim-api DEBUG=1"
clear

openim::util::desc "Run tidy to format and fix imports"
openim::util::run "make tidy"
clear

openim::util::desc "Vendor go.mod dependencies"
openim::util::run "make vendor"
clear

openim::util::desc "Run unit tests"
openim::util::run "make test"
clear

openim::util::desc "Run unit tests and get test coverage"
openim::util::run "make cover"
clear

openim::util::desc "Check for updates to go.mod dependencies"
openim::util::run "make updates"
clear

openim::util::desc "Clean all generated files"
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
