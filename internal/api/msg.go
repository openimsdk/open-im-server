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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/rpccall"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/idutil"
	"github.com/openimsdk/tools/utils/jsonutil"
	"github.com/openimsdk/tools/utils/timeutil"
)

type MessageApi struct {
	validate      *validator.Validate
	imAdminUserID []string
}

func NewMessageApi(imAdminUserID []string) MessageApi {
	return MessageApi{validate: validator.New(), imAdminUserID: imAdminUserID}
}

func (*MessageApi) SetOptions(options map[string]bool, value bool) {
	datautil.SetSwitchFromOptions(options, constant.IsHistory, value)
	datautil.SetSwitchFromOptions(options, constant.IsPersistent, value)
	datautil.SetSwitchFromOptions(options, constant.IsSenderSync, value)
	datautil.SetSwitchFromOptions(options, constant.IsConversationUpdate, value)
}

func (m *MessageApi) newUserSendMsgReq(_ *gin.Context, params *apistruct.SendMsg) *msg.SendMsgReq {
	var newContent string
	options := make(map[string]bool, 5)
	switch params.ContentType {
	case constant.OANotification:
		notification := sdkws.NotificationElem{}
		notification.Detail = jsonutil.StructToJsonString(params.Content)
		newContent = jsonutil.StructToJsonString(&notification)
	case constant.Text:
		fallthrough
	case constant.Picture:
		fallthrough
	case constant.Custom:
		fallthrough
	case constant.Voice:
		fallthrough
	case constant.Video:
		fallthrough
	case constant.File:
		fallthrough
	default:
		newContent = jsonutil.StructToJsonString(params.Content)
	}
	if params.IsOnlineOnly {
		m.SetOptions(options, false)
	}
	if params.NotOfflinePush {
		datautil.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	}
	pbData := msg.SendMsgReq{
		MsgData: &sdkws.MsgData{
			SendID:           params.SendID,
			GroupID:          params.GroupID,
			ClientMsgID:      idutil.GetMsgIDByMD5(params.SendID),
			SenderPlatformID: params.SenderPlatformID,
			SenderNickname:   params.SenderNickname,
			SenderFaceURL:    params.SenderFaceURL,
			SessionType:      params.SessionType,
			MsgFrom:          constant.SysMsgType,
			ContentType:      params.ContentType,
			Content:          []byte(newContent),
			CreateTime:       timeutil.GetCurrentTimestampByMill(),
			SendTime:         params.SendTime,
			Options:          options,
			OfflinePushInfo:  params.OfflinePushInfo,
			Ex:               params.Ex,
		},
	}
	return &pbData
}

func (m *MessageApi) GetSeq(c *gin.Context) {
	a2r.CallV2(c, msg.GetMaxSeqCaller.Invoke)
}

func (m *MessageApi) PullMsgBySeqs(c *gin.Context) {
	a2r.CallV2(c, msg.PullMessageBySeqsCaller.Invoke)
}

func (m *MessageApi) RevokeMsg(c *gin.Context) {
	a2r.CallV2(c, msg.RevokeMsgCaller.Invoke)
}

func (m *MessageApi) MarkMsgsAsRead(c *gin.Context) {
	a2r.CallV2(c, msg.MarkMsgsAsReadCaller.Invoke)
}

func (m *MessageApi) MarkConversationAsRead(c *gin.Context) {
	a2r.CallV2(c, msg.MarkConversationAsReadCaller.Invoke)
}

func (m *MessageApi) GetConversationsHasReadAndMaxSeq(c *gin.Context) {
	a2r.CallV2(c, msg.GetConversationsHasReadAndMaxSeqCaller.Invoke)
}

func (m *MessageApi) SetConversationHasReadSeq(c *gin.Context) {
	a2r.CallV2(c, msg.SetConversationHasReadSeqCaller.Invoke)
}

func (m *MessageApi) ClearConversationsMsg(c *gin.Context) {
	a2r.CallV2(c, msg.ClearConversationsMsgCaller.Invoke)
}

func (m *MessageApi) UserClearAllMsg(c *gin.Context) {
	a2r.CallV2(c, msg.UserClearAllMsgCaller.Invoke)
}

func (m *MessageApi) DeleteMsgs(c *gin.Context) {
	a2r.CallV2(c, msg.DeleteMsgsCaller.Invoke)
}

