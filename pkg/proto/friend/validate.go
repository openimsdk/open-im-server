package friend

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func (m *ApplyToAddFriendReq) Check() error {
	if m.GetToUserID() == "" {
		return errs.ErrArgs.Wrap("get toUserID is empty")
	}
	return nil
}
