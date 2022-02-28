package call_back_struct

type singleMsg struct {
	SendID      string `json:"sendID"`
	RecvID      string `json:"recvID"`
	ClientMsgID string `json:"clientMsgID"`
	ServerMsgID string `json:"serverMsgId"`
	SendTime    int64  `json:"sendTime"`
	CreateTime  int64  `json:"createTime"`
}

type CallbackBeforeSendSingleMsgReq struct {
	CommonCallbackReq
	singleMsg
}

type CallbackBeforeSendSingleMsgResp struct {
	CommonCallbackResp
}

type CallbackAfterSendSingleMsgReq struct {
	CommonCallbackReq
	singleMsg
}

type CallbackAfterSendSingleMsgResp struct {
	CommonCallbackResp
}

type groupMsg struct {

}

type CallbackBeforeSendGroupMsgReq struct {
	CommonCallbackReq
}

type CallbackBeforeSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackAfterSendGroupMsgReq struct {
	CommonCallbackReq
}

type CallbackAfterSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallBackWordFilterReq struct {
	CommonCallbackReq
	Content []byte `json:"content"`
	SendID  string `json:"SendID"`
	RecvID  string `json:"RecvID"`
	GroupID string `json:"GroupID"`
}

type CallBackWordFilterResp struct {
	CommonCallbackResp
	Content []byte `json:"content"`
}
