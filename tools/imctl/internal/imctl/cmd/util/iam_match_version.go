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

package util

import (
	"context"
	"fmt"
	"sync"

	"github.com/marmotedu/marmotedu-sdk-go/rest"
	"github.com/marmotedu/marmotedu-sdk-go/tools/clientcmd"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/scheme"
	"github.com/openim-sigs/component-base/pkg/version"
	"github.com/spf13/pflag"

	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	flagMatchBinaryVersion = "match-server-version"
)

// MatchVersionFlags is for setting the "match server version" function.
type MatchVersionFlags struct {
	Delegate genericclioptions.RESTClientGetter

	RequireMatchedServerVersion bool
	checkServerVersion          sync.Once
	matchesServerVersionErr     error
}

var _ genericclioptions.RESTClientGetter = &MatchVersionFlags{}

func (f *MatchVersionFlags) checkMatchingServerVersion() error {
	f.checkServerVersion.Do(func() {
		if !f.RequireMatchedServerVersion {
			return
		}

		clientConfig, err := f.Delegate.ToRESTConfig()
		if err != nil {
			f.matchesServerVersionErr = err
			return
		}

		setIAMDefaults(clientConfig)
		restClient, err := rest.RESTClientFor(clientConfig)
		if err != nil {
			f.matchesServerVersionErr = err
			return
		}

		var sVer *version.Info
		if err := restClient.Get().AbsPath("/version").Do(context.TODO()).Into(&sVer); err != nil {
			f.matchesServerVersionErr = err
			return
		}

		clientVersion := version.Get()

		// GitVersion includes GitCommit and GitTreeState, but best to be safe?
		if clientVersion.GitVersion != sVer.GitVersion || clientVersion.GitCommit != sVer.GitCommit ||
			clientVersion.GitTreeState != sVer.GitTreeState {
			f.matchesServerVersionErr = fmt.Errorf(
				"server version (%#v) differs from client version (%#v)",
				sVer,
				version.Get(),
			)
		}
	})

	return f.matchesServerVersionErr
}

// ToRESTConfig implements RESTClientGetter.
// Returns a REST client configuration based on a provided path
// to a .iamconfig file, loading rules, and config flag overrides.
// Expects the AddFlags method to have been called.
func (f *MatchVersionFlags) ToRESTConfig() (*rest.Config, error) {
	if err := f.checkMatchingServerVersion(); err != nil {
		return nil, err
	}
	clientConfig, err := f.Delegate.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	// TODO we should not have to do this.  It smacks of something going wrong.
	setIAMDefaults(clientConfig)
	return clientConfig, nil
}

func (f *MatchVersionFlags) ToRawIAMConfigLoader() clientcmd.ClientConfig {
	return f.Delegate.ToRawIAMConfigLoader()
}

func (f *MatchVersionFlags) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(
		&f.RequireMatchedServerVersion,
		flagMatchBinaryVersion,
		f.RequireMatchedServerVersion,
		"Require server version to match client version",
	)
}

func NewMatchVersionFlags(delegate genericclioptions.RESTClientGetter) *MatchVersionFlags {
	return &MatchVersionFlags{
		Delegate: delegate,
	}
}

// setIAMDefaults sets default values on the provided client config for accessing the
// IAM API or returns an error if any of the defaults are impossible or invalid.
// TODO this isn't what we want.  Each iamclient should be setting defaults as it sees fit.
func setIAMDefaults(config *rest.Config) error {
	// TODO remove this hack.  This is allowing the GetOptions to be serialized.
	config.GroupVersion = &scheme.GroupVersion{Group: "iam.api", Version: "v1"}

	if config.APIPath == "" {
		config.APIPath = "/api"
	}
	if config.Negotiator == nil {
		// This codec factory ensures the resources are not converted. Therefore, resources
		// will not be round-tripped through internal versions. Defaulting does not happen
		// on the client.
		config.Negotiator = runtime.NewSimpleClientNegotiator()
	}
	return rest.SetIAMDefaults(config)
}
