// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

const (
	FileName                         = "config.yaml"
	NotificationFileName             = "notification.yaml"
	ShareFileName                    = "share.yaml"
	WebhooksConfigFileName           = "webhooks.yml"
	KafkaConfigFileName              = "kafka.yml"
	RedisConfigFileName              = "redis.yml"
	ZookeeperConfigFileName          = "zookeeper.yml"
	MongodbConfigFileName            = "mongodb.yml"
	MinioConfigFileName              = "minio.yml"
	LogConfigFileName                = "log.yml"
	OpenIMAPICfgFileName             = "openim-api.yml"
	OpenIMCronTaskCfgFileName        = "openim-crontask.yml"
	OpenIMMsgGatewayCfgFileName      = "openim-msggateway.yml"
	OpenIMMsgTransferCfgFileName     = "openim-msgtransfer.yml"
	OpenIMPushCfgFileName            = "openim-push.yml"
	OpenIMRPCAuthCfgFileName         = "openim-rpc-auth.yml"
	OpenIMRPCConversationCfgFileName = "openim-rpc-conversation.yml"
	OpenIMRPCFriendCfgFileName       = "openim-rpc-friend.yml"
	OpenIMRPCGroupCfgFileName        = "openim-rpc-group.yml"
	OpenIMRPCMsgCfgFileName          = "openim-rpc-msg.yml"
	OpenIMRPCThirdCfgFileName        = "openim-rpc-third.yml"
	OpenIMRPCUserCfgFileName         = "openim-rpc-user.yml"
)

const (
	notificationEnvPrefix = "IMENV_NOTIFICATION"
	shareEnvPrefix        = "IMENV_SHARE"
	webhooksEnvPrefix     = "IMENV_WEBHOOKS"
	kafkaEnvPrefix        = "IMENV_KAFKA"
	redisEnvPrefix        = "IMENV_REDIS"
	zookeeperEnvPrefix    = "IMENV_ZOOKEEPER"
	mongodbEnvPrefix      = "IMENV_MONGODB"
	minioEnvPrefix        = "IMENV_MINIO"
	logEnvPrefix          = "IMENV_LOG"
	apiEnvPrefix          = "IMENV_OPENIM_API"
	cronTaskEnvPrefix     = "IMENV_OPENIM_CRONTASK"
	msgGatewayEnvPrefix   = "IMENV_OPENIM_MSGGATEWAY"
	msgTransferEnvPrefix  = "IMENV_OPENIM_MSGTRANSFER"
	pushEnvPrefix         = "IMENV_OPENIM_PUSH"
	authEnvPrefix         = "IMENV_OPENIM_RPC_AUTH"
	conversationEnvPrefix = "IMENV_OPENIM_RPC_CONVERSATION"
	friendEnvPrefix       = "IMENV_OPENIM_RPC_FRIEND"
	groupEnvPrefix        = "IMENV_OPENIM_RPC_GROUP"
	msgEnvPrefix          = "IMENV_OPENIM_RPC_MSG"
	thirdEnvPrefix        = "IMENV_OPENIM_RPC_THIRD"
	userEnvPrefix         = "IMENV_OPENIM_RPC_USER"
)

const (
	FlagConf = "config_folder_path"

	FlagTransferIndex = "index"
)
