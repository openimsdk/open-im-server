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

package msg

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/tools/errs"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/constant"
	pbchat "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/stringutil"
	"google.golang.org/protobuf/proto"
)

func toCommonCallback(ctx context.Context, msg *pbchat.SendMsgReq, command string) cbapi.CommonCallbackReq {
	return cbapi.CommonCallbackReq{
		SendID:           msg.MsgData.SendID,
		ServerMsgID:      msg.MsgData.ServerMsgID,
		CallbackCommand:  command,
		ClientMsgID:      msg.MsgData.ClientMsgID,
		OperationID:      mcontext.GetOperationID(ctx),
		SenderPlatformID: msg.MsgData.SenderPlatformID,
		SenderNickname:   msg.MsgData.SenderNickname,
		SessionType:      msg.MsgData.SessionType,
		MsgFrom:          msg.MsgData.MsgFrom,
		ContentType:      msg.MsgData.ContentType,
		Status:           msg.MsgData.Status,
		SendTime:         msg.MsgData.SendTime,
		CreateTime:       msg.MsgData.CreateTime,
		AtUserIDList:     msg.MsgData.AtUserIDList,
		SenderFaceURL:    msg.MsgData.SenderFaceURL,
		Content:          GetContent(msg.MsgData),
		Seq:              uint32(msg.MsgData.Seq),
		Ex:               msg.MsgData.Ex,
	}
}

func GetContent(msg *sdkws.MsgData) string {
	if msg.ContentType >= constant.NotificationBegin && msg.ContentType <= constant.NotificationEnd {
		var tips sdkws.TipsComm
		_ = proto.Unmarshal(msg.Content, &tips)
		content := tips.JsonDetail
		return content
	} else {
		return string(msg.Content)
	}
}

func (m *msgServer) webhookBeforeSendSingleMsg(ctx context.Context, before *config.BeforeConfig, msg *pbchat.SendMsgReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		if msg.MsgData.ContentType == constant.Typing {
			return nil
		}
		if !filterBeforeMsg(msg, before) {
			return nil
		}
		cbReq := &cbapi.CallbackBeforeSendSingleMsgReq{
			CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackBeforeSendSingleMsgCommand),
			RecvID:            msg.MsgData.RecvID,
		}
		resp := &cbapi.CallbackBeforeSendSingleMsgResp{}
		if err := m.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		return nil
	})
}

// Move to msgtransfer
func (m *msgServer) webhookAfterSendSingleMsg(ctx context.Context, after *config.AfterConfig, msg *pbchat.SendMsgReq) {
	if msg.MsgData.ContentType == constant.Typing {
		return
	}
	if !filterAfterMsg(msg, after) {
		return
	}
	cbReq := &cbapi.CallbackAfterSendSingleMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackAfterSendSingleMsgCommand),
		RecvID:            msg.MsgData.RecvID,
	}
	m.webhookClient.AsyncPostWithQuery(ctx, cbReq.GetCallbackCommand(), cbReq, &cbapi.CallbackAfterSendSingleMsgResp{}, after, buildKeyMsgDataQuery(msg.MsgData))
}

func (m *msgServer) webhookBeforeSendGroupMsg(ctx context.Context, before *config.BeforeConfig, msg *pbchat.SendMsgReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		if !filterBeforeMsg(msg, before) {
			return nil
		}
		if msg.MsgData.ContentType == constant.Typing {
			return nil
		}
		cbReq := &cbapi.CallbackBeforeSendGroupMsgReq{
			CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackBeforeSendGroupMsgCommand),
			GroupID:           msg.MsgData.GroupID,
		}
		resp := &cbapi.CallbackBeforeSendGroupMsgResp{}
		if err := m.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}
		return nil
	})
}

func (m *msgServer) webhookAfterSendGroupMsg(ctx context.Context, after *config.AfterConfig, msg *pbchat.SendMsgReq) {
	if msg.MsgData.ContentType == constant.Typing {
		return
	}
	if !filterAfterMsg(msg, after) {
		return
	}
	cbReq := &cbapi.CallbackAfterSendGroupMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackAfterSendGroupMsgCommand),
		GroupID:           msg.MsgData.GroupID,
	}

	m.webhookClient.AsyncPostWithQuery(ctx, cbReq.GetCallbackCommand(), cbReq, &cbapi.CallbackAfterSendGroupMsgResp{}, after, buildKeyMsgDataQuery(msg.MsgData))
}

