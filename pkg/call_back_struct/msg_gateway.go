package call_back_struct

type CallbackUserOnlineReq struct {
	UserStatusCallbackReq
	Token string `json:"token"`
	Seq   int    `json:"seq"`
}

type CallbackUserOnlineResp struct {
	CommonCallbackResp
}

type CallbackUserOfflineReq struct {
	UserStatusCallbackReq
	Seq int `json:"seq"`
}

type CallbackUserOfflineResp struct {
	CommonCallbackResp
}
