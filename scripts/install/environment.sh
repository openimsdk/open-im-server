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

# This is a file that initializes variables for the automation script that initializes the config file
# You need to supplement the script according to the specification.
# Read: https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/init-config.md
# 格式化 bash 注释：https://tool.lu/shell/
# 配置中心文档：https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/environment.md

OPENIM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

# 生成文件存放目录
LOCAL_OUTPUT_ROOT=""${OPENIM_ROOT}"/${OUT_DIR:-_output}"
source "${OPENIM_ROOT}/scripts/lib/init.sh"

#TODO: Access to the OPENIM_IP networks outside, or you want to use the OPENIM_IP network
# OPENIM_IP=127.0.0.1
if [ -z "${OPENIM_IP}" ]; then
	OPENIM_IP=$(openim::util::get_server_ip)
fi

# config.gateway custom bridge modes
# if [ -z "{IP_GATEWAY}" ] then
#     IP_GATEWAY=$(openim::util::get_local_ip)
# fi

function def() {
	local var_name="$1"
	local default_value="${2:-}"
	eval "readonly $var_name=\"\${$var_name:-$(printf '%q' "$default_value")}\""
}

# OpenIM Docker Compose 数据存储的默认路径
def "DATA_DIR" "${OPENIM_ROOT}"

# 设置统一的用户名，方便记忆
def "OPENIM_USER" "root"

# 设置统一的密码，方便记忆
readonly PASSWORD=${PASSWORD:-'openIM123'}

# 设置统一的数据库名称，方便管理
def "DATABASE_NAME" "openIM_v3"

# Linux系统 openim 用户
def "LINUX_USERNAME" "openim"
readonly LINUX_PASSWORD=${LINUX_PASSWORD:-"${PASSWORD}"}

# 设置安装目录
def "INSTALL_DIR" "${LOCAL_OUTPUT_ROOT}/installs"
mkdir -p ${INSTALL_DIR}

def "ENV_FILE" ""${OPENIM_ROOT}"/scripts/install/environment.sh"

###################### Docker compose ###################
# OPENIM AND CHAT
def "CHAT_BRANCH" "main"
def "SERVER_BRANCH" "main"

# Choose the appropriate image address, the default is GITHUB image,
# you can choose docker hub, for Chinese users can choose Ali Cloud
# export IMAGE_REGISTRY="ghcr.io/openimsdk"
# export IMAGE_REGISTRY="openim"
# export IMAGE_REGISTRY="registry.cn-hangzhou.aliyuncs.com/openimsdk"
def "IMAGE_REGISTRY" "ghcr.io/openimsdk"
# def "IMAGE_REGISTRY" "openim"
# def "IMAGE_REGISTRY" "registry.cn-hangzhou.aliyuncs.com/openimsdk"

# Choose the appropriate image tag, the default is the latest version
def "SERVER_IMAGE_TAG" "latest"

###################### OpenIM Docker Network ######################
# 设置 Docker 网络的网段
readonly DOCKER_BRIDGE_SUBNET=${DOCKER_BRIDGE_SUBNET:-'172.28.0.0/16'}
IP_PREFIX=$(echo $DOCKER_BRIDGE_SUBNET | cut -d '/' -f 1)
SUBNET=$(echo $DOCKER_BRIDGE_SUBNET | cut -d '/' -f 2)
LAST_OCTET=$(echo $IP_PREFIX | cut -d '.' -f 4)

generate_ip() {
	local NEW_IP="$(echo $IP_PREFIX | cut -d '.' -f 1-3).$((LAST_OCTET++))"
	echo $NEW_IP
}
LAST_OCTET=$((LAST_OCTET + 1))
DOCKER_BRIDGE_GATEWAY=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
MYSQL_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
MONGO_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
REDIS_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
KAFKA_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
ZOOKEEPER_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
MINIO_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
OPENIM_WEB_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
OPENIM_SERVER_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
OPENIM_CHAT_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
PROMETHEUS_NETWORK_ADDRESS=$(generate_ip)
LAST_OCTET=$((LAST_OCTET + 1))
GRAFANA_NETWORK_ADDRESS=$(generate_ip)

###################### openim 配置 ######################
# read: https://github.com/openimsdk/open-im-server/blob/main/deployment/README.md
def "OPENIM_DATA_DIR" "/data/openim"
def "OPENIM_INSTALL_DIR" "/opt/openim"
def "OPENIM_CONFIG_DIR" "/etc/openim/config"
def "OPENIM_LOG_DIR" "/var/log/openim"
def "CA_FILE" "${OPENIM_CONFIG_DIR}/cert/ca.pem"

def "OPNEIM_CONFIG" ""${OPENIM_ROOT}"/config"
def "OPENIM_SERVER_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}" # OpenIM服务地址

# OpenIM Websocket端口
readonly OPENIM_WS_PORT=${OPENIM_WS_PORT:-'10001'}

# OpenIM API端口
readonly API_OPENIM_PORT=${API_OPENIM_PORT:-'10002'}
def "API_LISTEN_IP" "0.0.0.0" # API的监听IP

