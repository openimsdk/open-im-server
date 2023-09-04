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

// Package jwt can be used to sign/show/verify jwt token with given secretID and secretKey.
package jwt

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

var jwtLong = templates.LongDesc(`
	JWT command.

	This commands is used to sigin/show/verify jwt token.`)

// NewCmdJWT returns new initialized instance of 'jwt' sub command.
func NewCmdJWT(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "jwt SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "JWT command-line tool",
		Long:                  jwtLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.ErrOut),
	}

	// add subcommands
	cmd.AddCommand(NewCmdSign(f, ioStreams))
	cmd.AddCommand(NewCmdShow(f, ioStreams))
	cmd.AddCommand(NewCmdVerify(f, ioStreams))

	return cmd
}
