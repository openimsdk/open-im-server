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

import (
	"time"
)

type Applet struct {
	ID         string    `gorm:"column:id;primary_key;size:64"`
	Name       string    `gorm:"column:name;size:64"`
	AppID      string    `gorm:"column:app_id;uniqueIndex;size:255"`
	Icon       string    `gorm:"column:icon;size:255"`
	URL        string    `gorm:"column:url;size:255"`
	MD5        string    `gorm:"column:md5;size:255"`
	Size       int64     `gorm:"column:size"`
	Version    string    `gorm:"column:version;size:64"`
	Priority   uint32    `gorm:"column:priority;size:64"`
	Status     uint8     `gorm:"column:status"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime;size:64"`
}

func (Applet) TableName() string {
	return "applets"
}