###################### openim-chat 配置信息 ######################
def "OPENIM_CHAT_DATA_DIR" "./openim-chat/${CHAT_BRANCH}"
def "OPENIM_CHAT_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}" # OpenIM服务地址
def "OPENIM_CHAT_API_PORT" "10008"                   # OpenIM API端口
def "CHAT_API_LISTEN_IP" ""                          # OpenIM API的监听IP

def "OPENIM_ADMIN_API_PORT" "10009" # OpenIM Admin API端口
def "ADMIN_API_LISTEN_IP" ""        # OpenIM Admin API的监听IP

def "OPENIM_ADMIN_PORT" "30200" # OpenIM chat Admin端口
def "OPENIM_CHAT_PORT" "30300"  # OpenIM chat Admin的监听IP

def "OPENIM_ADMIN_NAME" "admin" # openim-chat Admin用户名
def "OPENIM_CHAT_NAME" "chat"   # openim-chat chat用户名

# TODO 注意： 一般的配置都可以使用 def 函数来定义，如果是包含特殊字符，比如说:
# TODO readonly MSG_DESTRUCT_TIME=${MSG_DESTRUCT_TIME:-'0 2 * * *'}
# TODO 使用 readonly 来定义合适，负责无法正常解析, 并且 yaml 模板需要加 "" 来包裹
###################### Env 配置信息 ######################
def "ENVS_DISCOVERY" "zookeeper"

###################### Zookeeper 配置信息 ######################
def "ZOOKEEPER_SCHEMA" "openim"                    # Zookeeper的模式
def "ZOOKEEPER_PORT" "12181"                       # Zookeeper的端口
def "ZOOKEEPER_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}" # Zookeeper的地址
def "ZOOKEEPER_USERNAME" ""                        # Zookeeper的用户名
def "ZOOKEEPER_PASSWORD" ""                        # Zookeeper的密码

###################### MySQL 配置信息 ######################
def "MYSQL_PORT" "13306"                       # MySQL的端口
def "MYSQL_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}" # MySQL的地址
def "MYSQL_USERNAME" "${OPENIM_USER}"                 # MySQL的用户名
# MySQL的密码
readonly MYSQL_PASSWORD=${MYSQL_PASSWORD:-"${PASSWORD}"}
def "MYSQL_DATABASE" "${DATABASE_NAME}"        # MySQL的数据库名
def "MYSQL_MAX_OPEN_CONN" "1000"               # 最大打开的连接数
def "MYSQL_MAX_IDLE_CONN" "100"                # 最大空闲连接数
def "MYSQL_MAX_LIFETIME" "60"                  # 连接可以重用的最大生命周期（秒）
def "MYSQL_LOG_LEVEL" "4"                      # 日志级别
def "MYSQL_SLOW_THRESHOLD" "500"               # 慢查询阈值（毫秒）

###################### MongoDB 配置信息 ######################
def "MONGO_URI"                                # MongoDB的URI
def "MONGO_PORT" "37017"                       # MongoDB的端口
def "MONGO_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}" # MongoDB的地址
def "MONGO_DATABASE" "${DATABASE_NAME}"        # MongoDB的数据库名
def "MONGO_USERNAME" "${OPENIM_USER}"                 # MongoDB的用户名
# MongoDB的密码
readonly MONGO_PASSWORD=${MONGO_PASSWORD:-"${PASSWORD}"}
def "MONGO_MAX_POOL_SIZE" "100"                # 最大连接池大小

###################### Object 配置信息 ######################
# app要能访问到此ip和端口或域名
readonly API_URL=${API_URL:-"http://${OPENIM_IP}:${API_OPENIM_PORT}"}

def "OBJECT_ENABLE" "minio" # 对象是否启用
# 对象的API地址
readonly OBJECT_APIURL=${OBJECT_APIURL:-"${API_URL}"}
def "MINIO_BUCKET" "openim" # MinIO的存储桶名称
def "MINIO_PORT" "10005"    # MinIO的端口
# MinIO的端点URL
def MINIO_ADDRESS "${DOCKER_BRIDGE_GATEWAY}"
readonly MINIO_ENDPOINT=${MINIO_ENDPOINT:-"http://${MINIO_ADDRESS}:${MINIO_PORT}"}
def "MINIO_ACCESS_KEY" "${OPENIM_USER}"                                                  # MinIO的访问密钥ID
readonly MINIO_SECRET_KEY=${MINIO_SECRET_KEY:-"${PASSWORD}"}
def "MINIO_SESSION_TOKEN"                                                         # MinIO的会话令牌
readonly MINIO_SIGN_ENDPOINT=${MINIO_SIGN_ENDPOINT:-"http://${OPENIM_IP}:${MINIO_PORT}"} # signEndpoint为minio公网地址
def "MINIO_PUBLIC_READ" "false"                                                   # 公有读

