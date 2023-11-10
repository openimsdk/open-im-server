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
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type MessageApi struct {
	*rpcclient.Message
	validate      *validator.Validate
	userRpcClient *rpcclient.UserRpcClient
}

func NewMessageApi(msgRpcClient *rpcclient.Message, userRpcClient *rpcclient.User) MessageApi {
	return MessageApi{Message: msgRpcClient, validate: validator.New(), userRpcClient: rpcclient.NewUserRpcClientByUser(userRpcClient)}
}

func (MessageApi) SetOptions(options map[string]bool, value bool) {
	utils.SetSwitchFromOptions(options, constant.IsHistory, value)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, value)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, value)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, value)
}

func (m MessageApi) newUserSendMsgReq(_ *gin.Context, params *apistruct.SendMsg) *msg.SendMsgReq {
	var newContent string
	options := make(map[string]bool, 5)
	switch params.ContentType {
	case constant.OANotification:
		notification := sdkws.NotificationElem{}
		notification.Detail = utils.StructToJsonString(params.Content)
		newContent = utils.StructToJsonString(&notification)
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
		newContent = utils.StructToJsonString(params.Content)
	}
	if params.IsOnlineOnly {
		m.SetOptions(options, false)
	}
	if params.NotOfflinePush {
		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	}
	pbData := msg.SendMsgReq{
		MsgData: &sdkws.MsgData{
			SendID:           params.SendID,
			GroupID:          params.GroupID,
			ClientMsgID:      utils.GetMsgID(params.SendID),
			SenderPlatformID: params.SenderPlatformID,
			SenderNickname:   params.SenderNickname,
			SenderFaceURL:    params.SenderFaceURL,
			SessionType:      params.SessionType,
			MsgFrom:          constant.SysMsgType,
			ContentType:      params.ContentType,
			Content:          []byte(newContent),
			CreateTime:       utils.GetCurrentTimestampByMill(),
			SendTime:         params.SendTime,
			Options:          options,
			OfflinePushInfo:  params.OfflinePushInfo,
		},
	}
	return &pbData
}

func (m *MessageApi) GetSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetMaxSeq, m.Client, c)
}

func (m *MessageApi) PullMsgBySeqs(c *gin.Context) {
	a2r.Call(msg.MsgClient.PullMessageBySeqs, m.Client, c)
}

func (m *MessageApi) RevokeMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.RevokeMsg, m.Client, c)
}

func (m *MessageApi) MarkMsgsAsRead(c *gin.Context) {
	a2r.Call(msg.MsgClient.MarkMsgsAsRead, m.Client, c)
}

func (m *MessageApi) MarkConversationAsRead(c *gin.Context) {
	a2r.Call(msg.MsgClient.MarkConversationAsRead, m.Client, c)
}

func (m *MessageApi) GetConversationsHasReadAndMaxSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetConversationsHasReadAndMaxSeq, m.Client, c)
}

func (m *MessageApi) SetConversationHasReadSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.SetConversationHasReadSeq, m.Client, c)
}

func (m *MessageApi) ClearConversationsMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.ClearConversationsMsg, m.Client, c)
}

func (m *MessageApi) UserClearAllMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.UserClearAllMsg, m.Client, c)
}

func (m *MessageApi) DeleteMsgs(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMsgs, m.Client, c)
}

func (m *MessageApi) DeleteMsgPhysicalBySeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMsgPhysicalBySeq, m.Client, c)
}

func (m *MessageApi) DeleteMsgPhysical(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMsgPhysical, m.Client, c)
}

func (m *MessageApi) getSendMsgReq(c *gin.Context, req apistruct.SendMsg) (sendMsgReq *msg.SendMsgReq, err error) {
	var data interface{}
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
	case constant.Custom:
		data = apistruct.CustomElem{}
	case constant.OANotification:
		data = apistruct.OANotificationElem{}
		req.SessionType = constant.NotificationChatType
		if !authverify.IsManagerUserID(req.SendID) {
			return nil, errs.ErrNoPermission.
				Wrap("only app manager can as sender send OANotificationElem")
		}
	default:
		return nil, errs.ErrArgs.WithDetail("not support err contentType")
	}
	if err := mapstructure.WeakDecode(req.Content, &data); err != nil {
		return nil, err
	}
	log.ZDebug(c, "getSendMsgReq", "req", req.Content)
	if err := m.validate.Struct(data); err != nil {
		return nil, err
	}
	return m.newUserSendMsgReq(c, &req), nil
}

