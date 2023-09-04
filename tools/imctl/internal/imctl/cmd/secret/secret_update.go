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

package secret

import (
	"context"
	"fmt"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam"
	metav1 "github.com/openim-sigs/component-base/pkg/meta/v1"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/cli/genericclioptions"
)

const (
	updateUsageStr = "update SECRET_NAME"
)

// UpdateOptions is an options struct to support update subcommands.
type UpdateOptions struct {
	Description string
	Expires     int64

	Secret *v1.Secret

	iamclient iam.IamInterface
	genericclioptions.IOStreams
}

var (
	updateExample = templates.Examples(`
		# Update a secret resource
		iamctl secret update foo --expires=4h --description="new description"`)

	updateUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nSECRET_NAME is required arguments for the update command",
		updateUsageStr,
	)
)

// NewUpdateOptions returns an initialized UpdateOptions instance.
func NewUpdateOptions(ioStreams genericclioptions.IOStreams) *UpdateOptions {
	return &UpdateOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdUpdate returns new initialized instance of update sub command.
func NewCmdUpdate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewUpdateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "update SECRET_NAME",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Update a secret resource",
		TraverseChildren:      true,
		Long:                  "Update a secret resource.",
		Example:               updateExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.Description, "description", o.Description, "The description of the secret.")
	cmd.Flags().Int64Var(&o.Expires, "expires", o.Expires, "The expires of the secret.")

	return cmd
}

// Complete completes all the required options.
func (o *UpdateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, updateUsageErrStr)
	}

	o.Secret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: args[0],
		},
		Description: o.Description,
		Expires:     o.Expires,
	}

	o.iamclient, err = f.IAMClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *UpdateOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a update subcommand using the specified options.
func (o *UpdateOptions) Run(args []string) error {
	secret, err := o.iamclient.APIV1().Secrets().Update(context.TODO(), o.Secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "secret/%s updated\n", secret.Name)

	return nil
}