# 腾讯云COS的存储桶URL
def "COS_BUCKET_URL" "https://temp-1252357374.cos.ap-chengdu.myqcloud.com"
def "COS_SECRET_ID"                                                     # 腾讯云COS的密钥ID
def "COS_SECRET_KEY"                                                    # 腾讯云COS的密钥
def "COS_SESSION_TOKEN"                                                 # 腾讯云COS的会话令牌
def "COS_PUBLIC_READ" "false"                                           # 公有读
def "OSS_ENDPOINT" "https://oss-cn-chengdu.aliyuncs.com"                # 阿里云OSS的端点URL
def "OSS_BUCKET" "demo-9999999"                                         # 阿里云OSS的存储桶名称
def "OSS_BUCKET_URL" "https://demo-9999999.oss-cn-chengdu.aliyuncs.com" # 阿里云OSS的存储桶URL
def "OSS_ACCESS_KEY_ID"                                                 # 阿里云OSS的访问密钥ID
def "OSS_ACCESS_KEY_SECRET"                                             # 阿里云OSS的密钥
def "OSS_SESSION_TOKEN"                                                 # 阿里云OSS的会话令牌
def "OSS_PUBLIC_READ" "false"                                           # 公有读

###################### Redis 配置信息 ######################
def "REDIS_PORT" "16379"                                    # Redis的端口
def "REDIS_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}"              # Redis的地址
def "REDIS_USERNAME"                                        # Redis的用户名
readonly REDIS_PASSWORD=${REDIS_PASSWORD:-"${PASSWORD}"}

###################### Kafka 配置信息 ######################
def "KAFKA_USERNAME"                                        # `Kafka` 的用户名
def "KAFKA_PASSWORD"                                        # `Kafka` 的密码
def "KAFKA_PORT" "19094"                                    # `Kafka` 的端口
def "KAFKA_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}"              # `Kafka` 的地址
def "KAFKA_LATESTMSG_REDIS_TOPIC" "latestMsgToRedis"        # `Kafka` 的最新消息到Redis的主题
def "KAFKA_OFFLINEMSG_MONGO_TOPIC" "offlineMsgToMongoMysql" # `Kafka` 的离线消息到Mongo的主题
def "KAFKA_MSG_PUSH_TOPIC" "msgToPush"                      # `Kafka` 的消息到推送的主题
def "KAFKA_CONSUMERGROUPID_REDIS" "redis"                   # `Kafka` 的消费组ID到Redis
def "KAFKA_CONSUMERGROUPID_MONGO" "mongo"                   # `Kafka` 的消费组ID到Mongo
def "KAFKA_CONSUMERGROUPID_MYSQL" "mysql"                   # `Kafka` 的消费组ID到MySql
def "KAFKA_CONSUMERGROUPID_PUSH" "push"                     # `Kafka` 的消费组ID到推送

###################### openim-web 配置信息 ######################
def "OPENIM_WEB_PORT" "11001"                       # openim-web的端口
def "OPENIM_WEB_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}" # openim-web的地址
def "OPENIM_WEB_DIST_PATH" "/app/dist"              # openim-web的dist路径

###################### RPC 配置信息 ######################
def "RPC_REGISTER_IP"                               # RPC的注册IP
def "RPC_LISTEN_IP" "0.0.0.0"                       # RPC的监听IP

###################### prometheus 配置 ######################
def "PROMETHEUS_PORT" "19090"                       # Prometheus的端口
def "PROMETHEUS_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}" # Prometheus的地址

###################### Grafana 配置信息 ######################
def "GRAFANA_PORT" "3000"                        # Grafana的端口
def "GRAFANA_ADDRESS" "${DOCKER_BRIDGE_GATEWAY}" # Grafana的地址

###################### RPC Port Configuration Variables ######################
# For launching multiple programs, just fill in multiple ports separated by commas
# For example:
# readonly OPENIM_USER_PORT=${OPENIM_USER_PORT:-'10110, 10111, 10112'} #Try not to have Spaces

# OpenIM用户服务端口
readonly OPENIM_USER_PORT=${OPENIM_USER_PORT:-'10110'}
# OpenIM朋友服务端口
readonly OPENIM_FRIEND_PORT=${OPENIM_FRIEND_PORT:-'10120'}
# OpenIM消息服务端口
readonly OPENIM_MESSAGE_PORT=${OPENIM_MESSAGE_PORT:-'10130'}
# OpenIM消息网关服务端口
readonly OPENIM_MESSAGE_GATEWAY_PORT=${OPENIM_MESSAGE_GATEWAY_PORT:-'10140'}
# OpenIM组服务端口
readonly OPENIM_GROUP_PORT=${OPENIM_GROUP_PORT:-'10150'}
# OpenIM授权服务端口
readonly OPENIM_AUTH_PORT=${OPENIM_AUTH_PORT:-'10160'}
# OpenIM推送服务端口
readonly OPENIM_PUSH_PORT=${OPENIM_PUSH_PORT:-'10170'}
# OpenIM对话服务端口
readonly OPENIM_CONVERSATION_PORT=${OPENIM_CONVERSATION_PORT:-'10180'}
# OpenIM第三方服务端口
readonly OPENIM_THIRD_PORT=${OPENIM_THIRD_PORT:-'10190'}

