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
	"reflect"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
)

func TestGetDefaultConfigPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultConfigPath(); got != tt.want {
				t.Errorf("GetDefaultConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProjectRoot(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjectRoot(); got != tt.want {
				t.Errorf("GetProjectRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOptionsByNotification(t *testing.T) {
	type args struct {
		cfg NotificationConf
	}
	tests := []struct {
		name string
		args args
		want msgprocessor.Options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOptionsByNotification(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOptionsByNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initConfig(t *testing.T) {
	type args struct {
		config           interface{}
		configName       string
		configFolderPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initConfig(tt.args.config, tt.args.configName, tt.args.configFolderPath); (err != nil) != tt.wantErr {
				t.Errorf("initConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	type args struct {
		configFolderPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitConfig(tt.args.configFolderPath); (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
