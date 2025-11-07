package convert

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/utils/datautil"
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
