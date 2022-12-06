package gate

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	http2 "net/http"
	"time"
)

func callbackUserOnline(operationID, userID string, platformID int, token string, isAppBackground bool, connID string) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackUserOnline.Enable {
		return callbackResp
	}
	callbackUserOnlineReq := cbApi.CallbackUserOnlineReq{
		Token: token,
		UserStatusCallbackReq: cbApi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserOnlineCommand,
				OperationID:     operationID,
				PlatformID:      int32(platformID),
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:             int(time.Now().UnixNano() / 1e6),
		IsAppBackground: isAppBackground,
		ConnID:          connID,
	}
	callbackUserOnlineResp := &cbApi.CallbackUserOnlineResp{CommonCallbackResp: &callbackResp}
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackUserOnlineCommand, callbackUserOnlineReq, callbackUserOnlineResp, config.Config.Callback.CallbackUserOnline.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return callbackResp
}

func callbackUserOffline(operationID, userID string, platformID int, connID string) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackUserOffline.Enable {
		return callbackResp
	}
	callbackOfflineReq := cbApi.CallbackUserOfflineReq{
		UserStatusCallbackReq: cbApi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserOfflineCommand,
				OperationID:     operationID,
				PlatformID:      int32(platformID),
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:    int(time.Now().UnixNano() / 1e6),
		ConnID: connID,
	}
	callbackUserOfflineResp := &cbApi.CallbackUserOfflineResp{CommonCallbackResp: &callbackResp}
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackUserOfflineCommand, callbackOfflineReq, callbackUserOfflineResp, config.Config.Callback.CallbackUserOffline.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return callbackResp
}

func callbackUserKickOff(operationID string, userID string, platformID int) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackUserKickOff.Enable {
		return callbackResp
	}
	callbackUserKickOffReq := cbApi.CallbackUserKickOffReq{
		UserStatusCallbackReq: cbApi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserKickOffCommand,
				OperationID:     operationID,
				PlatformID:      int32(platformID),
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq: int(time.Now().UnixNano() / 1e6),
	}
	callbackUserKickOffResp := &cbApi.CallbackUserKickOffResp{CommonCallbackResp: &callbackResp}
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackUserKickOffCommand, callbackUserKickOffReq, callbackUserKickOffResp, config.Config.Callback.CallbackUserOffline.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return callbackResp
}
