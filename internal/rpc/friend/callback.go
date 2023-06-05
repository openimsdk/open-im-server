package friend

import (
	"context"
	cbapi "github.com/OpenIMSDK/Open-IM-Server/pkg/callbackstruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/http"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	pbfriend "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
)

func CallbackBeforeAddFriend(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) error {
	log.ZInfo(ctx, "CallbackBeforeAddFriend", "in", req)
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
	defer log.ZInfo(ctx, "CallbackBeforeAddFriend", "out", &resp)
	return http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriend)
}
