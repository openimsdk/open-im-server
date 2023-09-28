#!/bin/bash

IMAGEHUB="registry.cn-shenzhen.aliyuncs.com/huanglin_hub"
PROJECT=$1
ALLPRO="all"
servers=(openim-api openim-crontask openim-msggateway openim-msgtransfer openim-push openim-rpc-auth openim-rpc-conversation openim-rpc-friend openim-rpc-group openim-rpc-msg openim-rpc-third openim-rpc-user)


if [ "$1" != "" ]
then
    if [[ "${servers[@]}"  =~ "${1}" ]]
    then
        echo "building ${PROJECT}"
        DOCKER_PUSHIMG=${IMAGEHUB}/${PROJECT}:dev
        docker rmi  ${DOCKER_PUSHIMG}
        docker build -f manifest/dockerfiles/${PROJECT}/Dockerfile -t ${DOCKER_PUSHIMG} .
        docker push ${DOCKER_PUSHIMG}
    elif [[ ! "${servers[@]}"  =~ "${1}" ]]
    then
        if [ ${PROJECT} == ${ALLPRO} ]
        then
            echo "building allproject"
            for element in ${servers[@]}
            do
                SUB_IMG=${element}
                SUB_PUSHIMG=${IMAGEHUB}/${element}:dev
                docker rmi  ${SUB_PUSHIMG}
                docker build -f manifest/dockerfiles/${SUB_IMG}/Dockerfile -t ${SUB_PUSHIMG} .
                docker push ${SUB_PUSHIMG}
            done
        else
          echo "输入的项目名称不正确"
        fi
    fi
else
    echo "请传入一个参数"
fi