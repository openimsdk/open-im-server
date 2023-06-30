#!/usr/bin/env bash
version=errcode
repository=${1}
if [[ -z ${repository} ]]
then
  echo "repository is empty"
  exit 0
fi

set +e
echo "repository: ${repository}"
source ./path_info.cfg
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