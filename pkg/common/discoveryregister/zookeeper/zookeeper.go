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

package zookeeper

import (
	"os"
	"strings"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/zookeeper"
)

// NewZookeeperDiscoveryRegister creates a new instance of ZookeeperDR for Zookeeper service discovery and registration.
func NewZookeeperDiscoveryRegister(zkConf *config.Zookeeper) (discovery.SvcDiscoveryRegistry, error) {
	schema := getEnv("ZOOKEEPER_SCHEMA", zkConf.Schema)
	zkAddr := getZkAddrFromEnv(zkConf.ZkAddr)
	username := getEnv("ZOOKEEPER_USERNAME", zkConf.Username)
	password := getEnv("ZOOKEEPER_PASSWORD", zkConf.Password)

	zk, err := zookeeper.NewZkClient(
		zkAddr,
		schema,
		zookeeper.WithFreq(time.Hour),
		zookeeper.WithUserNameAndPassword(username, password),
		zookeeper.WithRoundRobin(),
		zookeeper.WithTimeout(10),
	)
	if err != nil {
		return nil, err
	}
	return zk, nil
}

// getEnv returns the value of an environment variable if it exists, otherwise it returns the fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getZkAddrFromEnv returns the Zookeeper addresses combined from the ZOOKEEPER_ADDRESS and ZOOKEEPER_PORT environment variables.
// If the environment variables are not set, it returns the fallback value.
func getZkAddrFromEnv(fallback []string) []string {
	address, addrExists := os.LookupEnv("ZOOKEEPER_ADDRESS")
	port, portExists := os.LookupEnv("ZOOKEEPER_PORT")

	if addrExists && portExists {
		addresses := strings.Split(address, ",")
		for i, addr := range addresses {
			addresses[i] = addr + ":" + port
		}
		return addresses
	}
	return fallback
}
