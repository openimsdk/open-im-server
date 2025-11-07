package controller

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/tools/db/pagination"
)

type ThirdDatabase interface {
	FcmUpdateToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) error
	SetAppBadge(ctx context.Context, userID string, value int) error
	// about log for debug
	UploadLogs(ctx context.Context, logs []*model.Log) error
	DeleteLogs(ctx context.Context, logID []string, userID string) error
	SearchLogs(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*model.Log, error)
	GetLogs(ctx context.Context, LogIDs []string, userID string) ([]*model.Log, error)
}

type thirdDatabase struct {
	cache cache.ThirdCache
	logdb database.Log
}

// DeleteLogs implements ThirdDatabase.
func (t *thirdDatabase) DeleteLogs(ctx context.Context, logID []string, userID string) error {
	return t.logdb.Delete(ctx, logID, userID)
}

// GetLogs implements ThirdDatabase.
func (t *thirdDatabase) GetLogs(ctx context.Context, LogIDs []string, userID string) ([]*model.Log, error) {
	return t.logdb.Get(ctx, LogIDs, userID)
}

// SearchLogs implements ThirdDatabase.
func (t *thirdDatabase) SearchLogs(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*model.Log, error) {
	return t.logdb.Search(ctx, keyword, start, end, pagination)
}

// UploadLogs implements ThirdDatabase.
func (t *thirdDatabase) UploadLogs(ctx context.Context, logs []*model.Log) error {
	return t.logdb.Create(ctx, logs)
}

func NewThirdDatabase(cache cache.ThirdCache, logdb database.Log) ThirdDatabase {
	return &thirdDatabase{cache: cache, logdb: logdb}
}

func (t *thirdDatabase) FcmUpdateToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) error {
	return t.cache.SetFcmToken(ctx, account, platformID, fcmToken, expireTime)
}

func (t *thirdDatabase) SetAppBadge(ctx context.Context, userID string, value int) error {
	return t.cache.SetUserBadgeUnreadCountSum(ctx, userID, value)
}
