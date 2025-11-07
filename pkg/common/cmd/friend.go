package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type FriendRpcCmd struct {
	*RootCmd
	ctx            context.Context
	configMap      map[string]any
	relationConfig *relation.Config
}

func NewFriendRpcCmd() *FriendRpcCmd {
	var relationConfig relation.Config
	ret := &FriendRpcCmd{relationConfig: &relationConfig}
	ret.configMap = map[string]any{
		config.OpenIMRPCFriendCfgFileName: &relationConfig.RpcConfig,
		config.RedisConfigFileName:        &relationConfig.RedisConfig,
		config.MongodbConfigFileName:      &relationConfig.MongodbConfig,
		config.ShareFileName:              &relationConfig.Share,
		config.NotificationFileName:       &relationConfig.NotificationConfig,
		config.WebhooksConfigFileName:     &relationConfig.WebhooksConfig,
		config.LocalCacheConfigFileName:   &relationConfig.LocalCacheConfig,
		config.DiscoveryConfigFilename:    &relationConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *FriendRpcCmd) Exec() error {
	return a.Execute()
}

func (a *FriendRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.relationConfig.Discovery, &a.relationConfig.RpcConfig.CircuitBreaker, &a.relationConfig.RpcConfig.RateLimiter, &a.relationConfig.RpcConfig.Prometheus, a.relationConfig.RpcConfig.RPC.ListenIP,
		a.relationConfig.RpcConfig.RPC.RegisterIP, a.relationConfig.RpcConfig.RPC.AutoSetPorts, a.relationConfig.RpcConfig.RPC.Ports,
		a.Index(), a.relationConfig.Discovery.RpcService.Friend, &a.relationConfig.NotificationConfig, a.relationConfig,
		[]string{
			a.relationConfig.RpcConfig.GetConfigFileName(),
			a.relationConfig.RedisConfig.GetConfigFileName(),
			a.relationConfig.MongodbConfig.GetConfigFileName(),
			a.relationConfig.NotificationConfig.GetConfigFileName(),
			a.relationConfig.Share.GetConfigFileName(),
			a.relationConfig.WebhooksConfig.GetConfigFileName(),
			a.relationConfig.LocalCacheConfig.GetConfigFileName(),
			a.relationConfig.Discovery.GetConfigFileName(),
		}, nil,
		relation.Start)
}
