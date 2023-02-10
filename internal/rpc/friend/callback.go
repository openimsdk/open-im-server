package friend

import (
	cbapi "Open_IM/pkg/callbackstruct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	pbfriend "Open_IM/pkg/proto/friend"
	"context"

	//"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	http2 "net/http"
)

func callbackBeforeAddFriendV1(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) error {
	resp := callbackBeforeAddFriend(ctx, req)
	if resp.ErrCode != 0 {
		return (&constant.ErrInfo{
			ErrCode: resp.ErrCode,
			ErrMsg:  resp.ErrMsg,
		}).Wrap()
	}
	return nil
}

func callbackBeforeAddFriend(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) cbapi.CommonCallbackResp {
	callbackResp := cbapi.CommonCallbackResp{OperationID: tracelog.GetOperationID(ctx)}
	if !config.Config.Callback.CallbackBeforeAddFriend.Enable {
		return callbackResp
	}

	commonCallbackReq := &cbapi.CallbackBeforeAddFriendReq{
		CallbackCommand: constant.CallbackBeforeAddFriendCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		ReqMsg:          req.ReqMsg,
		OperationID:     tracelog.GetOperationID(ctx),
	}
	resp := &cbapi.CallbackBeforeAddFriendResp{
		CommonCallbackResp: &callbackResp,
	}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), commonCallbackReq, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeAddFriendCommand, commonCallbackReq, resp, config.Config.Callback.CallbackBeforeAddFriend); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !*config.Config.Callback.CallbackBeforeAddFriend.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	return callbackResp
}
