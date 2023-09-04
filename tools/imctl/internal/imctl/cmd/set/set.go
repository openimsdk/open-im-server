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

// Package set used to set specific features on objects.
package set

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd/util"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/cli/genericclioptions"
)

var setLong = templates.LongDesc(`
	Configure objects.

	These commands help you make changes to existing objects.`)

// NewCmdSet returns an initialized Command instance for 'set' sub command.
func NewCmdSet(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Set specific features on objects",
		Long:                  setLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.ErrOut),
	}

	// add subcommands
	// cmd.AddCommand(NewCmdDB(f, ioStreams))

	return cmd
}
