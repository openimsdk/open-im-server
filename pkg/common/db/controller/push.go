package controller

import (
	"OpenIM/pkg/common/db/cache"
	"OpenIM/pkg/proto/sdkws"
	"context"
)

type PushDatabase interface {
	DelFcmToken(ctx context.Context, userID string, platformID int) error
	HandleSignalInvite(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error)
}

type pushDataBase struct {
	cache cache.Model
}

func NewPushDatabase(cache cache.Model) PushDatabase {
	return &pushDataBase{cache: cache}
}

func (p *pushDataBase) DelFcmToken(ctx context.Context, userID string, platformID int) error {
	return p.cache.DelFcmToken(ctx, userID, platformID)
}

func (p *pushDataBase) HandleSignalInvite(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error) {
	return p.cache.HandleSignalInvite(ctx, msg, pushToUserID)
}
