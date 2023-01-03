package call_back_struct

import commonPb "Open_IM/pkg/proto/sdk_ws"

type CallbackBeforePushReq struct {
	UserStatusBatchCallbackReq
	*commonPb.OfflinePushInfo
	ClientMsgID  string   `json:"clientMsgID"`
	SendID       string   `json:"sendID"`
	GroupID      string   `json:"groupID"`
	ContentType  int32    `json:"contentType"`
	SessionType  int32    `json:"sessionType"`
	AtUserIDList []string `json:"atUserIDList"`
	Content      string   `json:"content"`
}

type CallbackBeforePushResp struct {
	*CommonCallbackResp
	UserIDList      []string                  `json:"userIDList"`
	OfflinePushInfo *commonPb.OfflinePushInfo `json:"offlinePushInfo"`
}

type CallbackBeforeSuperGroupOnlinePushReq struct {
	//*commonPb.OfflinePushInfo
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
	*CommonCallbackResp
	UserIDList      []string                  `json:"userIDList"`
	OfflinePushInfo *commonPb.OfflinePushInfo `json:"offlinePushInfo"`
}
