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
# Makefile helper functions for generate necessary files and docs
# https://cloud.redhat.com/blog/kubernetes-deep-dive-code-generation-customresources
# ! The stock of code generated by `make gen` should be idempotent
#
# Questions about go mod instead of go path: https://github.com/kubernetes/kubernetes/issues/117181
# ==============================================================================
# Makefile helper functions for generate necessary files
#

## gen.init: Initialize openim server project ✨
.PHONY: gen.init
gen.init:
	@echo "===========> Initializing openim server project"
	@${ROOT_DIR}/scripts/init-config.sh

## gen.run: Generate necessary files and docs ✨
.PHONY: gen.run
#gen.run: gen.errcode gen.docgo
gen.run: gen.clean gen.errcode gen.docgo.doc

## gen.errcode: Generate necessary files and docs ✨
.PHONY: gen.errcode
gen.errcode: gen.errcode.code gen.errcode.doc

## gen.errcode.code: Generate openim error code go source files ✨
.PHONY: gen.errcode.code
gen.errcode.code: tools.verify.codegen
	@echo "===========> Generating openim error code go source files"
	@codegen -type=int ${ROOT_DIR}/internal/pkg/code

## gen.errcode.doc: Generate openim error code markdown documentation ✨
.PHONY: gen.errcode.doc
gen.errcode.doc: tools.verify.codegen
	@echo "===========> Generating error code markdown documentation"
	@codegen -type=int -doc \
		-output ${ROOT_DIR}/docs/guide/zh-CN/api/error_code_generated.md ${ROOT_DIR}/internal/pkg/code

## gen.docgo: Generate missing doc.go for go packages ✨
.PHONY: gen.ca.%
gen.ca.%:
	$(eval CA := $(word 1,$(subst ., ,$*)))
	@echo "===========> Generating CA files for $(CA)"
	@${ROOT_DIR}/scripts/gencerts.sh generate-openim-cert $(OUTPUT_DIR)/cert $(CA)

## gen.ca: Generate CA files for all certificates ✨
.PHONY: gen.ca
gen.ca: $(addprefix gen.ca., $(CERTIFICATES))

## gen.docgo: Generate missing doc.go for go packages ✨
.PHONY: gen.docgo.doc
gen.docgo.doc:
	@echo "===========> Generating missing doc.go for go packages"
	@${ROOT_DIR}/scripts/gendoc.sh

## gen.docgo.check: Check if there are untracked doc.go files ✨
.PHONY: gen.docgo.check
gen.docgo.check: gen.docgo.doc
	@n="$$(git ls-files --others '*/doc.go' | wc -l)"; \
	if test "$$n" -gt 0; then \
		git ls-files --others '*/doc.go' | sed -e 's/^/  /'; \
		echo "$@: untracked doc.go file(s) exist in working directory" >&2 ; \
		false ; \
	fi

## gen.docgo.add: Add untracked doc.go files to git index ✨
.PHONY: gen.docgo.add
gen.docgo.add:
	@git ls-files --others '*/doc.go' | $(XARGS) -- git add

## gen.docgo: Generate missing doc.go for go packages ✨
.PHONY: gen.defaultconfigs
gen.defaultconfigs:
	@${ROOT_DIR}/scripts/gen_default_config.sh

## gen.docgo: Generate missing doc.go for go packages ✨
.PHONY: gen.clean
gen.clean:
	@rm -rf ./api/client/{clientset,informers,listers}
	@$(FIND) -type f -name '*_generated.go' -delete

## gen.help: show help for gen
.PHONY: gen.help
gen.help: scripts/make-rules/gen.mk
	$(call smallhelp)