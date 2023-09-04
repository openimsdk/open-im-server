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

package jwt

import (
	"fmt"
	"regexp"

	"github.com/golang-jwt/jwt/v4"
	"github.com/openim-sigs/component-base/pkg/json"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	showUsageStr = "show TOKEN"
)

// ShowOptions is an options struct to support show subcommands.
type ShowOptions struct {
	Compact bool

	genericclioptions.IOStreams
}

var (
	showExample = templates.Examples(`
		# Show header and Claims for a JWT token
		iamctl jwt show XXX.XXX.XXX`)

	showUsageErrStr = fmt.Sprintf("expected '%s'.\nTOKEN is required arguments for the show command", showUsageStr)
)

// NewShowOptions returns an initialized ShowOptions instance.
func NewShowOptions(ioStreams genericclioptions.IOStreams) *ShowOptions {
	return &ShowOptions{
		Compact: false,

		IOStreams: ioStreams,
	}
}

// NewCmdShow returns new initialized instance of show sub command.
func NewCmdShow(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewShowOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   showUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Show header and claims for a JWT token",
		Long:                  "Show header and claims for a JWT token",
		TraverseChildren:      true,
		Example:               showExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmdutil.UsageErrorf(cmd, showUsageErrStr)
			}

			return nil
		},
	}

	// mark flag as deprecated
	cmd.Flags().BoolVar(&o.Compact, "compact", o.Compact, "output compact JSON.")

	return cmd
}

// Complete completes all the required options.
func (o *ShowOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *ShowOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a show subcommand using the specified options.
func (o *ShowOptions) Run(args []string) error {
	// get the token
	tokenData := []byte(args[0])

	// trim possible whitespace from token
	tokenData = regexp.MustCompile(`\s*$`).ReplaceAll(tokenData, []byte{})

	token, err := jwt.Parse(string(tokenData), nil)
	if token == nil {
		return fmt.Errorf("malformed token: %w", err)
	}

	// Print the token details
	fmt.Println("Header:")
	if err := printJSON(o.Compact, token.Header); err != nil {
		return fmt.Errorf("failed to output header: %w", err)
	}

	fmt.Println("Claims:")
	if err := printJSON(o.Compact, token.Claims); err != nil {
		return fmt.Errorf("failed to output claims: %w", err)
	}

	return nil
}

// printJSON print a json object in accordance with the prophecy (or the command line options).
func printJSON(compact bool, j interface{}) error {
	var out []byte
	var err error

	if !compact {
		out, err = json.MarshalIndent(j, "", "    ")
	} else {
		out, err = json.Marshal(j)
	}

	if err == nil {
		fmt.Println(string(out))
	}

	return err
}
