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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

const (
	FileName             = "config.yaml"
	NotificationFileName = "notification.yaml"
	ENV                  = "CONFIG_NAME"
	DefaultFolderPath    = "../config/"
	ConfKey              = "conf"
)

func GetOptionsByNotification(cfg NotificationConf) utils.Options {
	opts := utils.NewOptions()
	if cfg.UnreadCount {
		opts = utils.WithOptions(opts, utils.WithUnreadCount(true))
	}
	if cfg.OfflinePush.Enable {
		opts = utils.WithOptions(opts, utils.WithOfflinePush(true))
	}
	switch cfg.ReliabilityLevel {
	case constant.UnreliableNotification:
	case constant.ReliableNotificationNoMsg:
		opts = utils.WithOptions(opts, utils.WithHistory(true), utils.WithPersistent())
	}
	opts = utils.WithOptions(opts, utils.WithSendMsg(cfg.IsSendMsg))
	return opts
}

func (c *config) unmarshalConfig(config interface{}, configPath string) error {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(bytes, config); err != nil {
		return err
	}
	return nil
}

func (c *config) initConfig(config interface{}, configName, configFolderPath string) error {
	if configFolderPath == "" {
		configFolderPath = DefaultFolderPath
	}
	configPath := filepath.Join(configFolderPath, configName)
	defer func() {
		fmt.Println("use config", configPath)
	}()
	_, err := os.Stat(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		configPath = filepath.Join(Root, "config", configName)
	} else {
		Root = filepath.Dir(configPath)
	}
	return c.unmarshalConfig(config, configPath)
}

func (c *config) RegisterConf2Registry(registry discoveryregistry.SvcDiscoveryRegistry) error {
	bytes, err := yaml.Marshal(Config)
	if err != nil {
		return err
	}
	return registry.RegisterConf2Registry(ConfKey, bytes)
}

func (c *config) GetConfFromRegistry(registry discoveryregistry.SvcDiscoveryRegistry) ([]byte, error) {
	return registry.GetConfFromRegistry(ConfKey)
}

func InitConfig(configFolderPath string) error {
	err := Config.initConfig(&Config, FileName, configFolderPath)
	if err != nil {
		return err
	}
	err = Config.initConfig(&Config.Notification, NotificationFileName, configFolderPath)
	if err != nil {
		return err
	}
	return nil
}

func EncodeConfig() []byte {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(Config); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
