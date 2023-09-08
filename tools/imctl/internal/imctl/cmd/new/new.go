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

// Package new used to generate demo command code.
// nolint: predeclared
package new

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/openim-sigs/component-base/pkg/util/fileutil"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/util/templates"
	"github.com/openimsdk/open-im-server/tools/imctl/pkg/cli/genericclioptions"
	cmdutil "github.com/openimsdk/open-im-server/tools/imctl/pkg/util"
)

const (
	newUsageStr = "new CMD_NAME | CMD_NAME CMD_DESCRIPTION"
)

var (
	newLong = templates.LongDesc(`Used to generate demo command source code file.

Can use this command generate a command template file, and do some modify based on your needs.
This can improve your R&D efficiency.`)

	newExample = templates.Examples(`
		# Create a default 'test' command file without a description
		imctl new test

		# Create a default 'test' command file in /tmp/
		imctl new test -d /tmp/

		# Create a default 'test' command file with a description
		imctl new test "This is a test command"

		# Create command 'test' with two subcommands
		imctl new -g test "This is a test command with two subcommands"`)

	newUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nat least CMD_NAME is a required argument for the new command",
		newUsageStr,
	)

	cmdTemplate = `package {{.CommandName}}

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/openimsdk/open-im-server/tools/imctl/pkg/util"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/util/templates"
	"github.com/openimsdk/open-im-server/tools/imctl/pkg/cli/genericclioptions"
)

const (
	{{.CommandName}}UsageStr    = "{{.CommandName}} USERNAME PASSWORD"
	maxStringLength = 17
)

// {{.CommandFunctionName}}Options is an options struct to support '{{.CommandName}}' sub command.
type {{.CommandFunctionName}}Options struct {
	// options
	StringOption      string
	StringSliceOption []string
	IntOption         int
	BoolOption        bool

	// args
	Username string
	Password string

	genericclioptions.IOStreams
}

var (
	{{.CommandName}}Long = templates.LongDesc({{.Dot}}A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.{{.Dot}})

	{{.CommandName}}Example = templates.Examples({{.Dot}}
		# Print all option values for {{.CommandName}} 
		imctl {{.CommandName}} marmotedu marmotedupass{{.Dot}})

	{{.CommandName}}UsageErrStr = fmt.Sprintf("expected '%s'.\nUSERNAME and PASSWORD are required arguments for the {{.CommandName}} command", {{.CommandName}}UsageStr)
)

// New{{.CommandFunctionName}}Options returns an initialized {{.CommandFunctionName}}Options instance.
func New{{.CommandFunctionName}}Options(ioStreams genericclioptions.IOStreams) *{{.CommandFunctionName}}Options {
	return &{{.CommandFunctionName}}Options{
		StringOption: "default",
		IOStreams:    ioStreams,
	}
}

// NewCmd{{.CommandFunctionName}} returns new initialized instance of '{{.CommandName}}' sub command.
func NewCmd{{.CommandFunctionName}}(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := New{{.CommandFunctionName}}Options(ioStreams)

	cmd := &cobra.Command{
		Use:                   {{.CommandName}}UsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "{{.CommandDescription}}",
		TraverseChildren:      true,
		Long:                  {{.CommandName}}Long,
		Example:               {{.CommandName}}Example,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
		Args: func(cmd *cobra.Command, args []string) error {
			// nolint: gomnd // no need
			if len(args) < 2 {
				return cmdutil.UsageErrorf(cmd, {{.CommandName}}UsageErrStr)
			}

			// if need args equal to zero, uncomment the following code
			/*
				if len(args) != 0 {
					return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
				}
			*/

			return nil
		},
	}

	// mark flag as deprecated
	_ = cmd.Flags().MarkDeprecated("deprecated-opt", "This flag is deprecated and will be removed in future.")
	cmd.Flags().StringVarP(&o.StringOption, "string", "", o.StringOption, "String option.")
	cmd.Flags().StringSliceVar(&o.StringSliceOption, "slice", o.StringSliceOption, "String slice option.")
	cmd.Flags().IntVarP(&o.IntOption, "int", "i", o.IntOption, "Int option.")
	cmd.Flags().BoolVarP(&o.BoolOption, "bool", "b", o.BoolOption, "Bool option.")

	return cmd
}

// Complete completes all the required options.
func (o *{{.CommandFunctionName}}Options) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if o.StringOption != "" {
		o.StringOption += "(complete)"
	}

	o.Username = args[0]
	o.Password = args[1]

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *{{.CommandFunctionName}}Options) Validate(cmd *cobra.Command, args []string) error {
	if len(o.StringOption) > maxStringLength {
		return cmdutil.UsageErrorf(cmd, "--string length must less than 18")
	}

	if o.IntOption < 0 {
		return cmdutil.UsageErrorf(cmd, "--int must be a positive integer: %v", o.IntOption)
	}

	return nil
}

// Run executes a {{.CommandName}} sub command using the specified options.
func (o *{{.CommandFunctionName}}Options) Run(args []string) error {
	fmt.Fprintf(o.Out, "The following is option values:\n")
	fmt.Fprintf(o.Out, "==> --string: %v\n==> --slice: %v\n==> --int: %v\n==> --bool: %v\n",
		o.StringOption, o.StringSliceOption, o.IntOption, o.BoolOption)

	fmt.Fprintf(o.Out, "\nThe following is args values:\n")
	fmt.Fprintf(o.Out, "==> username: %v\n==> password: %v\n", o.Username, o.Password)

	return nil
}
`

	maincmdTemplate = `package {{.CommandName}}

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cmdutil "github.com/openimsdk/open-im-server/tools/imctl/pkg/util"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/util/templates"
	"github.com/openimsdk/open-im-server/tools/imctl/pkg/cli/genericclioptions"
)

const maxStringLength = 17

var (
	{{.CommandName}}Long = templates.LongDesc({{.Dot}}
	Demo command.

	This commands show you how to implement a command with two sub commands.{{.Dot}})
)

// NewCmd{{.CommandFunctionName}} returns new initialized instance of '{{.CommandName}}' sub command.
func NewCmd{{.CommandFunctionName}}(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "{{.CommandName}} SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "{{.CommandDescription}}",
		Long:                  {{.CommandName}}Long,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.ErrOut),
	}

	// add subcommands
	cmd.AddCommand(NewCmdSubCmd1(f, ioStreams))
	cmd.AddCommand(NewCmdSubCmd2(f, ioStreams))

	// add persistent flags for '{{.CommandName}}'
	cmdutil.AddCleanFlags(cmd)

	// persistent flag, we can get the value in subcommand via {{.Dot}}viper.Get{{.Dot}}
	cmd.PersistentFlags().StringP("persistent", "p", "this is a persistent option", "Cobra persistent option.")

	// bind flags with viper
	viper.BindPFlag("persistent", cmd.PersistentFlags().Lookup("persistent"))

	return cmd
}
`
	subcmd1Template = `package {{.CommandName}}

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cmdutil "github.com/openimsdk/open-im-server/tools/imctl/pkg/util"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/util/templates"
	"github.com/openimsdk/open-im-server/tools/imctl/pkg/cli/genericclioptions"
)

const (
	subcmd1UsageStr = "subcmd1 USERNAME PASSWORD"
)

// SubCmd1Options is an options struct to support subcmd1 subcommands.
type SubCmd1Options struct {
	// options
	StringOption      string
	StringSliceOption []string
	IntOption         int
	BoolOption        bool
	PersistentOption  string

	// args
	Username string
	Password string

	genericclioptions.IOStreams
}

var (
	subcmd1Long = templates.LongDesc({{.Dot}}A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.{{.Dot}})

	subcmd1Example = templates.Examples({{.Dot}}
		# Print all option values for subcmd1
		imctl {{.CommandName}} subcmd1 marmotedu marmotedupass

		# Print all option values for subcmd1 with --persistent specified
		imctl {{.CommandName}} subcmd1 marmotedu marmotedupass --persistent="specified persistent option in command line"{{.Dot}})

	subcmd1UsageErrStr = fmt.Sprintf("expected '%s'.\nUSERNAME and PASSWORD are required arguments for the subcmd1 command", subcmd1UsageStr)
)

// NewSubCmd1Options returns an initialized SubCmd1Options instance.
func NewSubCmd1Options(ioStreams genericclioptions.IOStreams) *SubCmd1Options {
	return &SubCmd1Options{
		StringOption: "default",
		IOStreams:    ioStreams,
	}
}

// NewCmdSubCmd1 returns new initialized instance of subcmd1 sub command.
func NewCmdSubCmd1(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewSubCmd1Options(ioStreams)

	cmd := &cobra.Command{
		Use:                   "subcmd1 USERNAME PASSWORD",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"sub1"},
		Short:                 "A brief description of your command",
		TraverseChildren:      true,
		Long:                  subcmd1Long,
		Example:               subcmd1Example,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
		Args: func(cmd *cobra.Command, args []string) error {
			// nolint: gomnd // no need
			if len(args) < 2 {
				return cmdutil.UsageErrorf(cmd, subcmd1UsageErrStr)
			}

			return nil
		},
	}

	// mark flag as deprecated
	_ = cmd.Flags().MarkDeprecated("deprecated-opt", "This flag is deprecated and will be removed in future.")
	cmd.Flags().StringVarP(&o.StringOption, "string", "", o.StringOption, "String option.")
	cmd.Flags().StringSliceVar(&o.StringSliceOption, "slice", o.StringSliceOption, "String slice option.")
	cmd.Flags().IntVarP(&o.IntOption, "int", "i", o.IntOption, "Int option.")
	cmd.Flags().BoolVarP(&o.BoolOption, "bool", "b", o.BoolOption, "Bool option.")

	return cmd
}

// Complete completes all the required options.
func (o *SubCmd1Options) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if o.StringOption != "" {
		o.StringOption += "(complete)"
	}

	o.PersistentOption = viper.GetString("persistent")
	o.Username = args[0]
	o.Password = args[1]

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *SubCmd1Options) Validate(cmd *cobra.Command, args []string) error {
	if len(o.StringOption) > maxStringLength {
		return cmdutil.UsageErrorf(cmd, "--string length must less than 18")
	}

	if o.IntOption < 0 {
		return cmdutil.UsageErrorf(cmd, "--int must be a positive integer: %v", o.IntOption)
	}

	return nil
}

// Run executes a subcmd1 subcommand using the specified options.
func (o *SubCmd1Options) Run(args []string) error {
	fmt.Fprintf(o.Out, "The following is option values:\n")
	fmt.Fprintf(o.Out, "==> --string: %v\n==> --slice: %v\n==> --int: %v\n==> --bool: %v\n==> --persistent: %v\n",
		o.StringOption, o.StringSliceOption, o.IntOption, o.BoolOption, o.PersistentOption)

	fmt.Fprintf(o.Out, "\nThe following is args values:\n")
	fmt.Fprintf(o.Out, "==> username: %v\n==> password: %v\n", o.Username, o.Password)
	return nil
}
`
	subcmd2Template = `package {{.CommandName}}

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/openimsdk/open-im-server/tools/imctl/pkg/util"
	"github.com/openimsdk/open-im-server/tools/imctl/internal/imctl/util/templates"
	"github.com/openimsdk/open-im-server/tools/imctl/pkg/cli/genericclioptions"
)

// SubCmd2Options is an options struct to support subcmd2 subcommands.
type SubCmd2Options struct {
	StringOption      string
	StringSliceOption []string
	IntOption         int
	BoolOption        bool

	genericclioptions.IOStreams
}

var (
	subcmd2Long = templates.LongDesc({{.Dot}}A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.{{.Dot}})

	subcmd2Example = templates.Examples({{.Dot}}
		# Print all option values for subcmd2
		imctl {{.CommandName}} subcmd2{{.Dot}})
)

// NewSubCmd2Options returns an initialized SubCmd2Options instance.
func NewSubCmd2Options(ioStreams genericclioptions.IOStreams) *SubCmd2Options {
	return &SubCmd2Options{
		StringOption: "default",
		IOStreams:    ioStreams,
	}
}

// NewCmdSubCmd2 returns new initialized instance of subcmd2 sub command.
func NewCmdSubCmd2(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewSubCmd2Options(ioStreams)

	cmd := &cobra.Command{
		Use:                   "subcmd2",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"sub2"},
		Short:                 "A brief description of your command",
		TraverseChildren:      true,
		Long:                  subcmd2Long,
		Example:               subcmd2Example,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	// mark flag as deprecated
	cmd.Flags().StringVarP(&o.StringOption, "string", "", o.StringOption, "String option.")
	cmd.Flags().StringSliceVar(&o.StringSliceOption, "slice", o.StringSliceOption, "String slice option.")
	cmd.Flags().IntVarP(&o.IntOption, "int", "i", o.IntOption, "Int option.")
	cmd.Flags().BoolVarP(&o.BoolOption, "bool", "b", o.BoolOption, "Bool option.")

	return cmd
}

// Complete completes all the required options.
func (o *SubCmd2Options) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	if o.StringOption != "" {
		o.StringOption += "(complete)"
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *SubCmd2Options) Validate(cmd *cobra.Command, args []string) error {
	if len(o.StringOption) > maxStringLength {
		return cmdutil.UsageErrorf(cmd, "--string length must less than 18")
	}

	if o.IntOption < 0 {
		return cmdutil.UsageErrorf(cmd, "--int must be a positive integer: %v", o.IntOption)
	}

	return nil
}

// Run executes a subcmd2 subcommand using the specified options.
func (o *SubCmd2Options) Run(args []string) error {
	fmt.Fprintf(o.Out, "The following is option values:\n")
	fmt.Fprintf(o.Out, "==> --string: %v\n==> --slice: %v\n==> --int: %v\n==> --bool: %v\n",
		o.StringOption, o.StringSliceOption, o.IntOption, o.BoolOption)
	return nil
}
`
)

