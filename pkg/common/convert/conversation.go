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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/protocol/conversation"
	"github.com/OpenIMSDK/tools/utils"
)

func ConversationDB2Pb(conversationDB *relation.ConversationModel) *conversation.Conversation {
	conversationPB := &conversation.Conversation{}
	conversationPB.LatestMsgDestructTime = conversationDB.LatestMsgDestructTime.Unix()
	if err := utils.CopyStructFields(conversationPB, conversationDB); err != nil {
		return nil
	}
	return conversationPB
}

func ConversationsDB2Pb(conversationsDB []*relation.ConversationModel) (conversationsPB []*conversation.Conversation) {
	for _, conversationDB := range conversationsDB {
		conversationPB := &conversation.Conversation{}
		if err := utils.CopyStructFields(conversationPB, conversationDB); err != nil {
			continue
		}
		conversationPB.LatestMsgDestructTime = conversationDB.LatestMsgDestructTime.Unix()
		conversationsPB = append(conversationsPB, conversationPB)
	}
	return conversationsPB
}

func ConversationPb2DB(conversationPB *conversation.Conversation) *relation.ConversationModel {
	conversationDB := &relation.ConversationModel{}
	if err := utils.CopyStructFields(conversationDB, conversationPB); err != nil {
		return nil
	}
	return conversationDB
}

func ConversationsPb2DB(conversationsPB []*conversation.Conversation) (conversationsDB []*relation.ConversationModel) {
	for _, conversationPB := range conversationsPB {
		conversationDB := &relation.ConversationModel{}
		if err := utils.CopyStructFields(conversationDB, conversationPB); err != nil {
			continue
		}
		conversationsDB = append(conversationsDB, conversationDB)
	}
	return conversationsDB
}
