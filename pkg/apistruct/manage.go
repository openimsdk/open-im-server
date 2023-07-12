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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

// define a DeleteUsersReq struct
type DeleteUsersReq struct {
	OperationID      string   `json:"operationID"      binding:"required"`
	DeleteUserIDList []string `json:"deleteUserIDList" binding:"required"`
}

// define a DeleteUsersResp struct
type DeleteUsersResp struct {
	FailedUserIDList []string `json:"data"`
}

// define a GetAllUsersUidReq struct
type GetAllUsersUidReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

// define a GetAllUsersUidResp struct
type GetAllUsersUidResp struct {
	UserIDList []string `json:"data"`
}

// define a GetUsersOnlineStatusReq struct
type GetUsersOnlineStatusReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList"  binding:"required,lte=200"`
}

// define a GetUsersOnlineStatusResp struct
type GetUsersOnlineStatusResp struct {

	//SuccessResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult `json:"data"`
}

// define a AccountCheckReq struct
type AccountCheckReq struct {
	OperationID     string   `json:"operationID"     binding:"required"`
	CheckUserIDList []string `json:"checkUserIDList" binding:"required,lte=100"`
}

// define a AccountCheckResp struct
type AccountCheckResp struct {
}

// define a ManagementSendMsg struct
type ManagementSendMsg struct {
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
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}

// define a ManagementSendMsgReq struct
type ManagementSendMsgReq struct {
	SendID           string                 `json:"sendID"           binding:"required"`
	RecvID           string                 `json:"recvID"           binding:"required_if" message:"recvID is required if sessionType is SingleChatType or NotificationChatType"`
	GroupID          string                 `json:"groupID"          binding:"required_if" message:"groupID is required if sessionType is GroupChatType or SuperGroupChatType"`
	SenderNickname   string                 `json:"senderNickname"`
	SenderFaceURL    string                 `json:"senderFaceURL"`
	SenderPlatformID int32                  `json:"senderPlatformID"`
	Content          map[string]interface{} `json:"content"          binding:"required"                                                                                          swaggerignore:"true"`
	ContentType      int32                  `json:"contentType"      binding:"required"`
	SessionType      int32                  `json:"sessionType"      binding:"required"`
	IsOnlineOnly     bool                   `json:"isOnlineOnly"`
	NotOfflinePush   bool                   `json:"notOfflinePush"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}

// define a ManagementSendMsgResp struct
type ManagementSendMsgResp struct {
	ResultList sdkws.UserSendMsgResp `json:"data"`
}

// define a ManagementBatchSendMsgReq struct
type ManagementBatchSendMsgReq struct {
	ManagementSendMsg
	IsSendAll  bool     `json:"isSendAll"`
	RecvIDList []string `json:"recvIDList"`
}

// define a ManagementBatchSendMsgResp struct
type ManagementBatchSendMsgResp struct {
	Data struct {
		ResultList   []*SingleReturnResult `json:"resultList"`
		FailedIDList []string
	} `json:"data"`
}

// define a SingleReturnResult struct
type SingleReturnResult struct {
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
	RecvID      string `json:"recvID"`
}

// define a CheckMsgIsSendSuccessReq struct
type CheckMsgIsSendSuccessReq struct {
	OperationID string `json:"operationID"`
}

// define a CheckMsgIsSendSuccessResp struct
type CheckMsgIsSendSuccessResp struct {
	Status int32 `json:"status"`
}
