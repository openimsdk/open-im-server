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

package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/openimsdk/open-im-server/v3/pkg/common/tls"
)

type TLSConfig struct {
	CACrt              string
	ClientCrt          string
	ClientKey          string
	ClientKeyPwd       string
	InsecureSkipVerify bool
}

// SetupTLSConfig set up the TLS config from config file.
func SetupTLSConfig(cfg *sarama.Config, tlsConfig *TLSConfig) error {
	if tlsConfig != nil {
		cfg.Net.TLS.Enable = true
		tlsConfig, err := tls.NewTLSConfig(
			tlsConfig.ClientCrt,
			tlsConfig.ClientKey,
			tlsConfig.CACrt,
			[]byte(tlsConfig.ClientKeyPwd),
			tlsConfig.InsecureSkipVerify,
		)
		if err != nil {
			return err
		}
		cfg.Net.TLS.Config = tlsConfig
	}
	return nil
}

// getEnvOrConfig returns the value of the environment variable if it exists,
// otherwise, it returns the value from the configuration file.
func getEnvOrConfig(envName string, configValue string) string {
	if value, exists := os.LookupEnv(envName); exists {
		return value
	}
	return configValue
}

// getKafkaAddrFromEnv returns the Kafka addresses combined from the KAFKA_ADDRESS and KAFKA_PORT environment variables.
// If the environment variables are not set, it returns the fallback value.
func getKafkaAddrFromEnv(fallback []string) []string {
	envAddr := os.Getenv("KAFKA_ADDRESS")
	envPort := os.Getenv("KAFKA_PORT")

	if envAddr != "" && envPort != "" {
		addresses := strings.Split(envAddr, ",")
		for i, addr := range addresses {
			addresses[i] = fmt.Sprintf("%s:%s", addr, envPort)
		}
		return addresses
	}

	return fallback
}
