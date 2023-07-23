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
# Makefile helper functions for tools(https://github.com/avelino/awesome-go) -> DIR: {TOOT_DIR}/tools | (go >= 1.19)
# Why download to the tools directory, thinking we might often switch Go versions using gvm.
#

# openim build use BUILD_TOOLS
BUILD_TOOLS ?= golangci-lint goimports addlicense deepcopy-gen conversion-gen ginkgo go-junit-report go-gitlint
# Code analysis tools
ANALYSIS_TOOLS = golangci-lint goimports golines go-callvis kube-score
# Code generation tools
GENERATION_TOOLS = deepcopy-gen conversion-gen protoc-gen-go cfssl rts codegen
# Testing tools
TEST_TOOLS = ginkgo go-junit-report gotests
# tenxun cos tools
COS_TOOLS = coscli coscmd
# Version control tools
VERSION_CONTROL_TOOLS = addlicense go-gitlint git-chglog github-release gsemver
# Utility tools
UTILITY_TOOLS = go-mod-outdated mockgen gothanks richgo kubeconform
# All tools
ALL_TOOLS ?= $(ANALYSIS_TOOLS) $(GENERATION_TOOLS) $(TEST_TOOLS) $(VERSION_CONTROL_TOOLS) $(UTILITY_TOOLS) $(COS_TOOLS)

## tools.install: Install a must tools
.PHONY: tools.install
tools.install: $(addprefix tools.verify., $(BUILD_TOOLS))
 
## tools.install-all: Install all tools
.PHONY: tools.install-all
tools.install-all: $(addprefix tools.install-all., $(ALL_TOOLS))

## tools.install.%: Install a single tool in $GOBIN/
.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $,The default installation path is $(GOBIN)/$*"
	@$(MAKE) install.$*

## tools.install-all.%: Parallelism install a single tool in ./tools/*
.PHONY: tools.install-all.%
tools.install-all.%:
	@echo "===========> Installing $,The default installation path is $(TOOLS_DIR)/$*"
	@$(MAKE) -j $(nproc) install.$*

## tools.verify.%: Check if a tool is installed and install it
.PHONY: tools.verify.%
tools.verify.%:
	@echo "===========> Verifying $* is installed"
	@if [ ! -f $(TOOLS_DIR)/$* ]; then GOBIN=$(TOOLS_DIR) $(MAKE) tools.install.$*; fi
	@echo "===========> $* is install in $(TOOLS_DIR)/$*"

## install.golangci-lint: Install golangci-lint
.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## install.goimports: Install goimports, used to format go source files
.PHONY: install.goimports
install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

## install.addlicense: Install addlicense, used to add license header to source files
.PHONY: install.addlicense
install.addlicense:
	@$(GO) install github.com/google/addlicense@latest

## install.deepcopy-gen: Install deepcopy-gen, used to generate deep copy functions
.PHONY: install.deepcopy-gen
install.deepcopy-gen:
	@$(GO) install k8s.io/code-generator/cmd/deepcopy-gen@latest

## install.conversion-gen: Install conversion-gen, used to generate conversion functions
.PHONY: install.conversion-gen
install.conversion-gen:
	@$(GO) install k8s.io/code-generator/cmd/conversion-gen@latest

## install.ginkgo: Install ginkgo to run a single test or set of tests
.PHONY: install.ginkgo
install.ginkgo:
	@$(GO) install github.com/onsi/ginkgo/ginkgo@v1.16.2

## Install go-gitlint: Install go-gitlint, used to check git commit message
.PHONY: install.go-gitlint
install.go-gitlint:
	@$(GO) install github.com/marmotedu/go-gitlint/cmd/go-gitlint@latest

## install.go-junit-report: Install go-junit-report, used to convert go test output to junit xml
.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) install github.com/jstemmer/go-junit-report@latest

## install.gotests: Install gotests, used to generate go tests
.PHONY: install.swagger
install.swagger:
	@$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest

# ==============================================================================
# Tools that might be used include go gvm
#

## install.kube-score: Install kube-score, used to check kubernetes yaml files
.PHONY: install.kube-score
install.kube-score:
	@$(GO) install github.com/zegl/kube-score/cmd/kube-score@latest

## install.kubeconform: Install kubeconform, used to check kubernetes yaml files
.PHONY: install.kubeconform
install.kubeconform:
	@$(GO) install github.com/yannh/kubeconform/cmd/kubeconform@latest

