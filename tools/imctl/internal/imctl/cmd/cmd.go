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
	"flag"
	"io"
	"os"

	cliflag "github.com/openimsdk/component-base/pkg/cli/flag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openimsdk/open-im-server/tools/imctl/pkg/cli/genericclioptions"

	cmdutil "github.com/openimsdk/open-im-server/tools/imctl/inernal/iamctl/cmd/util"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/color"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/completion"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/info"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/jwt"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/new"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/options"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/policy"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/secret"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/set"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/user"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/validate"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/cmd/version"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/util/templates"
	genericapiserver "github.com/openimsdk/open-im-server/tools/imctl/internal/pkg/server"
	"github.com/openimsdk/open-im-server/tools/imctl/pkg/cli/genericclioptions"
)

// NewDefaultIAMCtlCommand creates the `imctl` command with default arguments.
func NewDefaultIMCtlCommand() *cobra.Command {
	return NewIMCtlCommand(os.Stdin, os.Stdout, os.Stderr)
}

// NewIAMCtlCommand returns new initialized instance of 'imctl' root command.
func NewIMCtlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "imctl",
		Short: "imctl controls the IM platform",
		Long: templates.LongDesc(`
		imctl controls the IM platform, is the client side tool for IM platform.

		Find more information at:
			// TODO: add link to docs, from auto scripts and gendocs
			https://github.com/openimsdk/open-im-server/tree/main/docs`),
		Run: runHelp,
		// Hook before and after Run initialize and write profiles to disk,
		// respectively.
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return flushProfiling()
		},
	}

	flags := cmds.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	addProfilingFlags(flags)

	iamConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDeprecatedSecretFlag()
	iamConfigFlags.AddFlags(flags)
	matchVersionIAMConfigFlags := cmdutil.NewMatchVersionFlags(iamConfigFlags)
	matchVersionIAMConfigFlags.AddFlags(cmds.PersistentFlags())

	_ = viper.BindPFlags(cmds.PersistentFlags())
	cobra.OnInitialize(func() {
		genericapiserver.LoadConfig(viper.GetString(genericclioptions.FlagIAMConfig), "iamctl")
	})
	cmds.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	f := cmdutil.NewFactory(matchVersionIAMConfigFlags)

	// From this point and forward we get warnings on flags that contain "_" separators
	cmds.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands:",
			Commands: []*cobra.Command{
				info.NewCmdInfo(f, ioStreams),
				color.NewCmdColor(f, ioStreams),
				new.NewCmdNew(f, ioStreams),
				jwt.NewCmdJWT(f, ioStreams),
			},
		},
		{
			Message: "Identity and Access Management Commands:",
			Commands: []*cobra.Command{
				user.NewCmdUser(f, ioStreams),
				secret.NewCmdSecret(f, ioStreams),
				policy.NewCmdPolicy(f, ioStreams),
			},
		},
		{
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				validate.NewCmdValidate(f, ioStreams),
			},
		},
		{
			Message: "Settings Commands:",
			Commands: []*cobra.Command{
				set.NewCmdSet(f, ioStreams),
				completion.NewCmdCompletion(ioStreams.Out, ""),
			},
		},
	}
	groups.Add(cmds)

	filters := []string{"options"}
	templates.ActsAsRootCommand(cmds, filters, groups...)

	cmds.AddCommand(version.NewCmdVersion(f, ioStreams))
	cmds.AddCommand(options.NewCmdOptions(ioStreams.Out))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
