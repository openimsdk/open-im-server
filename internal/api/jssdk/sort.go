package jssdk

import "github.com/openimsdk/protocol/msg"

type sortActiveConversations struct {
	Conversation          []*msg.ActiveConversation
	PinnedConversationIDs map[string]struct{}
}

func (s sortActiveConversations) Top(limit int) []*msg.ActiveConversation {
	if limit > 0 && len(s.Conversation) > limit {
		return s.Conversation[:limit]
	}
	return s.Conversation
}

func (s sortActiveConversations) Len() int {
	return len(s.Conversation)
}

func (s sortActiveConversations) Less(i, j int) bool {
	iv, jv := s.Conversation[i], s.Conversation[j]
	_, ip := s.PinnedConversationIDs[iv.ConversationID]
	_, jp := s.PinnedConversationIDs[jv.ConversationID]
	if ip != jp {
		return ip
	}
	return iv.LastTime > jv.LastTime
}

func (s sortActiveConversations) Swap(i, j int) {
	s.Conversation[i], s.Conversation[j] = s.Conversation[j], s.Conversation[i]
}
