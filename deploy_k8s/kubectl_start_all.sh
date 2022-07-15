#!/usr/bin/env bash

source ./path_info.cfg

#mkdir -p /db/sdk #path for jssdk sqlite

for i in ${service[*]}
do
  kubectl -n openim apply -f ./${i}/deployment.yaml
done

