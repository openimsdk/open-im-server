package controller

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
)

type PushDatabase interface {
	DelFcmToken(ctx context.Context, userID string, platformID int) error
}

type pushDataBase struct {
	cache cache.MsgModel
}

func NewPushDatabase(cache cache.MsgModel) PushDatabase {
	return &pushDataBase{cache: cache}
}

func (p *pushDataBase) DelFcmToken(ctx context.Context, userID string, platformID int) error {
	return p.cache.DelFcmToken(ctx, userID, platformID)
}
