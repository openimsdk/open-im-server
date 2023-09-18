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

package apistruct

import (
	sdkws "github.com/OpenIMSDK/protocol/sdkws"
)

type SendMsg struct {
	SendID           string                 `json:"sendID"           binding:"required"`
	GroupID          string                 `json:"groupID"          binding:"required_if=SessionType 2|required_if=SessionType 3"`
	SenderNickname   string                 `json:"senderNickname"`
	SenderFaceURL    string                 `json:"senderFaceURL"`
	SenderPlatformID int32                  `json:"senderPlatformID"`
	Content          map[string]interface{} `json:"content"          binding:"required"                                            swaggerignore:"true"`
	ContentType      int32                  `json:"contentType"      binding:"required"`
	SessionType      int32                  `json:"sessionType"      binding:"required"`
	IsOnlineOnly     bool                   `json:"isOnlineOnly"`
	NotOfflinePush   bool                   `json:"notOfflinePush"`
	SendTime         int64                  `json:"sendTime"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}

type SendMsgReq struct {
	RecvID string `json:"recvID" binding:"required_if" message:"recvID is required if sessionType is SingleChatType or NotificationChatType"`
	SendMsg
}

type BatchSendMsgReq struct {
	SendMsg
	IsSendAll bool     `json:"isSendAll"`
	RecvIDs   []string `json:"recvIDs"`
}

type BatchSendMsgResp struct {
	Results   []*SingleReturnResult `json:"results"`
	FailedIDs []string              `json:"failedUserIDs"`
}

type SingleReturnResult struct {
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
	RecvID      string `json:"recvID"`
}
