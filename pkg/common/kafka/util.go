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
