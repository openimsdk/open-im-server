package callbackstruct

import common "Open_IM/pkg/proto/sdkws"

type CallbackBeforePushReq struct {
	UserStatusBatchCallbackReq
	*common.OfflinePushInfo
	ClientMsgID  string   `json:"clientMsgID"`
	SendID       string   `json:"sendID"`
	GroupID      string   `json:"groupID"`
	ContentType  int32    `json:"contentType"`
	SessionType  int32    `json:"sessionType"`
	AtUserIDList []string `json:"atUserIDList"`
	Content      string   `json:"content"`
}

type CallbackBeforePushResp struct {
	CommonCallbackResp
	UserIDList      []string                `json:"userIDList"`
	OfflinePushInfo *common.OfflinePushInfo `json:"offlinePushInfo"`
}

type CallbackBeforeSuperGroupOnlinePushReq struct {
	//*common.OfflinePushInfo
	UserStatusBaseCallback
	ClientMsgID  string   `json:"clientMsgID"`
	SendID       string   `json:"sendID"`
	GroupID      string   `json:"groupID"`
	ContentType  int32    `json:"contentType"`
	SessionType  int32    `json:"sessionType"`
	AtUserIDList []string `json:"atUserIDList"`
	Content      string   `json:"content"`
	Seq          uint32   `json:"seq"`
}

type CallbackBeforeSuperGroupOnlinePushResp struct {
	CommonCallbackResp
	UserIDList      []string                `json:"userIDList"`
	OfflinePushInfo *common.OfflinePushInfo `json:"offlinePushInfo"`
}
