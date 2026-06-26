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
	"strings"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/utils/runtimeenv"
)

//func TestNewDiscoveryRegister(t *testing.T) {
//	setupTestEnvironment()
//	conf := config.NewGlobalConfig()
//	tests := []struct {
//		envType        string
//		gatewayName    string
//		expectedError  bool
//		expectedResult bool
//	}{
//		{"zookeeper", "MessageGateway", false, true},
//		{"k8s", "MessageGateway", false, true},
//		{"direct", "MessageGateway", false, true},
//		{"invalid", "MessageGateway", true, false},
//	}
//
//	for _, test := range tests {
//		conf.Envs.Discovery = test.envType
//		conf.RpcRegisterName.OpenImMessageGatewayName = test.gatewayName
//		client, err := NewDiscoveryRegister(conf)
//
//		if test.expectedError {
//			assert.Error(t, err)
//		} else {
//			assert.NoError(t, err)
//			if test.expectedResult {
//				assert.Implements(t, (*discovery.SvcDiscoveryRegistry)(nil), client)
//			} else {
//				assert.Nil(t, client)
//			}
//		}
//	}
//}

func TestNewDiscoveryRegisterRejectsKubernetesOutsideCluster(t *testing.T) {
	if runtimeenv.RuntimeEnvironment() == config.KUBERNETES {
		t.Skip("outside-cluster fallback is only relevant outside Kubernetes")
	}
	discovery := &config.Discovery{
		Enable: config.KUBERNETES,
		Kubernetes: config.Kubernetes{
			Namespace: "default",
		},
	}

	client, err := NewDiscoveryRegister(discovery, &config.Share{}, nil)
	if err == nil && client != nil {
		client.Close()
	}
	if err == nil || !strings.Contains(err.Error(), "unsupported discovery type") {
		t.Fatalf("%q outside Kubernetes should not select Kubernetes discovery: %v", config.KUBERNETES, err)
	}
}
