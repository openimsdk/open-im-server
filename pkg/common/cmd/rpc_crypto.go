package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/crypto"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type CryptoRpcCmd struct {
	*RootCmd
	ctx          context.Context
	configMap    map[string]any
	cryptoConfig *crypto.Config
}

func NewCryptoRpcCmd() *CryptoRpcCmd {
	var cryptoConfig crypto.Config
	ret := &CryptoRpcCmd{cryptoConfig: &cryptoConfig}
	ret.configMap = map[string]any{
		OpenIMRPCCryptoCfgFileName: &cryptoConfig.RpcConfig,
		MongodbConfigFileName:      &cryptoConfig.MongodbConfig,
		ShareFileName:              &cryptoConfig.Share,
		DiscoveryConfigFilename:    &cryptoConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (c *CryptoRpcCmd) Exec() error {
	return c.Execute()
}

func (c *CryptoRpcCmd) runE() error {
	return startrpc.Start(c.ctx, &c.cryptoConfig.Discovery, &c.cryptoConfig.RpcConfig.Prometheus, c.cryptoConfig.RpcConfig.RPC.ListenIP,
		c.cryptoConfig.RpcConfig.RPC.RegisterIP, c.cryptoConfig.RpcConfig.RPC.AutoSetPorts, c.cryptoConfig.RpcConfig.RPC.Ports,
		c.Index(), c.cryptoConfig.Share.RpcRegisterName.Crypto, &c.cryptoConfig.Share, c.cryptoConfig,
		nil,
		crypto.Start)
}
