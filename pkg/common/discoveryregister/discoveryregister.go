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
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/discovery/zookeeper"
	"github.com/openimsdk/tools/errs"
	"time"
)

// NewDiscoveryRegister creates a new service discovery and registry client based on the provided environment type.
func NewDiscoveryRegister(discovery *config.Discovery, share *config.Share) (discovery.SvcDiscoveryRegistry, error) {
	switch discovery.Enable {
	case "zookeeper":
		return zookeeper.NewZkClient(
			discovery.ZooKeeper.Address,
			discovery.ZooKeeper.Schema,
			zookeeper.WithFreq(time.Hour),
			zookeeper.WithUserNameAndPassword(discovery.ZooKeeper.Username, discovery.ZooKeeper.Password),
			zookeeper.WithRoundRobin(),
			zookeeper.WithTimeout(10),
		)
	case "k8s":
		return kubernetes.NewK8sDiscoveryRegister(share.RpcRegisterName.MessageGateway)
	case "etcd":
		return etcd.NewSvcDiscoveryRegistry(
			discovery.Etcd.RootDirectory,
			discovery.Etcd.Address,
			nil,
			etcd.WithDialTimeout(10*time.Second),
			etcd.WithMaxCallSendMsgSize(20*1024*1024),
			etcd.WithUsernameAndPassword(discovery.Etcd.Username, discovery.Etcd.Password))
	default:
		return nil, errs.New("unsupported discovery type", "type", discovery.Enable).Wrap()
	}
}
