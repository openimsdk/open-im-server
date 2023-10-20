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
	"runtime"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"gopkg.in/yaml.v3"
)

//go:embed version
var Version string

const (
	FileName             = "config.yaml"
	NotificationFileName = "notification.yaml"
	DefaultFolderPath    = "../config/"
)

// getProjectRoot returns the absolute path of the project root directory
func GetProjectRoot() string {
	// Program counter (PC): This represents the address of the function.
	// File path: The full path to the source file from which the function was called. (pkg/common/config/parse.go)
	// Line number: The line number in the source file from which the function was called.
	// Success flag: it will be true if the information was successfully fetched, false otherwise.
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "../../..")
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

func initConfig(config interface{}, configName, configFolderPath string) error {
	configFolderPath = filepath.Join(configFolderPath, configName)
	_, err := os.Stat(configFolderPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat config path error: %w", err)
		}
		configFolderPath = filepath.Join(GetProjectRoot(), "config", configName)
	}
	data, err := os.ReadFile(configFolderPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("unmarshal yaml error: %w", err)
	}
	fmt.Println("use config", configFolderPath)
	return nil
}

func InitConfig(configFolderPath string) error {
	if configFolderPath == "" {
		configFolderPath = DefaultFolderPath
	}
	if err := initConfig(&Config, FileName, configFolderPath); err != nil {
		return err
	}
	return initConfig(&Config.Notification, NotificationFileName, configFolderPath)
}
