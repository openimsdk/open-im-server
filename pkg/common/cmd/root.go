package cmd

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Command cobra.Command
}

func NewRootCmd() RootCmd {
	c := cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Long:  `Start the server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return getConfFromCmdAndInit(cmd)
		},
	}
	c.Flags()
	rootCmd := RootCmd{Command: c}
	rootCmd.init()
	return rootCmd
}

func (r RootCmd) init() {
	r.Command.Flags().StringP(constant.FlagConf, "c", "", "Path to config file folder")
}

func getConfFromCmdAndInit(cmdLines *cobra.Command) error {
	configFolderPath, _ := cmdLines.Flags().GetString(constant.FlagConf)
	return config.InitConfig(configFolderPath)
}
