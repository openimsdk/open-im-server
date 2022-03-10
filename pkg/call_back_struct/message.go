package call_back_struct


type CallbackBeforeSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID  string `json:"recvID"`
}

type CallbackBeforeSendSingleMsgResp struct {
	CommonCallbackResp
}

type CallbackAfterSendSingleMsgReq struct {
	CommonCallbackReq
	RecvID  string `json:"recvID"`
}

type CallbackAfterSendSingleMsgResp struct {
	CommonCallbackResp
}

type CallbackBeforeSendGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackBeforeSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackAfterSendGroupMsgReq struct {
	CommonCallbackReq
	GroupID string `json:"groupID"`
}

type CallbackAfterSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackWordFilterReq struct {
	CommonCallbackReq
}

type CallbackWordFilterResp struct {
	CommonCallbackResp
	Content string `json:"content"`
}
