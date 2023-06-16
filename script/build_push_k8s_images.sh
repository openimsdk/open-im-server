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
currentPwd=${pwd}
i=0
for path in  ${service_source_root[*]}
do
  cd ${path}
  make build
  image="${repository}/${image_names[${i}]}:$version"
  echo ${image}
  docker build -t $image . -f ${path}/deploy.Dockerfile
  echo "build ${image} success"
  docker push ${image}
  echo "push ${image} success"
  echo "=============================="
  i=$((i + 1))
  rm -rf ${service_names[${i}]}
  cd ${currentPwd}
done

echo "build all images success"