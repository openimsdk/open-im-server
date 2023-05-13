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


source ./proto_dir.cfg

for ((i = 0; i < ${#all_proto[*]}; i++)); do
  proto=${all_proto[$i]}
  protoc -I ../../../  -I ./ --go_out=plugins=grpc:. $proto
  echo "protoc --go_out=plugins=grpc:." $proto
done
echo "proto file generate success"


j=0
for file in $(find ./Open_IM -name   "*.go"); do # Not recommended, will break on whitespace
    filelist[j]=$file
    j=`expr $j + 1`
done


for ((i = 0; i < ${#filelist[*]}; i++)); do
  proto=${filelist[$i]}
  cp $proto  ${proto#*./Open_IM/pkg/proto/}
done
rm Open_IM -rf
#find ./ -type f -path "*.pb.go"|xargs sed -i 's/\".\/sdk_ws\"/\"Open_IM\/pkg\/proto\/sdk_ws\"/g'