###################### RPC Register Name Variables ######################
def "OPENIM_USER_NAME" "User"                       # OpenIM用户服务名称
def "OPENIM_FRIEND_NAME" "Friend"                   # OpenIM朋友服务名称
def "OPENIM_MSG_NAME" "Msg"                         # OpenIM消息服务名称
def "OPENIM_PUSH_NAME" "Push"                       # OpenIM推送服务名称
def "OPENIM_MESSAGE_GATEWAY_NAME" "MessageGateway"  # OpenIM消息网关服务名称
def "OPENIM_GROUP_NAME" "Group"                     # OpenIM组服务名称
def "OPENIM_AUTH_NAME" "Auth"                       # OpenIM授权服务名称
def "OPENIM_CONVERSATION_NAME" "Conversation"       # OpenIM对话服务名称
def "OPENIM_THIRD_NAME" "Third"                     # OpenIM第三方服务名称

###################### Log Configuration Variables ######################
def "LOG_STORAGE_LOCATION" ""${OPENIM_ROOT}"/logs/" # 日志存储位置
def "LOG_ROTATION_TIME" "24"                        # 日志轮替时间
def "LOG_REMAIN_ROTATION_COUNT" "2"                 # 保留的日志轮替数量
def "LOG_REMAIN_LOG_LEVEL" "6"                      # 保留的日志级别
def "LOG_IS_STDOUT" "false"                         # 是否将日志输出到标准输出
def "LOG_IS_JSON" "false"                           # 日志是否为JSON格式
def "LOG_WITH_STACK" "false"                        # 日志是否带有堆栈信息

###################### Variables definition ######################
def "WEBSOCKET_MAX_CONN_NUM" "100000" # Websocket最大连接数
def "WEBSOCKET_MAX_MSG_LEN" "4096"    # Websocket最大消息长度
def "WEBSOCKET_TIMEOUT" "10"          # Websocket超时
def "PUSH_ENABLE" "getui"             # 推送是否启用
# GeTui推送URL
readonly GETUI_PUSH_URL=${GETUI_PUSH_URL:-'https://restapi.getui.com/v2/$appId'}
def "GETUI_MASTER_SECRET" ""          # GeTui主密钥
def "GETUI_APP_KEY" ""                # GeTui应用密钥
def "GETUI_INTENT" ""                 # GeTui推送意图
def "GETUI_CHANNEL_ID" ""             # GeTui渠道ID
def "GETUI_CHANNEL_NAME" ""           # GeTui渠道名称
def "FCM_SERVICE_ACCOUNT" "x.json"    # FCM服务账户
def "JPNS_APP_KEY" ""                 # JPNS应用密钥
def "JPNS_MASTER_SECRET" ""           # JPNS主密钥
def "JPNS_PUSH_URL" ""                # JPNS推送URL
def "JPNS_PUSH_INTENT" ""             # JPNS推送意图
def "MANAGER_USERID_1" "openIM123456" # 管理员ID 1
def "MANAGER_USERID_2" "openIM654321" # 管理员ID 2
def "MANAGER_USERID_3" "openIMAdmin"  # 管理员ID 3
def "NICKNAME_1" "system1"            # 昵称1
def "NICKNAME_2" "system2"            # 昵称2
def "NICKNAME_3" "system3"            # 昵称3
def "MULTILOGIN_POLICY" "1"           # 多登录策略
def "CHAT_PERSISTENCE_MYSQL" "true"   # 聊天持久化MySQL
def "MSG_CACHE_TIMEOUT" "86400"       # 消息缓存超时
def "GROUP_MSG_READ_RECEIPT" "true"   # 群消息已读回执启用
def "SINGLE_MSG_READ_RECEIPT" "true"  # 单一消息已读回执启用
def "RETAIN_CHAT_RECORDS" "365"       # 保留聊天记录
# 聊天记录清理时间
readonly CHAT_RECORDS_CLEAR_TIME=${CHAT_RECORDS_CLEAR_TIME:-'0 2 * * 3'}
# 消息销毁时间
readonly MSG_DESTRUCT_TIME=${MSG_DESTRUCT_TIME:-'0 2 * * *'}
# 密钥
readonly SECRET=${SECRET:-"${PASSWORD}"}
def "TOKEN_EXPIRE" "90"         # Token到期时间
def "FRIEND_VERIFY" "false"     # 朋友验证
def "IOS_PUSH_SOUND" "xxx"      # IOS推送声音
def "IOS_BADGE_COUNT" "true"    # IOS徽章计数
def "IOS_PRODUCTION" "false"    # IOS生产

