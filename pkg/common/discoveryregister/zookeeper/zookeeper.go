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

	"github.com/OpenIMSDK/tools/discoveryregistry"
	openkeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

// NewZookeeperDiscoveryRegister creates a new instance of ZookeeperDR for Zookeeper service discovery and registration.
func NewZookeeperDiscoveryRegister() (discoveryregistry.SvcDiscoveryRegistry, error) {
	schema := getEnv("ZOOKEEPER_SCHEMA", config.Config.Zookeeper.Schema)
	zkAddr := getZkAddrFromEnv(config.Config.Zookeeper.ZkAddr)
	username := getEnv("ZOOKEEPER_USERNAME", config.Config.Zookeeper.Username)
	password := getEnv("ZOOKEEPER_PASSWORD", config.Config.Zookeeper.Password)

	return openkeeper.NewClient(
		zkAddr,
		schema,
		openkeeper.WithFreq(time.Hour),
		openkeeper.WithUserNameAndPassword(username, password),
		openkeeper.WithRoundRobin(),
		openkeeper.WithTimeout(10),
		openkeeper.WithLogger(log.NewZkLogger()),
	)
}

// getEnv returns the value of an environment variable if it exists, otherwise it returns the fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getZkAddrFromEnv returns the value of an environment variable if it exists, otherwise it returns the fallback value.
func getZkAddrFromEnv(fallback []string) []string {
	if value, exists := os.LookupEnv("ZOOKEEPER_ADDRESS"); exists {
		return strings.Split(value, ",")
	}
	return fallback
}
