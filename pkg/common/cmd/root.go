package cmd

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Command  cobra.Command
	port     int
	portFlag bool

	prometheusPort     int
	prometheusPortFlag bool
}

func NewRootCmd() RootCmd {
	c := cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Long:  `Start the server`,
	}
	rootCmd := RootCmd{}
	c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if rootCmd.portFlag {
			rootCmd.port = rootCmd.getPortFlag(cmd)
		}
		if rootCmd.prometheusPortFlag {
			rootCmd.prometheusPort = rootCmd.GetPrometheusPortFlag(cmd)
		}
		return rootCmd.getConfFromCmdAndInit(cmd)
	}
	rootCmd.init()
	rootCmd.Command = c
	return rootCmd
}

func (r RootCmd) AddRunE(f func(cmd RootCmd) error) {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f(r)
	}
}

func (r RootCmd) init() {
	r.Command.Flags().StringP(constant.FlagConf, "c", "", "Path to config file folder")
}

func (r RootCmd) AddPortFlag() {
	r.Command.Flags().StringP(constant.FlagPort, "p", "", "server listen port")
	r.portFlag = true
}

func (r RootCmd) getPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.FlagPort)
	return port
}

func (r RootCmd) GetPortFlag() int {
	return r.port
}

func (r RootCmd) AddPrometheusPortFlag() {
	r.Command.Flags().StringP(constant.PrometheusPort, "pp", "", "server listen port")
	r.prometheusPortFlag = true
}

func (r RootCmd) GetPrometheusPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.PrometheusPort)
	return port
}

func (r RootCmd) getConfFromCmdAndInit(cmdLines *cobra.Command) error {
	configFolderPath, _ := cmdLines.Flags().GetString(constant.FlagConf)
	return config.InitConfig(configFolderPath)
}

func (r RootCmd) Execute() error {
	return r.Command.Execute()
}
