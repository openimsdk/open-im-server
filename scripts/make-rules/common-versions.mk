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
# OpenIM Makefile Versions used
#
# Define the latest version for each tool to ensure consistent versioning across installations
GOLANGCI_LINT_VERSION ?= latest
GOIMPORTS_VERSION ?= latest
ADDLICENSE_VERSION ?= latest
DEEPCOPY_GEN_VERSION ?= latest
CONVERSION_GEN_VERSION ?= latest
GINKGO_VERSION ?= v1.16.2
GO_GITLINT_VERSION ?= latest
GO_JUNIT_REPORT_VERSION ?= latest
GOTESTS_VERSION ?= latest
SWAGGER_VERSION ?= latest
KUBE_SCORE_VERSION ?= latest
KUBECONFORM_VERSION ?= latest
GSEMVER_VERSION ?= latest
GIT_CHGLOG_VERSION ?= latest
KO_VERSION ?= latest
GITHUB_RELEASE_VERSION ?= latest
COSCLI_VERSION ?= v0.19.0-beta
MINIO_VERSION ?= latest
DELVE_VERSION ?= latest
AIR_VERSION ?= latest
GOLINES_VERSION ?= latest
GO_MOD_OUTDATED_VERSION ?= latest
CFSSL_VERSION ?= latest
DEPTH_VERSION ?= latest
GO_CALLVIS_VERSION ?= latest
MISSPELL_VERSION ?= latest
GOTHANKS_VERSION ?= latest
RICHGO_VERSION ?= latest
RTS_VERSION ?= latest
TYPECHECK_VERSION ?= latest
COMMENT_LANG_DETECTOR_VERSION ?= latest
STANDARDIZER_VERSION ?= latest
GO_TESTS_VERSION ?= v1.6.0
GO_APIDIFF_VERSION ?= v0.8.2
KAFKACTL_VERSION ?= latest
GOTESTSUM_VERSION ?= latest

WIRE_VERSION ?= latest
# WIRE_VERSION ?= $(call get_go_version,github.com/google/wire)
MOCKGEN_VERSION ?= $(call get_go_version,github.com/golang/mock)
PROTOC_GEN_GO_VERSION ?= $(call get_go_version,github.com/golang/protobuf/protoc-gen-go)