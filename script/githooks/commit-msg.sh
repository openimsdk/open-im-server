
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
# ==============================================================================
#
# Store this file as .git/hooks/commit-msg in your repository in order to
# enforce checking for proper commit message format before actual commits.
# You may need to make the script executable by 'chmod +x .git/hooks/commit-msg'.

# commit-msg use go-gitlint tool, install go-gitlint via `go get github.com/llorllale/go-gitlint/cmd/go-gitlint`
go-gitlint --msg-file="$1"
