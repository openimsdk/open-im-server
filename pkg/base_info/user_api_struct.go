package base_info

import (
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
)

type GetUsersOnlineStatusReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
	Secret      string   `json:"secret" binding:"required,max=32"`
}
type OnlineStatus struct {
	UserID string `json:"userID"`
	Status string `json:"status"`
}
type GetUsersOnlineStatusResp struct {
	CommResp
	OnlineStatusList []*OnlineStatus `json:"data"`
}

type GetUserInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}

type GetUserInfoResp struct {
	CommResp
	UserInfoList []*open_im_sdk.UserInfo `json:"data"`
}
