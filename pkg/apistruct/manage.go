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

// SendMsg defines the structure for sending messages with various metadata.
type SendMsg struct {
	// SendID uniquely identifies the sender.
	SendID string `json:"sendID" binding:"required"`
	
	// GroupID is the identifier for the group, required if SessionType is 2 or 3.
	GroupID string `json:"groupID" binding:"required_if=SessionType 2|required_if=SessionType 3"`
	
	// SenderNickname is the nickname of the sender.
	SenderNickname string `json:"senderNickname"`
	
	// SenderFaceURL is the URL to the sender's avatar.
	SenderFaceURL string `json:"senderFaceURL"`
	
	// SenderPlatformID is an integer identifier for the sender's platform.
	SenderPlatformID int32 `json:"senderPlatformID"`
	
	// Content is the actual content of the message, required and excluded from Swagger documentation.
	Content map[string]interface{} `json:"content" binding:"required" swaggerignore:"true"`
	
	// ContentType is an integer that represents the type of the content.
	ContentType int32 `json:"contentType" binding:"required"`
	
	// SessionType is an integer that represents the type of session for the message.
	SessionType int32 `json:"sessionType" binding:"required"`
	
	// IsOnlineOnly specifies if the message is only sent when the receiver is online.
	IsOnlineOnly bool `json:"isOnlineOnly"`
	
	// NotOfflinePush specifies if the message should not trigger offline push notifications.
	NotOfflinePush bool `json:"notOfflinePush"`
	
	// SendTime is a timestamp indicating when the message was sent.
	SendTime int64 `json:"sendTime"`
	
	// OfflinePushInfo contains information for offline push notifications.
	OfflinePushInfo *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}

// SendMsgReq extends SendMsg with the requirement of RecvID when SessionType indicates a one-on-one or notification chat.
type SendMsgReq struct {
	// RecvID uniquely identifies the receiver and is required for one-on-one or notification chat types.
	RecvID string `json:"recvID" binding:"required_if" message:"recvID is required if sessionType is SingleChatType or NotificationChatType"`
	SendMsg
}

// BatchSendMsgReq defines the structure for sending a message to multiple recipients.
type BatchSendMsgReq struct {
	SendMsg
	
	// IsSendAll indicates whether the message should be sent to all users.
	IsSendAll bool `json:"isSendAll"`
	
	// RecvIDs is a slice of receiver identifiers to whom the message will be sent, required field.
	RecvIDs []string `json:"recvIDs" binding:"required"`
}

// BatchSendMsgResp contains the results of a batch message send operation.
type BatchSendMsgResp struct {
	// Results is a slice of SingleReturnResult, representing the outcome of each message sent.
	Results []*SingleReturnResult `json:"results"`
	
	// FailedIDs is a slice of user IDs for whom the message send failed.
	FailedIDs []string `json:"failedUserIDs"`
}

// SingleReturnResult encapsulates the result of a single message send attempt.
type SingleReturnResult struct {
	// ServerMsgID is the message identifier on the server-side.
	ServerMsgID string `json:"serverMsgID"`
	
	// ClientMsgID is the message identifier on the client-side.
	ClientMsgID string `json:"clientMsgID"`
	
	// SendTime is the timestamp of when the message was sent.
	SendTime int64 `json:"sendTime"`
	
	// RecvID uniquely identifies the receiver of the message.
	RecvID string `json:"recvID"`
}
