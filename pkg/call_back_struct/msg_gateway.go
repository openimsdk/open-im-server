package call_back_struct

type CallbackUserOnlineReq struct {
	UserStatusCallbackReq
	Token                        string `json:"token"`
	Seq                          int    `json:"seq"`
	IsAppBackgroundStatusChanged bool   `json:"isAppBackgroundStatusChanged"`
	ConnID                       string `json:"connID"`
}

type CallbackUserOnlineResp struct {
	*CommonCallbackResp
}

type CallbackUserOfflineReq struct {
	UserStatusCallbackReq
	Seq                          int    `json:"seq"`
	IsAppBackgroundStatusChanged bool   `json:"isAppBackgroundStatusChanged"`
	ConnID                       string `json:"connID"`
}

type CallbackUserOfflineResp struct {
	*CommonCallbackResp
}

type CallbackUserKickOffReq struct {
	UserStatusCallbackReq
	Seq int `json:"seq"`
}

type CallbackUserKickOffResp struct {
	*CommonCallbackResp
}
