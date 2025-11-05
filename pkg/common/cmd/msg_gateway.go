package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/msggateway"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"

	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type MsgGatewayCmd struct {
	*RootCmd
	ctx              context.Context
	configMap        map[string]any
	msgGatewayConfig *msggateway.Config
}

func NewMsgGatewayCmd() *MsgGatewayCmd {
	var msgGatewayConfig msggateway.Config
	ret := &MsgGatewayCmd{msgGatewayConfig: &msgGatewayConfig}
	ret.configMap = map[string]any{
		config.OpenIMMsgGatewayCfgFileName: &msgGatewayConfig.MsgGateway,
		config.ShareFileName:               &msgGatewayConfig.Share,
		config.RedisConfigFileName:         &msgGatewayConfig.RedisConfig,
		config.WebhooksConfigFileName:      &msgGatewayConfig.WebhooksConfig,
		config.DiscoveryConfigFilename:     &msgGatewayConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (m *MsgGatewayCmd) Exec() error {
	return m.Execute()
}

func (m *MsgGatewayCmd) runE() error {
	m.msgGatewayConfig.Index = config.Index(m.Index())
	rpc := m.msgGatewayConfig.MsgGateway.RPC
	var prometheus config.Prometheus
	return startrpc.Start(
		m.ctx, &m.msgGatewayConfig.Discovery,
		&prometheus,
		rpc.ListenIP, rpc.RegisterIP,
		rpc.AutoSetPorts,
		rpc.Ports, int(m.msgGatewayConfig.Index),
		m.msgGatewayConfig.Discovery.RpcService.MessageGateway,
		nil,
		m.msgGatewayConfig,
		[]string{},
		[]string{},
		msggateway.Start,
	)
}