###################### Prometheus 配置信息 ######################
def "PROMETHEUS_ENABLE" "false" # 是否启用 Prometheus
def "PROMETHEUS_URL" "/prometheus"
# Api 服务的 Prometheus 端口
readonly API_PROM_PORT=${API_PROM_PORT:-'20100'}
# User 服务的 Prometheus 端口
readonly USER_PROM_PORT=${USER_PROM_PORT:-'20110'}
# Friend 服务的 Prometheus 端口
readonly FRIEND_PROM_PORT=${FRIEND_PROM_PORT:-'20120'}
# Message 服务的 Prometheus 端口
readonly MESSAGE_PROM_PORT=${MESSAGE_PROM_PORT:-'20130'}
# Message Gateway 服务的 Prometheus 端口
readonly MSG_GATEWAY_PROM_PORT=${MSG_GATEWAY_PROM_PORT:-'20140'}
# Group 服务的 Prometheus 端口
readonly GROUP_PROM_PORT=${GROUP_PROM_PORT:-'20150'}
# Auth 服务的 Prometheus 端口
readonly AUTH_PROM_PORT=${AUTH_PROM_PORT:-'20160'}
# Push 服务的 Prometheus 端口
readonly PUSH_PROM_PORT=${PUSH_PROM_PORT:-'20170'}
# Conversation 服务的 Prometheus 端口
readonly CONVERSATION_PROM_PORT=${CONVERSATION_PROM_PORT:-'20230'}
# RTC 服务的 Prometheus 端口
readonly RTC_PROM_PORT=${RTC_PROM_PORT:-'21300'}
# Third 服务的 Prometheus 端口
readonly THIRD_PROM_PORT=${THIRD_PROM_PORT:-'21301'}

# Message Transfer 服务的 Prometheus 端口列表
readonly MSG_TRANSFER_PROM_PORT=${MSG_TRANSFER_PROM_PORT:-'21400, 21401, 21402, 21403'}

###################### OpenIM openim-api ######################
def "OPENIM_API_HOST" "127.0.0.1"
def "OPENIM_API_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-api" # OpenIM openim-api 二进制文件路径
def "OPENIM_API_CONFIG" ""${OPENIM_ROOT}"/config/"            # OpenIM openim-api 配置文件路径
def "OPENIM_API_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-api" # OpenIM openim-api 日志存储路径
def "OPENIM_API_LOG_LEVEL" "info"                             # OpenIM openim-api 日志级别
def "OPENIM_API_LOG_MAX_SIZE" "100"                           # OpenIM openim-api 日志最大大小（MB）
def "OPENIM_API_LOG_MAX_BACKUPS" "7"                          # OpenIM openim-api 日志最大备份数
def "OPENIM_API_LOG_MAX_AGE" "7"                              # OpenIM openim-api 日志最大保存时间（天）
def "OPENIM_API_LOG_COMPRESS" "false"                         # OpenIM openim-api 日志是否压缩
def "OPENIM_API_LOG_WITH_STACK" "${LOG_WITH_STACK}"           # OpenIM openim-api 日志是否带有堆栈信息

###################### OpenIM openim-cmdutils ######################
def "OPENIM_CMDUTILS_HOST" "127.0.0.1"
def "OPENIM_CMDUTILS_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-cmdutils" # OpenIM openim-cmdutils 二进制文件路径
def "OPENIM_CMDUTILS_CONFIG" ""${OPENIM_ROOT}"/config/"                 # OpenIM openim-cmdutils 配置文件路径
def "OPENIM_CMDUTILS_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-cmdutils" # OpenIM openim-cmdutils 日志存储路径
def "OPENIM_CMDUTILS_LOG_LEVEL" "info"                                  # OpenIM openim-cmdutils 日志级别
def "OPENIM_CMDUTILS_LOG_MAX_SIZE" "100"                                # OpenIM openim-cmdutils 日志最大大小（MB）
def "OPENIM_CMDUTILS_LOG_MAX_BACKUPS" "7"                               # OpenIM openim-cmdutils 日志最大备份数
def "OPENIM_CMDUTILS_LOG_MAX_AGE" "7"                                   # OpenIM openim-cmdutils 日志最大保存时间（天）
def "OPENIM_CMDUTILS_LOG_COMPRESS" "false"                              # OpenIM openim-cmdutils 日志是否压缩
def "OPENIM_CMDUTILS_LOG_WITH_STACK" "${LOG_WITH_STACK}"                # OpenIM openim-cmdutils 日志是否带有堆栈信息

###################### OpenIM openim-crontask ######################
def "OPENIM_CRONTASK_HOST" "127.0.0.1"
def "OPENIM_CRONTASK_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-crontask" # OpenIM openim-crontask 二进制文件路径
def "OPENIM_CRONTASK_CONFIG" ""${OPENIM_ROOT}"/config/"                 # OpenIM openim-crontask 配置文件路径
def "OPENIM_CRONTASK_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-crontask" # OpenIM openim-crontask 日志存储路径
def "OPENIM_CRONTASK_LOG_LEVEL" "info"                                  # OpenIM openim-crontask 日志级别
def "OPENIM_CRONTASK_LOG_MAX_SIZE" "100"                                # OpenIM openim-crontask 日志最大大小（MB）
def "OPENIM_CRONTASK_LOG_MAX_BACKUPS" "7"                               # OpenIM openim-crontask 日志最大备份数
def "OPENIM_CRONTASK_LOG_MAX_AGE" "7"                                   # OpenIM openim-crontask 日志最大保存时间（天）
def "OPENIM_CRONTASK_LOG_COMPRESS" "false"                              # OpenIM openim-crontask 日志是否压缩
def "OPENIM_CRONTASK_LOG_WITH_STACK" "${LOG_WITH_STACK}"                # OpenIM openim-crontask 日志是否带有堆栈信息

