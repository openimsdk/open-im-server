package call_back_struct

type Msg struct {
	SendID      string `json:"sendID"`
	CreateTime  int64  `json:"createTime"`
	Content    	string `json:"content"`
}

type SingleMsg struct {
	Msg
	RecvID      string `json:"recvID"`
}

type CallbackBeforeSendSingleMsgReq struct {
	CommonCallbackReq
	SingleMsg
}

type CallbackBeforeSendSingleMsgResp struct {
	CommonCallbackResp
	SingleMsg
}

type CallbackAfterSendSingleMsgReq struct {
	CommonCallbackReq
	SingleMsg
}

type CallbackAfterSendSingleMsgResp struct {
	CommonCallbackResp
}

type GroupMsg struct {
	Msg
	GroupID string `json:"groupID"`
}

type CallbackBeforeSendGroupMsgReq struct {
	CommonCallbackReq
	GroupMsg
}

type CallbackBeforeSendGroupMsgResp struct {
	CommonCallbackResp
	GroupMsg
}

type CallbackAfterSendGroupMsgReq struct {
	GroupMsg
	CommonCallbackReq
}

type CallbackAfterSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackWordFilterReq struct {
	CommonCallbackReq
	Content string `json:"content"`
}

type CallbackWordFilterResp struct {
	CommonCallbackResp
	Content string `json:"content"`
}
