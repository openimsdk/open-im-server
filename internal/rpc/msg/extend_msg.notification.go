package msg

import (
	"Open_IM/pkg/common/constant"
	pbFriend "Open_IM/pkg/proto/friend"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
)

func ExtendMessageUpdatedNotification(operationID, changedUserID string, needNotifiedUserID string, opUserID string) {
	selfInfoUpdatedTips := open_im_sdk.UserInfoUpdatedTips{UserID: changedUserID}
	commID := pbFriend.CommID{FromUserID: opUserID, ToUserID: needNotifiedUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(&commID, constant.ReactionMessageModifierNotification, &selfInfoUpdatedTips)
}