###################### OpenIM openim-msggateway ######################
def "OPENIM_MSGGATEWAY_HOST" "127.0.0.1"
def "OPENIM_MSGGATEWAY_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-msggateway"
def "OPENIM_MSGGATEWAY_CONFIG" ""${OPENIM_ROOT}"/config/"
def "OPENIM_MSGGATEWAY_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-msggateway"
def "OPENIM_MSGGATEWAY_LOG_LEVEL" "info"
def "OPENIM_MSGGATEWAY_LOG_MAX_SIZE" "100"
def "OPENIM_MSGGATEWAY_LOG_MAX_BACKUPS" "7"
def "OPENIM_MSGGATEWAY_LOG_MAX_AGE" "7"
def "OPENIM_MSGGATEWAY_LOG_COMPRESS" "false"
def "OPENIM_MSGGATEWAY_LOG_WITH_STACK" "${LOG_WITH_STACK}"

# 定义 openim-msggateway 的数量, 默认为 4
readonly OPENIM_MSGGATEWAY_NUM=${OPENIM_MSGGATEWAY_NUM:-'4'}

###################### OpenIM openim-msgtransfer ######################
def "OPENIM_MSGTRANSFER_HOST" "127.0.0.1"
def "OPENIM_MSGTRANSFER_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-msgtransfer" # OpenIM openim-msgtransfer 二进制文件路径
def "OPENIM_MSGTRANSFER_CONFIG" ""${OPENIM_ROOT}"/config/"                    # OpenIM openim-msgtransfer 配置文件路径
def "OPENIM_MSGTRANSFER_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-msgtransfer" # OpenIM openim-msgtransfer 日志存储路径
def "OPENIM_MSGTRANSFER_LOG_LEVEL" "info"                                     # OpenIM openim-msgtransfer 日志级别
def "OPENIM_MSGTRANSFER_LOG_MAX_SIZE" "100"                                   # OpenIM openim-msgtransfer 日志最大大小（MB）
def "OPENIM_MSGTRANSFER_LOG_MAX_BACKUPS" "7"                                  # OpenIM openim-msgtransfer 日志最大备份数
def "OPENIM_MSGTRANSFER_LOG_MAX_AGE" "7"                                      # OpenIM openim-msgtransfer 日志最大保存时间（天）
def "OPENIM_MSGTRANSFER_LOG_COMPRESS" "false"                                 # OpenIM openim-msgtransfer 日志是否压缩
def "OPENIM_MSGTRANSFER_LOG_WITH_STACK" "${LOG_WITH_STACK}"                   # OpenIM openim-msgtransfer 日志是否带有堆栈信息

###################### OpenIM openim-push ######################
def "OPENIM_PUSH_HOST" "127.0.0.1"
def "OPENIM_PUSH_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-push" # OpenIM openim-push 二进制文件路径
def "OPENIM_PUSH_CONFIG" ""${OPENIM_ROOT}"/config/"             # OpenIM openim-push 配置文件路径
def "OPENIM_PUSH_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-push" # OpenIM openim-push 日志存储路径
def "OPENIM_PUSH_LOG_LEVEL" "info"                              # OpenIM openim-push 日志级别
def "OPENIM_PUSH_LOG_MAX_SIZE" "100"                            # OpenIM openim-push 日志最大大小（MB）
def "OPENIM_PUSH_LOG_MAX_BACKUPS" "7"                           # OpenIM openim-push 日志最大备份数
def "OPENIM_PUSH_LOG_MAX_AGE" "7"                               # OpenIM openim-push 日志最大保存时间（天）
def "OPENIM_PUSH_LOG_COMPRESS" "false"                          # OpenIM openim-push 日志是否压缩
def "OPENIM_PUSH_LOG_WITH_STACK" "${LOG_WITH_STACK}"            # OpenIM openim-push 日志是否带有堆栈信息

###################### OpenIM openim-rpc-auth ######################
def "OPENIM_RPC_AUTH_HOST" "127.0.0.1"
def "OPENIM_RPC_AUTH_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-rpc-auth" # OpenIM openim-rpc-auth 二进制文件路径
def "OPENIM_RPC_AUTH_CONFIG" ""${OPENIM_ROOT}"/config/"                 # OpenIM openim-rpc-auth 配置文件路径
def "OPENIM_RPC_AUTH_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-rpc-auth" # OpenIM openim-rpc-auth 日志存储路径
def "OPENIM_RPC_AUTH_LOG_LEVEL" "info"                                  # OpenIM openim-rpc-auth 日志级别
def "OPENIM_RPC_AUTH_LOG_MAX_SIZE" "100"                                # OpenIM openim-rpc-auth 日志最大大小（MB）
def "OPENIM_RPC_AUTH_LOG_MAX_BACKUPS" "7"                               # OpenIM openim-rpc-auth 日志最大备份数
def "OPENIM_RPC_AUTH_LOG_MAX_AGE" "7"                                   # OpenIM openim-rpc-auth 日志最大保存时间（天）
def "OPENIM_RPC_AUTH_LOG_COMPRESS" "false"                              # OpenIM openim-rpc-auth 日志是否压缩
def "OPENIM_RPC_AUTH_LOG_WITH_STACK" "${LOG_WITH_STACK}"                # OpenIM openim-rpc-auth 日志是否带有堆栈信息

