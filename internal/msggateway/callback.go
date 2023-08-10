// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msggateway

import (
	"context"
	"time"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/mcontext"

	cbapi "github.com/OpenIMSDK/Open-IM-Server/pkg/callbackstruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/http"
)

func url() string {
	return config.Config.Callback.CallbackUrl
}

func CallbackUserOnline(ctx context.Context, userID string, platformID int, isAppBackground bool, connID string) error {
	if !config.Config.Callback.CallbackUserOnline.Enable {
		return nil
	}
	req := cbapi.CallbackUserOnlineReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserOnlineCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:             time.Now().UnixMilli(),
		IsAppBackground: isAppBackground,
		ConnID:          connID,
	}
	resp := cbapi.CommonCallbackResp{}
	return http.CallBackPostReturn(ctx, url(), &req, &resp, config.Config.Callback.CallbackUserOnline)
}

func CallbackUserOffline(ctx context.Context, userID string, platformID int, connID string) error {
	if !config.Config.Callback.CallbackUserOffline.Enable {
		return nil
	}
	req := &cbapi.CallbackUserOfflineReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserOfflineCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:    time.Now().UnixMilli(),
		ConnID: connID,
	}
	resp := &cbapi.CallbackUserOfflineResp{}
	return http.CallBackPostReturn(ctx, url(), req, resp, config.Config.Callback.CallbackUserOffline)
}

func CallbackUserKickOff(ctx context.Context, userID string, platformID int) error {
	if !config.Config.Callback.CallbackUserKickOff.Enable {
		return nil
	}
	req := &cbapi.CallbackUserKickOffReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserKickOffCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq: time.Now().UnixMilli(),
	}
	resp := &cbapi.CommonCallbackResp{}
	return http.CallBackPostReturn(ctx, url(), req, resp, config.Config.Callback.CallbackUserOffline)
}

// func callbackUserOnline(operationID, userID string, platformID int, token string, isAppBackground bool, connID
// string) cbApi.CommonCallbackResp {
//	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
//	if !config.Config.Callback.CallbackUserOnline.Enable {
//		return callbackResp
//	}
//	callbackUserOnlineReq := cbApi.CallbackUserOnlineReq{
//		Token: token,
//		UserStatusCallbackReq: cbApi.UserStatusCallbackReq{
//			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
//				CallbackCommand: constant.CallbackUserOnlineCommand,
//				OperationID:     operationID,
//				PlatformID:      int32(platformID),
//				Platform:        constant.PlatformIDToName(platformID),
//			},
//			UserID: userID,
//		},
//		Seq:             int(time.Now().UnixNano() / 1e6),
//		IsAppBackground: isAppBackground,
//		ConnID:          connID,
//	}
//	callbackUserOnlineResp := &cbApi.CallbackUserOnlineResp{CommonCallbackResp: &callbackResp}
// 	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, constant.CallbackUserOnlineCommand,
// callbackUserOnlineReq, callbackUserOnlineResp, config.Config.Callback.CallbackUserOnline.CallbackTimeOut); err != nil
// {
//		callbackResp.ErrCode = http2.StatusInternalServerError
//		callbackResp.ErrMsg = err.Error()
//	}
//	return callbackResp
//}
//func callbackUserOffline(operationID, userID string, platformID int, connID string) cbApi.CommonCallbackResp {
//	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
//	if !config.Config.Callback.CallbackUserOffline.Enable {
//		return callbackResp
//	}
//	callbackOfflineReq := cbApi.CallbackUserOfflineReq{
//		UserStatusCallbackReq: cbApi.UserStatusCallbackReq{
//			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
//				CallbackCommand: constant.CallbackUserOfflineCommand,
//				OperationID:     operationID,
//				PlatformID:      int32(platformID),
//				Platform:        constant.PlatformIDToName(platformID),
//			},
//			UserID: userID,
//		},
//		Seq:    int(time.Now().UnixNano() / 1e6),
//		ConnID: connID,
//	}
//	callbackUserOfflineResp := &cbApi.CallbackUserOfflineResp{CommonCallbackResp: &callbackResp}
// 	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, constant.CallbackUserOfflineCommand,
// callbackOfflineReq, callbackUserOfflineResp, config.Config.Callback.CallbackUserOffline.CallbackTimeOut); err != nil
// {
//		callbackResp.ErrCode = http2.StatusInternalServerError
//		callbackResp.ErrMsg = err.Error()
//	}
//	return callbackResp
//}
//func callbackUserKickOff(operationID string, userID string, platformID int) cbApi.CommonCallbackResp {
//	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
//	if !config.Config.Callback.CallbackUserKickOff.Enable {
//		return callbackResp
//	}
//	callbackUserKickOffReq := cbApi.CallbackUserKickOffReq{
//		UserStatusCallbackReq: cbApi.UserStatusCallbackReq{
//			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
//				CallbackCommand: constant.CallbackUserKickOffCommand,
//				OperationID:     operationID,
//				PlatformID:      int32(platformID),
//				Platform:        constant.PlatformIDToName(platformID),
//			},
//			UserID: userID,
//		},
//		Seq: int(time.Now().UnixNano() / 1e6),
//	}
//	callbackUserKickOffResp := &cbApi.CallbackUserKickOffResp{CommonCallbackResp: &callbackResp}
// 	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, constant.CallbackUserKickOffCommand,
// callbackUserKickOffReq, callbackUserKickOffResp, config.Config.Callback.CallbackUserOffline.CallbackTimeOut); err !=
// nil {
//		callbackResp.ErrCode = http2.StatusInternalServerError
//		callbackResp.ErrMsg = err.Error()
//	}
//	return callbackResp
//}
