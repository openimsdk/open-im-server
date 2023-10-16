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

	"gorm.io/gorm"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	dbimpl "github.com/openimsdk/open-im-server/v3/pkg/common/db/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type ThirdDatabase interface {
	FcmUpdateToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) error
	SetAppBadge(ctx context.Context, userID string, value int) error
	// about log for debug
	UploadLogs(ctx context.Context, logs []*relation.Log) error
	DeleteLogs(ctx context.Context, logID []string, userID string) error
	SearchLogs(ctx context.Context, keyword string, start time.Time, end time.Time, pageNumber int32, showNumber int32) (uint32, []*relation.Log, error)
	GetLogs(ctx context.Context, LogIDs []string, userID string) ([]*relation.Log, error)
	FindUsers(ctx context.Context, userIDs []string) ([]*relation.UserModel, error)
}

type thirdDatabase struct {
	cache  cache.MsgModel
	logdb  relation.LogInterface
	userdb relation.UserModelInterface
}

// FindUsers implements ThirdDatabase.
func (t *thirdDatabase) FindUsers(ctx context.Context, userIDs []string) ([]*relation.UserModel, error) {
	return t.userdb.Find(ctx, userIDs)
}

// DeleteLogs implements ThirdDatabase.
func (t *thirdDatabase) DeleteLogs(ctx context.Context, logID []string, userID string) error {
	return t.logdb.Delete(ctx, logID, userID)
}

// GetLogs implements ThirdDatabase.
func (t *thirdDatabase) GetLogs(ctx context.Context, LogIDs []string, userID string) ([]*relation.Log, error) {
	return t.logdb.Get(ctx, LogIDs, userID)
}

// SearchLogs implements ThirdDatabase.
func (t *thirdDatabase) SearchLogs(ctx context.Context, keyword string, start time.Time, end time.Time, pageNumber int32, showNumber int32) (uint32, []*relation.Log, error) {
	return t.logdb.Search(ctx, keyword, start, end, pageNumber, showNumber)
}

// UploadLogs implements ThirdDatabase.
func (t *thirdDatabase) UploadLogs(ctx context.Context, logs []*relation.Log) error {
	return t.logdb.Create(ctx, logs)
}

func NewThirdDatabase(cache cache.MsgModel, db *gorm.DB) ThirdDatabase {
	return &thirdDatabase{cache: cache, logdb: dbimpl.NewLogGorm(db), userdb: dbimpl.NewUserGorm(db)}
}

func (t *thirdDatabase) FcmUpdateToken(
	ctx context.Context,
	account string,
	platformID int,
	fcmToken string,
	expireTime int64,
) error {
	return t.cache.SetFcmToken(ctx, account, platformID, fcmToken, expireTime)
}

func (t *thirdDatabase) SetAppBadge(ctx context.Context, userID string, value int) error {
	return t.cache.SetUserBadgeUnreadCountSum(ctx, userID, value)
}
