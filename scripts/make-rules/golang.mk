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
# Build management helpers.  These functions help to set, save and load the
#

GO := go
GO_SUPPORTED_VERSIONS ?= 1.19|1.20|1.21|1.22

GO_LDFLAGS += -X $(VERSION_PACKAGE).gitVersion=$(GIT_TAG) \
	-X $(VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).gitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).buildDate=$(BUILD_DATE) \
	-s -w		# -s -w deletes debugging information and symbol tables
ifeq ($(DEBUG), 1)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	GO_LDFLAGS=
endif

GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk, ->/Makefile)
endif

GOPATH ?= $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

# COMMANDS is Specify all files under ${ROOT_DIR}/cmd/ and ${ROOT_DIR}/tools/ except those ending in.md
COMMANDS ?= $(filter-out %.md, $(wildcard ${ROOT_DIR}/cmd/* ${ROOT_DIR}/tools/* ${ROOT_DIR}/tools/data-conversion/chat/cmd/* ${ROOT_DIR}/tools/data-conversion/openim/cmd/* ${ROOT_DIR}/cmd/openim-rpc/*))
ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif

# BINS is the name of each file in ${COMMANDS}, excluding the directory path
# If there are no files in ${COMMANDS}, or if all files end in.md, ${BINS} will be empty
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

ifeq ($(OS),Windows_NT)
  NULL :=
  SPACE := $(NULL) $(NULL)
  ROOT_DIR := $(subst $(SPACE),\$(SPACE),$(shell cd))
else
  ROOT_DIR := $(shell pwd)
endif

ifeq ($(strip $(COMMANDS)),)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq ($(strip $(BINS)),)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

# TODO: EXCLUDE_TESTS variable, which contains the name of the package to be excluded from the test
EXCLUDE_TESTS=github.com/openimsdk/open-im-server/test github.com/openimsdk/open-im-server/v3/pkg/log github.com/openimsdk/open-im-server/db github.com/openimsdk/open-im-server/scripts github.com/openimsdk/open-im-server/config

# ==============================================================================
# ❯ tree -L 1 cmd
# cmd
# ├── openim-sdk-core/ - main.go
# ├── openim-api	
# ├── openim_cms_api
# ├── openim-crontask
# ├── openim_demo
# ├── openim-rpc-msg_gateway
# ├── openim-msgtransfer
# ├── openim-push
# ├── rpc/openim_admin_cms/ - main.go
# └── test/ - main.go
# COMMAND=openim
# PLATFORM=linux_amd64
# OS=linux
# ARCH=amd64
# BINS=openim-api openim_cms_api openim-crontask openim_demo openim-rpc-msg_gateway openim-msgtransfer openim-push 
# BIN_DIR=/root/workspaces/OpenIM/_output/bin
# ==============================================================================

## go.build: Build binaries
.PHONY: go.build
go.build: go.build.verify $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))
	@echo "===========> Building binary $(BINS) $(VERSION) for $(PLATFORM)"

## go.start: Start openim
.PHONY: go.start
go.start:
	@echo "===========> Starting openim"
	@$(ROOT_DIR)/scripts/start-all.sh

## go.stop: Stop openim
.PHONY: go.stop
go.stop:
	@echo "===========> Stopping openim"
	@$(ROOT_DIR)/scripts/stop-all.sh

## go.check: Check openim
.PHONY: go.check
go.check:
	@echo "===========> Checking openim"
	@$(ROOT_DIR)/scripts/check-all.sh

## go.check-component: Check openim component
.PHONY: go.check-component
go.check-component:
	@echo "===========> Checking openim component"
	@$(ROOT_DIR)/scripts/install/openim-tools.sh openim::tools::pre-start

## go.versionchecker: Design, detect some environment variables and versions
go.versionchecker:
	@$(ROOT_DIR)/scripts/install/openim-tools.sh openim::tools::post-start

## go.build.verify: Verify that a suitable version of Go exists
.PHONY: go.build.verify
go.build.verify:
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif

## go.build.%: Build binaries for a specific platform
# CGO_ENABLED=0 https://wiki.musl-libc.org/functional-differences-from-glibc.html
.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "=====> COMMAND=$(COMMAND)"
	@echo "=====> PLATFORM=$(PLATFORM)"
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS)_$(ARCH)"
	@mkdir -p $(BIN_DIR)/platforms/$(OS)/$(ARCH)
	@if [ "$(COMMAND)" == "openim-sdk-core" ]; then \
		echo "===========> DEBUG: OpenIM-SDK-Core It is no longer supported for openim-server $(COMMAND)"; \
	elif [ -d $(ROOT_DIR)/cmd/openim-rpc/$(COMMAND) ]; then \
		CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o \
		$(BIN_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_DIR)/cmd/openim-rpc/$(COMMAND)/main.go; \
	else \
		if [ -f $(ROOT_DIR)/cmd/$(COMMAND)/main.go ]; then \
			CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o \
			$(BIN_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_DIR)/cmd/$(COMMAND)/main.go; \
		elif [ -f $(ROOT_DIR)/tools/$(COMMAND)/$(COMMAND).go ]; then \
			CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o \
			$(BIN_TOOLS_DIR)/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_DIR)/tools/$(COMMAND)/$(COMMAND).go; \
			chmod +x $(BIN_TOOLS_DIR)/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT); \
		elif [ -f $(ROOT_DIR)/tools/data-conversion/openim/cmd/$(COMMAND)/$(COMMAND).go ]; then \
			CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o \
			$(BIN_TOOLS_DIR)/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_DIR)/tools/data-conversion/openim/cmd/$(COMMAND)/$(COMMAND).go; \
			chmod +x $(BIN_TOOLS_DIR)/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT); \
		elif [ -f $(ROOT_DIR)/tools/data-conversion/chat/cmd/$(COMMAND)/$(COMMAND).go ]; then \
			CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o \
			$(BIN_TOOLS_DIR)/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_DIR)/tools/data-conversion/chat/cmd/$(COMMAND)/$(COMMAND).go; \
			chmod +x $(BIN_TOOLS_DIR)/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT); \
		fi \
	fi

## go.install: Install deployment openim
.PHONY: go.install
go.install:
	@echo "===========> Installing deployment openim"
	@$(ROOT_DIR)/scripts/install-im-server.sh

## go.multiarch: Build multi-arch binaries
.PHONY: go.build.multiarch
go.build.multiarch: go.build.verify $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))

## go.lint: Run golangci to lint source codes
.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@$(TOOLS_DIR)/golangci-lint run --color always -c $(ROOT_DIR)/.golangci.yml $(ROOT_DIR)/... 

## go.test: Run unit test
.PHONY: go.test
go.test:
	@$(GO) test ./...

## go.test.api: Run api test
.PHONY: go.test.api
go.test.api:
	@echo "===========> Run api test"
	@$(ROOT_DIR)/scripts/install/test.sh openim::test::test

## go.test.e2e: Run e2e test
.PHONY: go.test.e2e
go.test.e2e: tools.verify.ginkgo
	@echo "===========> Run e2e test"
	@$(TOOLS_DIR)/ginkgo -v $(ROOT_DIR)/test/e2e

## go.demo: Run demo
.PHONY: go.demo
go.demo:
	@echo "===========> Run demo"
	@$(ROOT_DIR)/scripts/demo.sh

## go.test.junit-report: Run unit test
.PHONY: go.test.junit-report
go.test.junit-report: tools.verify.go-junit-report
	@touch $(TMP_DIR)/coverage.out
	@echo "===========> Run unit test > $(TMP_DIR)/report.xml"
# 	@$(GO) test -v -coverprofile=$(TMP_DIR)/coverage.out 2>&1 $(GO_BUILD_FLAGS) ./... | $(TOOLS_DIR)/go-junit-report -set-exit-code > $(TMP_DIR)/report.xml
	@$(GO) test -v -coverprofile=$(TMP_DIR)/coverage.out 2>&1 ./... | $(TOOLS_DIR)/go-junit-report -set-exit-code > $(TMP_DIR)/report.xml
	@sed -i '/mock_.*.go/d' $(TMP_DIR)/coverage.out
	@echo "===========> Test coverage of Go code is reported to $(TMP_DIR)/coverage.html by generating HTML"
	@$(GO) tool cover -html=$(TMP_DIR)/coverage.out -o $(TMP_DIR)/coverage.html

## go.test.cover: Run unit test with coverage
.PHONY: go.test.cover
go.test.cover: go.test.junit-report
	@$(GO) tool cover -func=$(TMP_DIR)/coverage.out | \
		awk -v target=$(COVERAGE) -f $(ROOT_DIR)/scripts/coverage.awk

## go.format: Run unit test and format codes
.PHONY: go.format
go.format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formatting codes"
	@$(FIND) -type f -name '*.go' -not -name '*pb*' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' -not -name '*pb*' | $(XARGS) $(TOOLS_DIR)/goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' -not -name '*pb*' | $(XARGS) $(TOOLS_DIR)/golines -w --max-len=200 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt

## go.imports: task to automatically handle import packages in Go files using goimports tool
.PHONY: go.imports
go.imports: tools.verify.goimports
	@$(TOOLS_DIR)/goimports -l -w $(SRC)

## go.verify: execute all verity scripts.
.PHONY: go.verify
go.verify:
	@echo "Starting verification..."
	@scripts_list=$$(find $(ROOT_DIR)/scripts -type f -name 'verify-*' | sort); \
	for script in $$scripts_list; do \
		echo "Executing $$script..."; \
		$$script || exit 1; \
		echo "$$script completed successfully"; \
	done
	@echo "All verification scripts executed successfully."

## go.updates: Check for updates to go.mod dependencies
.PHONY: go.updates
go.updates: tools.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct

## go.clean: Clean all builds directories and files
.PHONY: go.clean
go.clean:
	@echo "===========> Cleaning all builds tmp, bin, logs directories and files"
	@-rm -vrf $(TMP_DIR) $(BIN_DIR) $(BIN_TOOLS_DIR) $(LOGS_DIR)
	@echo "===========> End clean..."

## go.help: Show go tools help
.PHONY: go.help
go.help: scripts/make-rules/golang.mk
	$(call smallhelp)
