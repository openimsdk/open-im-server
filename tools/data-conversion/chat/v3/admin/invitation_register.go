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

// 邀请码被注册使用.
type InvitationRegister struct {
	InvitationCode string    `gorm:"column:invitation_code;primary_key;type:char(32)"`
	UsedByUserID   string    `gorm:"column:user_id;index:userID;type:char(64)"`
	CreateTime     time.Time `gorm:"column:create_time"`
}

func (InvitationRegister) TableName() string {
	return "invitation_registers"
}
