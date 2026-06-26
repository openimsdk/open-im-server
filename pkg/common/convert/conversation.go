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

package convert

import (
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/conversation"
)

func ConversationDB2Pb(conversationDB *model.Conversation) *conversation.Conversation {
	if conversationDB == nil {
		return nil
	}
	return &conversation.Conversation{
		OwnerUserID:           conversationDB.OwnerUserID,
		ConversationID:        conversationDB.ConversationID,
		RecvMsgOpt:            conversationDB.RecvMsgOpt,
		ConversationType:      conversationDB.ConversationType,
		UserID:                conversationDB.UserID,
		GroupID:               conversationDB.GroupID,
		IsPinned:              conversationDB.IsPinned,
		AttachedInfo:          conversationDB.AttachedInfo,
		IsPrivateChat:         conversationDB.IsPrivateChat,
		GroupAtType:           conversationDB.GroupAtType,
		Ex:                    conversationDB.Ex,
		BurnDuration:          conversationDB.BurnDuration,
		MinSeq:                conversationDB.MinSeq,
		MaxSeq:                conversationDB.MaxSeq,
		MsgDestructTime:       conversationDB.MsgDestructTime,
		LatestMsgDestructTime: conversationDB.LatestMsgDestructTime.UnixMilli(),
		IsMsgDestruct:         conversationDB.IsMsgDestruct,
	}
}

func ConversationsDB2Pb(conversationsDB []*model.Conversation) (conversationsPB []*conversation.Conversation) {
	for _, conversationDB := range conversationsDB {
		if conversationDB == nil {
			continue
		}
		conversationPB := &conversation.Conversation{
			OwnerUserID:           conversationDB.OwnerUserID,
			ConversationID:        conversationDB.ConversationID,
			RecvMsgOpt:            conversationDB.RecvMsgOpt,
			ConversationType:      conversationDB.ConversationType,
			UserID:                conversationDB.UserID,
			GroupID:               conversationDB.GroupID,
			IsPinned:              conversationDB.IsPinned,
			AttachedInfo:          conversationDB.AttachedInfo,
			IsPrivateChat:         conversationDB.IsPrivateChat,
			GroupAtType:           conversationDB.GroupAtType,
			Ex:                    conversationDB.Ex,
			BurnDuration:          conversationDB.BurnDuration,
			MinSeq:                conversationDB.MinSeq,
			MaxSeq:                conversationDB.MaxSeq,
			MsgDestructTime:       conversationDB.MsgDestructTime,
			LatestMsgDestructTime: conversationDB.LatestMsgDestructTime.UnixMilli(),
			IsMsgDestruct:         conversationDB.IsMsgDestruct,
		}
		conversationsPB = append(conversationsPB, conversationPB)
	}
	return conversationsPB
}

func ConversationPb2DB(conversationPB *conversation.Conversation) *model.Conversation {
	if conversationPB == nil {
		return nil
	}
	conversationDB := &model.Conversation{
		OwnerUserID:      conversationPB.OwnerUserID,
		ConversationID:   conversationPB.ConversationID,
		RecvMsgOpt:       conversationPB.RecvMsgOpt,
		ConversationType: conversationPB.ConversationType,
		UserID:           conversationPB.UserID,
		GroupID:          conversationPB.GroupID,
		IsPinned:         conversationPB.IsPinned,
		AttachedInfo:     conversationPB.AttachedInfo,
		IsPrivateChat:    conversationPB.IsPrivateChat,
		GroupAtType:      conversationPB.GroupAtType,
		Ex:               conversationPB.Ex,
		BurnDuration:     conversationPB.BurnDuration,
		MinSeq:           conversationPB.MinSeq,
		MaxSeq:           conversationPB.MaxSeq,
		MsgDestructTime:  conversationPB.MsgDestructTime,
		IsMsgDestruct:    conversationPB.IsMsgDestruct,
	}
	if conversationPB.LatestMsgDestructTime != 0 {
		conversationDB.LatestMsgDestructTime = time.UnixMilli(conversationPB.LatestMsgDestructTime)
	}
	return conversationDB
}

func ConversationsPb2DB(conversationsPB []*conversation.Conversation) (conversationsDB []*model.Conversation) {
	for _, conversationPB := range conversationsPB {
		if conversationPB == nil {
			continue
		}
		conversationDB := &model.Conversation{
			OwnerUserID:      conversationPB.OwnerUserID,
			ConversationID:   conversationPB.ConversationID,
			RecvMsgOpt:       conversationPB.RecvMsgOpt,
			ConversationType: conversationPB.ConversationType,
			UserID:           conversationPB.UserID,
			GroupID:          conversationPB.GroupID,
			IsPinned:         conversationPB.IsPinned,
			AttachedInfo:     conversationPB.AttachedInfo,
			IsPrivateChat:    conversationPB.IsPrivateChat,
			GroupAtType:      conversationPB.GroupAtType,
			Ex:               conversationPB.Ex,
			BurnDuration:     conversationPB.BurnDuration,
			MinSeq:           conversationPB.MinSeq,
			MaxSeq:           conversationPB.MaxSeq,
			MsgDestructTime:  conversationPB.MsgDestructTime,
			IsMsgDestruct:    conversationPB.IsMsgDestruct,
		}
		if conversationPB.LatestMsgDestructTime != 0 {
			conversationDB.LatestMsgDestructTime = time.UnixMilli(conversationPB.LatestMsgDestructTime)
		}
		conversationsDB = append(conversationsDB, conversationDB)
	}
	return conversationsDB
}
