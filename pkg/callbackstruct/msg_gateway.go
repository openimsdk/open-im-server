package callbackstruct

type CallbackUserOnlineReq struct {
	UserStatusCallbackReq
	// Token           string `json:"token"`
	Seq             int64  `json:"seq"`
	IsAppBackground bool   `json:"isAppBackground"`
	ConnID          string `json:"connID"`
}

type CallbackUserOnlineResp struct {
	CommonCallbackResp
}

type CallbackUserOfflineReq struct {
	UserStatusCallbackReq
	Seq    int64  `json:"seq"`
	ConnID string `json:"connID"`
}

type CallbackUserOfflineResp struct {
	CommonCallbackResp
}

type CallbackUserKickOffReq struct {
	UserStatusCallbackReq
	Seq int64 `json:"seq"`
}

type CallbackUserKickOffResp struct {
	CommonCallbackResp
}
