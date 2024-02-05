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

package callbackstruct

import (
	sdkws "github.com/OpenIMSDK/protocol/sdkws"
)

type CallbackBeforeSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID string `json:"recvID"`
}

type CallbackBeforeSendSingleMsgResp struct {
	CommonCallbackResp
}

type CallbackAfterSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID string `json:"recvID"`
}

type CallbackAfterSendSingleMsgResp struct {
	CommonCallbackResp
}

type CallbackBeforeSendGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackBeforeSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackAfterSendGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackAfterSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackMsgModifyCommandReq struct {
	CommonCallbackReq
}

type CallbackMsgModifyCommandResp struct {
	CommonCallbackResp
	Content          *string                `json:"content"`
	RecvID           *string                `json:"recvID"`
	GroupID          *string                `json:"groupID"`
	ClientMsgID      *string                `json:"clientMsgID"`
	ServerMsgID      *string                `json:"serverMsgID"`
	SenderPlatformID *int32                 `json:"senderPlatformID"`
	SenderNickname   *string                `json:"senderNickname"`
	SenderFaceURL    *string                `json:"senderFaceURL"`
	SessionType      *int32                 `json:"sessionType"`
	MsgFrom          *int32                 `json:"msgFrom"`
	ContentType      *int32                 `json:"contentType"`
	Status           *int32                 `json:"status"`
	Options          *map[string]bool       `json:"options"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
	AtUserIDList     *[]string              `json:"atUserIDList"`
	MsgDataList      *[]byte                `json:"msgDataList"`
	AttachedInfo     *string                `json:"attachedInfo"`
	Ex               *string                `json:"ex"`
}

type CallbackGroupMsgReadReq struct {
	CallbackCommand `json:"callbackCommand"`
	SendID          string `json:"sendID"`
	ReceiveID       string `json:"receiveID"`
	UnreadMsgNum    int64  `json:"unreadMsgNum"`
	ContentType     int64  `json:"contentType"`
}

type CallbackGroupMsgReadResp struct {
	CommonCallbackResp
}

type CallbackSingleMsgReadReq struct {
	CallbackCommand `json:"callbackCommand"`
	ConversationID  string  `json:"conversationID"`
	UserID          string  `json:"userID"`
	Seqs            []int64 `json:"Seqs"`
	ContentType     int32   `json:"contentType"`
}

type CallbackSingleMsgReadResp struct {
	CommonCallbackResp
}
