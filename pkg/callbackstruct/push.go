package callbackstruct

import common "github.com/openimsdk/protocol/sdkws"

type CallbackBeforePushReq struct {
	UserStatusBatchCallbackReq
	*common.OfflinePushInfo
	ClientMsgID string   `json:"clientMsgID"`
	SendID      string   `json:"sendID"`
	GroupID     string   `json:"groupID"`
	ContentType int32    `json:"contentType"`
	SessionType int32    `json:"sessionType"`
	AtUserIDs   []string `json:"atUserIDList"`
	Content     string   `json:"content"`
}

type CallbackBeforePushResp struct {
	CommonCallbackResp
	UserIDs         []string                `json:"userIDList"`
	OfflinePushInfo *common.OfflinePushInfo `json:"offlinePushInfo"`
}

type CallbackBeforeSuperGroupOnlinePushReq struct {
	UserStatusBaseCallback
	ClientMsgID string   `json:"clientMsgID"`
	SendID      string   `json:"sendID"`
	GroupID     string   `json:"groupID"`
	ContentType int32    `json:"contentType"`
	SessionType int32    `json:"sessionType"`
	AtUserIDs   []string `json:"atUserIDList"`
	Content     string   `json:"content"`
	Seq         int64    `json:"seq"`
}

type CallbackBeforeSuperGroupOnlinePushResp struct {
	CommonCallbackResp
	UserIDs         []string                `json:"userIDList"`
	OfflinePushInfo *common.OfflinePushInfo `json:"offlinePushInfo"`
}
