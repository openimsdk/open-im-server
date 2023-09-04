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

package user

import (
	"context"
	"fmt"

	v1 "github.com/marmotedu/api/apiserver/v1"
	apiclientv1 "github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam/apiserver/v1"
	metav1 "github.com/openim-sigs/component-base/pkg/meta/v1"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	createUsageStr = "create USERNAME PASSWORD EMAIL"
)

// CreateOptions is an options struct to support create subcommands.
type CreateOptions struct {
	Nickname string
	Phone    string

	User *v1.User

	Client apiclientv1.APIV1Interface
	genericclioptions.IOStreams
}

var (
	createLong = templates.LongDesc(`Create a user on iam platform.
If nickname not specified, username will be used.`)

	createExample = templates.Examples(`
		# Create user with given input
		iamctl user create foo Foo@2020 foo@foxmail.com

		# Create user wt 
		iamctl user create foo Foo@2020 foo@foxmail.com --phone=18128845xxx --nickname=colin`)

	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nUSERNAME, PASSWORD and EMAIL are required arguments for the create command",
		createUsageStr,
	)
)

// NewCreateOptions returns an initialized CreateOptions instance.
func NewCreateOptions(ioStreams genericclioptions.IOStreams) *CreateOptions {
	return &CreateOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdCreate returns new initialized instance of create sub command.
func NewCmdCreate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   createUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Create a user resource",
		TraverseChildren:      true,
		Long:                  createLong,
		Example:               createExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	// mark flag as deprecated
	cmd.Flags().StringVar(&o.Nickname, "nickname", o.Nickname, "The nickname of the user.")
	cmd.Flags().StringVar(&o.Phone, "phone", o.Phone, "The phone number of the user.")

	return cmd
}

// Complete completes all the required options.
func (o *CreateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) < 3 {
		return cmdutil.UsageErrorf(cmd, createUsageErrStr)
	}

	if o.Nickname == "" {
		o.Nickname = args[0]
	}

	o.User = &v1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: args[0],
		},
		Nickname: o.Nickname,
		Password: args[1],
		Email:    args[2],
		Phone:    o.Phone,
	}

	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	o.Client, err = apiclientv1.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if errs := o.User.Validate(); len(errs) != 0 {
		return errs.ToAggregate()
	}

	return nil
}

// Run executes a create subcommand using the specified options.
func (o *CreateOptions) Run(args []string) error {
	ret, err := o.Client.Users().Create(context.TODO(), o.User, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "user/%s created\n", ret.Name)

	return nil
}