## install.gsemver: Install gsemver, used to generate semver
.PHONY: install.gsemver
install.gsemver:
	@$(GO) install github.com/arnaud-deprez/gsemver@latest

## install.git-chglog: Install git-chglog, used to generate changelog
.PHONY: install.git-chglog
install.git-chglog:
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

## install.ko: Install ko, used to build go program into container images
.PHONY: install.ko
install.ko:
	@$(GO) install github.com/google/ko@latest

## install.github-release: Install github-release, used to create github release
.PHONY: install.github-release
install.github-release:
	@$(GO) install github.com/github-release/github-release@latest

## install.coscli: Install coscli, used to upload files to cos
# example: ./coscli  cp/sync -r  /home/off-line/docker-off-line/ cos://openim-1306374445/openim/image/amd/off-line/off-line/ -e cos.ap-guangzhou.myqcloud.com
# https://cloud.tencent.com/document/product/436/71763
# amd64
.PHONY: install.coscli
install.coscli:
	@wget -q https://ghproxy.com/https://github.com/tencentyun/coscli/releases/download/v0.13.0-beta/coscli-linux -O ${TOOLS_DIR}/coscli
	@chmod +x ${TOOLS_DIR}/coscli

## install.coscmd: Install coscmd, used to upload files to cos
.PHONY: install.coscmd
install.coscmd:
	@if which pip &>/dev/null; then pip install coscmd; else pip3 install coscmd; fi

## install.minio: Install minio, used to upload files to minio
.PHONY: install.minio
install.minio:
	@$(GO) install github.com/minio/minio@latest

## install.delve: Install delve, used to debug go program
.PHONY: install.delve
install.delve:
	@$(GO) install github.com/go-delve/delve/cmd/dlv@latest

## install.air: Install air, used to hot reload go program
.PHONY: install.air
install.air:
	@$(GO) install github.com/cosmtrek/air@latest

## install.gvm: Install gvm, gvm is a Go version manager, built on top of the official go tool.
# github: https://github.com/moovweb/gvm
.PHONY: install.gvm
install.gvm:
	@echo "===========> Installing gvm,The default installation path is ~/.gvm/scripts/gvm"
	@bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
	@$(shell source /root/.gvm/scripts/gvm)

## install.golines: Install golines, used to format long lines
.PHONY: install.golines
install.golines:
	@$(GO) install github.com/segmentio/golines@latest

## install.go-mod-outdated: Install go-mod-outdated, used to check outdated dependencies
.PHONY: install.go-mod-outdated
install.go-mod-outdated:
	@$(GO) install github.com/psampaz/go-mod-outdated@latest

## install.mockgen: Install mockgen, used to generate mock functions
.PHONY: install.mockgen
install.mockgen:
	@$(GO) install github.com/golang/mock/mockgen@latest

## install.gotests: Install gotests, used to generate test functions
.PHONY: install.gotests
install.gotests:
	@$(GO) install github.com/cweill/gotests/gotests@latest

## install.protoc-gen-go: Install protoc-gen-go, used to generate go source files from protobuf files
.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	@$(GO) install github.com/golang/protobuf/protoc-gen-go@latest

## install.cfssl: Install cfssl, used to generate certificates
.PHONY: install.cfssl
install.cfssl:
	@$(ROOT_DIR)/scripts/install/install.sh openim::install::install_cfssl

## install.depth: Install depth, used to check dependency tree
.PHONY: install.depth
install.depth:
	@$(GO) install github.com/KyleBanks/depth/cmd/depth@latest

## install.go-callvis: Install go-callvis, used to visualize call graph
.PHONY: install.go-callvis
install.go-callvis:
	@$(GO) install github.com/ofabry/go-callvis@latest

## install.gothanks: Install gothanks, used to thank go dependencies
.PHONY: install.gothanks
install.gothanks:
	@$(GO) install github.com/psampaz/gothanks@latest

## install.richgo: Install richgo
.PHONY: install.richgo
install.richgo:
	@$(GO) install github.com/kyoh86/richgo@latest

## install.rts: Install rts
.PHONY: install.rts
install.rts:
	@$(GO) install github.com/galeone/rts/cmd/rts@latest

## tools.help: Display help information about the tools package
.PHONY: tools.help
tools.help: scripts/make-rules/tools.mk
	$(call smallhelp)