#!/usr/bin/env bash

# Copyright © 2023 OpenIMSDK.
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

LC_ALL=C

local_branch="$(git rev-parse --abbrev-ref HEAD)"

valid_branch_regex="^(master|develop)$|(feature|release|hotfix)\/[a-z0-9._-]+$|^HEAD$"

message="There is something wrong with your branch name. Branch names in this project must adhere to this contract: $valid_branch_regex. 
Your commit will be rejected. You should rename your branch to a valid name and try again."
message2="If you're not familiar with the contribution process, please read our contributor documentation again
-> https://github.com/OpenIMSDK/Open-IM-Server/blob/main/CONTRIBUTING.md"


if [[ ! $local_branch =~ $valid_branch_regex ]]
then
    echo "$message"
    echo "$message2"
    exit 1
fi

exit 0
