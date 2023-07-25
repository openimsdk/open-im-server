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

package push

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/callbackstruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/http"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/config"
	"github.com/OpenIMSDK/tools/constant"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/utils"
)

func url() string {
	return config.Config.Callback.CallbackUrl
}

func callbackOfflinePush(
	ctx context.Context,
	userIDs []string,
	msg *sdkws.MsgData,
	offlinePushUserIDs *[]string,
) error {
	if !config.Config.Callback.CallbackOfflinePush.Enable {
		return nil
	}
	req := &callbackstruct.CallbackBeforePushReq{
		UserStatusBatchCallbackReq: callbackstruct.UserStatusBatchCallbackReq{
			UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackOfflinePushCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      int(msg.SenderPlatformID),
				Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
			},
			UserIDList: userIDs,
		},
		OfflinePushInfo: msg.OfflinePushInfo,
		ClientMsgID:     msg.ClientMsgID,
		SendID:          msg.SendID,
		GroupID:         msg.GroupID,
		ContentType:     msg.ContentType,
		SessionType:     msg.SessionType,
		AtUserIDs:       msg.AtUserIDList,
		Content:         GetContent(msg),
	}
	resp := &callbackstruct.CallbackBeforePushResp{}
	if err := http.CallBackPostReturn(ctx, url(), req, resp, config.Config.Callback.CallbackOfflinePush); err != nil {
		if err == errs.ErrCallbackContinue {
			return nil
		}
		return err
	}
	if len(resp.UserIDs) != 0 {
		*offlinePushUserIDs = resp.UserIDs
	}
	if resp.OfflinePushInfo != nil {
		msg.OfflinePushInfo = resp.OfflinePushInfo
	}
	return nil
}

func callbackOnlinePush(ctx context.Context, userIDs []string, msg *sdkws.MsgData) error {
	if !config.Config.Callback.CallbackOnlinePush.Enable || utils.Contain(msg.SendID, userIDs...) {
		return nil
	}
	req := callbackstruct.CallbackBeforePushReq{
		UserStatusBatchCallbackReq: callbackstruct.UserStatusBatchCallbackReq{
			UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackOnlinePushCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      int(msg.SenderPlatformID),
				Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
			},
			UserIDList: userIDs,
		},
		ClientMsgID: msg.ClientMsgID,
		SendID:      msg.SendID,
		GroupID:     msg.GroupID,
		ContentType: msg.ContentType,
		SessionType: msg.SessionType,
		AtUserIDs:   msg.AtUserIDList,
		Content:     GetContent(msg),
	}
	resp := &callbackstruct.CallbackBeforePushResp{}
	return http.CallBackPostReturn(ctx, url(), req, resp, config.Config.Callback.CallbackOnlinePush)
}

func callbackBeforeSuperGroupOnlinePush(
	ctx context.Context,
	groupID string,
	msg *sdkws.MsgData,
	pushToUserIDs *[]string,
) error {
	if !config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.Enable {
		return nil
	}
	req := callbackstruct.CallbackBeforeSuperGroupOnlinePushReq{
		UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
			CallbackCommand: constant.CallbackSuperGroupOnlinePushCommand,
			OperationID:     mcontext.GetOperationID(ctx),
			PlatformID:      int(msg.SenderPlatformID),
			Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
		},
		ClientMsgID: msg.ClientMsgID,
		SendID:      msg.SendID,
		GroupID:     groupID,
		ContentType: msg.ContentType,
		SessionType: msg.SessionType,
		AtUserIDs:   msg.AtUserIDList,
		Content:     GetContent(msg),
		Seq:         msg.Seq,
	}
	resp := &callbackstruct.CallbackBeforeSuperGroupOnlinePushResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackBeforeSuperGroupOnlinePush); err != nil {
		if err == errs.ErrCallbackContinue {
			return nil
		}
		return err
	}
	if len(resp.UserIDs) != 0 {
		*pushToUserIDs = resp.UserIDs
	}
	return nil
}
