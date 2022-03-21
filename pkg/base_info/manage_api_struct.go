package base_info

import (
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
)

type DeleteUsersReq struct {
	OperationID      string   `json:"operationID" binding:"required"`
	DeleteUserIDList []string `json:"deleteUserIDList" binding:"required"`
}
type DeleteUsersResp struct {
	CommResp
	FailedUserIDList []string `json:"data"`
}
type GetAllUsersUidReq struct {
	OperationID string `json:"operationID" binding:"required"`
}
type GetAllUsersUidResp struct {
	CommResp
	UserIDList []string `json:"data"`
}
type GetUsersOnlineStatusReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
}
type GetUsersOnlineStatusResp struct {
	CommResp
	SuccessResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult `json:"data"`
}
type AccountCheckReq struct {
	OperationID     string   `json:"operationID" binding:"required"`
	CheckUserIDList []string `json:"checkUserIDList" binding:"required,lte=100"`
}
type AccountCheckResp struct {
	CommResp
	ResultList []*pbUser.AccountCheckResp_SingleUserStatus `json:"data"`
}

type ManagementSendMsgResp struct {
	CommResp
	ResultList server_api_params.UserSendMsgResp `json:"data"`
}
