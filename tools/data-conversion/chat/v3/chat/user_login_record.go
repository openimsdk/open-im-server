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

package chat

import (
	"time"
)

// 用户登录信息表.
type UserLoginRecord struct {
	UserID    string    `gorm:"column:user_id;size:64"`
	LoginTime time.Time `gorm:"column:login_time"`
	IP        string    `gorm:"column:ip;type:varchar(32)"`
	DeviceID  string    `gorm:"column:device_id;type:varchar(255)"`
	Platform  string    `gorm:"column:platform;type:varchar(32)"`
}

func (UserLoginRecord) TableName() string {
	return "user_login_records"
}
