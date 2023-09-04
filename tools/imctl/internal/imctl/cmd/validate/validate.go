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

// Package validate validate the basic environment for iamctl to run.
package validate

import (
	"fmt"
	"net"
	"net/url"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

// ValidateOptions is an options struct to support 'validate' sub command.
type ValidateOptions struct {
	genericclioptions.IOStreams
}

// ValidateInfo defines the validate information.
type ValidateInfo struct {
	ItemName string
	Status   string
	Message  string
}

var validateExample = templates.Examples(`
		# Validate the basic environment for iamctl to run
		iamctl validate`)

// NewValidateOptions returns an initialized ValidateOptions instance.
func NewValidateOptions(ioStreams genericclioptions.IOStreams) *ValidateOptions {
	return &ValidateOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdValidate returns new initialized instance of 'validate' sub command.
func NewCmdValidate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewValidateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "validate",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Validate the basic environment for iamctl to run",
		TraverseChildren:      true,
		Long:                  "Validate the basic environment for iamctl to run.",
		Example:               validateExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	return cmd
}

// Complete completes all the required options.
func (o *ValidateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *ValidateOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a validate sub command using the specified options.
func (o *ValidateOptions) Run(args []string) error {
	data := [][]string{}
	FAIL := color.RedString("FAIL")
	PASS := color.GreenString("PASS")
	validateInfo := ValidateInfo{}

	// check if can access db
	validateInfo.ItemName = "iam-apiserver"
	target, err := url.Parse(viper.GetString("server.address"))
	if err != nil {
		return err
	}
	_, err = net.Dial("tcp", target.Host)
	// defer client.Close()
	if err != nil {
		validateInfo.Status = FAIL
		validateInfo.Message = fmt.Sprintf("%v", err)
	} else {
		validateInfo.Status = PASS
		validateInfo.Message = ""
	}

	data = append(data, []string{validateInfo.ItemName, validateInfo.Status, validateInfo.Message})

	table := tablewriter.NewWriter(o.Out)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(iamctl.TableWidth)
	table.SetHeader([]string{"ValidateItem", "Result", "Message"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()

	return nil
}
