package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/auth"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type AuthRpcCmd struct {
	*RootCmd
	ctx        context.Context
	configMap  map[string]any
	authConfig *auth.Config
}

func NewAuthRpcCmd() *AuthRpcCmd {
	var authConfig auth.Config
	ret := &AuthRpcCmd{authConfig: &authConfig}
	ret.configMap = map[string]any{
		config.OpenIMRPCAuthCfgFileName: &authConfig.RpcConfig,
		config.RedisConfigFileName:      &authConfig.RedisConfig,
		config.MongodbConfigFileName:    &authConfig.MongoConfig,
		config.ShareFileName:            &authConfig.Share,
		config.LocalCacheConfigFileName: &authConfig.LocalCacheConfig,
		config.DiscoveryConfigFilename:  &authConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}

	return ret
}

func (a *AuthRpcCmd) Exec() error {
	return a.Execute()
}

func (a *AuthRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.authConfig.Discovery, &a.authConfig.RpcConfig.CircuitBreaker, &a.authConfig.RpcConfig.RateLimiter, &a.authConfig.RpcConfig.Prometheus, a.authConfig.RpcConfig.RPC.ListenIP,
		a.authConfig.RpcConfig.RPC.RegisterIP, a.authConfig.RpcConfig.RPC.AutoSetPorts, a.authConfig.RpcConfig.RPC.Ports,
		a.Index(), a.authConfig.Discovery.RpcService.Auth, nil, a.authConfig,
		[]string{
			a.authConfig.RpcConfig.GetConfigFileName(),
			a.authConfig.Share.GetConfigFileName(),
			a.authConfig.RedisConfig.GetConfigFileName(),
			a.authConfig.Discovery.GetConfigFileName(),
		},
		[]string{
			a.authConfig.Discovery.RpcService.MessageGateway,
		},
		auth.Start)
}
