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
	"github.com/openimsdk/open-im-server/v3/internal/rpc/third"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type ThirdRpcCmd struct {
	*RootCmd
	ctx         context.Context
	configMap   map[string]StructEnvPrefix
	thirdConfig ThirdConfig
}
type ThirdConfig struct {
	RpcConfig          config.Third
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	ZookeeperConfig    config.ZooKeeper
	NotificationConfig config.Notification
	Share              config.Share
	MinioConfig        config.Minio
}

func NewThirdRpcCmd() *ThirdRpcCmd {
	var thirdConfig ThirdConfig
	ret := &ThirdRpcCmd{thirdConfig: thirdConfig}
	ret.configMap = map[string]StructEnvPrefix{
		OpenIMRPCThirdCfgFileName: {EnvPrefix: thridEnvPrefix, ConfigStruct: &thirdConfig.RpcConfig},
		RedisConfigFileName:       {EnvPrefix: redisEnvPrefix, ConfigStruct: &thirdConfig.RedisConfig},
		ZookeeperConfigFileName:   {EnvPrefix: zoopkeeperEnvPrefix, ConfigStruct: &thirdConfig.ZookeeperConfig},
		MongodbConfigFileName:     {EnvPrefix: mongodbEnvPrefix, ConfigStruct: &thirdConfig.MongodbConfig},
		ShareFileName:             {EnvPrefix: shareEnvPrefix, ConfigStruct: &thirdConfig.Share},
		NotificationFileName:      {EnvPrefix: notificationEnvPrefix, ConfigStruct: &thirdConfig.NotificationConfig},
		MinioConfigFileName:       {EnvPrefix: minioEnvPrefix, ConfigStruct: &thirdConfig.MinioConfig},
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return ret.preRunE()
	}
	return ret
}

func (a *ThirdRpcCmd) Exec() error {
	return a.Execute()
}

func (a *ThirdRpcCmd) preRunE() error {
	return startrpc.Start(a.ctx, &a.thirdConfig.ZookeeperConfig, &a.thirdConfig.RpcConfig.Prometheus, a.thirdConfig.RpcConfig.RPC.ListenIP,
		a.thirdConfig.RpcConfig.RPC.RegisterIP, a.thirdConfig.RpcConfig.RPC.Ports,
		a.Index(), a.thirdConfig.Share.RpcRegisterName.Auth, &a.thirdConfig, third.Start)
}
