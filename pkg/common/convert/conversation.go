package convert

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func ConversationDB2Pb(conversationDB *relation.ConversationModel) *conversation.Conversation {
	conversationPB := &conversation.Conversation{}
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
