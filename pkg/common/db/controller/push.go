package controller

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type PushDatabase interface {
	DelFcmToken(ctx context.Context, userID string, platformID int) error
	HandleSignalInvite(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error)
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

func (p *pushDataBase) HandleSignalInvite(ctx context.Context, msg *sdkws.MsgData, pushToUserID string) (isSend bool, err error) {
	return p.cache.HandleSignalInvite(ctx, msg, pushToUserID)
}
