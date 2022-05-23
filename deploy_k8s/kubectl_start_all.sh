#!/usr/bin/env bash

source ./path_info.cfg

for i in ${service[*]}
do
  kubectl -n openim apply -f ./${i}/deployment.yaml
done

