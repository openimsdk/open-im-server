package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/tools/cron"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type CronTaskCmd struct {
	*RootCmd
	ctx            context.Context
	configMap      map[string]any
	cronTaskConfig *cron.Config
}

func NewCronTaskCmd() *CronTaskCmd {
	var cronTaskConfig cron.Config
	ret := &CronTaskCmd{cronTaskConfig: &cronTaskConfig}
	ret.configMap = map[string]any{
		config.OpenIMCronTaskCfgFileName: &cronTaskConfig.CronTask,
		config.ShareFileName:             &cronTaskConfig.Share,
		config.DiscoveryConfigFilename:   &cronTaskConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *CronTaskCmd) Exec() error {
	return a.Execute()
}

func (a *CronTaskCmd) runE() error {
	var prometheus config.Prometheus
	return startrpc.Start(
		a.ctx, &a.cronTaskConfig.Discovery,
		&prometheus,
		"", "",
		true,
		nil, 0,
		"",
		nil,
		a.cronTaskConfig,
		[]string{},
		[]string{},
		cron.Start,
	)
}
