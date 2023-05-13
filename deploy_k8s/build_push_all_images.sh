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

#/bin/sh
source ./path_info.cfg

# images version
version=v2.3.8
git pull
cd ../script/; ./build_all_service.sh
cd ../deploy_k8s/

for i in  ${service[*]}
do
  mv ../bin/open_im_${i} ./${i}/
done

echo "move success"

echo "start to build images"

for i in ${service[*]}
do
	echo "start to build images" $i
	cd $i
	image="openim/${i}:$version"
	docker build -t $image . -f ./${i}.Dockerfile
	echo "build ${dockerfile} success"
	docker push $image
	echo "push ${image} success "
	cd ..
done

