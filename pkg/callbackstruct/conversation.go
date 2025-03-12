package callbackstruct

type CallbackBeforeCreateSingleChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	RecvID           string `json:"recvID"`
	SendID           string `json:"sendID"`
	ConversationID   string `json:"conversationID"`
	ConversationType int32  `json:"conversationType"`
}

type CallbackBeforeCreateSingleChatConversationsResp struct {
	CommonCallbackResp
	RecvID           *string `json:"recvID"`
	SendID           *string `json:"sendID"`
	ConversationID   *string `json:"conversationID"`
	ConversationType *int32  `json:"conversationType"`
}

type CallbackAfterCreateSingleChatConversationsReq struct {
	CallbackCommand  `json:"callbackCommand"`
	RecvID           string `json:"recvID"`
	SendID           string `json:"sendID"`
	ConversationID   string `json:"conversationID"`
	ConversationType int32  `json:"conversationType"`
}

type CallbackAfterCreateSingleChatConversationsResp struct {
	CommonCallbackResp
}

type CallbackBeforeCreateGroupChatConversationsReq struct {
	CallbackCommand `json:"callbackCommand"`
	UserIDs         []string `json:"userIDs"`
	GroupID         string   `json:"groupID"`
}

type CallbackBeforeCreateGroupChatConversationsResp struct {
	CommonCallbackResp
	UserIDs *[]string `json:"userIDs"`
	GroupID *string   `json:"groupID"`
}

type CallbackAfterCreateGroupChatConversationsReq struct {
	CallbackCommand `json:"callbackCommand"`
	UserIDs         []string `json:"userIDs"`
	GroupID         string   `json:"groupID"`
}

type CallbackAfterCreateGroupChatConversationsResp struct {
	CommonCallbackResp
}
