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

# PLATFORMS defines the target platforms for  the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: test ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- docker buildx create --name project-v3-builder
	docker buildx use project-v3-builder
	- docker buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- docker buildx rm project-v3-builder
	rm Dockerfile.cross

DOCKER := docker
DOCKER_SUPPORTED_API_VERSION ?= 1.32|1.40|1.41

REGISTRY_PREFIX ?= cubxxw
BASE_IMAGE = centos:centos8

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

# ifeq (${IMAGES},)
#   $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
# endif

# ==============================================================================
# Image targets
# ==============================================================================

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
image.build: image.verify go.build.verify $(addprefix image.build., $(addprefix $(IMAGE_PLAT)., $(IMAGES)))

## image.build.multiarch: Build docker images for all platforms
.PHONY: image.build.multiarch
image.build.multiarch: image.verify go.build.verify $(foreach p,$(PLATFORMS),$(addprefix image.build., $(addprefix $(p)., $(IMAGES))))

## image.build.%: Build docker image for a specific platform
.PHONY: image.build.%
image.build.%: go.build.%
	$(eval IMAGE := $(COMMAND))
	$(eval IMAGE_PLAT := $(subst _,/,$(PLATFORM)))
	@echo "===========> Building docker image $(IMAGE) $(VERSION) for $(IMAGE_PLAT)"
	@mkdir -p $(TMP_DIR)/$(IMAGE)
	@cat $(ROOT_DIR)/build/docker/$(IMAGE)/Dockerfile\
		| sed "s#BASE_IMAGE#$(BASE_IMAGE)#g" >$(TMP_DIR)/$(IMAGE)/Dockerfile
	@cp $(OUTPUT_DIR)/platforms/$(IMAGE_PLAT)/$(IMAGE) $(TMP_DIR)/$(IMAGE)/
	@DST_DIR=$(TMP_DIR)/$(IMAGE) $(ROOT_DIR)/build/docker/$(IMAGE)/build.sh 2>/dev/null || true
	$(eval BUILD_SUFFIX := $(_DOCKER_BUILD_EXTRA_ARGS) --pull -t $(REGISTRY_PREFIX)/$(IMAGE)-$(ARCH):$(VERSION) $(TMP_DIR)/$(IMAGE))
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
	REGISTRY_PREFIX=$(REGISTRY_PREFIX) PLATFROMS="$(PLATFORMS)" IMAGE=$* VERSION=$(VERSION) DOCKER_CLI_EXPERIMENTAL=enabled \
	  $(ROOT_DIR)/build/lib/create-manifest.sh

## image.help: Print help for image targets
.PHONY: image.help
image.help: scripts/make-rules/image.mk
	$(call smallhelp)