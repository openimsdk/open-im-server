package sdkws

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

func (x *MsgData) Check() error {
	if x.SendID == "" {
		return errs.ErrArgs.Wrap("sendID is empty")
	}
	if x.Content == nil {
		return errs.ErrArgs.Wrap("content is empty")
	}
	if x.ContentType < 101 || x.ContentType > 203 {
		return errs.ErrArgs.Wrap("content is empty")
	}
	if x.SessionType < 1 || x.SessionType > 4 {
		return errs.ErrArgs.Wrap("sessionType is invalid")
	}
	if x.SessionType == constant.SingleChatType || x.SessionType == constant.NotificationChatType {
		if x.RecvID == "" {
			return errs.ErrArgs.Wrap("recvID is empty")
		}
	}
	if x.SessionType == constant.GroupChatType || x.SessionType == constant.SuperGroupChatType {
		if x.GroupID == "" {
			return errs.ErrArgs.Wrap("GroupID is empty")
		}
	}
	return nil
}
