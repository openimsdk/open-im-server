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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/encryption"
	"gopkg.in/yaml.v3"
)

func main() {
	var configPath = flag.String("config", "./config", "path to config directory")
	flag.Parse()

	// Load encryption service config independently
	config, err := loadEncryptionConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load encryption config: %v", err)
	}

	// Start the encryption service
	ctx := context.Background()
	if err := encryption.Start(ctx, config); err != nil {
		log.Fatalf("Failed to start encryption service: %v", err)
	}
}

// loadEncryptionConfig loads configuration from multiple files following OpenIM pattern
func loadEncryptionConfig(configDir string) (*encryption.Config, error) {
	config := &encryption.Config{}

	// Load main encryption config (only contains rpc, prometheus, and encryption-specific settings)
	encryptionConfigFile := filepath.Join(configDir, "openim-rpc-encryption.yml")
	if err := loadYAMLFile(encryptionConfigFile, config); err != nil {
		return nil, fmt.Errorf("failed to load encryption config: %w", err)
	}

	// Load shared MongoDB configuration
	mongoConfigFile := filepath.Join(configDir, "mongodb.yml")
	if err := loadYAMLFile(mongoConfigFile, &config.MongodbConfig); err != nil {
		return nil, fmt.Errorf("failed to load mongodb config: %w", err)
	}

	// Load shared Discovery configuration
	discoveryConfigFile := filepath.Join(configDir, "discovery.yml")
	if err := loadYAMLFile(discoveryConfigFile, &config.Discovery); err != nil {
		return nil, fmt.Errorf("failed to load discovery config: %w", err)
	}

	return config, nil
}

// loadYAMLFile loads a YAML file into the given struct
func loadYAMLFile(filename string, out interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}
