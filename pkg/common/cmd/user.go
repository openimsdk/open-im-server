package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/user"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type UserRpcCmd struct {
	*RootCmd
	ctx        context.Context
	configMap  map[string]any
	userConfig *user.Config
}

func NewUserRpcCmd() *UserRpcCmd {
	var userConfig user.Config
	ret := &UserRpcCmd{userConfig: &userConfig}
	ret.configMap = map[string]any{
		config.OpenIMRPCUserCfgFileName: &userConfig.RpcConfig,
		config.RedisConfigFileName:      &userConfig.RedisConfig,
		config.MongodbConfigFileName:    &userConfig.MongodbConfig,
		config.KafkaConfigFileName:      &userConfig.KafkaConfig,
		config.ShareFileName:            &userConfig.Share,
		config.NotificationFileName:     &userConfig.NotificationConfig,
		config.WebhooksConfigFileName:   &userConfig.WebhooksConfig,
		config.LocalCacheConfigFileName: &userConfig.LocalCacheConfig,
		config.DiscoveryConfigFilename:  &userConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *UserRpcCmd) Exec() error {
	return a.Execute()
}

func (a *UserRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.userConfig.Discovery, &a.userConfig.RpcConfig.CircuitBreaker, &a.userConfig.RpcConfig.RateLimiter, &a.userConfig.RpcConfig.Prometheus, a.userConfig.RpcConfig.RPC.ListenIP,
		a.userConfig.RpcConfig.RPC.RegisterIP, a.userConfig.RpcConfig.RPC.AutoSetPorts, a.userConfig.RpcConfig.RPC.Ports,
		a.Index(), a.userConfig.Discovery.RpcService.User, &a.userConfig.NotificationConfig, a.userConfig,
		[]string{
			a.userConfig.RpcConfig.GetConfigFileName(),
			a.userConfig.RedisConfig.GetConfigFileName(),
			a.userConfig.MongodbConfig.GetConfigFileName(),
			a.userConfig.KafkaConfig.GetConfigFileName(),
			a.userConfig.NotificationConfig.GetConfigFileName(),
			a.userConfig.Share.GetConfigFileName(),
			a.userConfig.WebhooksConfig.GetConfigFileName(),
			a.userConfig.LocalCacheConfig.GetConfigFileName(),
			a.userConfig.Discovery.GetConfigFileName(),
		}, nil,
		user.Start)
}
