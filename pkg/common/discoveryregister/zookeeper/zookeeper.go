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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	openkeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

// NewZookeeperDiscoveryRegister creates a new instance of ZookeeperDR for Zookeeper service discovery and registration.
func NewZookeeperDiscoveryRegister() (discoveryregistry.SvcDiscoveryRegistry, error) {
	schema := getEnv("ZOOKEEPER_SCHEMA", config.Config.Zookeeper.Schema)
	zkAddr := getZkAddrFromEnv(config.Config.Zookeeper.ZkAddr)
	username := getEnv("ZOOKEEPER_USERNAME", config.Config.Zookeeper.Username)
	password := getEnv("ZOOKEEPER_PASSWORD", config.Config.Zookeeper.Password)

	zk, err := openkeeper.NewClient(
		zkAddr,
		schema,
		openkeeper.WithFreq(time.Hour),
		openkeeper.WithUserNameAndPassword(username, password),
		openkeeper.WithRoundRobin(),
		openkeeper.WithTimeout(10),
		openkeeper.WithLogger(log.NewZkLogger()),
	)
	if err != nil {
		uriFormat := "address:%s, username:%s, password:%s, schema:%s."
		errInfo := fmt.Sprintf(uriFormat,
			config.Config.Zookeeper.ZkAddr,
			config.Config.Zookeeper.Username,
			config.Config.Zookeeper.Password,
			config.Config.Zookeeper.Schema)
		return nil, errs.Wrap(err, errInfo)
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
