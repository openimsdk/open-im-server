package cms_api_struct

import server_api_params "Open_IM/pkg/proto/sdk_ws"

type AdminLoginRequest struct {
	AdminName string `json:"admin_name" binding:"required"`
	Secret    string `json:"secret" binding:"required"`
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
	Operation   int32    `json:"operation" binding:"required"`
}

type ReduceUserRegisterAddFriendIDListResponse struct {
}

type GetUserRegisterAddFriendIDListRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	RequestPagination
}

type GetUserRegisterAddFriendIDListResponse struct {
	Users []*server_api_params.UserInfo `json:"Users"`
	ResponsePagination
}
