package callbackstruct

type CallbackBeforeCreateSingleChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	OwnerUserID      string `json:"owner_user_id"`
	ConversationID   string `json:"conversation_id"`
	ConversationType int32  `json:"conversation_type"`
	UserID           string `json:"user_id"`
	RecvMsgOpt       int32  `json:"recv_msg_opt"`
	IsPinned         bool   `json:"is_pinned"`
	IsPrivateChat    bool   `json:"is_private_chat"`
	BurnDuration     int32  `json:"burn_duration"`
	GroupAtType      int32  `json:"group_at_type"`
	AttachedInfo     string `json:"attached_info"`
	Ex               string `json:"ex"`
}

type CallbackBeforeCreateSingleChatConversationsResp struct {
	CommonCallbackResp
	RecvMsgOpt    *int32  `json:"recv_msg_opt"`
	IsPinned      *bool   `json:"is_pinned"`
	IsPrivateChat *bool   `json:"is_private_chat"`
	BurnDuration  *int32  `json:"burn_duration"`
	GroupAtType   *int32  `json:"group_at_type"`
	AttachedInfo  *string `json:"attached_info"`
	Ex            *string `json:"ex"`
}

type CallbackAfterCreateSingleChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	OwnerUserID      string `json:"owner_user_id"`
	ConversationID   string `json:"conversation_id"`
	ConversationType int32  `json:"conversation_type"`
	UserID           string `json:"user_id"`
	RecvMsgOpt       int32  `json:"recv_msg_opt"`
	IsPinned         bool   `json:"is_pinned"`
	IsPrivateChat    bool   `json:"is_private_chat"`
	BurnDuration     int32  `json:"burn_duration"`
	GroupAtType      int32  `json:"group_at_type"`
	AttachedInfo     string `json:"attached_info"`
	Ex               string `json:"ex"`
}

type CallbackAfterCreateSingleChatConversationsResp struct {
	CommonCallbackResp
}

type CallbackBeforeCreateGroupChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	OwnerUserID      string `json:"owner_user_id"`
	ConversationID   string `json:"conversation_id"`
	ConversationType int32  `json:"conversation_type"`
	GroupID          string `json:"group_id"`
	RecvMsgOpt       int32  `json:"recv_msg_opt"`
	IsPinned         bool   `json:"is_pinned"`
	IsPrivateChat    bool   `json:"is_private_chat"`
	BurnDuration     int32  `json:"burn_duration"`
	GroupAtType      int32  `json:"group_at_type"`
	AttachedInfo     string `json:"attached_info"`
	Ex               string `json:"ex"`
}

type CallbackBeforeCreateGroupChatConversationsResp struct {
	CommonCallbackResp
	RecvMsgOpt    *int32  `json:"recv_msg_opt"`
	IsPinned      *bool   `json:"is_pinned"`
	IsPrivateChat *bool   `json:"is_private_chat"`
	BurnDuration  *int32  `json:"burn_duration"`
	GroupAtType   *int32  `json:"group_at_type"`
	AttachedInfo  *string `json:"attached_info"`
	Ex            *string `json:"ex"`
}

type CallbackAfterCreateGroupChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	OwnerUserID      string `json:"owner_user_id"`
	ConversationID   string `json:"conversation_id"`
	ConversationType int32  `json:"conversation_type"`
	GroupID          string `json:"group_id"`
	RecvMsgOpt       int32  `json:"recv_msg_opt"`
	IsPinned         bool   `json:"is_pinned"`
	IsPrivateChat    bool   `json:"is_private_chat"`
	BurnDuration     int32  `json:"burn_duration"`
	GroupAtType      int32  `json:"group_at_type"`
	AttachedInfo     string `json:"attached_info"`
	Ex               string `json:"ex"`
}

type CallbackAfterCreateGroupChatConversationsResp struct {
	CommonCallbackResp
}
