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
# TODO: For the time being only used for compilation, it can be arm or amd, please do not delete it, it can be extended with new functions
# ==============================================================================
# Path: scripts/make-rules/image.mk
# docker registry: registry.example.com/namespace/image:tag as: registry.hub.docker.com/cubxxw/<image-name>:<tag>
# https://docs.docker.com/build/building/multi-platform/
#

DOCKER := docker
DOCKER_SUPPORTED_API_VERSION ?= 1.32|1.40|1.41

REGISTRY_PREFIX ?= ghcr.io/OpenIMSDK
IMAGES ?= lvscare
IMAGE_PLAT ?= $(subst $(SPACE),$(COMMA),$(subst _,/,$(PLATFORMS)))

EXTRA_ARGS ?= --no-cache
_DOCKER_BUILD_EXTRA_ARGS :=

ifdef HTTP_PROXY
_DOCKER_BUILD_EXTRA_ARGS += --build-arg HTTP_PROXY=${HTTP_PROXY}
endif

ifneq ($(EXTRA_ARGS), )
_DOCKER_BUILD_EXTRA_ARGS += $(EXTRA_ARGS)
endif

# Determine image files by looking into build/docker/*/Dockerfile
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/docker/*)
# Determine images names by stripping out the dir names
IMAGES ?= $(filter-out tools,$(foreach image,${IMAGES_DIR},$(notdir ${image})))

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
## Build and push docker image for the manager for cross-platform support
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
.PHONY: docker-buildx
docker-buildx:
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-v3-builder
	$(CONTAINER_TOOL) buildx use project-v3-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMAGES} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm project-v3-builder
	rm Dockerfile.cross

## image.verify: Verify docker version
.PHONY: image.verify
image.verify:
	$(eval API_VERSION := $(shell $(DOCKER) version | grep -E 'API version: {1,6}[0-9]' | head -n1 | awk '{print $$3} END { if (NR==0) print 0}' ))
	$(eval PASS := $(shell echo "$(API_VERSION) > $(DOCKER_SUPPORTED_API_VERSION)" | bc))
	@if [ $(PASS) -ne 1 ]; then \
		$(DOCKER) -v ;\
		echo "Unsupported docker version. Docker API version should be greater than $(DOCKER_SUPPORTED_API_VERSION)"; \
		exit 1; \
	fi

## image.daemon.verify: Verify docker daemon experimental features
.PHONY: image.daemon.verify
image.daemon.verify:
	$(eval PASS := $(shell $(DOCKER) version | grep -q -E 'Experimental: {1,5}true' && echo 1 || echo 0))
	@if [ $(PASS) -ne 1 ]; then \
		echo "Experimental features of Docker daemon is not enabled. Please add \"experimental\": true in '/etc/docker/daemon.json' and then restart Docker daemon."; \
		exit 1; \
	fi

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
image.build.%: go.bin.%
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building LOCAL docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)/$(PLATFORM)
	@cat $(ROOT_DIR)/docker/$(IMAGE)/Dockerfile\
		>$(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp $(BIN_DIR)/$(PLATFORM)/$(IMAGE) $(TMP_DIR)/$(IMAGE)/$(PLATFORM)

	$(eval BUILD_SUFFIX := --load --pull -t $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION) $(TMP_DIR)/$(IMAGE))
	$(eval BUILD_SUFFIX_ARM := --load --pull -t $(REGISTRY_PREFIX)/$(IMAGE).$(ARCH):$(VERSION) $(TMP_DIR)/$(IMAGE))
	@if [ "$(ARCH)" == "amd64" ]; then \
		echo "===========> Creating LOCAL docker image tag $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION) for $(ARCH)"; \
		$(DOCKER) buildx build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX); \
	else \
		echo "===========> Creating LOCAL docker image tag $(REGISTRY_PREFIX)/$(IMAGE).$(ARCH):$(VERSION) for $(ARCH)"; \
		$(DOCKER) buildx build --platform $(IMAGE_PLAT) $(BUILD_SUFFIX_ARM); \
	fi

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
	REGISTRY_PREFIX=$(REGISTRY_PREFIX) PLATFROMS="$(PLATFORMS)" IMAGE=$* VERSION=$(VERSION) DOCKER_CLI_EXPERIMENTAL=enabled \
	  $(ROOT_DIR)/build/lib/create-manifest.sh

## image.help: Print help for image targets
.PHONY: image.help
image.help: scripts/make-rules/image.mk
	$(call smallhelp)