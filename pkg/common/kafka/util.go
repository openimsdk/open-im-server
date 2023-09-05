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
	"github.com/Shopify/sarama"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tls"
)

// SetupTLSConfig set up the TLS config from config file.
func SetupTLSConfig(cfg *sarama.Config) {
	if config.Config.Kafka.TLS != nil {
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = tls.NewTLSConfig(
			config.Config.Kafka.TLS.ClientCrt,
			config.Config.Kafka.TLS.ClientKey,
			config.Config.Kafka.TLS.CACrt,
			[]byte(config.Config.Kafka.TLS.ClientKeyPwd),
		)
	}
}
