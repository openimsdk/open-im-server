package call_back_struct

type msg struct {
	SendID      string `json:"sendID"`
	ClientMsgID string `json:"clientMsgID"`
	ServerMsgID string `json:"serverMsgId"`
	CreateTime  int64  `json:"createTime"`
	Content    	[]byte `json:"content"`
}

type singleMsg struct {
	msg
	RecvID      string `json:"recvID"`
}

type CallbackBeforeSendSingleMsgReq struct {
	CommonCallbackReq
	singleMsg
}

type CallbackBeforeSendSingleMsgResp struct {
	CommonCallbackResp
	singleMsg
}

type CallbackAfterSendSingleMsgReq struct {
	CommonCallbackReq
	singleMsg
}

type CallbackAfterSendSingleMsgResp struct {
	CommonCallbackResp
}

type groupMsg struct {
	msg
	GroupID string `json:"groupID"`
}

type CallbackBeforeSendGroupMsgReq struct {
	CommonCallbackReq
	groupMsg
}

type CallbackBeforeSendGroupMsgResp struct {
	CommonCallbackResp
	groupMsg
}

type CallbackAfterSendGroupMsgReq struct {
	groupMsg
	CommonCallbackReq
}

type CallbackAfterSendGroupMsgResp struct {
	CommonCallbackResp
}

type CallbackWordFilterReq struct {
	CommonCallbackReq
	Content []byte `json:"content"`
}

type CallbackWordFilterResp struct {
	CommonCallbackResp
	Content []byte `json:"content"`
}
