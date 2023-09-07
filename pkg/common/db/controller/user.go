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

	"github.com/OpenIMSDK/protocol/user"

	unrelationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/unrelation"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/tx"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type UserDatabase interface {
	// FindWithError Get the information of the specified user. If the userID is not found, it will also return an error
	FindWithError(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	// Find Get the information of the specified user If the userID is not found, no error will be returned
	Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the db
	Create(ctx context.Context, users []*relation.UserModel) (err error)
	// Update update (non-zero value) external guarantee userID exists
	Update(ctx context.Context, user *relation.UserModel) (err error)
	// UpdateByMap update (zero value) external guarantee userID exists
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	// Page If not found, no error is returned
	Page(ctx context.Context, pageNumber, showNumber int32) (users []*relation.UserModel, count int64, err error)
	// IsExist true as long as one exists
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
	// GetAllUserID Get all user IDs
	GetAllUserID(ctx context.Context, pageNumber, showNumber int32) ([]string, error)
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
	SetUserStatus(ctx context.Context, list []*user.OnlineStatus) error
}

type userDatabase struct {
	userDB  relation.UserModelInterface
	cache   cache.UserCache
	tx      tx.Tx
	mongoDB unrelationtb.UserModelInterface
}

func NewUserDatabase(userDB relation.UserModelInterface, cache cache.UserCache, tx tx.Tx, mongoDB unrelationtb.UserModelInterface) UserDatabase {
	return &userDatabase{userDB: userDB, cache: cache, tx: tx, mongoDB: mongoDB}
}

func (u *userDatabase) InitOnce(ctx context.Context, users []*relation.UserModel) (err error) {
	userIDs := utils.Slice(users, func(e *relation.UserModel) string {
		return e.UserID
	})
	result, err := u.userDB.Find(ctx, userIDs)
	if err != nil {
		return err
	}
	miss := utils.SliceAnySub(users, result, func(e *relation.UserModel) string { return e.UserID })
	if len(miss) > 0 {
		_ = u.userDB.Create(ctx, miss)
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
		err = errs.ErrRecordNotFound.Wrap("userID not found")
	}
	return
}

// Find Get the information of the specified user. If the userID is not found, no error will be returned.
func (u *userDatabase) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	users, err = u.cache.GetUsersInfo(ctx, userIDs)
	return
}

// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the db.
func (u *userDatabase) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	if err := u.tx.Transaction(func(tx any) error {
		err = u.userDB.Create(ctx, users)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	var userIDs []string
	for _, user := range users {
		userIDs = append(userIDs, user.UserID)
	}
	return u.cache.DelUsersInfo(userIDs...).ExecDel(ctx)
}

// Update (non-zero value) externally guarantees that userID exists.
func (u *userDatabase) Update(ctx context.Context, user *relation.UserModel) (err error) {
	if err := u.userDB.Update(ctx, user); err != nil {
		return err
	}
	return u.cache.DelUsersInfo(user.UserID).ExecDel(ctx)
}

// UpdateByMap update (zero value) externally guarantees that userID exists.
func (u *userDatabase) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	if err := u.userDB.UpdateByMap(ctx, userID, args); err != nil {
		return err
	}
	return u.cache.DelUsersInfo(userID).ExecDel(ctx)
}

// Page Gets, returns no error if not found.
func (u *userDatabase) Page(
	ctx context.Context,
	pageNumber, showNumber int32,
) (users []*relation.UserModel, count int64, err error) {
	return u.userDB.Page(ctx, pageNumber, showNumber)
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
func (u *userDatabase) GetAllUserID(ctx context.Context, pageNumber, showNumber int32) (userIDs []string, err error) {
	return u.userDB.GetAllUserID(ctx, pageNumber, showNumber)
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
func (u *userDatabase) SetUserStatus(ctx context.Context, list []*user.OnlineStatus) error {
	return u.cache.SetUserStatus(ctx, list)
}
