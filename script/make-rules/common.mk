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
# Makefile helper functions for common tasks
#




# ==============================================================================
# Makefile helper functions for common tasks

# Help information for the makefile package
define makehelp
	@printf "\n\033[1mUsage: make <TARGETS> <OPTIONS> ...\033[0m\n\n\\033[1mTargets:\\033[0m\n\n"
	@sed -n 's/^##//p' $< | awk -F':' '{printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' | sed -e 's/^/ /'
	@printf "\n\033[1m$$USAGE_OPTIONS\033[0m\n"
endef

# Here are some examples of builds
define MAKEFILE_EXAMPLE
# make build BINS=imctl                                          Only a single imctl binary is built.
# make -j (nproc) all                                            Run tidy gen add-copyright format lint cover build concurrently.
# make gen                                                       Generate all necessary files.
# make linux.arm64                                               imctl is compiled on arm64 platform.
# make verify-copyright                                          Verify the license headers for all files.
# make install-deepcopy-gen                                      Install deepcopy-gen tools if the license is missing.
# make build BINS=imctl V=1 DEBUG=1                             Build debug binaries for only imctl.
# make multiarch PLATFORMS="linux_arm64 linux_amd64" V=1   Build binaries for both platforms.
endef
export MAKEFILE_EXAMPLE

# Define all help functions	@printf "\n\033[1mCurrent imctl version information: $(shell imctl version):\033[0m\n\n"
define makeallhelp
	@printf "\n\033[1mMake example:\033[0m\n\n"
	$(call MAKEFILE_EXAMPLE)
	@printf "\n\033[1mAriables:\033[0m\n\n"
	@echo "  DEBUG: $(DEBUG)"
	@echo "  BINS: $(BINS)"
	@echo "  PLATFORMS: $(PLATFORMS)"
	@echo "  V: $(V)"
endef

# Help information for other makefile packages
CUT_OFF?="---------------------------------------------------------------------------------"
HELP_NAME:=$(shell basename $(MAKEFILE_LIST))
define smallhelp
	@sed -n 's/^##//p' $< | awk -F':' '{printf "\033[36m%-35s\033[0m %s\n", $$1, $$2}' | sed -e 's/^/ /'
	@echo $(CUT_OFF)
endef