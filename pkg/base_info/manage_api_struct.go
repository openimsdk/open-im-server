package base_info

import (
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/proto/sdk_ws"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
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

type ManagementSendMsg struct {
	OperationID         string `json:"operationID" binding:"required"`
	BusinessOperationID string `json:"businessOperationID"`
	SendID              string `json:"sendID" binding:"required"`
	GroupID             string `json:"groupID" `
	SenderNickname      string `json:"senderNickname" `
	SenderFaceURL       string `json:"senderFaceURL" `
	SenderPlatformID    int32  `json:"senderPlatformID"`
	//ForceList        []string                     `json:"forceList" `
	Content         map[string]interface{}       `json:"content" binding:"required" swaggerignore:"true"`
	ContentType     int32                        `json:"contentType" binding:"required"`
	SessionType     int32                        `json:"sessionType" binding:"required"`
	IsOnlineOnly    bool                         `json:"isOnlineOnly"`
	NotOfflinePush  bool                         `json:"notOfflinePush"`
	OfflinePushInfo *open_im_sdk.OfflinePushInfo `json:"offlinePushInfo"`
}

type ManagementSendMsgReq struct {
	ManagementSendMsg
	RecvID string `json:"recvID" `
}

type ManagementSendMsgResp struct {
	CommResp
	ResultList server_api_params.UserSendMsgResp `json:"data"`
}

type ManagementBatchSendMsgReq struct {
	ManagementSendMsg
	RecvIDList []string `json:"recvIDList"`
}

type ManagementBatchSendMsgResp struct {
	CommResp
	Data struct {
		ResultList   []*SingleReturnResult `json:"resultList"`
		FailedIDList []string
	} `json:"data"`
}
type SingleReturnResult struct {
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
	RecvID      string `json:"recvID" `
}

type CheckMsgIsSendSuccessReq struct {
	OperationID string `json:"operationID"`
}

type CheckMsgIsSendSuccessResp struct {
	CommResp
	Status int32 `json:"status"`
}
