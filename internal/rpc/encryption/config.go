// Copyright © 2024 OpenIM. All rights reserved.
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

package encryption

import (
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

// PrometheusConfig contains Prometheus configuration
type PrometheusConfig struct {
	Enable bool  `yaml:"enable"`
	Ports  []int `yaml:"ports"`
}

// Config represents the configuration for encryption service
type Config struct {
	RpcConfig        config.RPC       `yaml:"rpc"`
	MongodbConfig    config.Mongo     `yaml:"mongo"`
	Discovery        config.Discovery `yaml:"discovery"`
	PrometheusConfig PrometheusConfig `yaml:"prometheus"`
	EncryptionConfig EncryptionConfig `yaml:"encryption"`
}

// EncryptionConfig contains encryption-specific configuration
type EncryptionConfig struct {
	Mode   string       `yaml:"mode"` // "aes", "signal", "hybrid"
	AES    AESConfig    `yaml:"aes"`
	Signal SignalConfig `yaml:"signal"`
}

// AESConfig contains AES encryption configuration
type AESConfig struct {
	Enabled bool `yaml:"enabled"`
}

// SignalConfig contains Signal Protocol configuration
type SignalConfig struct {
	Enabled                bool          `yaml:"enabled"`
	PreKeyBatch            int           `yaml:"preKeyBatch"`
	KeyRotationInterval    time.Duration `yaml:"keyRotationInterval"`
	SessionCleanupInterval time.Duration `yaml:"sessionCleanupInterval"`
	PrekeyCleanupInterval  time.Duration `yaml:"prekeyCleanupInterval"`

	// Security settings
	MaxOneTimePreKeys    int `yaml:"maxOneTimePreKeys"`
	MaxSessionsPerDevice int `yaml:"maxSessionsPerDevice"`

	// Validation settings
	ValidateSignatures  bool `yaml:"validateSignatures"`
	RequireIdentityKeys bool `yaml:"requireIdentityKeys"`
}

// GetEncryptionMode returns the current encryption mode
func (c *Config) GetEncryptionMode() string {
	if c.EncryptionConfig.Mode == "" {
		return "aes" // default to AES for compatibility
	}
	return c.EncryptionConfig.Mode
}

// IsSignalEnabled returns true if Signal Protocol is enabled
func (c *Config) IsSignalEnabled() bool {
	return c.EncryptionConfig.Signal.Enabled &&
		(c.EncryptionConfig.Mode == "signal" || c.EncryptionConfig.Mode == "hybrid")
}

// IsAESEnabled returns true if AES encryption is enabled
func (c *Config) IsAESEnabled() bool {
	return c.EncryptionConfig.AES.Enabled ||
		c.EncryptionConfig.Mode == "aes" ||
		c.EncryptionConfig.Mode == "hybrid"
}

// GetSignalConfig returns Signal Protocol configuration
func (c *Config) GetSignalConfig() *SignalConfig {
	// Set defaults if not specified
	if c.EncryptionConfig.Signal.PreKeyBatch == 0 {
		c.EncryptionConfig.Signal.PreKeyBatch = 100
	}
	if c.EncryptionConfig.Signal.KeyRotationInterval == 0 {
		c.EncryptionConfig.Signal.KeyRotationInterval = 7 * 24 * time.Hour // 7 days
	}
	if c.EncryptionConfig.Signal.SessionCleanupInterval == 0 {
		c.EncryptionConfig.Signal.SessionCleanupInterval = 30 * 24 * time.Hour // 30 days
	}
	if c.EncryptionConfig.Signal.PrekeyCleanupInterval == 0 {
		c.EncryptionConfig.Signal.PrekeyCleanupInterval = 7 * 24 * time.Hour // 7 days
	}
	if c.EncryptionConfig.Signal.MaxOneTimePreKeys == 0 {
		c.EncryptionConfig.Signal.MaxOneTimePreKeys = 100
	}
	if c.EncryptionConfig.Signal.MaxSessionsPerDevice == 0 {
		c.EncryptionConfig.Signal.MaxSessionsPerDevice = 1000
	}

	return &c.EncryptionConfig.Signal
}
