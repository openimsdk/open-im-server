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
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
	"github.com/openimsdk/tools/utils/datautil"
	"time"

	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type UserDatabase interface {
	// FindWithError Get the information of the specified user. If the userID is not found, it will also return an error
	FindWithError(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	// Find Get the information of the specified user If the userID is not found, no error will be returned
	Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	// Find userInfo By Nickname
	FindByNickname(ctx context.Context, nickname string) (users []*relation.UserModel, err error)
	// Find notificationAccounts
	FindNotification(ctx context.Context, level int64) (users []*relation.UserModel, err error)
	// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the db
	Create(ctx context.Context, users []*relation.UserModel) (err error)
	// UpdateByMap update (zero value) external guarantee userID exists
	UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error)
	// FindUser
	PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error)
	// FindUser with keyword
	PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID string, nickName string, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error)
	// Page If not found, no error is returned
	Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error)
	// IsExist true as long as one exists
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
	// GetAllUserID Get all user IDs
	GetAllUserID(ctx context.Context, pagination pagination.Pagination) (int64, []string, error)
	// Get user by userID
	GetUserByID(ctx context.Context, userID string) (user *relation.UserModel, err error)
	// InitOnce Inside the function, first query whether it exists in the db, if it exists, do nothing; if it does not exist, insert it
	InitOnce(ctx context.Context, users []*relation.UserModel) (err error)
	// CountTotal Get the total number of users
	CountTotal(ctx context.Context, before *time.Time) (int64, error)
	// CountRangeEverydayTotal Get the user increment in the range
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
	// SubscribeUsersStatus Subscribe a user's presence status
	SubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error
	// UnsubscribeUsersStatus unsubscribe a user's presence status
	UnsubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error
	// GetAllSubscribeList Get a list of all subscriptions
	GetAllSubscribeList(ctx context.Context, userID string) ([]string, error)
	// GetSubscribedList Get all subscribed lists
	GetSubscribedList(ctx context.Context, userID string) ([]string, error)
	// GetUserStatus Get the online status of the user
	GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error)
	// SetUserStatus Set the user status and store the user status in redis
	SetUserStatus(ctx context.Context, userID string, status, platformID int32) error

	// CRUD user command
	AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error
	DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error
	UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error
	GetUserCommands(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error)
	GetAllUserCommands(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error)
}

type userDatabase struct {
	tx      tx.Tx
	userDB  relation.UserModelInterface
	cache   cache.UserCache
	mongoDB relation.SubscribeUserModelInterface
}

func NewUserDatabase(userDB relation.UserModelInterface, cache cache.UserCache, tx tx.Tx, mongoDB relation.SubscribeUserModelInterface) UserDatabase {
	return &userDatabase{userDB: userDB, cache: cache, tx: tx, mongoDB: mongoDB}
}

func (u *userDatabase) InitOnce(ctx context.Context, users []*relation.UserModel) error {
	// Extract user IDs from the given user models.
	userIDs := datautil.Slice(users, func(e *relation.UserModel) string {
		return e.UserID
	})

	// Find existing users in the database.
	existingUsers, err := u.userDB.Find(ctx, userIDs)
	if err != nil {
		return err
	}

	// Determine which users are missing from the database.
	missingUsers := datautil.SliceAnySub(users, existingUsers, func(e *relation.UserModel) string {
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
func (u *userDatabase) FindWithError(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
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
func (u *userDatabase) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	return u.cache.GetUsersInfo(ctx, userIDs)
}

func (u *userDatabase) FindByNickname(ctx context.Context, nickname string) (users []*relation.UserModel, err error) {
	return u.userDB.TakeByNickname(ctx, nickname)
}

func (u *userDatabase) FindNotification(ctx context.Context, level int64) (users []*relation.UserModel, err error) {
	return u.userDB.TakeNotification(ctx, level)
}

// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the db.
func (u *userDatabase) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = u.userDB.Create(ctx, users); err != nil {
			return err
		}
		return u.cache.DelUsersInfo(datautil.Slice(users, func(e *relation.UserModel) string {
			return e.UserID
		})...).ExecDel(ctx)
	})
}

// UpdateByMap update (zero value) externally guarantees that userID exists.
func (u *userDatabase) UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := u.userDB.UpdateByMap(ctx, userID, args); err != nil {
			return err
		}
		return u.cache.DelUsersInfo(userID).ExecDel(ctx)
	})
}

// Page Gets, returns no error if not found.
func (u *userDatabase) Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	return u.userDB.Page(ctx, pagination)
}

func (u *userDatabase) PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	return u.userDB.PageFindUser(ctx, level1, level2, pagination)
}

func (u *userDatabase) PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID, nickName string, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
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

func (u *userDatabase) GetUserByID(ctx context.Context, userID string) (user *relation.UserModel, err error) {
	return u.userDB.Take(ctx, userID)
}

// CountTotal Get the total number of users.
func (u *userDatabase) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	return u.userDB.CountTotal(ctx, before)
}

// CountRangeEverydayTotal Get the user increment in the range.
func (u *userDatabase) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	return u.userDB.CountRangeEverydayTotal(ctx, start, end)
}

// SubscribeUsersStatus Subscribe or unsubscribe a user's presence status.
func (u *userDatabase) SubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error {
	err := u.mongoDB.AddSubscriptionList(ctx, userID, userIDs)
	return err
}

// UnsubscribeUsersStatus unsubscribe a user's presence status.
func (u *userDatabase) UnsubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error {
	err := u.mongoDB.UnsubscriptionList(ctx, userID, userIDs)
	return err
}

// GetAllSubscribeList Get a list of all subscriptions.
func (u *userDatabase) GetAllSubscribeList(ctx context.Context, userID string) ([]string, error) {
	list, err := u.mongoDB.GetAllSubscribeList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetSubscribedList Get all subscribed lists.
func (u *userDatabase) GetSubscribedList(ctx context.Context, userID string) ([]string, error) {
	list, err := u.mongoDB.GetSubscribedList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetUserStatus get user status.
func (u *userDatabase) GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error) {
	onlineStatusList, err := u.cache.GetUserStatus(ctx, userIDs)
	return onlineStatusList, err
}

// SetUserStatus Set the user status and save it in redis.
func (u *userDatabase) SetUserStatus(ctx context.Context, userID string, status, platformID int32) error {
	return u.cache.SetUserStatus(ctx, userID, status, platformID)
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
