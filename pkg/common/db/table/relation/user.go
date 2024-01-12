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
	"github.com/OpenIMSDK/protocol/user"
	"time"

	"github.com/OpenIMSDK/tools/pagination"
)

type UserModel struct {
	UserID           string    `bson:"user_id"`
	Nickname         string    `bson:"nickname"`
	FaceURL          string    `bson:"face_url"`
	Ex               string    `bson:"ex"`
	AppMangerLevel   int32     `bson:"app_manger_level"`
	GlobalRecvMsgOpt int32     `bson:"global_recv_msg_opt"`
	CreateTime       time.Time `bson:"create_time"`
}

func (u *UserModel) GetNickname() string {
	return u.Nickname
}

func (u *UserModel) GetFaceURL() string {
	return u.FaceURL
}

func (u UserModel) GetUserID() string {
	return u.UserID
}

func (u UserModel) GetEx() string {
	return u.Ex
}

type UserModelInterface interface {
	Create(ctx context.Context, users []*UserModel) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error)
	Find(ctx context.Context, userIDs []string) (users []*UserModel, err error)
	Take(ctx context.Context, userID string) (user *UserModel, err error)
	TakeNotification(ctx context.Context, level int64) (user []*UserModel, err error)
	TakeByNickname(ctx context.Context, nickname string) (user []*UserModel, err error)
	Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*UserModel, err error)
	PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*UserModel, err error)
	PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID, userName string, pagination pagination.Pagination) (count int64, users []*UserModel, err error)
	Exist(ctx context.Context, userID string) (exist bool, err error)
	GetAllUserID(ctx context.Context, pagination pagination.Pagination) (count int64, userIDs []string, err error)
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	// 获取用户总数
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// 获取范围内用户增量
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
	//CRUD user command
	AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error
	DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error
	UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error
	GetUserCommand(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error)
	GetAllUserCommand(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error)
}
