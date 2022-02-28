package call_back_struct

type CommonCallbackReq struct {
	CallbackCommand string `json:"callbackCommand"`
	ServerMsgID string `json:"serverID"`
	ClientID string `json:"clientID"`
	OperationID string `json:"operationID"`
}

type CommonCallbackResp struct {
	ActionCode int `json:"actionCode"`
	ErrCode int `json:"errCode"`
	ErrMsg string `json:"errMsg"`
	OperationID string `json:"operationID"`
}


