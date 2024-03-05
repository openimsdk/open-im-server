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

package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/open-im-server/v3/pkg/util/genutil"
	"gopkg.in/yaml.v3"
)

//go:embed version
var Version string

const (
	FileName             = "config.yaml"
	NotificationFileName = "notification.yaml"
	DefaultFolderPath    = "../config/"
)

// GetDefaultConfigPath returns the absolute path to the default configuration directory
// relative to the executable's location. It is intended for use in Kubernetes container configurations.
// Errors are returned to the caller to allow for flexible error handling.
func GetDefaultConfigPath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", errs.Wrap(err, "failed to get executable path")
	}

	// Calculate the config path as a directory relative to the executable's location
	configPath, err := genutil.OutDir(filepath.Join(filepath.Dir(executablePath), "../config/"))
	if err != nil {
		return "", errs.Wrap(err, "failed to get output directory")
	}
	return configPath, nil
}

// GetProjectRoot returns the absolute path of the project root directory by navigating up from the directory
// containing the executable. It provides a detailed error if the path cannot be determined.
func GetProjectRoot() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", errs.Wrap(err, "failed to retrieve executable path")
	}

	// Attempt to compute the project root by navigating up from the executable's directory
	projectRoot, err := genutil.OutDir(filepath.Join(filepath.Dir(executablePath), "../../../../.."))
	if err != nil {
		return "", err
	}

	return projectRoot, nil
}

func GetOptionsByNotification(cfg NotificationConf) msgprocessor.Options {
	opts := msgprocessor.NewOptions()

	if cfg.UnreadCount {
		opts = msgprocessor.WithOptions(opts, msgprocessor.WithUnreadCount(true))
	}
	if cfg.OfflinePush.Enable {
		opts = msgprocessor.WithOptions(opts, msgprocessor.WithOfflinePush(true))
	}
	switch cfg.ReliabilityLevel {
	case constant.UnreliableNotification:
	case constant.ReliableNotificationNoMsg:
		opts = msgprocessor.WithOptions(opts, msgprocessor.WithHistory(true), msgprocessor.WithPersistent())
	}
	opts = msgprocessor.WithOptions(opts, msgprocessor.WithSendMsg(cfg.IsSendMsg))

	return opts
}

// initConfig loads configuration from a specified path into the provided config structure.
// If the specified config file does not exist, it attempts to load from the project's default "config" directory.
// It logs informative messages regarding the configuration path being used.
func initConfig(config any, configName, configFolderPath string) error {
	configFilePath := filepath.Join(configFolderPath, configName)
	_, err := os.Stat(configFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errs.Wrap(err, fmt.Sprintf("failed to check existence of config file at path: %s", configFilePath))
		}
		var projectRoot string
		projectRoot, err = GetProjectRoot()
		if err != nil {
			return err
		}
		configFilePath = filepath.Join(projectRoot, "config", configName)
		fmt.Printf("Configuration file not found at specified path. Falling back to project path: %s\n", configFilePath)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		// Wrap and return the error if reading the configuration file fails.
		return errs.Wrap(err, fmt.Sprintf("failed to read configuration file at path: %s", configFilePath))
	}

	if err = yaml.Unmarshal(data, config); err != nil {
		// Wrap and return the error if unmarshalling the YAML configuration fails.
		return errs.Wrap(err, "failed to unmarshal YAML configuration")
	}

	fmt.Printf("Configuration file loaded successfully from path: %s\n", configFilePath)
	return nil
}

// InitConfig initializes the application configuration by loading it from a specified folder path.
// If the folder path is not provided, it attempts to use the OPENIMCONFIG environment variable,
// and as a fallback, it uses the default configuration path. It loads both the main configuration
// and notification configuration, wrapping errors for better context.
func InitConfig(configFolderPath string) error {
	// Use the provided config folder path, or fallback to environment variable or default path
	if configFolderPath == "" {
		configFolderPath = os.Getenv("OPENIMCONFIG")
		if configFolderPath == "" {
			var err error
			configFolderPath, err = GetDefaultConfigPath()
			if err != nil {
				return err
			}
		}
	}

	// Initialize the main configuration
	if err := initConfig(&Config, FileName, configFolderPath); err != nil {
		return err
	}

	// Initialize the notification configuration
	if err := initConfig(&Config.Notification, NotificationFileName, configFolderPath); err != nil {
		return err
	}

	return nil
}
