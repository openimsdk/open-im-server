package call_back_struct

import commonPb "Open_IM/pkg/proto/sdk_ws"

type CallbackOfflinePushReq struct {
	UserStatusCallbackReq
	*commonPb.OfflinePushInfo
}

type CallbackOfflinePushResp struct {
	CommonCallbackResp
}
