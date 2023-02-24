package friend

import (
	cbapi "OpenIM/pkg/callbackstruct"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/http"
	"OpenIM/pkg/common/tracelog"
	pbfriend "OpenIM/pkg/proto/friend"
	"context"
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
		OperationID:     tracelog.GetOperationID(ctx),
	}
	resp := &cbapi.CallbackBeforeAddFriendResp{}
	return http.CallBackPostReturn(config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriend)
}
