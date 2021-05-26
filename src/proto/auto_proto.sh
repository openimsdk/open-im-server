#!/usr/bin/env bash

source ./proto_dir.cfg

for ((i = 0; i < ${#all_proto[*]}; i++)); do
  proto=${all_proto[$i]}
  protoc --go_out=plugins=grpc:. $proto
  echo "protoc --go_out=plugins=grpc:." $proto
done
echo "proto file generate success..."
