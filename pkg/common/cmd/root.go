package cmd

import (
	"fmt"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Command        cobra.Command
	Name           string
	port           int
	prometheusPort int
}

func NewRootCmd(name string) (rootCmd *RootCmd) {
	rootCmd = &RootCmd{Name: name}
	c := cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf(`Start %s server`, name),
		Long:  fmt.Sprintf(`Start %s server`, name),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := rootCmd.getConfFromCmdAndInit(cmd); err != nil {
				panic(err)
			}
			if err := log.InitFromConfig("OpenIM.log.all", name, config.Config.Log.RemainLogLevel, config.Config.Log.IsStdout, config.Config.Log.IsJson, config.Config.Log.StorageLocation, config.Config.Log.RemainRotationCount); err != nil {
				panic(err)
			}
			return nil
		},
	}
	rootCmd.Command = c
	rootCmd.addConfFlag()
	return rootCmd
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
	r.Command.Flags().IntP(constant.FlagPrometheusPort, "", 0, "server prometheus listen port")
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
