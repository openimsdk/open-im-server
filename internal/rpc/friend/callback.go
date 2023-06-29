package friend

import (
	"context"

	cbapi "github.com/OpenIMSDK/Open-IM-Server/pkg/callbackstruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/http"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbfriend "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
)

func CallbackBeforeAddFriend(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) error {
	if !config.Config.Callback.CallbackBeforeAddFriend.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeAddFriendReq{
		CallbackCommand: constant.CallbackBeforeAddFriendCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		ReqMsg:          req.ReqMsg,
		OperationID:     mcontext.GetOperationID(ctx),
	}
	resp := &cbapi.CallbackBeforeAddFriendResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriend); err != nil {
		if err == errs.ErrCallbackContinue {
			return nil
		}
		return err
	}
	return nil
}
