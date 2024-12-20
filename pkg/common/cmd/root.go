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

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"github.com/spf13/cobra"
)

type RootCmd struct {
	Command        cobra.Command
	processName    string
	port           int
	prometheusPort int
	log            config.Log
	index          int
	configPath     string
}

func (r *RootCmd) ConfigPath() string {
	return r.configPath
}

func (r *RootCmd) Index() int {
	return r.index
}

func (r *RootCmd) Port() int {
	return r.port
}

type CmdOpts struct {
	loggerPrefixName string
	configMap        map[string]any
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
func WithConfigMap(configMap map[string]any) func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.configMap = configMap
	}
}

func NewRootCmd(processName string, opts ...func(*CmdOpts)) *RootCmd {
	rootCmd := &RootCmd{processName: processName}
	cmd := cobra.Command{
		Use:  "Start openIM application",
		Long: fmt.Sprintf(`Start %s `, processName),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.persistentPreRun(cmd, opts...)
		},
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	cmd.Flags().StringP(FlagConf, "c", "", "path of config directory")
	cmd.Flags().IntP(FlagTransferIndex, "i", 0, "process startup sequence number")

	rootCmd.Command = cmd
	return rootCmd
}

func (r *RootCmd) persistentPreRun(cmd *cobra.Command, opts ...func(*CmdOpts)) error {
	cmdOpts := r.applyOptions(opts...)
	if err := r.initializeConfiguration(cmd, cmdOpts); err != nil {
		return err
	}

	if err := r.initializeLogger(cmdOpts); err != nil {
		return errs.WrapMsg(err, "failed to initialize logger")
	}

	return nil
}

func (r *RootCmd) initializeConfiguration(cmd *cobra.Command, opts *CmdOpts) error {
	configDirectory, _, err := r.getFlag(cmd)
	if err != nil {
		return err
	}

	runtimeEnv := runtimeenv.PrintRuntimeEnvironment()

	// Load common configuration file
	//opts.configMap[ShareFileName] = StructEnvPrefix{EnvPrefix: shareEnvPrefix, ConfigStruct: &r.share}
	for configFileName, configStruct := range opts.configMap {
		err := config.Load(configDirectory, configFileName, ConfigEnvPrefixMap[configFileName], runtimeEnv, configStruct)
		if err != nil {
			return err
		}
	}
	// Load common log configuration file
	return config.Load(configDirectory, config.LogConfigFileName, ConfigEnvPrefixMap[config.LogConfigFileName], runtimeEnv, &r.log)
}

func (r *RootCmd) applyOptions(opts ...func(*CmdOpts)) *CmdOpts {
	cmdOpts := defaultCmdOpts()
	for _, opt := range opts {
		opt(cmdOpts)
	}

	return cmdOpts
}

func (r *RootCmd) initializeLogger(cmdOpts *CmdOpts) error {
	err := log.InitLoggerFromConfig(

		cmdOpts.loggerPrefixName,
		r.processName,
		"", "",
		r.log.RemainLogLevel,
		r.log.IsStdout,
		r.log.IsJson,
		r.log.StorageLocation,
		r.log.RemainRotationCount,
		r.log.RotationTime,
		version.Version,
		r.log.IsSimplify,
	)
	if err != nil {
		return errs.Wrap(err)
	}
	return errs.Wrap(log.InitConsoleLogger(r.processName, r.log.RemainLogLevel, r.log.IsJson, version.Version))

}

func defaultCmdOpts() *CmdOpts {
	return &CmdOpts{
		loggerPrefixName: "openim-service-log",
	}
}

func (r *RootCmd) getFlag(cmd *cobra.Command) (string, int, error) {
	configDirectory, err := cmd.Flags().GetString(FlagConf)
	if err != nil {
		return "", 0, errs.Wrap(err)
	}
	r.configPath = configDirectory
	index, err := cmd.Flags().GetInt(FlagTransferIndex)
	if err != nil {
		return "", 0, errs.Wrap(err)
	}
	r.index = index
	return configDirectory, index, nil
}

func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}
