package api

import (
	"github.com/openimsdk/protocol/msg"
	"sort"
	"testing"
)

func TestName(t *testing.T) {
	val := sortActiveConversations{
		Conversation: []*msg.ActiveConversation{
			{
				ConversationID: "100",
				LastTime:       100,
			},
			{
				ConversationID: "200",
				LastTime:       200,
			},
			{
				ConversationID: "300",
				LastTime:       300,
			},
			{
				ConversationID: "400",
				LastTime:       400,
			},
		},
		//PinnedConversationIDs: map[string]struct{}{
		//	"100": {},
		//	"300": {},
		//},
	}
	sort.Sort(&val)
	t.Log(val)

}
