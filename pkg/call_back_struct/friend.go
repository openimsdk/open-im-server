package call_back_struct

type CallbackBeforeAddFriendReq struct {
	CallbackCommand string `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
	OperationID     string `json:"operationID"`
}

type CallbackBeforeAddFriendResp struct {
	*CommonCallbackResp
}
