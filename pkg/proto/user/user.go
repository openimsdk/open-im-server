package user

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func (x *GetAllUserIDReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	return nil
}

func (x *AccountCheckReq) Check() error {
	if x.CheckUserIDs == nil {
		return errs.ErrArgs.Wrap("CheckUserIDs is empty")
	}
	return nil
}

func (x *GetDesignateUsersReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("UserIDs is empty")
	}
	return nil
}

func (x *UpdateUserInfoReq) Check() error {
	if x.UserInfo == nil {
		return errs.ErrArgs.Wrap("UserInfo is empty")
	}
	if x.UserInfo.UserID == "" {
		return errs.ErrArgs.Wrap("UserID is empty")
	}
	return nil
}

func (x *SetGlobalRecvMessageOptReq) Check() error {
	if x.GlobalRecvMsgOpt > 2 || x.GlobalRecvMsgOpt < 0 {
		return errs.ErrArgs.Wrap("GlobalRecvMsgOpt is invalid")
	}
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("UserID is empty")
	}
	return nil
}

func (x *SetConversationReq) Check() error {
	if err := x.Conversation.Check(); err != nil {
		return err
	}
	if x.NotificationType < 1 || x.NotificationType > 3 {
		return errs.ErrArgs.Wrap("NotificationType is invalid")
	}
	return nil
}

func (x *SetRecvMsgOptReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("OwnerUserID is empty")
	}
	if x.ConversationID == "" {
		return errs.ErrArgs.Wrap("ConversationID is empty")
	}
	if x.RecvMsgOpt < 0 || x.RecvMsgOpt > 2 {
		return errs.ErrArgs.Wrap("RecvMsgOpt is invalid")
	}
	if x.NotificationType < 1 || x.NotificationType > 3 {
		return errs.ErrArgs.Wrap("NotificationType is invalid")
	}
	return nil
}

func (x *GetConversationReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("OwnerUserID is empty")
	}
	if x.ConversationID == "" {
		return errs.ErrArgs.Wrap("ConversationID is empty")
	}
	return nil
}

func (x *GetConversationsReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("OwnerUserID is empty")
	}
	if x.ConversationIDs == nil {
		return errs.ErrArgs.Wrap("ConversationIDs is empty")
	}
	return nil
}

func (x *GetAllConversationsReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("OwnerUserID is empty")
	}
	return nil
}

func (x *BatchSetConversationsReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("OwnerUserID is empty")
	}
	if x.Conversations == nil {
		return errs.ErrArgs.Wrap("ConversationIDs is empty")
	}
	if x.NotificationType < 1 || x.NotificationType > 3 {
		return errs.ErrArgs.Wrap("NotificationType is invalid")
	}
	return nil
}

func (x *GetPaginationUsersReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	return nil
}

func (x *UserRegisterReq) Check() error {
	if x.Secret == "" {
		return errs.ErrArgs.Wrap("Secret is empty")
	}
	if x.Users == nil {
		return errs.ErrArgs.Wrap("Users is empty")
	}
	return nil
}

func (x *GetGlobalRecvMessageOptReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("UserID is empty")
	}
	return nil
}

func (x *UserRegisterCountReq) Check() error {
	if x.Start <= 0 {
		return errs.ErrArgs.Wrap("start is invalid")
	}
	if x.End <= 0 {
		return errs.ErrArgs.Wrap("end is invalid")
	}
	return nil
}
