package callbackstruct

import (
	"time"
)

type Conversation struct {
	OwnerUserID           string    `json:"owner_user_id"`
	ConversationID        string    `json:"conversation_id"`
	ConversationType      int32     `json:"conversation_type"`
	UserID                string    `json:"user_id"`
	GroupID               string    `json:"group_id"`
	RecvMsgOpt            int32     `json:"recv_msg_opt"`
	IsPinned              bool      `json:"is_pinned"`
	IsPrivateChat         bool      `json:"is_private_chat"`
	BurnDuration          int32     `json:"burn_duration"`
	GroupAtType           int32     `json:"group_at_type"`
	AttachedInfo          string    `json:"attached_info"`
	Ex                    string    `json:"ex"`
	MaxSeq                int64     `json:"max_seq"`
	MinSeq                int64     `json:"min_seq"`
	CreateTime            time.Time `json:"create_time"`
	IsMsgDestruct         bool      `json:"is_msg_destruct"`
	MsgDestructTime       int64     `json:"msg_destruct_time"`
	LatestMsgDestructTime time.Time `json:"latest_msg_destruct_time"`
}

type CallbackBeforeCreateSingleChatConversationsReq struct {
	CallbackCommand       `json:"callbackCommand"`
	OwnerUserID           string    `json:"owner_user_id"`
	ConversationID        string    `json:"conversation_id"`
	ConversationType      int32     `json:"conversation_type"`
	UserID                string    `json:"user_id"`
	GroupID               string    `json:"group_id"`
	RecvMsgOpt            int32     `json:"recv_msg_opt"`
	IsPinned              bool      `json:"is_pinned"`
	IsPrivateChat         bool      `json:"is_private_chat"`
	BurnDuration          int32     `json:"burn_duration"`
	GroupAtType           int32     `json:"group_at_type"`
	AttachedInfo          string    `json:"attached_info"`
	Ex                    string    `json:"ex"`
	MaxSeq                int64     `json:"max_seq"`
	MinSeq                int64     `json:"min_seq"`
	CreateTime            time.Time `json:"create_time"`
	IsMsgDestruct         bool      `json:"is_msg_destruct"`
	MsgDestructTime       int64     `json:"msg_destruct_time"`
	LatestMsgDestructTime time.Time `json:"latest_msg_destruct_time"`
}

type CallbackBeforeCreateSingleChatConversationsResp struct {
	CommonCallbackResp
	OwnerUserID           *string    `json:"owner_user_id"`
	ConversationID        *string    `json:"conversation_id"`
	ConversationType      *int32     `json:"conversation_type"`
	UserID                *string    `json:"user_id"`
	GroupID               *string    `json:"group_id"`
	RecvMsgOpt            *int32     `json:"recv_msg_opt"`
	IsPinned              *bool      `json:"is_pinned"`
	IsPrivateChat         *bool      `json:"is_private_chat"`
	BurnDuration          *int32     `json:"burn_duration"`
	GroupAtType           *int32     `json:"group_at_type"`
	AttachedInfo          *string    `json:"attached_info"`
	Ex                    *string    `json:"ex"`
	MaxSeq                *int64     `json:"max_seq"`
	MinSeq                *int64     `json:"min_seq"`
	CreateTime            *time.Time `json:"create_time"`
	IsMsgDestruct         *bool      `json:"is_msg_destruct"`
	MsgDestructTime       *int64     `json:"msg_destruct_time"`
	LatestMsgDestructTime *time.Time `json:"latest_msg_destruct_time"`
}

type CallbackAfterCreateSingleChatConversationsReq struct {
	CallbackCommand       `json:"callbackCommand"`
	OwnerUserID           string    `json:"owner_user_id"`
	ConversationID        string    `json:"conversation_id"`
	ConversationType      int32     `json:"conversation_type"`
	UserID                string    `json:"user_id"`
	GroupID               string    `json:"group_id"`
	RecvMsgOpt            int32     `json:"recv_msg_opt"`
	IsPinned              bool      `json:"is_pinned"`
	IsPrivateChat         bool      `json:"is_private_chat"`
	BurnDuration          int32     `json:"burn_duration"`
	GroupAtType           int32     `json:"group_at_type"`
	AttachedInfo          string    `json:"attached_info"`
	Ex                    string    `json:"ex"`
	MaxSeq                int64     `json:"max_seq"`
	MinSeq                int64     `json:"min_seq"`
	CreateTime            time.Time `json:"create_time"`
	IsMsgDestruct         bool      `json:"is_msg_destruct"`
	MsgDestructTime       int64     `json:"msg_destruct_time"`
	LatestMsgDestructTime time.Time `json:"latest_msg_destruct_time"`
}

