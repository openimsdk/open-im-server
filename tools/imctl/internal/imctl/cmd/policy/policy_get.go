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

package policy

import (
	"bytes"
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam"
	"github.com/openim-sigs/component-base/pkg/json"
	metav1 "github.com/openim-sigs/component-base/pkg/meta/v1"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	getUsageStr = "get POLICY_NAME"
)

// GetOptions is an options struct to support get subcommands.
type GetOptions struct {
	Name string

	iamclient iam.IamInterface
	genericclioptions.IOStreams
}

var (
	getExample = templates.Examples(`
		# Display a policy resource
		iamctl policy get foo`)

	getUsageErrStr = fmt.Sprintf("expected '%s'.\nPOLICY_NAME is required arguments for the get command", getUsageStr)
)

// NewGetOptions returns an initialized GetOptions instance.
func NewGetOptions(ioStreams genericclioptions.IOStreams) *GetOptions {
	return &GetOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdGet returns new initialized instance of get sub command.
func NewCmdGet(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewGetOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   getUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Display a authorization policy resource",
		TraverseChildren:      true,
		Long:                  "Display a authorization policy resource.",
		Example:               getExample,
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
func (o *GetOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, getUsageErrStr)
	}

	o.Name = args[0]

	o.iamclient, err = f.IAMClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *GetOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a get subcommand using the specified options.
func (o *GetOptions) Run(args []string) error {
	policy, err := o.iamclient.APIV1().Policies().Get(context.TODO(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(policy.Policy); err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "%12s %s\n", color.RedString(policy.Name+":"), bf.String())

	return nil
}
