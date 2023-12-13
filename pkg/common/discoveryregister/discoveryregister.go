package discoveryregister

import (
	"errors"

	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/kubernetes"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/zookeeper"

	"github.com/OpenIMSDK/tools/discoveryregistry"
)

// NewDiscoveryRegister creates a new service discovery and registry client based on the provided environment type.
func NewDiscoveryRegister(envType string) (discoveryregistry.SvcDiscoveryRegistry, error) {
	switch envType {
	case "zookeeper":
		return zookeeper.NewZookeeperDiscoveryRegister()
	case "k8s":
		return kubernetes.NewK8sDiscoveryRegister()
	default:
		return nil, errors.New("envType not correct")
	}
}
