#!/bin/bash

PROJECT=$1
DOCKER_IMG=${PROJECT}
docker rmi  ${DOCKER_IMG}
docker build -f manifest/dockerfiles/${DOCKER_IMG}/Dockerfile -t ${DOCKER_IMG} .