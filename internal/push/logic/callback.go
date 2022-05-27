package logic

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	http2 "net/http"
)

func callbackOfflinePush(operationID, userID string, info *commonPb.OfflinePushInfo, platformID int32) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackOfflinePush.Enable {
		return callbackResp
	}
	callbackOfflinePushReq := cbApi.CallbackOfflinePushReq{
		UserStatusCallbackReq: cbApi.UserStatusCallbackReq{
			CallbackCommand: constant.CallbackOfflinePushCommand,
			OperationID:     operationID,
			UserID:          userID,
			PlatformID:      platformID,
			Platform:        constant.PlatformIDToName(platformID),
		},
		OfflinePushInfo: info,
	}
	callbackOfflinePushResp := &cbApi.CallbackOfflinePushResp{CommonCallbackResp: &callbackResp}
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, callbackOfflinePushReq, callbackOfflinePushResp, config.Config.Callback.CallbackOfflinePush.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackOfflinePush.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	return callbackResp
}
