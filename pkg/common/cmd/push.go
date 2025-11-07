package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/push"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type PushRpcCmd struct {
	*RootCmd
	ctx        context.Context
	configMap  map[string]any
	pushConfig *push.Config
}

func NewPushRpcCmd() *PushRpcCmd {
	var pushConfig push.Config
	ret := &PushRpcCmd{pushConfig: &pushConfig}
	ret.configMap = map[string]any{
		config.OpenIMPushCfgFileName:    &pushConfig.RpcConfig,
		config.RedisConfigFileName:      &pushConfig.RedisConfig,
		config.MongodbConfigFileName:    &pushConfig.MongoConfig,
		config.KafkaConfigFileName:      &pushConfig.KafkaConfig,
		config.ShareFileName:            &pushConfig.Share,
		config.NotificationFileName:     &pushConfig.NotificationConfig,
		config.WebhooksConfigFileName:   &pushConfig.WebhooksConfig,
		config.LocalCacheConfigFileName: &pushConfig.LocalCacheConfig,
		config.DiscoveryConfigFilename:  &pushConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		ret.pushConfig.FcmConfigPath = config.Path(ret.ConfigPath())
		return ret.runE()
	}
	return ret
}

func (a *PushRpcCmd) Exec() error {
	return a.Execute()
}

func (a *PushRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.pushConfig.Discovery, &a.pushConfig.RpcConfig.CircuitBreaker, &a.pushConfig.RpcConfig.RateLimiter, &a.pushConfig.RpcConfig.Prometheus, a.pushConfig.RpcConfig.RPC.ListenIP,
		a.pushConfig.RpcConfig.RPC.RegisterIP, a.pushConfig.RpcConfig.RPC.AutoSetPorts, a.pushConfig.RpcConfig.RPC.Ports,
		a.Index(), a.pushConfig.Discovery.RpcService.Push, &a.pushConfig.NotificationConfig, a.pushConfig,
		[]string{
			a.pushConfig.RpcConfig.GetConfigFileName(),
			a.pushConfig.RedisConfig.GetConfigFileName(),
			a.pushConfig.KafkaConfig.GetConfigFileName(),
			a.pushConfig.NotificationConfig.GetConfigFileName(),
			a.pushConfig.Share.GetConfigFileName(),
			a.pushConfig.WebhooksConfig.GetConfigFileName(),
			a.pushConfig.LocalCacheConfig.GetConfigFileName(),
			a.pushConfig.Discovery.GetConfigFileName(),
		},
		[]string{
			a.pushConfig.Discovery.RpcService.MessageGateway,
		},
		push.Start)
}