func (m *MessageApi) SendMessage(c *gin.Context) {
	req := apistruct.SendMsgReq{}
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	if !authverify.IsAppManagerUid(c) {
		apiresp.GinError(c, errs.ErrNoPermission.Wrap("only app manager can send message"))
		return
	}
	sendMsgReq, err := m.getSendMsgReq(c, req.SendMsg)
	if err != nil {
		log.ZError(c, "decodeData failed", err)
		apiresp.GinError(c, err)
		return
	}
	sendMsgReq.MsgData.RecvID = req.RecvID
	var status int
	respPb, err := m.Client.SendMsg(c, sendMsgReq)
	if err != nil {
		status = constant.MsgSendFailed
		log.ZError(c, "send message err", err)
		apiresp.GinError(c, err)
		return
	}
	status = constant.MsgSendSuccessed
	_, err = m.Client.SetSendMsgStatus(c, &msg.SetSendMsgStatusReq{
		Status: int32(status),
	})
	if err != nil {
		log.ZError(c, "SetSendMsgStatus failed", err)
	}
	apiresp.GinSuccess(c, respPb)
}

func (m *MessageApi) SendBusinessNotification(c *gin.Context) {
	req := struct {
		Key        string `json:"key"`
		Data       string `json:"data"`
		SendUserID string `json:"sendUserID"`
		RecvUserID string `json:"recvUserID"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	if !authverify.IsAppManagerUid(c) {
		apiresp.GinError(c, errs.ErrNoPermission.Wrap("only app manager can send message"))
		return
	}
	sendMsgReq := msg.SendMsgReq{
		MsgData: &sdkws.MsgData{
			SendID: req.SendUserID,
			RecvID: req.RecvUserID,
			Content: []byte(utils.StructToJsonString(&sdkws.NotificationElem{
				Detail: utils.StructToJsonString(&struct {
					Key  string `json:"key"`
					Data string `json:"data"`
				}{Key: req.Key, Data: req.Data}),
			})),
			MsgFrom:     constant.SysMsgType,
			ContentType: constant.BusinessNotification,
			SessionType: constant.SingleChatType,
			CreateTime:  utils.GetCurrentTimestampByMill(),
			ClientMsgID: utils.GetMsgID(mcontext.GetOpUserID(c)),
			Options: config.GetOptionsByNotification(config.NotificationConf{
				IsSendMsg:        false,
				ReliabilityLevel: 1,
				UnreadCount:      false,
			}),
		},
	}
	respPb, err := m.Client.SendMsg(c, &sendMsgReq)
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
		log.ZError(c, "BatchSendMsg BindJSON failed", err)
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	log.ZInfo(c, "BatchSendMsg", "req", req)
	if err := authverify.CheckAdmin(c); err != nil {
		apiresp.GinError(c, errs.ErrNoPermission.Wrap("only app manager can send message"))
		return
	}

	var recvIDs []string
	var err error
	if req.IsSendAll {
		pageNumber := 1
		showNumber := 500
		for {
			recvIDsPart, err := m.userRpcClient.GetAllUserIDs(c, int32(pageNumber), int32(showNumber))
			if err != nil {
				log.ZError(c, "GetAllUserIDs failed", err)
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
		log.ZError(c, "decodeData failed", err)
		apiresp.GinError(c, err)
		return
	}
	for _, recvID := range recvIDs {
		sendMsgReq.MsgData.RecvID = recvID
		rpcResp, err := m.Client.SendMsg(c, sendMsgReq)
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
	a2r.Call(msg.MsgClient.GetSendMsgStatus, m.Client, c)
}

func (m *MessageApi) GetUsersOnlineStatus(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetSendMsgStatus, m.Client, c)
}

func (m *MessageApi) GetActiveUser(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetActiveUser, m.Client, c)
}

func (m *MessageApi) GetActiveGroup(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetActiveGroup, m.Client, c)
}

func (m *MessageApi) SearchMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.SearchMessage, m.Client, c)
}

func (m *MessageApi) GetServerTime(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetServerTime, m.Client, c)
}