###################### OpenIM openim-rpc-conversation ######################
def "OPENIM_RPC_CONVERSATION_HOST" "127.0.0.1"
def "OPENIM_RPC_CONVERSATION_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-rpc-conversation" # OpenIM openim-rpc-conversation 二进制文件路径
def "OPENIM_RPC_CONVERSATION_CONFIG" ""${OPENIM_ROOT}"/config/"                         # OpenIM openim-rpc-conversation 配置文件路径
def "OPENIM_RPC_CONVERSATION_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-rpc-conversation" # OpenIM openim-rpc-conversation 日志存储路径
def "OPENIM_RPC_CONVERSATION_LOG_LEVEL" "info"                                          # OpenIM openim-rpc-conversation 日志级别
def "OPENIM_RPC_CONVERSATION_LOG_MAX_SIZE" "100"                                        # OpenIM openim-rpc-conversation 日志最大大小（MB）
def "OPENIM_RPC_CONVERSATION_LOG_MAX_BACKUPS" "7"                                       # OpenIM openim-rpc-conversation 日志最大备份数
def "OPENIM_RPC_CONVERSATION_LOG_MAX_AGE" "7"                                           # OpenIM openim-rpc-conversation 日志最大保存时间（天）
def "OPENIM_RPC_CONVERSATION_LOG_COMPRESS" "false"                                      # OpenIM openim-rpc-conversation 日志是否压缩
def "OPENIM_RPC_CONVERSATION_LOG_WITH_STACK" "${LOG_WITH_STACK}"                        # OpenIM openim-rpc-conversation 日志是否带有堆栈信息

###################### OpenIM openim-rpc-friend ######################
def "OPENIM_RPC_FRIEND_HOST" "127.0.0.1"
def "OPENIM_RPC_FRIEND_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-rpc-friend" # OpenIM openim-rpc-friend 二进制文件路径
def "OPENIM_RPC_FRIEND_CONFIG" ""${OPENIM_ROOT}"/config/"                   # OpenIM openim-rpc-friend 配置文件路径
def "OPENIM_RPC_FRIEND_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-rpc-friend" # OpenIM openim-rpc-friend 日志存储路径
def "OPENIM_RPC_FRIEND_LOG_LEVEL" "info"                                    # OpenIM openim-rpc-friend 日志级别
def "OPENIM_RPC_FRIEND_LOG_MAX_SIZE" "100"                                  # OpenIM openim-rpc-friend 日志最大大小（MB）
def "OPENIM_RPC_FRIEND_LOG_MAX_BACKUPS" "7"                                 # OpenIM openim-rpc-friend 日志最大备份数
def "OPENIM_RPC_FRIEND_LOG_MAX_AGE" "7"                                     # OpenIM openim-rpc-friend 日志最大保存时间（天）
def "OPENIM_RPC_FRIEND_LOG_COMPRESS" "false"                                # OpenIM openim-rpc-friend 日志是否压缩
def "OPENIM_RPC_FRIEND_LOG_WITH_STACK" "${LOG_WITH_STACK}"                  # OpenIM openim-rpc-friend 日志是否带有堆栈信息

###################### OpenIM openim-rpc-group ######################
def "OPENIM_RPC_GROUP_HOST" "127.0.0.1"
def "OPENIM_RPC_GROUP_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-rpc-group" # OpenIM openim-rpc-group 二进制文件路径
def "OPENIM_RPC_GROUP_CONFIG" ""${OPENIM_ROOT}"/config/"                  # OpenIM openim-rpc-group 配置文件路径
def "OPENIM_RPC_GROUP_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-rpc-group" # OpenIM openim-rpc-group 日志存储路径
def "OPENIM_RPC_GROUP_LOG_LEVEL" "info"                                   # OpenIM openim-rpc-group 日志级别
def "OPENIM_RPC_GROUP_LOG_MAX_SIZE" "100"                                 # OpenIM openim-rpc-group 日志最大大小（MB）
def "OPENIM_RPC_GROUP_LOG_MAX_BACKUPS" "7"                                # OpenIM openim-rpc-group 日志最大备份数
def "OPENIM_RPC_GROUP_LOG_MAX_AGE" "7"                                    # OpenIM openim-rpc-group 日志最大保存时间（天）
def "OPENIM_RPC_GROUP_LOG_COMPRESS" "false"                               # OpenIM openim-rpc-group 日志是否压缩
def "OPENIM_RPC_GROUP_LOG_WITH_STACK" "${LOG_WITH_STACK}"                 # OpenIM openim-rpc-group 日志是否带有堆栈信息

