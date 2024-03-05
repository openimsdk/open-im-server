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

// return absolude path join ../config/, this is k8s container config path.
func GetDefaultConfigPath() string {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("GetDefaultConfigPath error:", err.Error())
		return ""
	}

	configPath, err := genutil.OutDir(filepath.Join(filepath.Dir(executablePath), "../config/"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}
	return configPath
}

// getProjectRoot returns the absolute path of the project root directory.
func GetProjectRoot() string {
	executablePath, _ := os.Executable()

	projectRoot, err := genutil.OutDir(filepath.Join(filepath.Dir(executablePath), "../../../../.."))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}
	return projectRoot
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
	configFolderPath = filepath.Join(configFolderPath, configName)
	_, err := os.Stat(configFolderPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("stat config path error:", err.Error())
			return fmt.Errorf("stat config path error: %w", err)
		}
		configFolderPath = filepath.Join(GetProjectRoot(), "config", configName)
		fmt.Println("flag's path,enviment's path,default path all is not exist,using project path:", configFolderPath)
	}
	data, err := os.ReadFile(configFolderPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("unmarshal yaml error: %w", err)
	}
	fmt.Println("The path of the configuration file to start the process:", configFolderPath)

	return nil
}

func InitConfig(config *GlobalConfig, configFolderPath string) error {
	if configFolderPath == "" {
		envConfigPath := os.Getenv("OPENIMCONFIG")
		if envConfigPath != "" {
			configFolderPath = envConfigPath
		} else {
			configFolderPath = GetDefaultConfigPath()
		}
	}

	if err := initConfig(config, FileName, configFolderPath); err != nil {
		return err
	}

	return initConfig(&config.Notification, NotificationFileName, configFolderPath)
}
