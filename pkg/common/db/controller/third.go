package controller

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
)

type ThirdDatabase interface {
	FcmUpdateToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) error
	SetAppBadge(ctx context.Context, userID string, value int) error
}

type thirdDatabase struct {
	cache cache.MsgModel
}

func NewThirdDatabase(cache cache.MsgModel) ThirdDatabase {
	return &thirdDatabase{cache: cache}
}

func (t *thirdDatabase) FcmUpdateToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) error {
	return t.cache.SetFcmToken(ctx, account, platformID, fcmToken, expireTime)
}

func (t *thirdDatabase) SetAppBadge(ctx context.Context, userID string, value int) error {
	return t.cache.SetUserBadgeUnreadCountSum(ctx, userID, value)
}
