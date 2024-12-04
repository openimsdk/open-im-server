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
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func ConversationDB2Pb(conversationDB *model.Conversation) *conversation.Conversation {
	conversationPB := &conversation.Conversation{}
	conversationPB.LatestMsgDestructTime = conversationDB.LatestMsgDestructTime.UnixMilli()
	if err := datautil.CopyStructFields(conversationPB, conversationDB); err != nil {
		return nil
	}
	return conversationPB
}

func ConversationsDB2Pb(conversationsDB []*model.Conversation) (conversationsPB []*conversation.Conversation) {
	for _, conversationDB := range conversationsDB {
		conversationPB := &conversation.Conversation{}
		if err := datautil.CopyStructFields(conversationPB, conversationDB); err != nil {
			continue
		}
		conversationPB.LatestMsgDestructTime = conversationDB.LatestMsgDestructTime.UnixMilli()
		conversationsPB = append(conversationsPB, conversationPB)
	}
	return conversationsPB
}

func ConversationPb2DB(conversationPB *conversation.Conversation) *model.Conversation {
	conversationDB := &model.Conversation{}
	if err := datautil.CopyStructFields(conversationDB, conversationPB); err != nil {
		return nil
	}
	return conversationDB
}

func ConversationsPb2DB(conversationsPB []*conversation.Conversation) (conversationsDB []*model.Conversation) {
	for _, conversationPB := range conversationsPB {
		conversationDB := &model.Conversation{}
		if err := datautil.CopyStructFields(conversationDB, conversationPB); err != nil {
			continue
		}
		conversationsDB = append(conversationsDB, conversationDB)
	}
	return conversationsDB
}
