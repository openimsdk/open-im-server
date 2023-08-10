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

OPENIM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

# 生成文件存放目录
LOCAL_OUTPUT_ROOT="${OPENIM_ROOT}/${OUT_DIR:-_output}"

# app要能访问到此ip和端口或域名
readonly API_URL=${API_URL:-http://127.0.0.1:10002/object/}
readonly DATA_DIR=${DATA_DIR:-${OPENIM_ROOT}}

# 设置统一的用户名，方便记忆
readonly USER=${USER:-'root'} # Setting a username

# 设置统一的密码，方便记忆
readonly PASSWORD=${PASSWORD:-'openIM123'} # Setting a password

# Linux系统 going 用户
readonly LINUX_USERNAME=${LINUX_USERNAME:-going}
# Linux root & going 用户密码
readonly LINUX_PASSWORD=${LINUX_PASSWORD:-${PASSWORD}}

# 设置安装目录
readonly INSTALL_DIR=${INSTALL_DIR:-/tmp/installation}
mkdir -p ${INSTALL_DIR}
readonly ENV_FILE=${OPENIM_ROOT}/scripts/install/environment.sh

# MINIO 配置信息
readonly OBJECT_ENABLE=${OBJECT_ENABLE:-minio}
readonly OBJECT_APIURL=${OBJECT_APIURL:-http://127.0.0.1:10002/object/}
readonly MINIO_BUCKET=${MINIO_BUCKET:-openim}
readonly MINIO_ENDPOINT=${MINIO_ENDPOINT:-http://127.0.0.1:10005}
readonly MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY:-root}
readonly MINIO_SECRET_KEY=${MINIO_SECRET_KEY:-openIM123}
readonly COS_BUCKET_URL=${COS_BUCKET_URL:-https://temp-1252357374.cos.ap-chengdu.myqcloud.com}
readonly OSS_ENDPOINT=${OSS_ENDPOINT:-http://oss-cn-chengdu.aliyuncs.com}
readonly OSS_BUCKET=${OSS_BUCKET:-demo-9999999}
readonly OSS_BUCKET_URL=${OSS_BUCKET_URL:-https://demo-9999999.oss-cn-chengdu.aliyuncs.com}
readonly OSS_ACCESS_KEY_ID=${OSS_ACCESS_KEY_ID:-root}

# MariaDB 配置信息
readonly MARIADB_ADMIN_USERNAME=${MARIADB_ADMIN_USERNAME:-root} # MariaDB root 用户
readonly MARIADB_ADMIN_PASSWORD=${MARIADB_ADMIN_PASSWORD:-${PASSWORD}} # MariaDB root 用户密码
readonly MARIADB_HOST=${MARIADB_HOST:-127.0.0.1:3306} # MariaDB 主机地址
readonly MARIADB_DATABASE=${MARIADB_DATABASE:-openim} # MariaDB openim 应用使用的数据库名
readonly MARIADB_USERNAME=${MARIADB_USERNAME:-openim} # openim 数据库用户名
readonly MARIADB_PASSWORD=${MARIADB_PASSWORD:-${PASSWORD}} # openim 数据库密码

# Redis 配置信息
readonly REDIS_HOST=${REDIS_HOST:-127.0.0.1} # Redis 主机地址
readonly REDIS_PORT=${REDIS_PORT:-6379} # Redis 监听端口
readonly REDIS_USERNAME=${REDIS_USERNAME:-''} # Redis 用户名
readonly REDIS_PASSWORD=${REDIS_PASSWORD:-${PASSWORD}} # Redis 密码

# MongoDB 配置
readonly MONGO_ADMIN_USERNAME=${MONGO_ADMIN_USERNAME:-root} # MongoDB root 用户
readonly MONGO_ADMIN_PASSWORD=${MONGO_ADMIN_PASSWORD:-${PASSWORD}} # MongoDB root 用户密码
readonly MONGO_HOST=${MONGO_HOST:-127.0.0.1} # MongoDB 地址
readonly MONGO_PORT=${MONGO_PORT:-27017} # MongoDB 端口
readonly MONGO_USERNAME=${MONGO_USERNAME:-openim} # MongoDB 用户名
readonly MONGO_PASSWORD=${MONGO_PASSWORD:-${PASSWORD}} # MongoDB 密码

# openim 配置
readonly OPENIM_DATA_DIR=${OPENIM_DATA_DIR:-/data/openim} # openim 各组件数据目录
readonly OPENIM_INSTALL_DIR=${OPENIM_INSTALL_DIR:-/opt/openim} # openim 安装文件存放目录
readonly OPENIM_CONFIG_DIR=${OPENIM_CONFIG_DIR:-/etc/openim} # openim 配置文件存放目录
readonly OPENIM_LOG_DIR=${OPENIM_LOG_DIR:-/var/log/openim} # openim 日志文件存放目录
readonly CA_FILE=${CA_FILE:-${OPENIM_CONFIG_DIR}/cert/ca.pem} # CA

# openim-apiserver 配置
readonly OPENIM_APISERVER_HOST=${OPENIM_APISERVER_HOST:-127.0.0.1} # openim-apiserver 部署机器 IP 地址
readonly OPENIM_APISERVER_GRPC_BIND_ADDRESS=${OPENIM_APISERVER_GRPC_BIND_ADDRESS:-0.0.0.0}
readonly OPENIM_APISERVER_GRPC_BIND_PORT=${OPENIM_APISERVER_GRPC_BIND_PORT:-8081}
readonly OPENIM_APISERVER_INSECURE_BIND_ADDRESS=${OPENIM_APISERVER_INSECURE_BIND_ADDRESS:-127.0.0.1}
readonly OPENIM_APISERVER_INSECURE_BIND_PORT=${OPENIM_APISERVER_INSECURE_BIND_PORT:-8080}
readonly OPENIM_APISERVER_SECURE_BIND_ADDRESS=${OPENIM_APISERVER_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly OPENIM_APISERVER_SECURE_BIND_PORT=${OPENIM_APISERVER_SECURE_BIND_PORT:-8443}
readonly OPENIM_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${OPENIM_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE:-${OPENIM_CONFIG_DIR}/cert/openim-apiserver.pem}
readonly OPENIM_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${OPENIM_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE:-${OPENIM_CONFIG_DIR}/cert/openim-apiserver-key.pem}

# openim-authz-server 配置
readonly OPENIM_AUTHZ_SERVER_HOST=${OPENIM_AUTHZ_SERVER_HOST:-127.0.0.1} # openim-authz-server 部署机器 IP 地址
readonly OPENIM_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=${OPENIM_AUTHZ_SERVER_INSECURE_BIND_ADDRESS:-127.0.0.1}
readonly OPENIM_AUTHZ_SERVER_INSECURE_BIND_PORT=${OPENIM_AUTHZ_SERVER_INSECURE_BIND_PORT:-9090}
readonly OPENIM_AUTHZ_SERVER_SECURE_BIND_ADDRESS=${OPENIM_AUTHZ_SERVER_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly OPENIM_AUTHZ_SERVER_SECURE_BIND_PORT=${OPENIM_AUTHZ_SERVER_SECURE_BIND_PORT:-9443}
readonly OPENIM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${OPENIM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_CERT_FILE:-${OPENIM_CONFIG_DIR}/cert/openim-authz-server.pem}
readonly OPENIM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${OPENIM_AUTHZ_SERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE:-${OPENIM_CONFIG_DIR}/cert/openim-authz-server-key.pem}
readonly OPENIM_AUTHZ_SERVER_CLIENT_CA_FILE=${OPENIM_AUTHZ_SERVER_CLIENT_CA_FILE:-${CA_FILE}}
readonly OPENIM_AUTHZ_SERVER_RPCSERVER=${OPENIM_AUTHZ_SERVER_RPCSERVER:-${OPENIM_APISERVER_HOST}:${OPENIM_APISERVER_GRPC_BIND_PORT}}

# openim-pump 配置
readonly OPENIM_PUMP_HOST=${OPENIM_PUMP_HOST:-127.0.0.1} # openim-pump 部署机器 IP 地址
readonly OPENIM_PUMP_COLLECTION_NAME=${OPENIM_PUMP_COLLECTION_NAME:-openim_analytics}
readonly OPENIM_PUMP_MONGO_URL=${OPENIM_PUMP_MONGO_URL:-mongodb://${MONGO_USERNAME}:${MONGO_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/${OPENIM_PUMP_COLLECTION_NAME}?authSource=${OPENIM_PUMP_COLLECTION_NAME}}

# openim-watcher配置
readonly OPENIM_WATCHER_HOST=${OPENIM_WATCHER_HOST:-127.0.0.1} # openim-watcher 部署机器 IP 地址

# openimctl 配置
readonly CONFIG_USER_USERNAME=${CONFIG_USER_USERNAME:-admin}
readonly CONFIG_USER_PASSWORD=${CONFIG_USER_PASSWORD:-Admin@2021}
readonly CONFIG_USER_CLIENT_CERTIFICATE=${CONFIG_USER_CLIENT_CERTIFICATE:-${HOME}/.openim/cert/admin.pem}
readonly CONFIG_USER_CLIENT_KEY=${CONFIG_USER_CLIENT_KEY:-${HOME}/.openim/cert/admin-key.pem}
readonly CONFIG_SERVER_ADDRESS=${CONFIG_SERVER_ADDRESS:-${OPENIM_APISERVER_HOST}:${OPENIM_APISERVER_SECURE_BIND_PORT}}
readonly CONFIG_SERVER_CERTIFICATE_AUTHORITY=${CONFIG_SERVER_CERTIFICATE_AUTHORITY:-${CA_FILE}}
