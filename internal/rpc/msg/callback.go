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

	"github.com/OpenIMSDK/protocol/sdkws"
	"google.golang.org/protobuf/proto"

	"github.com/OpenIMSDK/protocol/constant"
	pbchat "github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/utils"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
)

func cbURL() string {
	return config.Config.Callback.CallbackUrl
}

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

func callbackBeforeSendSingleMsg(ctx context.Context, msg *pbchat.SendMsgReq) error {
	if !config.Config.Callback.CallbackBeforeSendSingleMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackBeforeSendSingleMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackBeforeSendSingleMsgCommand),
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbapi.CallbackBeforeSendSingleMsgResp{}
	if err := http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackBeforeSendSingleMsg); err != nil {
		return err
	}
	return nil
}

func callbackAfterSendSingleMsg(ctx context.Context, msg *pbchat.SendMsgReq) error {
	if !config.Config.Callback.CallbackAfterSendSingleMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackAfterSendSingleMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackAfterSendSingleMsgCommand),
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbapi.CallbackAfterSendSingleMsgResp{}
	if err := http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackAfterSendSingleMsg); err != nil {
		return err
	}
	return nil
}

func callbackBeforeSendGroupMsg(ctx context.Context, msg *pbchat.SendMsgReq) error {
	if !config.Config.Callback.CallbackBeforeSendSingleMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackBeforeSendGroupMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackBeforeSendGroupMsgCommand),
		GroupID:           msg.MsgData.GroupID,
	}
	resp := &cbapi.CallbackBeforeSendGroupMsgResp{}
	if err := http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackBeforeSendGroupMsg); err != nil {
		return err
	}
	return nil
}

func callbackAfterSendGroupMsg(ctx context.Context, msg *pbchat.SendMsgReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackAfterSendGroupMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackAfterSendGroupMsgCommand),
		GroupID:           msg.MsgData.GroupID,
	}
	resp := &cbapi.CallbackAfterSendGroupMsgResp{}
	if err := http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg); err != nil {
		return err
	}
	return nil
}

func callbackMsgModify(ctx context.Context, msg *pbchat.SendMsgReq) error {
	if !config.Config.Callback.CallbackMsgModify.Enable || msg.MsgData.ContentType != constant.Text {
		return nil
	}
	req := &cbapi.CallbackMsgModifyCommandReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackMsgModifyCommand),
	}
	resp := &cbapi.CallbackMsgModifyCommandResp{}
	if err := http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackMsgModify); err != nil {
		return err
	}
	if resp.Content != nil {
		msg.MsgData.Content = []byte(*resp.Content)
	}
	utils.NotNilReplace(msg.MsgData.OfflinePushInfo, resp.OfflinePushInfo)
	utils.NotNilReplace(&msg.MsgData.RecvID, resp.RecvID)
	utils.NotNilReplace(&msg.MsgData.GroupID, resp.GroupID)
	utils.NotNilReplace(&msg.MsgData.ClientMsgID, resp.ClientMsgID)
	utils.NotNilReplace(&msg.MsgData.ServerMsgID, resp.ServerMsgID)
	utils.NotNilReplace(&msg.MsgData.SenderPlatformID, resp.SenderPlatformID)
	utils.NotNilReplace(&msg.MsgData.SenderNickname, resp.SenderNickname)
	utils.NotNilReplace(&msg.MsgData.SenderFaceURL, resp.SenderFaceURL)
	utils.NotNilReplace(&msg.MsgData.SessionType, resp.SessionType)
	utils.NotNilReplace(&msg.MsgData.MsgFrom, resp.MsgFrom)
	utils.NotNilReplace(&msg.MsgData.ContentType, resp.ContentType)
	utils.NotNilReplace(&msg.MsgData.Status, resp.Status)
	utils.NotNilReplace(&msg.MsgData.Options, resp.Options)
	utils.NotNilReplace(&msg.MsgData.AtUserIDList, resp.AtUserIDList)
	utils.NotNilReplace(&msg.MsgData.AttachedInfo, resp.AttachedInfo)
	utils.NotNilReplace(&msg.MsgData.Ex, resp.Ex)
	log.ZDebug(ctx, "callbackMsgModify", "msg", msg.MsgData)
	return nil
}
func CallbackGroupMsgRead(ctx context.Context, req *cbapi.CallbackGroupMsgReadReq) error {
	if !config.Config.Callback.CallbackGroupMsgRead.Enable || req.ContentType != constant.Text {
		return nil
	}
	req.CallbackCommand = cbapi.CallbackGroupMsgReadCommand

	resp := &cbapi.CallbackGroupMsgReadResp{}
	if err := http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackMsgModify); err != nil {
		return err
	}
	return nil
}

func CallbackSingleMsgRead(ctx context.Context, req *cbapi.CallbackSingleMsgReadReq) error {
	if !config.Config.Callback.CallbackSingleMsgRead.Enable || req.ContentType != constant.Text {
		return nil
	}
	req.CallbackCommand = cbapi.CallbackSingleMsgRead

	resp := &cbapi.CallbackSingleMsgReadResp{}

	if err := http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackMsgModify); err != nil {
		return err
	}
	return nil
}
func CallbackAfterRevokeMsg(ctx context.Context, req *pbchat.RevokeMsgReq) error {
	if !config.Config.Callback.CallbackAfterRevokeMsg.Enable {
		return nil
	}
	callbackReq := &cbapi.CallbackAfterRevokeMsgReq{
		CallbackCommand: cbapi.CallbackAfterRevokeMsgCommand,
		ConversationID:  req.ConversationID,
		Seq:             req.Seq,
		UserID:          req.UserID,
	}
	resp := &cbapi.CallbackAfterRevokeMsgResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, callbackReq, resp, config.Config.Callback.CallbackAfterRevokeMsg); err != nil {
		return err
	}
	utils.StructFieldNotNilReplace(req, resp)
	return nil
}
