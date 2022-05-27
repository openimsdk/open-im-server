package logic

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	http2 "net/http"
)

func callbackOfflinePush(operationID, userID string, info *commonPb.OfflinePushInfo, platformID int32) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackUserOnline.Enable {
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
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), callbackOfflinePushResp)
	return callbackResp
}
