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

package relation

import (
	"context"
	"time"
)

const (
	UserModelTableName = "users"
)

type UserModel struct {
	UserID           string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname         string    `gorm:"column:name;size:255"`
	FaceURL          string    `gorm:"column:face_url;size:255"`
	Ex               string    `gorm:"column:ex;size:1024"`
	CreateTime       time.Time `gorm:"column:create_time;index:create_time;autoCreateTime"`
	AppMangerLevel   int32     `gorm:"column:app_manger_level;default:1"`
	GlobalRecvMsgOpt int32     `gorm:"column:global_recv_msg_opt"`
}

func (u *UserModel) GetNickname() string {
	return u.Nickname
}

func (u *UserModel) GetFaceURL() string {
	return u.FaceURL
}

func (u *UserModel) GetUserID() string {
	return u.UserID
}

func (u *UserModel) GetEx() string {
	return u.Ex
}

func (UserModel) TableName() string {
	return UserModelTableName
}

// UserModelInterface defines the operations available for managing user models.
type UserModelInterface interface {
	// Create inserts a new user or multiple users into the database.
	Create(ctx context.Context, users []*UserModel) (err error)

	// UpdateByMap updates a user's information based on a map of changes.
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)

	// Update modifies a user's information in the database.
	Update(ctx context.Context, user *UserModel) (err error)

	// Find retrieves information for a list of users by their IDs. If a user does not exist, it is simply skipped without returning an error.
	Find(ctx context.Context, userIDs []string) (users []*UserModel, err error)

	// Take retrieves a specific user's information by their ID. Returns an error if the user does not exist.
	Take(ctx context.Context, userID string) (user *UserModel, err error)

	// Page retrieves a paginated list of users and the total count of users. If no users exist, returns an empty list without an error.
	Page(ctx context.Context, pageNumber, showNumber int32) (users []*UserModel, count int64, err error)

	// GetAllUserID retrieves all user IDs in a paginated manner.
	GetAllUserID(ctx context.Context, pageNumber, showNumber int32) (userIDs []string, err error)

	// GetUserGlobalRecvMsgOpt retrieves a user's global message receiving option.
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)

	// CountTotal returns the total number of users before a specified time.
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)

	// CountRangeEverydayTotal calculates the daily increment of users within a specified time range.
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
}
