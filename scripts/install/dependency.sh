#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script will install the dependencies required for openim

set -o errexit
set +o nounset
set -o pipefail

OPENIM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)
[[ -z ${COMMON_SOURCED} ]] && source "${OPENIM_ROOT}"/scripts/install/common.sh

# Start MySQL service
docker run -d \
    --name mysql \
    -p 13306:3306 \
    -p 23306:33060 \
    -v "${DATA_DIR}/components/mysql/data:/var/lib/mysql" \
    -v "/etc/localtime:/etc/localtime" \
    -e MYSQL_ROOT_PASSWORD=${PASSWORD} \
    --restart always \
    mysql:5.7

# Start MongoDB service
docker run -d \
    --name mongo \
    -p 37017:27017 \
    -v "${DATA_DIR}/components/mongodb/data/db:/data/db" \
    -v "${DATA_DIR}/components/mongodb/data/logs:/data/logs" \
    -v "${DATA_DIR}/components/mongodb/data/conf:/etc/mongo" \
    -v "./scripts/mongo-init.sh:/docker-entrypoint-initdb.d/mongo-init.sh:ro" \
    -e TZ=Asia/Shanghai \
    -e wiredTigerCacheSizeGB=1 \
    -e MONGO_INITDB_ROOT_USERNAME=${USER} \
    -e MONGO_INITDB_ROOT_PASSWORD=${PASSWORD} \
    -e MONGO_INITDB_DATABASE=openIM \
    -e MONGO_USERNAME=${USER} \
    -e MONGO_PASSWORD=${PASSWORD} \
    --restart always \
    mongo:6.0.2 --wiredTigerCacheSizeGB 1 --auth

# Start Redis service
docker run -d \
    --name redis \
    -p 16379:6379 \
    -v "${DATA_DIR}/components/redis/data:/data" \
    -v "${DATA_DIR}/components/redis/config/redis.conf:/usr/local/redis/config/redis.conf" \
    -e TZ=Asia/Shanghai \
    --sysctl net.core.somaxconn=1024 \
    --restart always \
    redis:7.0.0 redis-server --requirepass ${PASSWORD} --appendonly yes

# Start Zookeeper service
docker run -d \
    --name zookeeper \
    -p 2181:2181 \
    -v "/etc/localtime:/etc/localtime" \
    -e TZ=Asia/Shanghai \
    --restart always \
    wurstmeister/zookeeper

# Start Kafka service
docker run -d \
    --name kafka \
    -p 9092:9092 \
    -e TZ=Asia/Shanghai \
    -e KAFKA_BROKER_ID=0 \
    -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
    -e KAFKA_CREATE_TOPICS="latestMsgToRedis:8:1,msgToPush:8:1,offlineMsgToMongoMysql:8:1" \
    -e KAFKA_ADVERTISED_LISTENERS="INSIDE://127.0.0.1:9092,OUTSIDE://103.116.45.174:9092" \
    -e KAFKA_LISTENERS="INSIDE://:9092,OUTSIDE://:9093" \
    -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP="INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT" \
    -e KAFKA_INTER_BROKER_LISTENER_NAME=INSIDE \
    --restart always \
    --link zookeeper \
    wurstmeister/kafka

# Start MinIO service
docker run -d \
    --name minio \
    -p 10005:9000 \
    -p 9090:9090 \
    -v "/mnt/data:/data" \
    -v "/mnt/config:/root/.minio" \
    -e MINIO_ROOT_USER=${USER} \
    -e MINIO_ROOT_PASSWORD=${PASSWORD} \
    --restart always \
    minio/minio server /data --console-address ':9090'