// NewOptions is an options struct to support 'new' sub command.
type NewOptions struct {
	Group  bool
	Outdir string

	// command template options, will render to command template
	CommandName         string
	CommandDescription  string
	CommandFunctionName string
	Dot                 string

	genericclioptions.IOStreams
}

// NewNewOptions returns an initialized NewOptions instance.
func NewNewOptions(ioStreams genericclioptions.IOStreams) *NewOptions {
	return &NewOptions{
		Group:              false,
		Outdir:             ".",
		CommandDescription: "A brief description of your command",
		Dot:                "`",
		IOStreams:          ioStreams,
	}
}

// NewCmdNew returns new initialized instance of 'new' sub command.
func NewCmdNew(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewNewOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   newUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate demo command code",
		Long:                  newLong,
		Example:               newExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate(cmd))
			cmdutil.CheckErr(o.Run(args))
		},
		Aliases:    []string{},
		SuggestFor: []string{},
	}

	cmd.Flags().BoolVarP(&o.Group, "group", "g", o.Group, "Generate two subcommands.")
	cmd.Flags().StringVarP(&o.Outdir, "outdir", "d", o.Outdir, "Where to create demo command files.")

	return cmd
}

// Complete completes all the required options.
func (o *NewOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, newUsageErrStr)
	}

	o.CommandName = strings.ToLower(args[0])
	if len(args) > 1 {
		o.CommandDescription = args[1]
	}

	o.CommandFunctionName = cases.Title(language.English).String(o.CommandName)

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *NewOptions) Validate(cmd *cobra.Command) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *NewOptions) Run(args []string) error {
	if o.Group {
		return o.CreateCommandWithSubCommands()
	}

	return o.CreateCommand()
}

// CreateCommand create the command with options.
func (o *NewOptions) CreateCommand() error {
	return o.GenerateGoCode(o.CommandName+".go", cmdTemplate)
}

// CreateCommandWithSubCommands create sub commands with options.
func (o *NewOptions) CreateCommandWithSubCommands() error {
	if err := o.GenerateGoCode(o.CommandName+".go", maincmdTemplate); err != nil {
		return err
	}

	if err := o.GenerateGoCode(o.CommandName+"_subcmd1.go", subcmd1Template); err != nil {
		return err
	}

	if err := o.GenerateGoCode(o.CommandName+"_subcmd2.go", subcmd2Template); err != nil {
		return err
	}

	return nil
}

// GenerateGoCode generate go source file.
func (o *NewOptions) GenerateGoCode(name, codeTemplate string) error {
	tmpl, err := template.New("cmd").Parse(codeTemplate)
	if err != nil {
		return err
	}

	err = fileutil.EnsureDirAll(o.Outdir)
	if err != nil {
		return err
	}

	filename := filepath.Join(o.Outdir, name)
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = tmpl.Execute(fd, o)
	if err != nil {
		return err
	}

	fmt.Printf("Command file generated: %s\n", filename)

	return nil
}
