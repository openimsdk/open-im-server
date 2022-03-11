#!/bin/bash
source ../setting.env
VERSION_TAG="2.0.2"
IMAGE_TAG="${DOCKER_REGISTRY_ADDR}open_im_sdk_server:v$VERSION_TAG"

if [ ! -f "./Open-IM-SDK-Core.tar.gz" ];then

rm -rf Open-IM-Server.tar.gz && \
wget https://github.91chi.fun//https://github.com//OpenIMSDK/Open-IM-Server/archive/refs/tags/v$VERSION_TAG.tar.gz && \
mv v$VERSION_TAG.tar.gz Open-IM-Server.tar.gz && \
wget https://github.91chi.fun//https://github.com//OpenIMSDK/Open-IM-SDK-Core/archive/refs/tags/v$VERSION_TAG.tar.gz && \
mv v$VERSION_TAG.tar.gz Open-IM-SDK-Core.tar.gz

fi
echo  "下载完成" && \
docker build . -t $IMAGE_TAG --no-cache && \
echo "构建完成" && \
echo "构建完成" && \
echo "构建完成" && \
echo "构建完成" && \
echo "构建完成" && \
echo "构建完成" && \
echo "构建完成" && \
echo "构建完成" && \
echo "构建完成" && \
echo "构建完成" && \
docker push $IMAGE_TAG && \
cp development.tmp.yaml development.yaml

if [[ `uname` == 'Darwin' ]]; then
    sed -i "" "s#IMAGE_TAG#${IMAGE_TAG}#g" ./development.yaml
    sed -i "" "s#NODE_PORT_API#${NODE_PORT_API}#g" development.yaml
    sed -i "" "s#NODE_PORT_MSG_GATEWAY#${NODE_PORT_MSG_GATEWAY}#g" development.yaml
    sed -i "" "s#NODE_PORT_SDK_SERVER#${NODE_PORT_SDK_SERVER}#g" development.yaml
    sed -i "" "s#NODE_PORT_DEMO#${NODE_PORT_DEMO}#g" development.yaml
elif [[ `uname` == 'Linux' ]]; then
  sed -i "s#IMAGE_TAG#${IMAGE_TAG}#g" development.yaml
  sed -i "s#NODE_PORT_API#${NODE_PORT_API}#g" development.yaml
  sed -i "s#NODE_PORT_MSG_GATEWAY#${NODE_PORT_MSG_GATEWAY}#g" development.yaml
  sed -i "s#NODE_PORT_SDK_SERVER#${NODE_PORT_SDK_SERVER}#g" development.yaml
  sed -i "s#NODE_PORT_DEMO#${NODE_PORT_DEMO}#g" development.yaml
fi
echo "推送镜像完成" && kubectl -n ${K8S_NAMESPACE} delete -f development.yaml
kubectl -n ${K8S_NAMESPACE} apply -f development.yaml
