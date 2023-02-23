package notification

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/proto/sdkws"
	"context"
)

// send to myself
func (c *Check) UserInfoUpdatedNotification(ctx context.Context, opUserID string, changedUserID string) {
	selfInfoUpdatedTips := sdkws.UserInfoUpdatedTips{UserID: changedUserID}
	c.friendNotification(ctx, opUserID, changedUserID, constant.UserInfoUpdatedNotification, &selfInfoUpdatedTips)
}
