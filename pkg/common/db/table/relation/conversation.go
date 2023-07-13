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

package relation

import (
	"context"
	"time"
)

const (
	conversationModelTableName = "conversations"
)

type ConversationModel struct {
	OwnerUserID           string    `gorm:"column:owner_user_id;primary_key;type:char(128)" json:"OwnerUserID"`
	ConversationID        string    `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType      int32     `gorm:"column:conversation_type" json:"conversationType"`
	UserID                string    `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID               string    `gorm:"column:group_id;type:char(128)" json:"groupID"`
	RecvMsgOpt            int32     `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	IsPinned              bool      `gorm:"column:is_pinned" json:"isPinned"`
	IsPrivateChat         bool      `gorm:"column:is_private_chat" json:"isPrivateChat"`
	BurnDuration          int32     `gorm:"column:burn_duration;default:30" json:"burnDuration"`
	GroupAtType           int32     `gorm:"column:group_at_type" json:"groupAtType"`
	AttachedInfo          string    `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                    string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	MaxSeq                int64     `gorm:"column:max_seq" json:"maxSeq"`
	MinSeq                int64     `gorm:"column:min_seq" json:"minSeq"`
	CreateTime            time.Time `gorm:"column:create_time;index:create_time;autoCreateTime"`
	IsMsgDestruct         bool      `gorm:"column:is_msg_destruct;default:false"`
	MsgDestructTime       int64     `gorm:"column:msg_destruct_time;default:604800"`
	LatestMsgDestructTime time.Time `gorm:"column:latest_msg_destruct_time;autoCreateTime"`
}

func (ConversationModel) TableName() string {
	return conversationModelTableName
}

type ConversationModelInterface interface {
	Create(ctx context.Context, conversations []*ConversationModel) (err error)
	Delete(ctx context.Context, groupIDs []string) (err error)
	UpdateByMap(ctx context.Context, userIDs []string, conversationID string, args map[string]interface{}) (rows int64, err error)
	Update(ctx context.Context, conversation *ConversationModel) (err error)
	Find(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []*ConversationModel, err error)
	FindUserID(ctx context.Context, userIDs []string, conversationIDs []string) ([]string, error)
	FindUserIDAllConversationID(ctx context.Context, userID string) ([]string, error)
	Take(ctx context.Context, userID, conversationID string) (conversation *ConversationModel, err error)
	FindConversationID(ctx context.Context, userID string, conversationIDs []string) (existConversationID []string, err error)
	FindUserIDAllConversations(ctx context.Context, userID string) (conversations []*ConversationModel, err error)
	FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
	GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error)
	FindSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
	GetAllConversationIDs(ctx context.Context) ([]string, error)
	GetUserAllHasReadSeqs(ctx context.Context, ownerUserID string) (hashReadSeqs map[string]int64, err error)
	GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*ConversationModel, error)
	GetConversationIDsNeedDestruct(ctx context.Context) ([]*ConversationModel, error)
	NewTx(tx any) ConversationModelInterface
}
