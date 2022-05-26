package call_back_struct

type CallbackUserOnlineReq struct {
	UserStatusCallbackReq
	Token string `json:"token"`
}

type CallbackUserOnlineResp struct {
	CommonCallbackResp
}

type CallbackUserOfflineReq struct {
	UserStatusCallbackReq
}

type CallbackUserOfflineResp struct {
	CommonCallbackResp
}