func (m *MessageApi) DeleteMsgPhysicalBySeq(c *gin.Context) {
	a2r.CallV2(c, msg.DeleteMsgPhysicalBySeqCaller.Invoke)
}

func (m *MessageApi) DeleteMsgPhysical(c *gin.Context) {
	a2r.CallV2(c, msg.DeleteMsgPhysicalCaller.Invoke)
}

func (m *MessageApi) getSendMsgReq(c *gin.Context, req apistruct.SendMsg) (sendMsgReq *msg.SendMsgReq, err error) {
	var data any
	log.ZDebug(c, "getSendMsgReq", "req", req.Content)
	switch req.ContentType {
	case constant.Text:
		data = apistruct.TextElem{}
	case constant.Picture:
		data = apistruct.PictureElem{}
	case constant.Voice:
		data = apistruct.SoundElem{}
	case constant.Video:
		data = apistruct.VideoElem{}
	case constant.File:
		data = apistruct.FileElem{}
	case constant.AtText:
		data = apistruct.AtElem{}
	case constant.Custom:
		data = apistruct.CustomElem{}
	case constant.Quote:
		data = apistruct.QuoteElem{}
	case constant.Stream:
		data = apistruct.StreamMsgElem{}
	case constant.OANotification:
		data = apistruct.OANotificationElem{}
		req.SessionType = constant.NotificationChatType
		if err = user.GetNotificationAccountCaller.Execute(c, &user.GetNotificationAccountReq{UserID: req.SendID}); err != nil {
			return nil, err
		}
	default:
		return nil, errs.WrapMsg(errs.ErrArgs, "unsupported content type", "contentType", req.ContentType)
	}
	if err := mapstructure.WeakDecode(req.Content, &data); err != nil {
		return nil, errs.WrapMsg(err, "failed to decode message content")
	}
	log.ZDebug(c, "getSendMsgReq", "decodedContent", data)
	if err := m.validate.Struct(data); err != nil {
		return nil, errs.WrapMsg(err, "validation error")
	}
	return m.newUserSendMsgReq(c, &req), nil
}

// SendMessage handles the sending of a message. It's an HTTP handler function to be used with Gin framework.
func (m *MessageApi) SendMessage(c *gin.Context) {
	// Initialize a request struct for sending a message.
	req := apistruct.SendMsgReq{}

	// Bind the JSON request body to the request struct.
	if err := c.BindJSON(&req); err != nil {
		// Respond with an error if request body binding fails.
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}

	// Check if the user has the app manager role.
	if !authverify.IsAppManagerUid(c, m.imAdminUserID) {
		// Respond with a permission error if the user is not an app manager.
		apiresp.GinError(c, errs.ErrNoPermission.WrapMsg("only app manager can send message"))
		return
	}

	// Prepare the message request with additional required data.
	sendMsgReq, err := m.getSendMsgReq(c, req.SendMsg)
	if err != nil {
		// Log and respond with an error if preparation fails.
		apiresp.GinError(c, err)
		return
	}

	// Set the receiver ID in the message data.
	sendMsgReq.MsgData.RecvID = req.RecvID

	// Attempt to send the message using the client.
	respPb, err := msg.SendMsgCaller.Invoke(c, sendMsgReq)
	if err != nil {
		// Set the status to failed and respond with an error if sending fails.
		apiresp.GinError(c, err)
		return
	}

	// Set the status to successful if the message is sent.
	var status = constant.MsgSendSuccessed

	// Attempt to update the message sending status in the system.
	err = msg.SetSendMsgStatusCaller.Execute(c, &msg.SetSendMsgStatusReq{
		Status: int32(status),
	})

	if err != nil {
		// Log the error if updating the status fails.
		apiresp.GinError(c, err)
		return
	}

	// Respond with a success message and the response payload.
	apiresp.GinSuccess(c, respPb)
}

