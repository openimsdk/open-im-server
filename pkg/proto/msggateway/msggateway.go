package msggateway

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func (x *OnlinePushMsgReq) Check() error {
	if x.MsgData == nil {
		return errs.ErrArgs.Wrap("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	if x.PushToUserID == "" {
		return errs.ErrArgs.Wrap("PushToUserID is empty")
	}
	return nil
}

func (x *OnlineBatchPushOneMsgReq) Check() error {
	if x.MsgData == nil {
		return errs.ErrArgs.Wrap("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	if x.PushToUserIDs == nil {
		return errs.ErrArgs.Wrap("PushToUserIDs is empty")
	}
	return nil
}

func (x *GetUsersOnlineStatusReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("UserIDs is empty")
	}
	return nil
}

func (x *KickUserOfflineReq) Check() error {
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errs.ErrArgs.Wrap("PlatformID is invalid")
	}
	if x.KickUserIDList == nil {
		return errs.ErrArgs.Wrap("KickUserIDList is empty")
	}
	return nil
}

func (x *MultiTerminalLoginCheckReq) Check() error {
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errs.ErrArgs.Wrap("PlatformID is invalid")
	}
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("UserID is empty")
	}
	if x.Token == "" {
		return errs.ErrArgs.Wrap("Token is empty")
	}
	return nil
}
