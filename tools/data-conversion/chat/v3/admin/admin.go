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

// Admin 后台管理员.
type Admin struct {
	Account    string    `gorm:"column:account;primary_key;type:varchar(64)"`
	Password   string    `gorm:"column:password;type:varchar(64)"`
	FaceURL    string    `gorm:"column:face_url;type:varchar(255)"`
	Nickname   string    `gorm:"column:nickname;type:varchar(64)"`
	UserID     string    `gorm:"column:user_id;type:varchar(64)"` // openIM userID
	Level      int32     `gorm:"column:level;default:1"  `
	CreateTime time.Time `gorm:"column:create_time"`
}

func (Admin) TableName() string {
	return "admins"
}