func (m *MessageApi) SendBusinessNotification(c *gin.Context) {
	req := struct {
		Key        string `json:"key"`
		Data       string `json:"data"`
		SendUserID string `json:"sendUserID" binding:"required"`
		RecvUserID string `json:"recvUserID" binding:"required"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}

	if !authverify.IsAppManagerUid(c, m.imAdminUserID) {
		apiresp.GinError(c, errs.ErrNoPermission.WrapMsg("only app manager can send message"))
		return
	}
	sendMsgReq := msg.SendMsgReq{
		MsgData: &sdkws.MsgData{
			SendID: req.SendUserID,
			RecvID: req.RecvUserID,
			Content: []byte(jsonutil.StructToJsonString(&sdkws.NotificationElem{
				Detail: jsonutil.StructToJsonString(&struct {
					Key  string `json:"key"`
					Data string `json:"data"`
				}{Key: req.Key, Data: req.Data}),
			})),
			MsgFrom:     constant.SysMsgType,
			ContentType: constant.BusinessNotification,
			SessionType: constant.SingleChatType,
			CreateTime:  timeutil.GetCurrentTimestampByMill(),
			ClientMsgID: idutil.GetMsgIDByMD5(mcontext.GetOpUserID(c)),
			Options: config.GetOptionsByNotification(config.NotificationConfig{
				IsSendMsg:        false,
				ReliabilityLevel: 1,
				UnreadCount:      false,
			}),
		},
	}
	respPb, err := msg.SendMsgCaller.Invoke(c, &sendMsgReq)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, respPb)
}

func (m *MessageApi) BatchSendMsg(c *gin.Context) {
	var (
		req  apistruct.BatchSendMsgReq
		resp apistruct.BatchSendMsgResp
	)
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	if err := authverify.CheckAdmin(c, m.imAdminUserID); err != nil {
		apiresp.GinError(c, errs.ErrNoPermission.WrapMsg("only app manager can send message"))
		return
	}

	var recvIDs []string
	if req.IsSendAll {
		pageNumber := 1
		showNumber := 500
		for {
			recvIDsPart, err := rpccall.ExtractField(c, user.GetAllUserIDCaller.Invoke, &user.GetAllUserIDReq{Pagination: &sdkws.RequestPagination{
				PageNumber: int32(pageNumber),
				ShowNumber: int32(showNumber),
			}}, (*user.GetAllUserIDResp).GetUserIDs)
			if err != nil {
				apiresp.GinError(c, err)
				return
			}
			recvIDs = append(recvIDs, recvIDsPart...)
			if len(recvIDsPart) < showNumber {
				break
			}
			pageNumber++
		}
	} else {
		recvIDs = req.RecvIDs
	}
	log.ZDebug(c, "BatchSendMsg nums", "nums ", len(recvIDs))
	sendMsgReq, err := m.getSendMsgReq(c, req.SendMsg)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	for _, recvID := range recvIDs {
		sendMsgReq.MsgData.RecvID = recvID
		rpcResp, err := msg.SendMsgCaller.Invoke(c, sendMsgReq)
		if err != nil {
			resp.FailedIDs = append(resp.FailedIDs, recvID)
			continue
		}
		resp.Results = append(resp.Results, &apistruct.SingleReturnResult{
			ServerMsgID: rpcResp.ServerMsgID,
			ClientMsgID: rpcResp.ClientMsgID,
			SendTime:    rpcResp.SendTime,
			RecvID:      recvID,
		})
	}
	apiresp.GinSuccess(c, resp)
}

func (m *MessageApi) CheckMsgIsSendSuccess(c *gin.Context) {
	a2r.CallV2(c, msg.GetSendMsgStatusCaller.Invoke)
}

func (m *MessageApi) GetUsersOnlineStatus(c *gin.Context) {
	a2r.CallV2(c, msg.GetSendMsgStatusCaller.Invoke)
}

func (m *MessageApi) GetActiveUser(c *gin.Context) {
	a2r.CallV2(c, msg.GetActiveUserCaller.Invoke)
}

func (m *MessageApi) GetActiveGroup(c *gin.Context) {
	a2r.CallV2(c, msg.GetActiveGroupCaller.Invoke)
}

func (m *MessageApi) SearchMsg(c *gin.Context) {
	a2r.CallV2(c, msg.SearchMessageCaller.Invoke)
}

func (m *MessageApi) GetServerTime(c *gin.Context) {
	a2r.CallV2(c, msg.GetServerTimeCaller.Invoke)
}

func (m *MessageApi) GetStreamMsg(c *gin.Context) {
	a2r.CallV2(c, msg.GetStreamMsgCaller.Invoke)
}

func (m *MessageApi) AppendStreamMsg(c *gin.Context) {
	a2r.CallV2(c, msg.AppendStreamMsgCaller.Invoke)
}
