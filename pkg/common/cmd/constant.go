package cmd

import (
	"strings"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

var ConfigEnvPrefixMap map[string]string

func init() {
	ConfigEnvPrefixMap = make(map[string]string)
	fileNames := []string{
		config.FileName, config.NotificationFileName, config.ShareFileName, config.WebhooksConfigFileName,
		config.KafkaConfigFileName, config.RedisConfigFileName,
		config.MongodbConfigFileName, config.MinioConfigFileName, config.LogConfigFileName,
		config.OpenIMAPICfgFileName, config.OpenIMCronTaskCfgFileName, config.OpenIMMsgGatewayCfgFileName,
		config.OpenIMMsgTransferCfgFileName, config.OpenIMPushCfgFileName, config.OpenIMRPCAuthCfgFileName,
		config.OpenIMRPCConversationCfgFileName, config.OpenIMRPCFriendCfgFileName, config.OpenIMRPCGroupCfgFileName,
		config.OpenIMRPCMsgCfgFileName, config.OpenIMRPCThirdCfgFileName, config.OpenIMRPCUserCfgFileName, config.DiscoveryConfigFilename,
	}

	for _, fileName := range fileNames {
		envKey := strings.TrimSuffix(strings.TrimSuffix(fileName, ".yml"), ".yaml")
		envKey = "IMENV_" + envKey
		envKey = strings.ToUpper(strings.ReplaceAll(envKey, "-", "_"))
		ConfigEnvPrefixMap[fileName] = envKey
	}
}

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
