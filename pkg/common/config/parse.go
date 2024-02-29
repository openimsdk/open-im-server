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
	"gopkg.in/yaml.v3"

	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/open-im-server/v3/pkg/util/genutil"
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

func initConfig(config any, configName, configFolderPath string) error {
	configFolderPath = filepath.Join(configFolderPath, configName)
	_, err := os.Stat(configFolderPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errs.Wrap(err, "stat config path error")
		}
		configFolderPath = filepath.Join(GetProjectRoot(), "config", configName)
		fmt.Println("flag's path,enviment's path,default path all is not exist,using project path:", configFolderPath)
	}
	data, err := os.ReadFile(configFolderPath)
	if err != nil {
		return errs.Wrap(err, "read file error")
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return errs.Wrap(err, "unmarshal yaml error")
	}
	fmt.Println("The path of the configuration file to start the process:", configFolderPath)

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
				// Wrap and return the error if getting the default config path fails
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
