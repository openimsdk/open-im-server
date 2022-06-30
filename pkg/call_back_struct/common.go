package call_back_struct

type CommonCallbackReq struct {
	SendID           string   `json:"sendID"`
	CallbackCommand  string   `json:"callbackCommand"`
	ServerMsgID      string   `json:"serverMsgID"`
	ClientMsgID      string   `json:"clientMsgID"`
	OperationID      string   `json:"operationID"`
	SenderPlatformID int32    `json:"senderPlatformID"`
	SenderNickname   string   `json:"senderNickname"`
	SessionType      int32    `json:"sessionType"`
	MsgFrom          int32    `json:"msgFrom"`
	ContentType      int32    `json:"contentType"`
	Status           int32    `json:"status"`
	CreateTime       int64    `json:"createTime"`
	Content          string   `json:"content"`
	AtUserIDList     []string `json:"atUserList"`
	SenderFaceURL    string   `json:"faceURL"`
}

type CommonCallbackResp struct {
	ActionCode  int    `json:"actionCode"`
	ErrCode     int    `json:"errCode"`
	ErrMsg      string `json:"errMsg"`
	OperationID string `json:"operationID"`
}

type UserStatusCallbackReq struct {
	CallbackCommand string `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	UserID          string `json:"userID"`
	PlatformID      int32  `json:"platformID"`
	Platform        string `json:"platform"`
}
