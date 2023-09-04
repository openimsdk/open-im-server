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
	"github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam"
	"github.com/openim-sigs/component-base/pkg/json"
	metav1 "github.com/openim-sigs/component-base/pkg/meta/v1"
	"github.com/ory/ladon"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/cli/genericclioptions"
)

const (
	updateUsageStr = "update POLICY_NAME POLICY"
)

// UpdateOptions is an options struct to support update subcommands.
type UpdateOptions struct {
	Policy *v1.Policy

	iamclient iam.IamInterface
	genericclioptions.IOStreams
}

var (
	updateExample = templates.Examples(`
		# Update a authorization policy with new policy.
		iamctl policy update foo "{"description":"This is a updated policy","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}"`)

	updateUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nPOLICY_NAME and POLICY is required arguments for the update command",
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
		Use:                   updateUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Update a authorization policy resource",
		TraverseChildren:      true,
		Long:                  "Update a authorization policy resource.",
		Example:               updateExample,
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
func (o *UpdateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error

	if len(args) < 2 {
		return cmdutil.UsageErrorf(cmd, updateUsageErrStr)
	}

	var pol ladon.DefaultPolicy
	if err = json.Unmarshal([]byte(args[1]), &pol); err != nil {
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
	ret, err := o.iamclient.APIV1().Policies().Update(context.TODO(), o.Policy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "policy/%s updated\n", ret.Name)

	return nil
}
