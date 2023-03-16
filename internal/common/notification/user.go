package notification

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

// send to myself
func (c *Check) UserInfoUpdatedNotification(ctx context.Context, opUserID string, changedUserID string) {
	selfInfoUpdatedTips := sdkws.UserInfoUpdatedTips{UserID: changedUserID}
	c.friendNotification(ctx, opUserID, changedUserID, constant.UserInfoUpdatedNotification, &selfInfoUpdatedTips)
}
