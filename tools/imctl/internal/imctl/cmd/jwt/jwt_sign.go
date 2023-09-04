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
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/cli/genericclioptions"
	"github.com/marmotedu/iam/internal/pkg/middleware/auth"
)

const (
	signUsageStr = "sign SECRETID SECRETKEY"
)

// ErrSigningMethod defines invalid signing method error.
var ErrSigningMethod = errors.New("invalid signing method")

// SignOptions is an options struct to support sign subcommands.
type SignOptions struct {
	Timeout   time.Duration
	NotBefore time.Duration
	Algorithm string
	Audience  string
	Issuer    string
	Claims    ArgList
	Head      ArgList

	genericclioptions.IOStreams
}

var (
	signExample = templates.Examples(`
		# Sign a token with secretID and secretKey
		iamctl sign tgydj8d9EQSnFqKf iBdEdFNBLN1nR3fV

		# Sign a token with expires and sign method
		iamctl sign tgydj8d9EQSnFqKf iBdEdFNBLN1nR3fV --timeout=2h --algorithm=HS256`)

	signUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nSECRETID and SECRETKEY are required arguments for the sign command",
		signUsageStr,
	)
)

// NewSignOptions returns an initialized SignOptions instance.
func NewSignOptions(ioStreams genericclioptions.IOStreams) *SignOptions {
	return &SignOptions{
		Timeout:   2 * time.Hour,
		Algorithm: "HS256",
		Audience:  auth.AuthzAudience,
		Issuer:    "iamctl",
		Claims:    make(ArgList),
		Head:      make(ArgList),

		IOStreams: ioStreams,
	}
}

// NewCmdSign returns new initialized instance of sign sub command.
func NewCmdSign(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewSignOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   signUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Sign a jwt token with given secretID and secretKey",
		Long:                  "Sign a jwt token with given secretID and secretKey",
		TraverseChildren:      true,
		Example:               signExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return cmdutil.UsageErrorf(cmd, signUsageErrStr)
			}

			return nil
		},
	}

	// mark flag as deprecated
	cmd.Flags().DurationVar(&o.Timeout, "timeout", o.Timeout, "JWT token expires time.")
	cmd.Flags().DurationVar(
		&o.NotBefore,
		"not-before",
		o.NotBefore,
		"Identifies the time before which the JWT MUST NOT be accepted for processing.",
	)
	cmd.Flags().StringVar(
		&o.Algorithm,
		"algorithm",
		o.Algorithm,
		"Signing algorithm - possible values are HS256, HS384, HS512.",
	)
	cmd.Flags().StringVar(
		&o.Audience,
		"audience",
		o.Audience,
		"Identifies the recipients that the JWT is intended for.",
	)
	cmd.Flags().StringVar(&o.Issuer, "issuer", o.Issuer, "Identifies the principal that issued the JWT.")
	cmd.Flags().Var(&o.Claims, "claim", "Add additional claims. may be used more than once.")
	cmd.Flags().Var(&o.Head, "header", "Add additional header params. may be used more than once.")

	return cmd
}

// Complete completes all the required options.
func (o *SignOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *SignOptions) Validate(cmd *cobra.Command, args []string) error {
	switch o.Algorithm {
	case "HS256", "HS384", "HS512":
	default:
		return ErrSigningMethod
	}

	return nil
}

// Run executes a sign subcommand using the specified options.
func (o *SignOptions) Run(args []string) error {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(o.Timeout).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Add(o.NotBefore).Unix(),
		"aud": o.Audience,
		"iss": o.Issuer,
	}

	// add command line claims
	if len(o.Claims) > 0 {
		for k, v := range o.Claims {
			claims[k] = v
		}
	}

	// create a new token
	token := jwt.NewWithClaims(jwt.GetSigningMethod(o.Algorithm), claims)

	// add command line headers
	if len(o.Head) > 0 {
		for k, v := range o.Head {
			token.Header[k] = v
		}
	}
	token.Header["kid"] = args[0]

	tokenString, err := token.SignedString([]byte(args[1]))
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, tokenString+"\n")

	return nil
}
