package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/group"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/versionctx"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type GroupRpcCmd struct {
	*RootCmd
	ctx         context.Context
	configMap   map[string]any
	groupConfig *group.Config
}

func NewGroupRpcCmd() *GroupRpcCmd {
	var groupConfig group.Config
	ret := &GroupRpcCmd{groupConfig: &groupConfig}
	ret.configMap = map[string]any{
		config.OpenIMRPCGroupCfgFileName: &groupConfig.RpcConfig,
		config.RedisConfigFileName:       &groupConfig.RedisConfig,
		config.MongodbConfigFileName:     &groupConfig.MongodbConfig,
		config.ShareFileName:             &groupConfig.Share,
		config.NotificationFileName:      &groupConfig.NotificationConfig,
		config.WebhooksConfigFileName:    &groupConfig.WebhooksConfig,
		config.LocalCacheConfigFileName:  &groupConfig.LocalCacheConfig,
		config.DiscoveryConfigFilename:   &groupConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *GroupRpcCmd) Exec() error {
	return a.Execute()
}

func (a *GroupRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.groupConfig.Discovery, &a.groupConfig.RpcConfig.Prometheus, a.groupConfig.RpcConfig.RPC.ListenIP,
		a.groupConfig.RpcConfig.RPC.RegisterIP, a.groupConfig.RpcConfig.RPC.AutoSetPorts, a.groupConfig.RpcConfig.RPC.Ports,
		a.Index(), a.groupConfig.Discovery.RpcService.Group, &a.groupConfig.NotificationConfig, a.groupConfig,
		[]string{
			a.groupConfig.RpcConfig.GetConfigFileName(),
			a.groupConfig.RedisConfig.GetConfigFileName(),
			a.groupConfig.MongodbConfig.GetConfigFileName(),
			a.groupConfig.NotificationConfig.GetConfigFileName(),
			a.groupConfig.Share.GetConfigFileName(),
			a.groupConfig.WebhooksConfig.GetConfigFileName(),
			a.groupConfig.LocalCacheConfig.GetConfigFileName(),
			a.groupConfig.Discovery.GetConfigFileName(),
		}, nil,
		group.Start, versionctx.EnableVersionCtx())
}
