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
	"errors"
	"os"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/direct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/kubernetes"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/zookeeper"
)

// NewDiscoveryRegister creates a new service discovery and registry client based on the provided environment type.
func NewDiscoveryRegister(envType string) (discoveryregistry.SvcDiscoveryRegistry, error) {

	if os.Getenv("ENVS_DISCOVERY") != "" {
		envType = os.Getenv("ENVS_DISCOVERY")
	}

	switch envType {
	case "zookeeper":
		return zookeeper.NewZookeeperDiscoveryRegister()
	case "k8s":
		return kubernetes.NewK8sDiscoveryRegister()
	case "direct":
		return direct.NewConnDirect()
	default:
		return nil, errs.Wrap(errors.New("envType not correct"))
	}
}
