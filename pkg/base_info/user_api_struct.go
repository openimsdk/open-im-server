package base_info

import (
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
)

type GetUserInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}
type GetUserInfoResp struct {
	CommResp
	UserInfoList []*open_im_sdk.UserInfo `json:"data"`
}
