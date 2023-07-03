package push

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func (x *PushMsgReq) Check() error {
	if x.MsgData == nil {
		return errs.ErrArgs.Wrap("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	if x.ConversationID == "" {
		return errs.ErrArgs.Wrap("ConversationID is empty")
	}
	return nil
}

func (x *DelUserPushTokenReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("UserID is empty")
	}
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errs.ErrArgs.Wrap("PlatformID is invalid")
	}
	return nil
}
