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
	"time"

	"github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam"
	"github.com/olekukonko/tablewriter"
	metav1 "github.com/openim-sigs/component-base/pkg/meta/v1"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/cli/genericclioptions"
)

const (
	defaltLimit = 1000
)

// ListOptions is an options struct to support list subcommands.
type ListOptions struct {
	Offset int64
	Limit  int64

	iamclient iam.IamInterface
	genericclioptions.IOStreams
}

var listExample = templates.Examples(`
		# List all secrets
		iamctl secret list

		# List secrets with limit and offset 
		iamctl secret list --offset=0 --limit=5`)

// NewListOptions returns an initialized ListOptions instance.
func NewListOptions(ioStreams genericclioptions.IOStreams) *ListOptions {
	return &ListOptions{
		IOStreams: ioStreams,
		Offset:    0,
		Limit:     defaltLimit,
	}
}

// NewCmdList returns new initialized instance of list sub command.
func NewCmdList(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "list",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Display all secret resources",
		TraverseChildren:      true,
		Long:                  "Display all secret resources.",
		Example:               listExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().Int64VarP(&o.Offset, "offset", "o", o.Offset, "Specify the offset of the first row to be returned.")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "l", o.Limit, "Specify the amount records to be returned.")

	return cmd
}

// Complete completes all the required options.
func (o *ListOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.iamclient, err = f.IAMClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *ListOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *ListOptions) Run(args []string) error {
	secrets, err := o.iamclient.APIV1().Secrets().List(context.TODO(), metav1.ListOptions{
		Offset: &o.Offset,
		Limit:  &o.Limit,
	})
	if err != nil {
		return err
	}

	data := make([][]string, 0, len(secrets.Items))
	table := tablewriter.NewWriter(o.Out)

	for _, secret := range secrets.Items {
		data = append(data, []string{
			secret.Name,
			secret.SecretID,
			secret.SecretKey,
			time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05"),
			secret.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	table = setHeader(table)
	table = cmdutil.TableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()

	return nil
}
