package config

import "strings"

var EnvPrefixMap map[string]string

func init() {
	EnvPrefixMap = make(map[string]string)
	fileNames := []string{
		FileName, NotificationFileName, ShareFileName, WebhooksConfigFileName,
		KafkaConfigFileName, RedisConfigFileName,
		MongodbConfigFileName, MinioConfigFileName, LogConfigFileName,
		OpenIMAPICfgFileName, OpenIMCronTaskCfgFileName, OpenIMMsgGatewayCfgFileName,
		OpenIMMsgTransferCfgFileName, OpenIMPushCfgFileName, OpenIMRPCAuthCfgFileName,
		OpenIMRPCConversationCfgFileName, OpenIMRPCFriendCfgFileName, OpenIMRPCGroupCfgFileName,
		OpenIMRPCMsgCfgFileName, OpenIMRPCThirdCfgFileName, OpenIMRPCUserCfgFileName, DiscoveryConfigFilename,
	}

	for _, fileName := range fileNames {
		envKey := strings.TrimSuffix(strings.TrimSuffix(fileName, ".yml"), ".yaml")
		envKey = "IMENV_" + envKey
		envKey = strings.ToUpper(strings.ReplaceAll(envKey, "-", "_"))
		EnvPrefixMap[fileName] = envKey
	}
}

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
