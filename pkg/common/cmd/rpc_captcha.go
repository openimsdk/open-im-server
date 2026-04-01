package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/captcha"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type CaptchaRpcCmd struct {
	*RootCmd
	ctx           context.Context
	configMap     map[string]any
	captchaConfig *captcha.Config
}

func NewCaptchaRpcCmd() *CaptchaRpcCmd {
	var captchaConfig captcha.Config
	ret := &CaptchaRpcCmd{captchaConfig: &captchaConfig}
	ret.configMap = map[string]any{
		OpenIMRPCCaptchaCfgFileName: &captchaConfig.RpcConfig,
		MongodbConfigFileName:       &captchaConfig.MongodbConfig,
		ShareFileName:               &captchaConfig.Share,
		DiscoveryConfigFilename:     &captchaConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (c *CaptchaRpcCmd) Exec() error {
	return c.Execute()
}

func (c *CaptchaRpcCmd) runE() error {
	return startrpc.Start(c.ctx, &c.captchaConfig.Discovery, &c.captchaConfig.RpcConfig.Prometheus, c.captchaConfig.RpcConfig.RPC.ListenIP,
		c.captchaConfig.RpcConfig.RPC.RegisterIP, c.captchaConfig.RpcConfig.RPC.AutoSetPorts, c.captchaConfig.RpcConfig.RPC.Ports,
		c.Index(), c.captchaConfig.Share.RpcRegisterName.Captcha, &c.captchaConfig.Share, c.captchaConfig,
		nil,
		captcha.Start)
}
