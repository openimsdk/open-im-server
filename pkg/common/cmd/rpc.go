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
	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
)

type rpcInitFuc func(ctx context.Context, config *AuthConfig, disCov discovery.SvcDiscoveryRegistry, server *grpc.Server) error

type AuthRpcCmd struct {
	*RootCmd
	initFunc   rpcInitFuc
	ctx        context.Context
	configMap  map[string]StructEnvPrefix
	authConfig AuthConfig
}
type AuthConfig struct {
	RpcConfig       config2.Auth
	RedisConfig     config2.Redis
	ZookeeperConfig config2.ZooKeeper
}

func NewAuthRpcCmd(initFunc rpcInitFuc) *AuthRpcCmd {
	var authConfig AuthConfig
	ret := &AuthRpcCmd{initFunc: initFunc, authConfig: authConfig}
	ret.configMap = map[string]StructEnvPrefix{
		OpenIMRPCAuthCfgFileName: {EnvPrefix: authEnvPrefix, ConfigStruct: &authConfig.RpcConfig},
		RedisConfigFileName:      {EnvPrefix: redisEnvPrefix, ConfigStruct: &authConfig.RedisConfig},
		ZookeeperConfigFileName:  {EnvPrefix: zoopkeeperEnvPrefix, ConfigStruct: &authConfig.ZookeeperConfig},
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config2.Version)
	ret.RunE()
	return ret
}

func (a *AuthRpcCmd) Exec() error {
	return a.Execute()
}

func (a *AuthRpcCmd) RunE() error {
	return startrpc.Start(a.ctx, &a.authConfig.ZookeeperConfig, &a.authConfig.RpcConfig.Prometheus, a.authConfig.RpcConfig.RPC.ListenIP,
		a.authConfig.RpcConfig.RPC.RegisterIP, a.authConfig.RpcConfig.RPC.Ports,
		a.Index(), a.authConfig.ZookeeperConfig.RpcRegisterName.Auth, &a.authConfig, a.initFunc)
}
