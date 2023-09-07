// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/spf13/cobra"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type RootCmd struct {
	Command        cobra.Command
	Name           string
	port           int
	prometheusPort int
}

type CmdOpts struct {
	loggerPrefixName string
}

func WithCronTaskLogName() func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.loggerPrefixName = "OpenIM.CronTask.log.all"
	}
}

func WithLogName(logName string) func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.loggerPrefixName = logName
	}
}

func NewRootCmd(name string, opts ...func(*CmdOpts)) (rootCmd *RootCmd) {
	rootCmd = &RootCmd{Name: name}
	c := cobra.Command{
		Use:   "start openIM application",
		Short: fmt.Sprintf(`Start %s `, name),
		Long:  fmt.Sprintf(`Start %s `, name),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := rootCmd.getConfFromCmdAndInit(cmd); err != nil {
				panic(err)
			}
			cmdOpts := &CmdOpts{}
			for _, opt := range opts {
				opt(cmdOpts)
			}
			if cmdOpts.loggerPrefixName == "" {
				cmdOpts.loggerPrefixName = "OpenIM.log.all"
			}
			if err := log.InitFromConfig(cmdOpts.loggerPrefixName, name, config.Config.Log.RemainLogLevel, config.Config.Log.IsStdout, config.Config.Log.IsJson, config.Config.Log.StorageLocation, config.Config.Log.RemainRotationCount, config.Config.Log.RotationTime); err != nil {
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
	fmt.Println("configFolderPath:", configFolderPath)
	return config2.InitConfig(configFolderPath)
}

func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}

func (r *RootCmd) AddCommand(cmds ...*cobra.Command) {
	r.Command.AddCommand(cmds...)
}
