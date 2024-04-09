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

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

func callbackOfflinePush(ctx context.Context, callback *config.Webhooks, userIDs []string, msg *sdkws.MsgData, offlinePushUserIDs *[]string) error {
	if !callback.BeforeOfflinePush.Enable || msg.ContentType == constant.Typing {
		return nil
	}
	req := &callbackstruct.CallbackBeforePushReq{
		UserStatusBatchCallbackReq: callbackstruct.UserStatusBatchCallbackReq{
			UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
				CallbackCommand: callbackstruct.CallbackOfflinePushCommand,
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
	if err := http.CallBackPostReturn(ctx, callback.URL, req, resp, callback.BeforeOfflinePush); err != nil {
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

func callbackOnlinePush(ctx context.Context, callback *config.Webhooks, userIDs []string, msg *sdkws.MsgData) error {
	if !callback.BeforeOnlinePush.Enable || datautil.Contain(msg.SendID, userIDs...) || msg.ContentType == constant.Typing {
		return nil
	}
	req := callbackstruct.CallbackBeforePushReq{
		UserStatusBatchCallbackReq: callbackstruct.UserStatusBatchCallbackReq{
			UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
				CallbackCommand: callbackstruct.CallbackOnlinePushCommand,
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
	if err := http.CallBackPostReturn(ctx, callback.URL, req, resp, callback.BeforeOnlinePush); err != nil {
		return err
	}
	return nil
}

func callbackBeforeSuperGroupOnlinePush(
	ctx context.Context,
	callback *config.Webhooks,
	groupID string,
	msg *sdkws.MsgData,
	pushToUserIDs *[]string,
) error {
	if !callback.BeforeGroupOnlinePush.Enable || msg.ContentType == constant.Typing {
		return nil
	}
	req := callbackstruct.CallbackBeforeSuperGroupOnlinePushReq{
		UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
			CallbackCommand: callbackstruct.CallbackSuperGroupOnlinePushCommand,
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
	if err := http.CallBackPostReturn(ctx, callback.URL, req, resp, callback.BeforeGroupOnlinePush); err != nil {
		return err
	}

	if len(resp.UserIDs) != 0 {
		*pushToUserIDs = resp.UserIDs
	}
	return nil
}
