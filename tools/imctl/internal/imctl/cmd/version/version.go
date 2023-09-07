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

// Package version print the client and server version information.
package version

import (
	"context"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/openim-sigs/component-base/pkg/json"
	"github.com/openim-sigs/component-base/pkg/version"
	"github.com/spf13/cobra"

	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/util/templates"
	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/cli/genericclioptions"
	cmdutil "github.com/OpenIMSDK/Open-IM-Server/tools/imctl/pkg/util"
)

// Version is a struct for version information.
type Version struct {
	ClientVersion *version.Info `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"`
	ServerVersion *version.Info `json:"serverVersion,omitempty" yaml:"serverVersion,omitempty"`
}

var versionExample = templates.Examples(`
		# Print the client and server versions for the current context
		imctl version`)

// Options is a struct to support version command.
type Options struct {
	ClientOnly bool
	Short      bool
	Output     string

	client *restclient.RESTClient
	genericclioptions.IOStreams
}

// NewOptions returns initialized Options.
func NewOptions(ioStreams genericclioptions.IOStreams) *Options {
	return &Options{
		IOStreams: ioStreams,
	}
}

// NewCmdVersion returns a cobra command for fetching versions.
func NewCmdVersion(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewOptions(ioStreams)
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the client and server version information",
		Long:    "Print the client and server version information for the current context",
		Example: versionExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().BoolVar(
		&o.ClientOnly,
		"client",
		o.ClientOnly,
		"If true, shows client version only (no server required).",
	)
	cmd.Flags().BoolVar(&o.Short, "short", o.Short, "If true, print just the version number.")
	cmd.Flags().StringVarP(&o.Output, "output", "o", o.Output, "One of 'yaml' or 'json'.")

	return cmd
}

// Complete completes all the required options.
func (o *Options) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	var err error
	if o.ClientOnly {
		return nil
	}

	o.client, err = f.RESTClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate validates the provided options.
func (o *Options) Validate() error {
	if o.Output != "" && o.Output != "yaml" && o.Output != "json" {
		return errors.New(`--output must be 'yaml' or 'json'`)
	}

	return nil
}

// Run executes version command.
func (o *Options) Run() error {
	var (
		serverVersion *version.Info
		serverErr     error
		versionInfo   Version
	)

	clientVersion := version.Get()
	versionInfo.ClientVersion = &clientVersion

	if !o.ClientOnly && o.client != nil {
		// Always request fresh data from the server
		if err := o.client.Get().AbsPath("/version").Do(context.TODO()).Into(&serverVersion); err != nil {
			return err
		}
		versionInfo.ServerVersion = serverVersion
	}

	switch o.Output {
	case "":
		if o.Short {
			fmt.Fprintf(o.Out, "Client Version: %s\n", clientVersion.GitVersion)

			if serverVersion != nil {
				fmt.Fprintf(o.Out, "Server Version: %s\n", serverVersion.GitVersion)
			}
		} else {
			fmt.Fprintf(o.Out, "Client Version: %s\n", fmt.Sprintf("%#v", clientVersion))
			if serverVersion != nil {
				fmt.Fprintf(o.Out, "Server Version: %s\n", fmt.Sprintf("%#v", *serverVersion))
			}
		}
	case "yaml":
		marshaled, err := yaml.Marshal(&versionInfo)
		if err != nil {
			return err
		}

		fmt.Fprintln(o.Out, string(marshaled))
	case "json":
		marshaled, err := json.MarshalIndent(&versionInfo, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(o.Out, string(marshaled))
	default:
		// There is a bug in the program if we hit this case.
		// However, we follow a policy of never panicking.
		return fmt.Errorf("VersionOptions were not validated: --output=%q should have been rejected", o.Output)
	}

	return serverErr
}
