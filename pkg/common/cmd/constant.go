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
	RpcPushServer         = "push"
	RpcAuthServer         = "auth"
	RpcConversationServer = "conversation"
	RpcFriendServer       = "friend"
	RpcGroupServer        = "group"
	RpcMsgServer          = "msg"
	RpcThirdServer        = "third"
	RpcUserServer         = "user"
	ApiServer             = "api"
	CronTaskServer        = "cronTask"
	MsgGatewayServer      = "msgGateway"
	MsgTransferServer     = "msgTransfer"
)
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
	OpenIMMsgGatewayCfgFileName      = "openim-msg-gateway.yml"
	OpenIMMsgTransferCfgFileName     = "openim-msg-transfer.yml"
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
	notificationEnvPrefix = "openim-notification"
	shareEnvPrefix        = "openim-share"
	webhooksEnvPrefix     = "openim-webhooks"
	logEnvPrefix          = "openim-log"
	redisEnvPrefix        = "openim-redis"
	mongodbEnvPrefix      = "openim-mongodb"
	zoopkeeperEnvPrefix   = "openim-zookeeper"
	authEnvPrefix         = "openim-auth"
	conversationEnvPrefix = "openim-conversation"
	friendEnvPrefix       = "openim-friend"
	groupEnvPrefix        = "openim-group"
)

const (
	FlagConf = "config_folder_path"

	FlagTransferIndex = "index"
)
