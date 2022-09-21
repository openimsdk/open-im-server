#!/usr/bin/env bash

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




