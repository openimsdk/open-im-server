package callbackstruct

type CallbackAfterRevokeMsgReq struct {
	CallbackCommand `json:"callbackCommand"`
	ConversationID  string `json:"conversationID"`
	Seq             int64  `json:"seq"`
	UserID          string `json:"userID"`
	EventTime       int64  `json:"eventTime"`
}
type CallbackAfterRevokeMsgResp struct {
	CommonCallbackResp
}
