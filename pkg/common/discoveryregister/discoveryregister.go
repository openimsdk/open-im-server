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

package discoveryregister

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/kubernetes"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"time"
)

const (
	zookeeper = "zoopkeeper"
	kubenetes = "k8s"

	direct = "direct"
)

// NewDiscoveryRegister creates a new service discovery and registry client based on the provided environment type.
func NewDiscoveryRegister(zookeeperConfig *config.ZooKeeper) (discovery.SvcDiscoveryRegistry, error) {

	switch zookeeperConfig.Env {
	case "zookeeper":
		return zookeeper.NewZookeeperDiscoveryRegister(&config.Zookeeper)
		zk, err := zookeeper.NewZkClient(
			zookeeperConfig.zkAddr,
			schema,
			zookeeper.WithFreq(time.Hour),
			zookeeper.WithUserNameAndPassword(username, password),
			zookeeper.WithRoundRobin(),
			zookeeper.WithTimeout(10),
		)
		if err != nil {
			return nil, err
		}
	case "k8s":
		return kubernetes.NewK8sDiscoveryRegister(config.RpcRegisterName.OpenImMessageGatewayName)
	case "direct":
		return direct.NewConnDirect(config)
	default:
		return nil, errs.New("unsupported discovery type", "type", config.Envs.Discovery).Wrap()
	}
}
