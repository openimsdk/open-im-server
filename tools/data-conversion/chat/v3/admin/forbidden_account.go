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

// ForbiddenAccount 封号表.
type ForbiddenAccount struct {
	UserID         string    `gorm:"column:user_id;index:userID;primary_key;type:char(64)"`
	Reason         string    `gorm:"column:reason;type:varchar(255)" `
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(255)"`
	CreateTime     time.Time `gorm:"column:create_time" `
}

func (ForbiddenAccount) TableName() string {
	return "forbidden_accounts"
}
