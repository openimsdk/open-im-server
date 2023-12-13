package discoveryregister

import (
	"os"
	"testing"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/stretchr/testify/assert"
)

func setupTestEnvironment() {
	os.Setenv("ZOOKEEPER_SCHEMA", "openim")
	os.Setenv("ZOOKEEPER_ADDRESS", "172.28.0.1:12181")
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
