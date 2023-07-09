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

version=errcode
repository=${1}
if [[ -z ${repository} ]]
then
  echo "repository is empty"
  exit 0
fi

set +e
echo "repository: ${repository}"
source ./path_info.sh
echo "start to build docker images"
currentPwd=`pwd`
echo ${currentPwd}
i=0
for path in  ${service_source_root[*]}
do
  cd ${path}
  make build
  image="${repository}/${image_names[${i}]}:$version"
  echo ${image}
  docker build -t $image . -f ./deploy.Dockerfile
  echo "build ${image} success"
  docker push ${image}
  echo "push ${image} success"
  echo "=============================="
  i=$((i + 1))
  cd ${currentPwd}
done

echo "build all images success"