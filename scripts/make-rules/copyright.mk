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
# wget https://github.com/google/addlicense/releases/download/v1.0.0/addlicense_1.0.0_Linux_x86_64.tar.gz
# Makefile helper functions for copyright
#

LICENSE_TEMPLATE ?= $(ROOT_DIR)/scripts/template/LICENSE_TEMPLATES

## copyright.verify: Validate boilerplate headers for assign files
.PHONY: copyright.verify
copyright.verify: tools.verify.addlicense
	@echo "===========> Validate boilerplate headers for assign files starting in the $(ROOT_DIR) directory"
	@$(TOOLS_DIR)/addlicense -v -check -ignore **/test/** -ignore **pb**  -f $(LICENSE_TEMPLATE) $(CODE_DIRS)
	@echo "===========> End of boilerplate headers check..."

## copyright.add: Add the boilerplate headers for all files
.PHONY: copyright.add
copyright.add: tools.verify.addlicense
	@echo "===========> Adding $(LICENSE_TEMPLATE) the boilerplate headers for all files"
	@$(TOOLS_DIR)/addlicense -y $(shell date +"%Y") -ignore **pb** -v -c "OpenIM." -f $(LICENSE_TEMPLATE) $(CODE_DIRS)
	@echo "===========> End the copyright is added..."

# Addlicense Flags:
#   -c string
#         copyright holder (default "Google LLC")
#   -check
#         check only mode: verify presence of license headers and exit with non-zero code if missing
#   -f string
#         license file
#   -ignore value
#         file patterns to ignore, for example: -ignore **/*.go -ignore vendor/**
#   -l string
#         license type: apache, bsd, mit, mpl (default "apache")
#   -s    Include SPDX identifier in license header. Set -s=only to only include SPDX identifier.
#   -skip value
#         [deprecated: see -ignore] file extensions to skip, for example: -skip rb -skip go
#   -v    verbose mode: print the name of the files that are modified or were skipped
#   -y string
#         copyright year(s) (default "2023")

## copyright.advertise: Advertise the license of the project
.PHONY: copyright.advertise
copyright.advertise:
	@chmod +x $(ROOT_DIR)/scripts/advertise.sh
	@$(ROOT_DIR)/scripts/advertise.sh

## copyright.help: Show copyright help
.PHONY: copyright.help
copyright.help: scripts/make-rules/copyright.mk
	$(call smallhelp)