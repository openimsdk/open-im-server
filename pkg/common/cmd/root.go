package cmd

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Command cobra.Command
	port    int

	prometheusPort int
}

func NewRootCmd() (rootCmd *RootCmd) {
	rootCmd = &RootCmd{}
	c := cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Long:  `Start the server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			rootCmd.port = rootCmd.getPortFlag(cmd)
			rootCmd.prometheusPort = rootCmd.getPrometheusPortFlag(cmd)
			return rootCmd.getConfFromCmdAndInit(cmd)
		},
	}
	rootCmd.Command = c
	rootCmd.init()
	return rootCmd
}

func (r *RootCmd) AddRunE(f func(cmd RootCmd) error) {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f(*r)
	}
}

func (r *RootCmd) AddRpc(f func(port, prometheusPort int) error) {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f(r.port, r.prometheusPort)
	}
}

func (r *RootCmd) init() {
	r.Command.Flags().StringP(constant.FlagConf, "c", "", "Path to config file folder")
}

func (r *RootCmd) AddPortFlag() {
	r.Command.Flags().IntP(constant.FlagPort, "p", 0, "server listen port")
}

func (r *RootCmd) getPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.FlagPort)
	return port
}

func (r *RootCmd) GetPortFlag() int {
	return r.port
}

func (r *RootCmd) AddPrometheusPortFlag() {
	r.Command.Flags().StringP(constant.PrometheusPort, "pp", "", "server listen port")
}

func (r *RootCmd) getPrometheusPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.PrometheusPort)
	return port
}

func (r *RootCmd) GetPrometheusPortFlag() int {
	return r.prometheusPort
}

func (r *RootCmd) getConfFromCmdAndInit(cmdLines *cobra.Command) error {
	configFolderPath, _ := cmdLines.Flags().GetString(constant.FlagConf)
	return config.InitConfig(configFolderPath)
}

func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}

func (r *RootCmd) AddCommand(cmds ...*cobra.Command) {
	r.Command.AddCommand(cmds...)
}
