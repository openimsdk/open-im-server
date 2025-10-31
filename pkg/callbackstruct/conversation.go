package callbackstruct

type CallbackBeforeCreateSingleChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	OwnerUserID      string `json:"ownerUserId"`
	ConversationID   string `json:"conversationId"`
	ConversationType int32  `json:"conversationType"`
	UserID           string `json:"userId"`
	RecvMsgOpt       int32  `json:"recvMsgOpt"`
	IsPinned         bool   `json:"isPinned"`
	IsPrivateChat    bool   `json:"isPrivateChat"`
	BurnDuration     int32  `json:"burnDuration"`
	GroupAtType      int32  `json:"groupAtType"`
	AttachedInfo     string `json:"attachedInfo"`
	Ex               string `json:"ex"`
}

type CallbackBeforeCreateSingleChatConversationsResp struct {
	CommonCallbackResp
	RecvMsgOpt    *int32  `json:"recvMsgOpt"`
	IsPinned      *bool   `json:"isPinned"`
	IsPrivateChat *bool   `json:"isPrivateChat"`
	BurnDuration  *int32  `json:"burnDuration"`
	GroupAtType   *int32  `json:"groupAtType"`
	AttachedInfo  *string `json:"attachedInfo"`
	Ex            *string `json:"ex"`
}

type CallbackAfterCreateSingleChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	OwnerUserID      string `json:"ownerUserId"`
	ConversationID   string `json:"conversationId"`
	ConversationType int32  `json:"conversationType"`
	UserID           string `json:"userId"`
	RecvMsgOpt       int32  `json:"recvMsgOpt"`
	IsPinned         bool   `json:"isPinned"`
	IsPrivateChat    bool   `json:"isPrivateChat"`
	BurnDuration     int32  `json:"burnDuration"`
	GroupAtType      int32  `json:"groupAtType"`
	AttachedInfo     string `json:"attachedInfo"`
	Ex               string `json:"ex"`
}

type CallbackAfterCreateSingleChatConversationsResp struct {
	CommonCallbackResp
}

type CallbackBeforeCreateGroupChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	OwnerUserID      string `json:"ownerUserId"`
	ConversationID   string `json:"conversationId"`
	ConversationType int32  `json:"conversationType"`
	GroupID          string `json:"groupId"`
	RecvMsgOpt       int32  `json:"recvMsgOpt"`
	IsPinned         bool   `json:"isPinned"`
	IsPrivateChat    bool   `json:"isPrivateChat"`
	BurnDuration     int32  `json:"burnDuration"`
	GroupAtType      int32  `json:"groupAtType"`
	AttachedInfo     string `json:"attachedInfo"`
	Ex               string `json:"ex"`
}

type CallbackBeforeCreateGroupChatConversationsResp struct {
	CommonCallbackResp
	RecvMsgOpt    *int32  `json:"recvMsgOpt"`
	IsPinned      *bool   `json:"isPinned"`
	IsPrivateChat *bool   `json:"isPrivateChat"`
	BurnDuration  *int32  `json:"burnDuration"`
	GroupAtType   *int32  `json:"groupAtType"`
	AttachedInfo  *string `json:"attachedInfo"`
	Ex            *string `json:"ex"`
}

type CallbackAfterCreateGroupChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	OwnerUserID      string `json:"ownerUserId"`
	ConversationID   string `json:"conversationId"`
	ConversationType int32  `json:"conversationType"`
	GroupID          string `json:"groupId"`
	RecvMsgOpt       int32  `json:"recvMsgOpt"`
	IsPinned         bool   `json:"isPinned"`
	IsPrivateChat    bool   `json:"isPrivateChat"`
	BurnDuration     int32  `json:"burnDuration"`
	GroupAtType      int32  `json:"groupAtType"`
	AttachedInfo     string `json:"attachedInfo"`
	Ex               string `json:"ex"`
}

type CallbackAfterCreateGroupChatConversationsResp struct {
	CommonCallbackResp
}
