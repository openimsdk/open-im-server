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
#