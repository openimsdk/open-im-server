#!/usr/bin/env bash

service=(
  #api service file
  api
  #rpc service file
  user
  friend
  group
  auth
  conversation
  msg-gateway
  msg-transfer
  msg
  push
)

for i in ${service[*]}
do
    kubectl -n openim delete deployment "${i}-deployment"
done

kubectl -n openim delete service api
kubectl -n openim delete service msg-gateway

echo done