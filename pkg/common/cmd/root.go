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

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/spf13/cobra"
)

type RootCmdPt interface {
	GetPortFromConfig(portType string) int
}

type RootCmd struct {
	Command        cobra.Command
	Name           string
	port           int
	prometheusPort int
	cmdItf         RootCmdPt
	config         *config.GlobalConfig
}

func (rc *RootCmd) Port() int {
	return rc.port
}

type CmdOpts struct {
	loggerPrefixName string
}

func WithCronTaskLogName() func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.loggerPrefixName = "openim-crontask"
	}
}

func WithLogName(logName string) func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.loggerPrefixName = logName
	}
}

func NewRootCmd(name string, opts ...func(*CmdOpts)) *RootCmd {
	rootCmd := &RootCmd{Name: name, config: config.NewGlobalConfig()}
	cmd := cobra.Command{
		Use:   "Start openIM application",
		Short: fmt.Sprintf(`Start %s `, name),
		Long:  fmt.Sprintf(`Start %s `, name),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.persistentPreRun(cmd, opts...)
		},
	}
	rootCmd.Command = cmd
	rootCmd.addConfFlag()
	return rootCmd
}

func (rc *RootCmd) persistentPreRun(cmd *cobra.Command, opts ...func(*CmdOpts)) error {
	if err := rc.initializeConfiguration(cmd); err != nil {
		return fmt.Errorf("failed to get configuration from command: %w", err)
	}

	cmdOpts := rc.applyOptions(opts...)

	if err := rc.initializeLogger(cmdOpts); err != nil {
		return errs.Wrap(err, "failed to initialize logger")
	}

	return nil
}

func (rc *RootCmd) initializeConfiguration(cmd *cobra.Command) error {
	return rc.getConfFromCmdAndInit(cmd)
}

func (rc *RootCmd) applyOptions(opts ...func(*CmdOpts)) *CmdOpts {
	cmdOpts := defaultCmdOpts()
	for _, opt := range opts {
		opt(cmdOpts)
	}

	return cmdOpts
}

func (rc *RootCmd) initializeLogger(cmdOpts *CmdOpts) error {
	logConfig := rc.config.Log

	return log.InitFromConfig(

		cmdOpts.loggerPrefixName,
		rc.Name,
		logConfig.RemainLogLevel,
		logConfig.IsStdout,
		logConfig.IsJson,
		logConfig.StorageLocation,
		logConfig.RemainRotationCount,
		logConfig.RotationTime,
	)
}

func defaultCmdOpts() *CmdOpts {
	return &CmdOpts{
		loggerPrefixName: "openim-all",
	}
}

func (r *RootCmd) SetRootCmdPt(cmdItf RootCmdPt) {
	r.cmdItf = cmdItf
}

func (r *RootCmd) addConfFlag() {
	r.Command.Flags().StringP(constant.FlagConf, "c", "", "path to config file folder")
}

func (r *RootCmd) AddPortFlag() {
	r.Command.Flags().IntP(constant.FlagPort, "p", 0, "server listen port")
}

func (r *RootCmd) getPortFlag(cmd *cobra.Command) int {
	port, err := cmd.Flags().GetInt(constant.FlagPort)
	if err != nil {
		// Wrapping the error with additional context
		return 0
	}
	if port == 0 {
		port = r.PortFromConfig(constant.FlagPort)
	}
	return port
}

// // GetPortFlag returns the port flag.
func (r *RootCmd) GetPortFlag() int {
	return r.port
}

func (r *RootCmd) AddPrometheusPortFlag() {
	r.Command.Flags().IntP(constant.FlagPrometheusPort, "", 0, "server prometheus listen port")
}

func (r *RootCmd) getPrometheusPortFlag(cmd *cobra.Command) int {
	port, err := cmd.Flags().GetInt(constant.FlagPrometheusPort)
	if err != nil || port == 0 {
		port = r.PortFromConfig(constant.FlagPrometheusPort)
		if err != nil {
			return 0
		}
	}
	return port
}

func (r *RootCmd) GetPrometheusPortFlag() int {
	return r.prometheusPort
}

func (r *RootCmd) getConfFromCmdAndInit(cmdLines *cobra.Command) error {
	configFolderPath, _ := cmdLines.Flags().GetString(constant.FlagConf)
	fmt.Println("The directory of the configuration file to start the process:", configFolderPath)
	return config2.InitConfig(r.config, configFolderPath)
}

func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}

func (r *RootCmd) AddCommand(cmds ...*cobra.Command) {
	r.Command.AddCommand(cmds...)
}

func (r *RootCmd) PortFromConfig(portType string) int {
	// Retrieve the port and cache it
	port := r.cmdItf.GetPortFromConfig(portType)
	return port
}
