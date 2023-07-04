#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
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

image=openim/open_im_server:v1.0.5
rm Open-IM-Server -rf
git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive
cd Open-IM-Server
git checkout tuoyun
cd cmd/Open-IM-SDK-Core/
git checkout tuoyun
cd ../../
docker build -t  $image . -f deploy.Dockerfile
docker push $image