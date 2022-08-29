package cms_api_struct

import (
	"Open_IM/pkg/base_info"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
)

type AdminLoginRequest struct {
	AdminName   string `json:"adminID" binding:"required"`
	Secret      string `json:"secret" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type AddUserRegisterAddFriendIDListRequest struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}

type AddUserRegisterAddFriendIDListResponse struct {
}

type ReduceUserRegisterAddFriendIDListRequest struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
	Operation   int32    `json:"operation"`
}

type ReduceUserRegisterAddFriendIDListResponse struct {
}

type GetUserRegisterAddFriendIDListRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	base_info.RequestPagination
}

type GetUserRegisterAddFriendIDListResponse struct {
	Users []*server_api_params.UserInfo `json:"users"`
	base_info.ResponsePagination
}