func (m *msgServer) webhookBeforeMsgModify(ctx context.Context, before *config.BeforeConfig, msg *pbchat.SendMsgReq, beforeMsgData **sdkws.MsgData) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		//if msg.MsgData.ContentType != constant.Text {
		//	return nil
		//}
		if !filterBeforeMsg(msg, before) {
			return nil
		}
		cbReq := &cbapi.CallbackMsgModifyCommandReq{
			CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackBeforeMsgModifyCommand),
		}
		resp := &cbapi.CallbackMsgModifyCommandResp{}
		if err := m.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}
		if beforeMsgData != nil {
			*beforeMsgData = proto.Clone(msg.MsgData).(*sdkws.MsgData)
		}
		if resp.Content != nil {
			msg.MsgData.Content = []byte(*resp.Content)
			if err := json.Unmarshal(msg.MsgData.Content, &struct{}{}); err != nil {
				return errs.ErrArgs.WrapMsg("webhook msg modify content is not json", "content", string(msg.MsgData.Content))
			}
		}
		datautil.NotNilReplace(msg.MsgData.OfflinePushInfo, resp.OfflinePushInfo)
		datautil.NotNilReplace(&msg.MsgData.RecvID, resp.RecvID)
		datautil.NotNilReplace(&msg.MsgData.GroupID, resp.GroupID)
		datautil.NotNilReplace(&msg.MsgData.ClientMsgID, resp.ClientMsgID)
		datautil.NotNilReplace(&msg.MsgData.ServerMsgID, resp.ServerMsgID)
		datautil.NotNilReplace(&msg.MsgData.SenderPlatformID, resp.SenderPlatformID)
		datautil.NotNilReplace(&msg.MsgData.SenderNickname, resp.SenderNickname)
		datautil.NotNilReplace(&msg.MsgData.SenderFaceURL, resp.SenderFaceURL)
		datautil.NotNilReplace(&msg.MsgData.SessionType, resp.SessionType)
		datautil.NotNilReplace(&msg.MsgData.MsgFrom, resp.MsgFrom)
		datautil.NotNilReplace(&msg.MsgData.ContentType, resp.ContentType)
		datautil.NotNilReplace(&msg.MsgData.Status, resp.Status)
		datautil.NotNilReplace(&msg.MsgData.Options, resp.Options)
		datautil.NotNilReplace(&msg.MsgData.AtUserIDList, resp.AtUserIDList)
		datautil.NotNilReplace(&msg.MsgData.AttachedInfo, resp.AttachedInfo)
		datautil.NotNilReplace(&msg.MsgData.Ex, resp.Ex)
		return nil
	})
}

func (m *msgServer) webhookAfterGroupMsgRead(ctx context.Context, after *config.AfterConfig, req *cbapi.CallbackGroupMsgReadReq) {
	req.CallbackCommand = cbapi.CallbackAfterGroupMsgReadCommand
	m.webhookClient.AsyncPost(ctx, req.GetCallbackCommand(), req, &cbapi.CallbackGroupMsgReadResp{}, after)
}

func (m *msgServer) webhookAfterSingleMsgRead(ctx context.Context, after *config.AfterConfig, req *cbapi.CallbackSingleMsgReadReq) {

	req.CallbackCommand = cbapi.CallbackAfterSingleMsgReadCommand

	m.webhookClient.AsyncPost(ctx, req.GetCallbackCommand(), req, &cbapi.CallbackSingleMsgReadResp{}, after)

}

func (m *msgServer) webhookAfterRevokeMsg(ctx context.Context, after *config.AfterConfig, req *pbchat.RevokeMsgReq) {
	callbackReq := &cbapi.CallbackAfterRevokeMsgReq{
		CallbackCommand: cbapi.CallbackAfterRevokeMsgCommand,
		ConversationID:  req.ConversationID,
		Seq:             req.Seq,
		UserID:          req.UserID,
	}
	m.webhookClient.AsyncPost(ctx, callbackReq.GetCallbackCommand(), callbackReq, &cbapi.CallbackAfterRevokeMsgResp{}, after)
}

func buildKeyMsgDataQuery(msg *sdkws.MsgData) map[string]string {
	keyMsgData := apistruct.KeyMsgData{
		SendID:  msg.SendID,
		RecvID:  msg.RecvID,
		GroupID: msg.GroupID,
	}

	return map[string]string{
		webhook.Key: base64.StdEncoding.EncodeToString(stringutil.StructToJsonBytes(keyMsgData)),
	}
}
