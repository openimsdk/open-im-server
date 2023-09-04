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
	"context"
	"fmt"

	"github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam"
	metav1 "github.com/openim-sigs/component-base/pkg/meta/v1"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	deleteUsageStr = "delete POLICY_NAME"
)

// DeleteOptions is an options struct to support delete subcommands.
type DeleteOptions struct {
	Name string

	iamclient iam.IamInterface
	genericclioptions.IOStreams
}

var (
	deleteExample = templates.Examples(`
		# Delete a policy resource
		iamctl policy delete foo`)

	deleteUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nPOLICY_NAME is required arguments for the delete command",
		deleteUsageStr,
	)
)

// NewDeleteOptions returns an initialized DeleteOptions instance.
func NewDeleteOptions(ioStreams genericclioptions.IOStreams) *DeleteOptions {
	return &DeleteOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdDelete returns new initialized instance of delete sub command.
func NewCmdDelete(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewDeleteOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "delete",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Delete a authorization policy resource",
		TraverseChildren:      true,
		Long:                  "Delete a authorization policy resource.",
		Example:               deleteExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run())
		},
		SuggestFor: []string{},
	}

	return cmd
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, deleteUsageErrStr)
	}

	o.Name = args[0]

	o.iamclient, err = f.IAMClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *DeleteOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a delete subcommand using the specified options.
func (o *DeleteOptions) Run() error {
	if err := o.iamclient.APIV1().Policies().Delete(context.TODO(), o.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "policy/%s deleted\n", o.Name)

	return nil
}
