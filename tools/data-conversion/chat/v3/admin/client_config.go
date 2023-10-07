// Copyright © 2023 OpenIM open source community. All rights reserved.
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

package admin

// ClientConfig 客户端相关配置项.
type ClientConfig struct {
	Key   string `gorm:"column:key;primary_key;type:varchar(255)"`
	Value string `gorm:"column:value;not null;type:text"`
}

func (ClientConfig) TableName() string {
	return "client_config"
}
