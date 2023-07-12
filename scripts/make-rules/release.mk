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

# ==============================================================================
# Makefile helper functions for release
# Versions are used after merging
#

## release.run: release the project
.PHONY: release.run
release.run: release.verify release.ensure-tag
	@scripts/release.sh

## release.verify: Check if a tool is installed and install it
.PHONY: release.verify
release.verify: tools.verify.git-chglog tools.verify.github-release tools.verify.coscmd tools.verify.coscli

## release.tag: release the project
.PHONY: release.tag
release.tag: tools.verify.gsemver release.ensure-tag
	@git push origin `git describe --tags --abbrev=0`

## release.ensure-tag: ensure tag
.PHONY: release.ensure-tag
release.ensure-tag: tools.verify.gsemver
	@scripts/ensure_tag.sh

## release.help: Display help information about the release package
.PHONY: release.help
release.help: scripts/make-rules/release.mk
	$(call smallhelp)
