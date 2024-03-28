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

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
	"github.com/openimsdk/protocol/constant"
	pbchat "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
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

func callbackBeforeSendSingleMsg(ctx context.Context, globalConfig *config.GlobalConfig, msg *pbchat.SendMsgReq) error {
	if !globalConfig.Callback.CallbackBeforeSendSingleMsg.Enable || msg.MsgData.ContentType == constant.Typing {
		return nil
	}
	req := &cbapi.CallbackBeforeSendSingleMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackBeforeSendSingleMsgCommand),
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbapi.CallbackBeforeSendSingleMsgResp{}
	if err := http.CallBackPostReturn(ctx, globalConfig.Callback.CallbackUrl, req, resp, globalConfig.Callback.CallbackBeforeSendSingleMsg); err != nil {
		return err
	}
	return nil
}

func callbackAfterSendSingleMsg(ctx context.Context, globalConfig *config.GlobalConfig, msg *pbchat.SendMsgReq) error {
	if !globalConfig.Callback.CallbackAfterSendSingleMsg.Enable || msg.MsgData.ContentType == constant.Typing {
		return nil
	}
	req := &cbapi.CallbackAfterSendSingleMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackAfterSendSingleMsgCommand),
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbapi.CallbackAfterSendSingleMsgResp{}
	if err := http.CallBackPostReturn(ctx, globalConfig.Callback.CallbackUrl, req, resp, globalConfig.Callback.CallbackAfterSendSingleMsg); err != nil {
		return err
	}
	return nil
}

func callbackBeforeSendGroupMsg(ctx context.Context, globalConfig *config.GlobalConfig, msg *pbchat.SendMsgReq) error {
	if !globalConfig.Callback.CallbackBeforeSendGroupMsg.Enable || msg.MsgData.ContentType == constant.Typing {
		return nil
	}
	req := &cbapi.CallbackBeforeSendGroupMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackBeforeSendGroupMsgCommand),
		GroupID:           msg.MsgData.GroupID,
	}
	resp := &cbapi.CallbackBeforeSendGroupMsgResp{}
	if err := http.CallBackPostReturn(ctx, globalConfig.Callback.CallbackUrl, req, resp, globalConfig.Callback.CallbackBeforeSendGroupMsg); err != nil {
		return err
	}
	return nil
}

func callbackAfterSendGroupMsg(ctx context.Context, globalConfig *config.GlobalConfig, msg *pbchat.SendMsgReq) error {
	if !globalConfig.Callback.CallbackAfterSendGroupMsg.Enable || msg.MsgData.ContentType == constant.Typing {
		return nil
	}
	req := &cbapi.CallbackAfterSendGroupMsgReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackAfterSendGroupMsgCommand),
		GroupID:           msg.MsgData.GroupID,
	}
	resp := &cbapi.CallbackAfterSendGroupMsgResp{}
	if err := http.CallBackPostReturn(ctx, globalConfig.Callback.CallbackUrl, req, resp, globalConfig.Callback.CallbackAfterSendGroupMsg); err != nil {
		return err
	}
	return nil
}

func callbackMsgModify(ctx context.Context, globalConfig *config.GlobalConfig, msg *pbchat.SendMsgReq) error {
	if !globalConfig.Callback.CallbackMsgModify.Enable || msg.MsgData.ContentType != constant.Text {
		return nil
	}
	req := &cbapi.CallbackMsgModifyCommandReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackMsgModifyCommand),
	}
	resp := &cbapi.CallbackMsgModifyCommandResp{}
	if err := http.CallBackPostReturn(ctx, globalConfig.Callback.CallbackUrl, req, resp, globalConfig.Callback.CallbackMsgModify); err != nil {
		return err
	}
	if resp.Content != nil {
		msg.MsgData.Content = []byte(*resp.Content)
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
	log.ZDebug(ctx, "callbackMsgModify", "msg", msg.MsgData)
	return nil
}

func CallbackGroupMsgRead(ctx context.Context, globalConfig *config.GlobalConfig, req *cbapi.CallbackGroupMsgReadReq) error {
	if !globalConfig.Callback.CallbackGroupMsgRead.Enable {
		return nil
	}
	req.CallbackCommand = cbapi.CallbackGroupMsgReadCommand

	resp := &cbapi.CallbackGroupMsgReadResp{}
	if err := http.CallBackPostReturn(ctx, globalConfig.Callback.CallbackUrl, req, resp, globalConfig.Callback.CallbackGroupMsgRead); err != nil {
		return err
	}
	return nil
}

func CallbackSingleMsgRead(ctx context.Context, globalConfig *config.GlobalConfig, req *cbapi.CallbackSingleMsgReadReq) error {
	if !globalConfig.Callback.CallbackSingleMsgRead.Enable {
		return nil
	}
	req.CallbackCommand = cbapi.CallbackSingleMsgRead

	resp := &cbapi.CallbackSingleMsgReadResp{}

	if err := http.CallBackPostReturn(ctx, globalConfig.Callback.CallbackUrl, req, resp, globalConfig.Callback.CallbackSingleMsgRead); err != nil {
		return err
	}
	return nil
}
func CallbackAfterRevokeMsg(ctx context.Context, globalConfig *config.GlobalConfig, req *pbchat.RevokeMsgReq) error {
	if !globalConfig.Callback.CallbackAfterRevokeMsg.Enable {
		return nil
	}
	callbackReq := &cbapi.CallbackAfterRevokeMsgReq{
		CallbackCommand: cbapi.CallbackAfterRevokeMsgCommand,
		ConversationID:  req.ConversationID,
		Seq:             req.Seq,
		UserID:          req.UserID,
	}
	resp := &cbapi.CallbackAfterRevokeMsgResp{}
	if err := http.CallBackPostReturn(ctx, globalConfig.Callback.CallbackUrl, callbackReq, resp, globalConfig.Callback.CallbackAfterRevokeMsg); err != nil {
		return err
	}
	return nil
}
