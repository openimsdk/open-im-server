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

type OptResult struct {
	ConversationID string `json:"conversationID"`
	Result         *int32 `json:"result"`
}
type GetAllConversationMessageOptReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID"  binding:"required"`
}
type GetAllConversationMessageOptResp struct {
	ConversationOptResultList []*OptResult `json:"data"`
}
type GetReceiveMessageOptReq struct {
	ConversationIDList []string `json:"conversationIDList" binding:"required"`
	OperationID        string   `json:"operationID"        binding:"required"`
	FromUserID         string   `json:"fromUserID"         binding:"required"`
}
type GetReceiveMessageOptResp struct {
	ConversationOptResultList []*OptResult `json:"data"`
}
type SetReceiveMessageOptReq struct {
	FromUserID         string   `json:"fromUserID"         binding:"required"`
	OperationID        string   `json:"operationID"        binding:"required"`
	Opt                *int32   `json:"opt"                binding:"required"`
	ConversationIDList []string `json:"conversationIDList" binding:"required"`
}
type SetReceiveMessageOptResp struct {
	ConversationOptResultList []*OptResult `json:"data"`
}

type Conversation struct {
	OwnerUserID           string `json:"ownerUserID"           binding:"required"`
	ConversationID        string `json:"conversationID"        binding:"required"`
	ConversationType      int32  `json:"conversationType"      binding:"required"`
	UserID                string `json:"userID"`
	GroupID               string `json:"groupID"`
	RecvMsgOpt            int32  `json:"recvMsgOpt"            binding:"omitempty,oneof=0 1 2"`
	UnreadCount           int32  `json:"unreadCount"           binding:"omitempty"`
	DraftTextTime         int64  `json:"draftTextTime"`
	IsPinned              bool   `json:"isPinned"              binding:"omitempty"`
	IsPrivateChat         bool   `json:"isPrivateChat"`
	BurnDuration          int32  `json:"burnDuration"`
	GroupAtType           int32  `json:"groupAtType"`
	IsNotInGroup          bool   `json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `json:"updateUnreadCountTime"`
	AttachedInfo          string `json:"attachedInfo"`
	Ex                    string `json:"ex"`
}

type SetConversationReq struct {
	Conversation
	NotificationType int32  `json:"notificationType"`
	OperationID      string `json:"operationID"      binding:"required"`
}

type SetConversationResp struct {
}
type ModifyConversationFieldReq struct {
	Conversation
	FieldType   int32    `json:"fieldType"   binding:"required"`
	UserIDList  []string `json:"userIDList"  binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}
type ModifyConversationFieldResp struct {
}

type BatchSetConversationsReq struct {
	Conversations    []Conversation `json:"conversations"    binding:"required"`
	NotificationType int32          `json:"notificationType"`
	OwnerUserID      string         `json:"ownerUserID"      binding:"required"`
	OperationID      string         `json:"operationID"      binding:"required"`
}

type BatchSetConversationsResp struct {
	Data struct {
		Success []string `json:"success"`
		Failed  []string `json:"failed"`
	} `json:"data"`
}

type GetConversationReq struct {
	ConversationID string `json:"conversationID" binding:"required"`
	OwnerUserID    string `json:"ownerUserID"    binding:"required"`
	OperationID    string `json:"operationID"    binding:"required"`
}

type GetConversationResp struct {
	Conversation Conversation `json:"data"`
}

type GetAllConversationsReq struct {
	OwnerUserID string `json:"ownerUserID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type GetAllConversationsResp struct {
	Conversations []Conversation `json:"data"`
}

type GetConversationsReq struct {
	ConversationIDs []string `json:"conversationIDs" binding:"required"`
	OwnerUserID     string   `json:"ownerUserID"     binding:"required"`
	OperationID     string   `json:"operationID"     binding:"required"`
}

type GetConversationsResp struct {
	Conversations []Conversation `json:"data"`
}

type SetRecvMsgOptReq struct {
	OwnerUserID      string `json:"ownerUserID"      binding:"required"`
	ConversationID   string `json:"conversationID"`
	RecvMsgOpt       int32  `json:"recvMsgOpt"       binding:"omitempty,oneof=0 1 2"`
	OperationID      string `json:"operationID"      binding:"required"`
	NotificationType int32  `json:"notificationType"`
}

type SetRecvMsgOptResp struct {
}