###################### OpenIM openim-rpc-msg ######################
def "OPENIM_RPC_MSG_HOST" "127.0.0.1"
def "OPENIM_RPC_MSG_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-rpc-msg" # OpenIM openim-rpc-msg 二进制文件路径
def "OPENIM_RPC_MSG_CONFIG" ""${OPENIM_ROOT}"/config/"                # OpenIM openim-rpc-msg 配置文件路径
def "OPENIM_RPC_MSG_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-rpc-msg" # OpenIM openim-rpc-msg 日志存储路径
def "OPENIM_RPC_MSG_LOG_LEVEL" "info"                                 # OpenIM openim-rpc-msg 日志级别
def "OPENIM_RPC_MSG_LOG_MAX_SIZE" "100"                               # OpenIM openim-rpc-msg 日志最大大小（MB）
def "OPENIM_RPC_MSG_LOG_MAX_BACKUPS" "7"                              # OpenIM openim-rpc-msg 日志最大备份数
def "OPENIM_RPC_MSG_LOG_MAX_AGE" "7"                                  # OpenIM openim-rpc-msg 日志最大保存时间（天）
def "OPENIM_RPC_MSG_LOG_COMPRESS" "false"                             # OpenIM openim-rpc-msg 日志是否压缩
def "OPENIM_RPC_MSG_LOG_WITH_STACK" "${LOG_WITH_STACK}"               # OpenIM openim-rpc-msg 日志是否带有堆栈信息

###################### OpenIM openim-rpc-third ######################
def "OPENIM_RPC_THIRD_HOST" "127.0.0.1"
def "OPENIM_RPC_THIRD_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-rpc-third" # OpenIM openim-rpc-third 二进制文件路径
def "OPENIM_RPC_THIRD_CONFIG" ""${OPENIM_ROOT}"/config/"                  # OpenIM openim-rpc-third 配置文件路径
def "OPENIM_RPC_THIRD_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-rpc-third" # OpenIM openim-rpc-third 日志存储路径
def "OPENIM_RPC_THIRD_LOG_LEVEL" "info"                                   # OpenIM openim-rpc-third 日志级别
def "OPENIM_RPC_THIRD_LOG_MAX_SIZE" "100"                                 # OpenIM openim-rpc-third 日志最大大小（MB）
def "OPENIM_RPC_THIRD_LOG_MAX_BACKUPS" "7"                                # OpenIM openim-rpc-third 日志最大备份数
def "OPENIM_RPC_THIRD_LOG_MAX_AGE" "7"                                    # OpenIM openim-rpc-third 日志最大保存时间（天）
def "OPENIM_RPC_THIRD_LOG_COMPRESS" "false"                               # OpenIM openim-rpc-third 日志是否压缩
def "OPENIM_RPC_THIRD_LOG_WITH_STACK" "${LOG_WITH_STACK}"                 # OpenIM openim-rpc-third 日志是否带有堆栈信息

###################### OpenIM openim-rpc-user ######################
def "OPENIM_RPC_USER_HOST" "127.0.0.1"
def "OPENIM_RPC_USER_BINARY" "${OPENIM_OUTPUT_HOSTBIN}/openim-rpc-user" # OpenIM openim-rpc-user 二进制文件路径
def "OPENIM_RPC_USER_CONFIG" ""${OPENIM_ROOT}"/config/"                 # OpenIM openim-rpc-user 配置文件路径
def "OPENIM_RPC_USER_LOG_DIR" "${LOG_STORAGE_LOCATION}/openim-rpc-user" # OpenIM openim-rpc-user 日志存储路径
def "OPENIM_RPC_USER_LOG_LEVEL" "info"                                  # OpenIM openim-rpc-user 日志级别
def "OPENIM_RPC_USER_LOG_MAX_SIZE" "100"                                # OpenIM openim-rpc-user 日志最大大小（MB）
def "OPENIM_RPC_USER_LOG_MAX_BACKUPS" "7"                               # OpenIM openim-rpc-user 日志最大备份数
def "OPENIM_RPC_USER_LOG_MAX_AGE" "7"                                   # OpenIM openim-rpc-user 日志最大保存时间（天）
def "OPENIM_RPC_USER_LOG_COMPRESS" "false"                              # OpenIM openim-rpc-user 日志是否压缩
def "OPENIM_RPC_USER_LOG_WITH_STACK" "${LOG_WITH_STACK}"                # OpenIM openim-rpc-user 日志是否带有堆栈信息

###################### 设计中...暂时不需要######################

# openimctl 配置
def "CONFIG_USER_USERNAME" "admin"
def "CONFIG_USER_PASSWORD" "Admin@2021"
def "CONFIG_USER_CLIENT_CERTIFICATE" "${HOME}/.openim/cert/admin.pem"
def "CONFIG_USER_CLIENT_KEY" "${HOME}/.openim/cert/admin-key.pem"
def "CONFIG_SERVER_ADDRESS" "${OPENIM_APISERVER_HOST}:${OPENIM_APISERVER_SECURE_BIND_PORT}"
def "CONFIG_SERVER_CERTIFICATE_AUTHORITY" "${CA_FILE}"
