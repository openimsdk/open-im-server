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

package database

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type User interface {
	Create(ctx context.Context, users []*model.User) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error)
	Find(ctx context.Context, userIDs []string) (users []*model.User, err error)
	Take(ctx context.Context, userID string) (user *model.User, err error)
	TakeNotification(ctx context.Context, level int64) (user []*model.User, err error)
	TakeByNickname(ctx context.Context, nickname string) (user []*model.User, err error)
	Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*model.User, err error)
	PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*model.User, err error)
	PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID, nickName string, pagination pagination.Pagination) (count int64, users []*model.User, err error)
	Exist(ctx context.Context, userID string) (exist bool, err error)
	GetAllUserID(ctx context.Context, pagination pagination.Pagination) (count int64, userIDs []string, err error)
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	// Get user total quantity
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// Get user total quantity every day
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)

	SortQuery(ctx context.Context, userIDName map[string]string, asc bool) ([]*model.User, error)

	// CRUD user command
	AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error
	DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error
	UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error
	GetUserCommand(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error)
	GetAllUserCommand(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error)
}
