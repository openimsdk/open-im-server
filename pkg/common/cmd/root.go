package cmd

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	log "OpenIM/pkg/common/logger"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Command        cobra.Command
	port           int
	prometheusPort int
}

func NewRootCmd() (rootCmd *RootCmd) {
	rootCmd = &RootCmd{}
	c := cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Long:  `Start the server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.InitFromConfig("newlog")
			return rootCmd.getConfFromCmdAndInit(cmd)
		},
	}
	rootCmd.Command = c
	rootCmd.addConfFlag()
	return rootCmd
}

func (r *RootCmd) SetDesc(use, short, long string) {
	r.Command.Use = use
	r.Command.Short = short
	r.Command.Long = long
}

func (r *RootCmd) addConfFlag() {
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
	r.Command.Flags().String(constant.FlagPrometheusPort, "", "server prometheus listen port")
}

func (r *RootCmd) getPrometheusPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.FlagPrometheusPort)
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
