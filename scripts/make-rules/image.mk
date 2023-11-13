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
# Makefile helper functions for docker image
# ==============================================================================
# Path: scripts/make-rules/image.mk
# docker registry: registry.example.com/namespace/image:tag as: registry.hub.docker.com/cubxxw/<image-name>:<tag>
# https://docs.docker.com/build/building/multi-platform/
#

DOCKER := docker

# read: https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md
REGISTRY_PREFIX ?= registry.cn-hangzhou.aliyuncs.com/openimsdk
# REGISTRY_PREFIX ?= ghcr.io/openimsdk

BASE_IMAGE ?= ghcr.io/openim-sigs/openim-bash-image

IMAGE_PLAT ?= $(subst $(SPACE),$(COMMA),$(subst _,/,$(PLATFORMS)))

EXTRA_ARGS ?= --no-cache
_DOCKER_BUILD_EXTRA_ARGS :=

ifdef HTTP_PROXY
_DOCKER_BUILD_EXTRA_ARGS += --build-arg HTTP_PROXY=${HTTP_PROXY}
endif

ifneq ($(EXTRA_ARGS), )
_DOCKER_BUILD_EXTRA_ARGS += $(EXTRA_ARGS)
endif

# Determine image files by looking into build/images/*/Dockerfile
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/images/*)
# Determine images names by stripping out the dir names, and filter out the undesired directories
# IMAGES ?= $(filter-out Dockerfile,$(foreach image,${IMAGES_DIR},$(notdir ${image})))
IMAGES ?= $(filter-out Dockerfile openim-tools openim-cmdutils,$(foreach image,${IMAGES_DIR},$(notdir ${image})))

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

# ==============================================================================
# Image targets
# ==============================================================================

# PLATFORMS defines the target platforms for  the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
## image.docker-buildx: Build and push docker image for the manager for cross-platform support
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
.PHONY: image.docker-buildx
image.docker-buildx:
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-v3-builder
	$(CONTAINER_TOOL) buildx use project-v3-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMAGES} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm project-v3-builder
	rm Dockerfile.cross

## image.verify: Verify docker version
.PHONY: image.verify
image.verify:
	@$(ROOT_DIR)/scripts/lib/util.sh openim::util::check_docker_and_compose_versions

## image.daemon.verify: Verify docker daemon experimental features
.PHONY: image.daemon.verify
image.daemon.verify:
	@$(ROOT_DIR)/scripts/lib/util.sh openim::util::ensure_docker_daemon_connectivity
	@$(ROOT_DIR)/scripts/lib/util.sh openim::util::ensure-docker-buildx

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
## image.build: Build docker images
.PHONY: image.build
image.build: image.verify $(addprefix image.build., $(addprefix $(PLATFORM)., $(IMAGES)))

.PHONY: image.build.multiarch
image.build.multiarch: image.verify $(foreach p,$(PLATFORMS),$(addprefix image.build., $(addprefix $(p)., $(IMAGES))))

## image.build.%: Build docker image for a specific platform
.PHONY: image.build.%
image.build.%: go.build.%
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)/$(PLATFORM)
	@cat $(ROOT_DIR)/build/images/Dockerfile\
		| sed "s#BASE_IMAGE#$(BASE_IMAGE)#g" \
		| sed "s#BINARY_NAME#$(IMAGE)#g" >$(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp $(BIN_DIR)/platforms/$(IMAGE_PLAT)/$(IMAGE) $(TMP_DIR)/$(IMAGE)
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION) $(TMP_DIR)/$(IMAGE))
	@echo $(DOCKER) build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX)
	@if [ $(shell $(GO) env GOARCH) != $(ARCH) ] ; then \
		$(MAKE) image.daemon.verify ;\
		$(DOCKER) build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX) ; \
	else \
		$(DOCKER) build $(BUILD_SUFFIX) ; \
	fi
	@rm -rf $(TMP_DIR)/$(IMAGE)

# https://docs.docker.com/build/building/multi-platform/
# busybox image supports amd64, arm32v5, arm32v6, arm32v7, arm64v8, i386, ppc64le, and s390x
## image.buildx.%: Build docker images with buildx
.PHONY: image.buildx.%
image.buildx.%:
	$(eval IMAGE := $(word 1,$(subst ., ,$*)))
	echo "===========> Building docker image $(IMAGE) $(VERSION)"
	$(DOCKER) buildx build -f $(ROOT_DIR)/Dockerfile --pull --no-cache --platform=$(PLATFORMS) --push . -t $(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION)

## image.push: Push docker images
.PHONY: image.push
image.push: image.verify go.build.verify $(addprefix image.push., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

## image.push.multiarch: Push docker images for all platforms
.PHONY: image.push.multiarch
image.push.multiarch: image.verify go.build.verify $(foreach p,$(PLATFORMS),$(addprefix image.push., $(addprefix $(p)., $(IMAGES))))

## image.push.%: Push docker image for a specific platform
.PHONY: image.push.%
image.push.%: image.build.%
	@echo "===========> Pushing image $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX)"
	$(DOCKER) push $(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION)

## image.manifest.push: Push manifest list for multi-arch images
.PHONY: image.manifest.push
image.manifest.push: export DOCKER_CLI_EXPERIMENTAL := enabled
image.manifest.push: image.verify go.build.verify \
$(addprefix image.manifest.push., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

## image.manifest.push.%: Push manifest list for multi-arch images for a specific platform
.PHONY: image.manifest.push.%
image.manifest.push.%: image.push.% image.manifest.remove.%
	@echo "===========> Pushing manifest $(IMAGE) $(VERSION) to $(REGISTRY_PREFIX) and then remove the local manifest list"
	@$(DOCKER) manifest create $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION) \
		$(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION)
	@$(DOCKER) manifest annotate $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION) \
		$(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION) \
		--os $(OS) --arch ${ARCH}
	@$(DOCKER) manifest push --purge $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION)

# Docker cli has a bug: https://github.com/docker/cli/issues/954
# If you find your manifests were not updated,
# Please manually delete them in $HOME/.docker/manifests/
# and re-run.
## image.manifest.remove.%: Remove local manifest list
.PHONY: image.manifest.remove.%
image.manifest.remove.%:
	@rm -rf ${HOME}/.docker/manifests/docker.io_$(REGISTRY_PREFIX)_$(IMAGE)-$(VERSION)

## image.manifest.push.multiarch: Push manifest list for multi-arch images for all platforms
.PHONY: image.manifest.push.multiarch
image.manifest.push.multiarch: image.push.multiarch $(addprefix image.manifest.push.multiarch., $(IMAGES))

## image.manifest.push.multiarch.%: Push manifest list for multi-arch images for all platforms for a specific image
.PHONY: image.manifest.push.multiarch.%
image.manifest.push.multiarch.%:
	@echo "===========> Pushing manifest $* $(VERSION) to $(REGISTRY_PREFIX) and then remove the local manifest list"
	REGISTRY_PREFIX=$(REGISTRY_PREFIX) PLATFORMS="$(PLATFORMS)" IMAGE=$* VERSION=$(VERSION) DOCKER_CLI_EXPERIMENTAL=enabled \
	  $(ROOT_DIR)/build/lib/create-manifest.sh

## image.help: Print help for image targets
.PHONY: image.help
image.help: scripts/make-rules/image.mk
	$(call smallhelp)