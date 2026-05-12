package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/redpacket"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type RedPacketRpcCmd struct {
	*RootCmd
	ctx             context.Context
	configMap       map[string]any
	redPacketConfig *redpacket.Config
}

func NewRedPacketRpcCmd() *RedPacketRpcCmd {
	var redPacketConfig redpacket.Config
	ret := &RedPacketRpcCmd{redPacketConfig: &redPacketConfig}
	ret.configMap = map[string]any{
		OpenIMRPCRedPacketCfgFileName: &redPacketConfig.RpcConfig,
		MongodbConfigFileName:         &redPacketConfig.MongodbConfig,
		ShareFileName:                 &redPacketConfig.Share,
		DiscoveryConfigFilename:       &redPacketConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (c *RedPacketRpcCmd) Exec() error {
	return c.Execute()
}

func (c *RedPacketRpcCmd) runE() error {
	return startrpc.Start(c.ctx, &c.redPacketConfig.Discovery, &c.redPacketConfig.RpcConfig.Prometheus, c.redPacketConfig.RpcConfig.RPC.ListenIP,
		c.redPacketConfig.RpcConfig.RPC.RegisterIP, c.redPacketConfig.RpcConfig.RPC.AutoSetPorts, c.redPacketConfig.RpcConfig.RPC.Ports,
		c.Index(), c.redPacketConfig.Share.RpcRegisterName.RedPacket, &c.redPacketConfig.Share, c.redPacketConfig,
		nil,
		redpacket.Start)
}
