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
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	http2 "github.com/openimsdk/open-im-server/v3/pkg/common/http"
	pbmsg "github.com/openimsdk/open-im-server/v3/tools/data-conversion/openim/proto/msg"
	"net/http"
	"reflect"
	"strings"

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
	case constant.OANotification:
		data = apistruct.OANotificationElem{}
		req.SessionType = constant.NotificationChatType
		if err = m.userRpcClient.GetNotificationByID(c, req.SendID); err != nil {
			return nil, err
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
	if !authverify.IsAppManagerUid(c) {
		// Respond with a permission error if the user is not an app manager.
		apiresp.GinError(c, errs.ErrNoPermission.Wrap("only app manager can send message"))
		return
	}

	// Prepare the message request with additional required data.
	sendMsgReq, err := m.getSendMsgReq(c, req.SendMsg)
	if err != nil {
		// Log and respond with an error if preparation fails.
		log.ZError(c, "decodeData failed", err)
		apiresp.GinError(c, err)
		return
	}

	// Set the receiver ID in the message data.
	sendMsgReq.MsgData.RecvID = req.RecvID

	// Declare a variable to store the message sending status.
	var status int

	// Attempt to send the message using the client.
	respPb, err := m.Client.SendMsg(c, sendMsgReq)
	if err != nil {
		// Set the status to failed and respond with an error if sending fails.
		status = constant.MsgSendFailed
		log.ZError(c, "send message err", err)
		apiresp.GinError(c, err)
		return
	}

	// Set the status to successful if the message is sent.
	status = constant.MsgSendSuccessed

	// Attempt to update the message sending status in the system.
	_, err = m.Client.SetSendMsgStatus(c, &msg.SetSendMsgStatusReq{
		Status: int32(status),
	})
	if err != nil {
		// Log the error if updating the status fails.
		log.ZError(c, "SetSendMsgStatus failed", err)
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

func (m *MessageApi) CallbackExample(c *gin.Context) {
	// 1. Callback after sending a single chat message
	var req callbackstruct.CallbackAfterSendSingleMsgReq

	if err := c.BindJSON(&req); err != nil {
		log.ZError(c, "CallbackExample BindJSON failed", err)
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}

	resp := callbackstruct.CallbackAfterSendSingleMsgResp{
		CommonCallbackResp: callbackstruct.CommonCallbackResp{
			ActionCode: 0,
			ErrCode:    200,
			ErrMsg:     "success",
			ErrDlt:     "successful",
			NextCode:   0,
		},
	}
	c.JSON(http.StatusOK, resp)

	// 2. If the user receiving the message is a customer service bot, return the message.

	// UserID of the robot account
	robotics := "5078764102"
	// Administrator token
	imtoken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJpbUFkbWluIiwiUGxhdGZvcm1JRCI6MTAsImV4cCI6MTcxMzI1MjI0OSwibmJmIjoxNzA1NDc1OTQ5LCJpYXQiOjE3MDU0NzYyNDl9.Zi-uFre8zq6msT3mFOumgcfNKBJ92kTw9ewsKeRVbZ4"
	if req.SendID == robotics {
		return
	}
	// Processing text messages
	if req.ContentType == constant.Picture {
		user, err := m.userRpcClient.GetUserInfo(c, robotics)
		if err != nil {
			log.ZError(c, "CallbackExample get Sender failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}

		content := make(map[string]any, 1)

		// Handle message structures
		text := apistruct.PictureElem{}
		log.ZDebug(c, "callback", "contextCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC", req.Content)
		err = json.Unmarshal([]byte(req.Content), &text)
		if err != nil {
			log.ZError(c, "CallbackExample unmarshal failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}
		log.ZDebug(c, "callback", "text%TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT", text)
		content["content"] = text.SourcePath
		content["sourcePicture"] = text.SourcePicture
		content["bigPicture"] = text.BigPicture
		content["snapshotPicture"] = text.SnapshotPicture

		log.ZDebug(c, "callback", "contextAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", content)

		mapStruct := make(map[string]any)
		mapStruct1, err := convertStructToMap(text.SnapshotPicture)

		if err != nil {
			log.ZError(c, "CallbackExample struct to map failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}
		mapStruct["snapshotPicture"] = mapStruct1

		mapStruct2, err := convertStructToMap(text.BigPicture)
		if err != nil {
			log.ZError(c, "CallbackExample struct to map failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}
		mapStruct["bigPicture"] = mapStruct2

		mapStruct3, err := convertStructToMap(text.SourcePicture)
		if err != nil {
			log.ZError(c, "CallbackExample struct to map failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}
		mapStruct["sourcePicture"] = mapStruct3
		mapStruct["sourcePath"] = text.SourcePath

		input := &apistruct.SendMsgReq{
			RecvID: req.SendID,
			SendMsg: apistruct.SendMsg{
				SendID:           user.UserID,
				SenderNickname:   user.Nickname,
				SenderFaceURL:    user.FaceURL,
				SenderPlatformID: req.SenderPlatformID,
				Content:          mapStruct,
				ContentType:      req.ContentType,
				SessionType:      req.SessionType,
				SendTime:         utils.GetCurrentTimestampByMill(), // millisecond
			},
		}

		url := "http://127.0.0.1:10002/msg/send_msg"
		header := make(map[string]string, 2)
		header["token"] = imtoken
		type sendResp struct {
			ErrCode int               `json:"errCode"`
			ErrMsg  string            `json:"errMsg"`
			ErrDlt  string            `json:"errDlt"`
			Data    pbmsg.SendMsgResp `json:"data,omitempty"`
		}

		output := &sendResp{}

		// Initiate a post request that calls the interface that sends the message (the bot sends a message to user)
		b, err := http2.Post(c, url, header, input, 10)
		if err != nil {
			log.ZError(c, "CallbackExample send message failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}
		if err = json.Unmarshal(b, output); err != nil {
			log.ZError(c, "CallbackExample unmarshal failed", err)
			apiresp.GinError(c, errs.ErrInternalServer.WithDetail(err.Error()).Wrap())
			return
		}
		res := &msg.SendMsgResp{
			ServerMsgID: output.Data.ServerMsgID,
			ClientMsgID: output.Data.ClientMsgID,
			SendTime:    output.Data.SendTime,
		}

		apiresp.GinSuccess(c, res)
	}
}

func convertStructToMap(input interface{}) (map[string]interface{}, error) {
	// 使用反射创建一个空的 map
	result := make(map[string]interface{})

	// 获取结构体的类型信息
	inputType := reflect.TypeOf(input)

	// 获取结构体的值信息
	inputValue := reflect.ValueOf(input)

	// 确保输入是结构体类型
	if inputType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Input is not a struct")
	}

	// 遍历结构体的字段
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		fieldValue := inputValue.Field(i)

		// 获取mapstructure标签的值作为map的键
		mapKey := field.Tag.Get("mapstructure")
		fmt.Println(mapKey)
		// 如果没有mapstructure标签，则使用字段名作为键
		if mapKey == "" {
			mapKey = field.Name
		}

		// 转换为小写形式，以匹配JSON的命名约定
		mapKey = strings.ToLower(mapKey)

		// 将字段名作为 map 的键，字段值作为 map 的值
		result[mapKey] = fieldValue.Interface()
	}

	return result, nil
}
