package controller

import (
	"OpenIM/pkg/common/db/cache"
	"OpenIM/pkg/proto/sdkws"
	"context"
)

type ThirdDatabase interface {
	GetSignalInvitationInfoByClientMsgID(ctx context.Context, clientMsgID string) (invitationInfo *sdkws.SignalInviteReq, err error)
	GetAvailableSignalInvitationInfo(ctx context.Context, userID string) (invitationInfo *sdkws.SignalInviteReq, err error)
	FcmUpdateToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) error
	SetAppBadge(ctx context.Context, userID string, value int) error
}

type thirdDatabase struct {
	cache cache.Cache
}

func NewThirdDatabase(cache cache.Cache) ThirdDatabase {
	return &thirdDatabase{cache: cache}
}

func (t *thirdDatabase) GetSignalInvitationInfoByClientMsgID(ctx context.Context, clientMsgID string) (invitationInfo *sdkws.SignalInviteReq, err error) {
	return t.cache.GetSignalInvitationInfoByClientMsgID(ctx, clientMsgID)
}

func (t *thirdDatabase) GetAvailableSignalInvitationInfo(ctx context.Context, userID string) (invitationInfo *sdkws.SignalInviteReq, err error) {
	return t.cache.GetAvailableSignalInvitationInfo(ctx, userID)
}

func (t *thirdDatabase) FcmUpdateToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) error {
	return t.cache.SetFcmToken(ctx, account, platformID, fcmToken, expireTime)
}

func (t *thirdDatabase) SetAppBadge(ctx context.Context, userID string, value int) error {
	return t.cache.SetUserBadgeUnreadCountSum(ctx, userID, value)
}
