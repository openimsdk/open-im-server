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

package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/openimsdk/tools/system/program"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/auth"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
)

type AuthRpcCmd struct {
	*RootCmd
	ctx        context.Context
	configMap  map[string]any
	authConfig *auth.Config
}

func NewAuthRpcCmd() *AuthRpcCmd {
	var authConfig auth.Config
	ret := &AuthRpcCmd{authConfig: &authConfig}
	ret.configMap = map[string]any{
		OpenIMRPCAuthCfgFileName: &authConfig.RpcConfig,
		RedisConfigFileName:      &authConfig.RedisConfig,
		ShareFileName:            &authConfig.Share,
		DiscoveryConfigFilename:  &authConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}

	return ret
}

func (a *AuthRpcCmd) Exec() error {
	return a.Execute()
}

func (a *AuthRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.authConfig.Discovery, &a.authConfig.RpcConfig.Prometheus, a.authConfig.RpcConfig.RPC.ListenIP,
		a.Index(), a.authConfig.Share.RpcRegisterName.Auth, &a.authConfig.Share, a.authConfig, auth.Start)
}
