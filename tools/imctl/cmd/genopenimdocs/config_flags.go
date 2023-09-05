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

package genericclioptions

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/marmotedu/marmotedu-sdk-go/rest"
	"github.com/marmotedu/marmotedu-sdk-go/tools/clientcmd"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Defines flag for imctl.
const (
	FlagIMConfig      = "imconfig"
	FlagBearerToken   = "user.token"
	FlagUsername      = "user.username"
	FlagPassword      = "user.password"
	FlagSecretID      = "user.secret-id"
	FlagSecretKey     = "user.secret-key"
	FlagCertFile      = "user.client-certificate"
	FlagKeyFile       = "user.client-key"
	FlagTLSServerName = "server.tls-server-name"
	FlagInsecure      = "server.insecure-skip-tls-verify"
	FlagCAFile        = "server.certificate-authority"
	FlagAPIServer     = "server.address"
	FlagTimeout       = "server.timeout"
	FlagMaxRetries    = "server.max-retries"
	FlagRetryInterval = "server.retry-interval"
)

// RESTClientGetter is an interface that the ConfigFlags describe to provide an easier way to mock for commands
// and eliminate the direct coupling to a struct type.  Users may wish to duplicate this type in their own packages
// as per the golang type overlapping.
type RESTClientGetter interface {
	// ToRESTConfig returns restconfig
	ToRESTConfig() (*rest.Config, error)
	// ToRawIMConfigLoader return imconfig loader as-is
	ToRawIMConfigLoader() clientcmd.ClientConfig
}

var _ RESTClientGetter = &ConfigFlags{}

// ConfigFlags composes the set of values necessary
// for obtaining a REST client config.
type ConfigFlags struct {
	IMConfig *string

	BearerToken *string
	Username    *string
	Password    *string
	SecretID    *string
	SecretKey   *string

	Insecure      *bool
	TLSServerName *string
	CertFile      *string
	KeyFile       *string
	CAFile        *string

	APIServer     *string
	Timeout       *time.Duration
	MaxRetries    *int
	RetryInterval *time.Duration

	clientConfig clientcmd.ClientConfig
	lock         sync.Mutex
	// If set to true, will use persistent client config and
	// propagate the config to the places that need it, rather than
	// loading the config multiple times
	usePersistentConfig bool
}

// ToRESTConfig implements RESTClientGetter.
// Returns a REST client configuration based on a provided path
// to a .imconfig file, loading rules, and config flag overrides.
// Expects the AddFlags method to have been called.
func (f *ConfigFlags) ToRESTConfig() (*rest.Config, error) {
	return f.ToRawIMConfigLoader().ClientConfig()
}

// ToRawIMConfigLoader binds config flag values to config overrides
// Returns an interactive clientConfig if the password flag is enabled,
// or a non-interactive clientConfig otherwise.
func (f *ConfigFlags) ToRawIMConfigLoader() clientcmd.ClientConfig {
	if f.usePersistentConfig {
		return f.toRawIMPersistentConfigLoader()
	}

	return f.toRawIMConfigLoader()
}

func (f *ConfigFlags) toRawIMConfigLoader() clientcmd.ClientConfig {
	config := clientcmd.NewConfig()
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return clientcmd.NewClientConfigFromConfig(config)
}

// toRawIMPersistentConfigLoader binds config flag values to config overrides
// Returns a persistent clientConfig for propagation.
func (f *ConfigFlags) toRawIMPersistentConfigLoader() clientcmd.ClientConfig {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.clientConfig == nil {
		f.clientConfig = f.toRawIMConfigLoader()
	}

	return f.clientConfig
}

// AddFlags binds client configuration flags to a given flagset.
func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.IMConfig != nil {
		flags.StringVar(f.IMConfig, FlagIMConfig, *f.IMConfig,
			fmt.Sprintf("Path to the %s file to use for CLI requests", FlagIMConfig))
	}

	if f.BearerToken != nil {
		flags.StringVar(
			f.BearerToken,
			FlagBearerToken,
			*f.BearerToken,
			"Bearer token for authentication to the API server",
		)
	}

	if f.Username != nil {
		flags.StringVar(f.Username, FlagUsername, *f.Username, "Username for basic authentication to the API server")
	}

	if f.Password != nil {
		flags.StringVar(f.Password, FlagPassword, *f.Password, "Password for basic authentication to the API server")
	}

	if f.SecretID != nil {
		flags.StringVar(f.SecretID, FlagSecretID, *f.SecretID, "SecretID for JWT authentication to the API server")
	}

	if f.SecretKey != nil {
		flags.StringVar(f.SecretKey, FlagSecretKey, *f.SecretKey, "SecretKey for jwt authentication to the API server")
	}

	if f.CertFile != nil {
		flags.StringVar(f.CertFile, FlagCertFile, *f.CertFile, "Path to a client certificate file for TLS")
	}
	if f.KeyFile != nil {
		flags.StringVar(f.KeyFile, FlagKeyFile, *f.KeyFile, "Path to a client key file for TLS")
	}
	if f.TLSServerName != nil {
		flags.StringVar(f.TLSServerName, FlagTLSServerName, *f.TLSServerName, ""+
			"Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used")
	}
	if f.Insecure != nil {
		flags.BoolVar(f.Insecure, FlagInsecure, *f.Insecure, ""+
			"If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	}
	if f.CAFile != nil {
		flags.StringVar(f.CAFile, FlagCAFile, *f.CAFile, "Path to a cert file for the certificate authority")
	}

	if f.APIServer != nil {
		flags.StringVarP(f.APIServer, FlagAPIServer, "s", *f.APIServer, "The address and port of the IM API server")
	}

	if f.Timeout != nil {
		flags.DurationVar(
			f.Timeout,
			FlagTimeout,
			*f.Timeout,
			"The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.",
		)
	}

	if f.MaxRetries != nil {
		flag.IntVar(f.MaxRetries, FlagMaxRetries, *f.MaxRetries, "Maximum number of retries.")
	}

	if f.RetryInterval != nil {
		flags.DurationVar(
			f.RetryInterval,
			FlagRetryInterval,
			*f.RetryInterval,
			"The interval time between each attempt.",
		)
	}
}

// WithDeprecatedPasswordFlag enables the username and password config flags.
func (f *ConfigFlags) WithDeprecatedPasswordFlag() *ConfigFlags {
	f.Username = pointer.ToString("")
	f.Password = pointer.ToString("")

	return f
}

// WithDeprecatedSecretFlag enables the secretID and secretKey config flags.
func (f *ConfigFlags) WithDeprecatedSecretFlag() *ConfigFlags {
	f.SecretID = pointer.ToString("")
	f.SecretKey = pointer.ToString("")

	return f
}

// NewConfigFlags returns ConfigFlags with default values set.
func NewConfigFlags(usePersistentConfig bool) *ConfigFlags {
	return &ConfigFlags{
		IMConfig: pointer.ToString(""),

		BearerToken:   pointer.ToString(""),
		Insecure:      pointer.ToBool(false),
		TLSServerName: pointer.ToString(""),
		CertFile:      pointer.ToString(""),
		KeyFile:       pointer.ToString(""),
		CAFile:        pointer.ToString(""),

		APIServer:           pointer.ToString(""),
		Timeout:             pointer.ToDuration(30 * time.Second),
		MaxRetries:          pointer.ToInt(0),
		RetryInterval:       pointer.ToDuration(1 * time.Second),
		usePersistentConfig: usePersistentConfig,
	}
}