type CallbackAfterCreateSingleChatConversationsResp struct {
	CommonCallbackResp
}

type CallbackBeforeCreateGroupChatConversationsReq struct {
	CallbackCommand       `json:"callbackCommand"`
	OwnerUserID           string    `json:"owner_user_id"`
	ConversationID        string    `json:"conversation_id"`
	ConversationType      int32     `json:"conversation_type"`
	UserID                string    `json:"user_id"`
	GroupID               string    `json:"group_id"`
	RecvMsgOpt            int32     `json:"recv_msg_opt"`
	IsPinned              bool      `json:"is_pinned"`
	IsPrivateChat         bool      `json:"is_private_chat"`
	BurnDuration          int32     `json:"burn_duration"`
	GroupAtType           int32     `json:"group_at_type"`
	AttachedInfo          string    `json:"attached_info"`
	Ex                    string    `json:"ex"`
	MaxSeq                int64     `json:"max_seq"`
	MinSeq                int64     `json:"min_seq"`
	CreateTime            time.Time `json:"create_time"`
	IsMsgDestruct         bool      `json:"is_msg_destruct"`
	MsgDestructTime       int64     `json:"msg_destruct_time"`
	LatestMsgDestructTime time.Time `json:"latest_msg_destruct_time"`
}

type CallbackBeforeCreateGroupChatConversationsResp struct {
	CommonCallbackResp
	OwnerUserID           *string    `json:"owner_user_id"`
	ConversationID        *string    `json:"conversation_id"`
	ConversationType      *int32     `json:"conversation_type"`
	UserID                *string    `json:"user_id"`
	GroupID               *string    `json:"group_id"`
	RecvMsgOpt            *int32     `json:"recv_msg_opt"`
	IsPinned              *bool      `json:"is_pinned"`
	IsPrivateChat         *bool      `json:"is_private_chat"`
	BurnDuration          *int32     `json:"burn_duration"`
	GroupAtType           *int32     `json:"group_at_type"`
	AttachedInfo          *string    `json:"attached_info"`
	Ex                    *string    `json:"ex"`
	MaxSeq                *int64     `json:"max_seq"`
	MinSeq                *int64     `json:"min_seq"`
	CreateTime            *time.Time `json:"create_time"`
	IsMsgDestruct         *bool      `json:"is_msg_destruct"`
	MsgDestructTime       *int64     `json:"msg_destruct_time"`
	LatestMsgDestructTime *time.Time `json:"latest_msg_destruct_time"`
}

type CallbackAfterCreateGroupChatConversationsReq struct {
	CallbackCommand       `json:"callbackCommand"`
	OwnerUserID           string    `json:"owner_user_id"`
	ConversationID        string    `json:"conversation_id"`
	ConversationType      int32     `json:"conversation_type"`
	UserID                string    `json:"user_id"`
	GroupID               string    `json:"group_id"`
	RecvMsgOpt            int32     `json:"recv_msg_opt"`
	IsPinned              bool      `json:"is_pinned"`
	IsPrivateChat         bool      `json:"is_private_chat"`
	BurnDuration          int32     `json:"burn_duration"`
	GroupAtType           int32     `json:"group_at_type"`
	AttachedInfo          string    `json:"attached_info"`
	Ex                    string    `json:"ex"`
	MaxSeq                int64     `json:"max_seq"`
	MinSeq                int64     `json:"min_seq"`
	CreateTime            time.Time `json:"create_time"`
	IsMsgDestruct         bool      `json:"is_msg_destruct"`
	MsgDestructTime       int64     `json:"msg_destruct_time"`
	LatestMsgDestructTime time.Time `json:"latest_msg_destruct_time"`
}

type CallbackAfterCreateGroupChatConversationsResp struct {
	CommonCallbackResp
}
