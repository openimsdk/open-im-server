package call_back_struct

import commonPb "Open_IM/pkg/proto/sdk_ws"

type CallbackOfflinePushReq struct {
	UserStatusCallbackReq
	*commonPb.OfflinePushInfo
	//CommonCallbackReq
	SendID       string   `json:"sendID"`
	GroupID      string   `json:"groupID"`
	ContentType  int32    `json:"contentType"`
	SessionType  int32    `json:"sessionType"`
	AtUserIDList []string `json:"atUserIDList"`
	Content      string   `json`
}

type CallbackOfflinePushResp struct {
	*CommonCallbackResp
}
