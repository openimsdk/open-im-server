package callbackstruct

type CallbackAfterRevokeMsgReq struct {
	CallbackCommand `json:"callbackCommand"`
	ConversationID  string `json:"conversationID"`
	Seq             int64  `json:"seq"`
	UserID          string `json:"userID"`
}

type CallbackAfterRevokeMsgResp struct {
	CommonCallbackResp
}
