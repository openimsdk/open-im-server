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

	v1 "github.com/marmotedu/api/apiserver/v1"
	apiclientv1 "github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam/apiserver/v1"
	"github.com/openim-sigs/component-base/pkg/json"
	metav1 "github.com/openim-sigs/component-base/pkg/meta/v1"
	"github.com/ory/ladon"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/cli/genericclioptions"
)

const (
	createUsageStr = "create POLICY_NAME POLICY"
)

// CreateOptions is an options struct to support create subcommands.
type CreateOptions struct {
	Policy *v1.Policy

	Client apiclientv1.APIV1Interface
	genericclioptions.IOStreams
}

var (
	createExample = templates.Examples(`
		# Create a authorization policy
		iamctl policy create foo '{"description":"This is a updated policy","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}'`)

	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nPOLICY_NAME and POLICY are required arguments for the create command",
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
		Short:                 "Create a authorization policy resource",
		TraverseChildren:      true,
		Long:                  "Create a authorization policy resource.",
		Example:               createExample,
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
func (o *CreateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return cmdutil.UsageErrorf(cmd, createUsageErrStr)
	}

	var pol ladon.DefaultPolicy
	if err := json.Unmarshal([]byte(args[1]), &pol); err != nil {
		return err
	}

	o.Policy = &v1.Policy{
		ObjectMeta: metav1.ObjectMeta{
			Name: args[0],
		},
		Policy: v1.AuthzPolicy{
			DefaultPolicy: pol,
		},
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
	if errs := o.Policy.Validate(); len(errs) != 0 {
		return errs.ToAggregate()
	}

	return nil
}

// Run executes a create subcommand using the specified options.
func (o *CreateOptions) Run(args []string) error {
	ret, err := o.Client.Policies().Create(context.TODO(), o.Policy, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "policy/%s created\n", ret.Name)

	return nil
}
