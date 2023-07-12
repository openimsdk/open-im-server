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
	if x.ContentType <= constant.ContentTypeBegin || x.ContentType >= constant.NotificationEnd {
		return errs.ErrArgs.Wrap("content type is invalid")
	}
	if x.SessionType < constant.SingleChatType || x.SessionType > constant.NotificationChatType {
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
