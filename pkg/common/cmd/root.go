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

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/spf13/cobra"
)

// define a root command
type RootCmd struct {
	Command        cobra.Command
	Name           string
	port           int
	prometheusPort int
}

// create a new root command
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
			if err := log.InitFromConfig("OpenIM.log.all", name,
				config.Config.Log.RemainLogLevel, config.Config.Log.IsStdout,
				config.Config.Log.IsJson, config.Config.Log.StorageLocation,
				config.Config.Log.RemainRotationCount); err != nil {
				panic(err)
			}
			return nil
		},
	}
	rootCmd.Command = c
	rootCmd.addConfFlag()
	return rootCmd
}

// add config flag
func (r *RootCmd) addConfFlag() {
	r.Command.Flags().StringP(constant.FlagConf, "c", "", "Path to config file folder")
}

// add port flag
func (r *RootCmd) AddPortFlag() {
	r.Command.Flags().IntP(constant.FlagPort, "p", 0, "server listen port")
}

// get port flag
func (r *RootCmd) getPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.FlagPort)
	return port
}

// get port flag
func (r *RootCmd) GetPortFlag() int {
	return r.port
}

// add promethus port flag
func (r *RootCmd) AddPrometheusPortFlag() {
	r.Command.Flags().IntP(constant.FlagPrometheusPort, "", 0, "server prometheus listen port")
}

// // gey promethus port plag
func (r *RootCmd) getPrometheusPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.FlagPrometheusPort)
	return port
}

// gey promethus port plag
func (r *RootCmd) GetPrometheusPortFlag() int {
	return r.prometheusPort
}

// get config from cmd and init
func (r *RootCmd) getConfFromCmdAndInit(cmdLines *cobra.Command) error {
	configFolderPath, _ := cmdLines.Flags().GetString(constant.FlagConf)
	return config.InitConfig(configFolderPath)
}

// excute
func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}

// add command
func (r *RootCmd) AddCommand(cmds ...*cobra.Command) {
	r.Command.AddCommand(cmds...)
}
