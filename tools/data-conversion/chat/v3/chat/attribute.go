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

// Attribute 用户属性表.
type Attribute struct {
	UserID           string    `gorm:"column:user_id;primary_key;type:char(64)"`
	Account          string    `gorm:"column:account;type:char(64)"`
	PhoneNumber      string    `gorm:"column:phone_number;type:varchar(32)"`
	AreaCode         string    `gorm:"column:area_code;type:varchar(8)"`
	Email            string    `gorm:"column:email;type:varchar(64)" `
	Nickname         string    `gorm:"column:nickname;type:varchar(64)" `
	FaceURL          string    `gorm:"column:face_url;type:varchar(255)" `
	Gender           int32     `gorm:"column:gender"`
	CreateTime       time.Time `gorm:"column:create_time"`
	ChangeTime       time.Time `gorm:"column:change_time"`
	BirthTime        time.Time `gorm:"column:birth_time"`
	Level            int32     `gorm:"column:level;default:1"`
	AllowVibration   int32     `gorm:"column:allow_vibration;default:1"`
	AllowBeep        int32     `gorm:"column:allow_beep;default:1"`
	AllowAddFriend   int32     `gorm:"column:allow_add_friend;default:1"`
	GlobalRecvMsgOpt int32     `gorm:"column:global_recv_msg_opt;default:0"`
}

func (Attribute) TableName() string {
	return "attributes"
}
