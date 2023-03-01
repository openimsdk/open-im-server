package controller

import (
	"OpenIM/pkg/common/db/cache"
	"context"
)

type PushDatabase interface {
	DelFcmToken(ctx context.Context, userID string, platformID int) error
}

type pushDataBase struct {
	cache cache.Cache
}

func NewPushDatabase(cache cache.Cache) PushDatabase {
	return &pushDataBase{cache: cache}
}

func (p *pushDataBase) DelFcmToken(ctx context.Context, userID string, platformID int) error {
	return p.cache.DelFcmToken(ctx, userID, platformID)
}
