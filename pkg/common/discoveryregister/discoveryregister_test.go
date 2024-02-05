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
	"os"
	"testing"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/stretchr/testify/assert"
)

func setupTestEnvironment() {
	os.Setenv("ZOOKEEPER_SCHEMA", "openim")
	os.Setenv("ZOOKEEPER_ADDRESS", "172.28.0.1")
	os.Setenv("ZOOKEEPER_PORT", "12181")
	os.Setenv("ZOOKEEPER_USERNAME", "")
	os.Setenv("ZOOKEEPER_PASSWORD", "")
}

func TestNewDiscoveryRegister(t *testing.T) {
	setupTestEnvironment()

	tests := []struct {
		envType        string
		expectedError  bool
		expectedResult bool
	}{
		{"zookeeper", false, true},
		{"k8s", false, true}, // 假设 k8s 配置也已正确设置
		{"direct", false, true},
		{"invalid", true, false},
	}

	for _, test := range tests {
		client, err := NewDiscoveryRegister(test.envType)

		if test.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			if test.expectedResult {
				assert.Implements(t, (*discoveryregistry.SvcDiscoveryRegistry)(nil), client)
			} else {
				assert.Nil(t, client)
			}
		}
	}
}
