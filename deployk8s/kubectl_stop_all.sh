#!/usr/bin/env bash

service=(
  #api service file
  api
  cms-api
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
  sdk-server
)

for i in ${service[*]}
do
    kubectl -n openim delete deployment "${i}-deployment"
done

kubectl -n openim delete service api
kubectl -n openim delete service sdk-server
kubectl -n openim delete service msg-gateway

echo done