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

package controller

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
)

type UserDatabase interface {
	// FindWithError Get the information of the specified user. If the userID is not found, it will also return an error
	FindWithError(ctx context.Context, userIDs []string) (users []*model.User, err error)
	// Find Get the information of the specified user If the userID is not found, no error will be returned
	Find(ctx context.Context, userIDs []string) (users []*model.User, err error)
	// Find userInfo By Nickname
	FindByNickname(ctx context.Context, nickname string) (users []*model.User, err error)
	// FindNotification find system account by level
	FindNotification(ctx context.Context, level int64) (users []*model.User, err error)
	// FindSystemAccount find all system account
	FindSystemAccount(ctx context.Context) (users []*model.User, err error)
	// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the storage
	Create(ctx context.Context, users []*model.User) (err error)
	// UpdateByMap update (zero value) external guarantee userID exists
	UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error)
	// FindUser
	PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*model.User, err error)
	// FindUser with keyword
	PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID string, nickName string, pagination pagination.Pagination) (count int64, users []*model.User, err error)
	// Page If not found, no error is returned
	Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*model.User, err error)
	// IsExist true as long as one exists
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
	// GetAllUserID Get all user IDs
	GetAllUserID(ctx context.Context, pagination pagination.Pagination) (int64, []string, error)
	// Get user by userID
	GetUserByID(ctx context.Context, userID string) (user *model.User, err error)
	// InitOnce Inside the function, first query whether it exists in the storage, if it exists, do nothing; if it does not exist, insert it
	InitOnce(ctx context.Context, users []*model.User) (err error)
	// CountTotal Get the total number of users
	CountTotal(ctx context.Context, before *time.Time) (int64, error)
	// CountRangeEverydayTotal Get the user increment in the range
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)

	SortQuery(ctx context.Context, userIDName map[string]string, asc bool) ([]*model.User, error)

	// CRUD user command
	AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error
	DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error
	UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error
	GetUserCommands(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error)
	GetAllUserCommands(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error)
}

type userDatabase struct {
	tx     tx.Tx
	userDB database.User
	cache  cache.UserCache
}

func NewUserDatabase(userDB database.User, cache cache.UserCache, tx tx.Tx) UserDatabase {
	return &userDatabase{userDB: userDB, cache: cache, tx: tx}
}

func (u *userDatabase) InitOnce(ctx context.Context, users []*model.User) error {
	// Extract user IDs from the given user models.
	userIDs := datautil.Slice(users, func(e *model.User) string {
		return e.UserID
	})

	// Find existing users in the database.
	existingUsers, err := u.userDB.Find(ctx, userIDs)
	if err != nil {
		return err
	}

	// Determine which users are missing from the database.
	missingUsers := datautil.SliceAnySub(users, existingUsers, func(e *model.User) string {
		return e.UserID
	})

	// Create records for missing users.
	if len(missingUsers) > 0 {
		if err := u.userDB.Create(ctx, missingUsers); err != nil {
			return err
		}
	}

	return nil
}

// FindWithError Get the information of the specified user and return an error if the userID is not found.
func (u *userDatabase) FindWithError(ctx context.Context, userIDs []string) (users []*model.User, err error) {
	userIDs = datautil.Distinct(userIDs)

	// TODO: Add logic to identify which user IDs are distinct and which user IDs were not found.

	users, err = u.cache.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return
	}

	if len(users) != len(userIDs) {
		err = errs.ErrRecordNotFound.WrapMsg("userID not found")
	}
	return
}

// Find Get the information of the specified user. If the userID is not found, no error will be returned.
func (u *userDatabase) Find(ctx context.Context, userIDs []string) (users []*model.User, err error) {
	return u.cache.GetUsersInfo(ctx, userIDs)
}

func (u *userDatabase) FindByNickname(ctx context.Context, nickname string) (users []*model.User, err error) {
	return u.userDB.TakeByNickname(ctx, nickname)
}

func (u *userDatabase) FindNotification(ctx context.Context, level int64) (users []*model.User, err error) {
	return u.userDB.TakeNotification(ctx, level)
}

func (u *userDatabase) FindSystemAccount(ctx context.Context) (users []*model.User, err error) {
	return u.userDB.TakeGTEAppManagerLevel(ctx, constant.AppNotificationAdmin)
}

// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the storage.
func (u *userDatabase) Create(ctx context.Context, users []*model.User) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = u.userDB.Create(ctx, users); err != nil {
			return err
		}
		return u.cache.DelUsersInfo(datautil.Slice(users, func(e *model.User) string {
			return e.UserID
		})...).ChainExecDel(ctx)
	})
}

// UpdateByMap update (zero value) externally guarantees that userID exists.
func (u *userDatabase) UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := u.userDB.UpdateByMap(ctx, userID, args); err != nil {
			return err
		}
		return u.cache.DelUsersInfo(userID).ChainExecDel(ctx)
	})
}

// Page Gets, returns no error if not found.
func (u *userDatabase) Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*model.User, err error) {
	return u.userDB.Page(ctx, pagination)
}

func (u *userDatabase) PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*model.User, err error) {
	return u.userDB.PageFindUser(ctx, level1, level2, pagination)
}

func (u *userDatabase) PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID, nickName string, pagination pagination.Pagination) (count int64, users []*model.User, err error) {
	return u.userDB.PageFindUserWithKeyword(ctx, level1, level2, userID, nickName, pagination)
}

// IsExist Does userIDs exist? As long as there is one, it will be true.
func (u *userDatabase) IsExist(ctx context.Context, userIDs []string) (exist bool, err error) {
	users, err := u.userDB.Find(ctx, userIDs)
	if err != nil {
		return false, err
	}
	if len(users) > 0 {
		return true, nil
	}
	return false, nil
}

// GetAllUserID Get all user IDs.
func (u *userDatabase) GetAllUserID(ctx context.Context, pagination pagination.Pagination) (total int64, userIDs []string, err error) {
	return u.userDB.GetAllUserID(ctx, pagination)
}

func (u *userDatabase) GetUserByID(ctx context.Context, userID string) (user *model.User, err error) {
	return u.cache.GetUserInfo(ctx, userID)
}

// CountTotal Get the total number of users.
func (u *userDatabase) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	return u.userDB.CountTotal(ctx, before)
}

// CountRangeEverydayTotal Get the user increment in the range.
func (u *userDatabase) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	return u.userDB.CountRangeEverydayTotal(ctx, start, end)
}

func (u *userDatabase) SortQuery(ctx context.Context, userIDName map[string]string, asc bool) ([]*model.User, error) {
	return u.userDB.SortQuery(ctx, userIDName, asc)
}

func (u *userDatabase) AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error {
	return u.userDB.AddUserCommand(ctx, userID, Type, UUID, value, ex)
}

func (u *userDatabase) DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error {
	return u.userDB.DeleteUserCommand(ctx, userID, Type, UUID)
}

func (u *userDatabase) UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error {
	return u.userDB.UpdateUserCommand(ctx, userID, Type, UUID, val)
}

func (u *userDatabase) GetUserCommands(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error) {
	commands, err := u.userDB.GetUserCommand(ctx, userID, Type)
	return commands, err
}

func (u *userDatabase) GetAllUserCommands(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error) {
	commands, err := u.userDB.GetAllUserCommand(ctx, userID)
	return commands, err
}
