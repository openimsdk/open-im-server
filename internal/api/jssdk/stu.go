package jssdk

import (
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/sdkws"
)

type ActiveConversationsReq struct {
	Count int `json:"count"`
}

type ConversationMsg struct {
	Conversation *conversation.Conversation `json:"conversation"`
	LastMsg      *sdkws.MsgData             `json:"lastMsg"`
	User         *sdkws.UserInfo
	Friend       *sdkws.FriendInfo
	Group        *sdkws.GroupInfo
	MaxSeq       int64 `json:"maxSeq"`
	ReadSeq      int64 `json:"readSeq"`
}

type ConversationsResp struct {
	UnreadCount   int64             `json:"unreadCount"`
	Conversations []ConversationMsg `json:"conversations"`
}
