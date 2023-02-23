package controller

import (
	"OpenIM/pkg/common/db/cache"
	"context"
)

type PushInterface interface {
	DelFcmToken(ctx context.Context, userID string, platformID int) error
}

type PushDataBase struct {
	cache cache.MsgCache
}

func (p *PushDataBase) DelFcmToken(ctx context.Context, userID string, platformID int) error {
	return p.cache.DelFcmToken(ctx, userID, platformID)
}
