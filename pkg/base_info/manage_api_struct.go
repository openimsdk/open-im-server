package base_info

import (
	pbRelay "Open_IM/pkg/proto/relay"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
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
	Content         map[string]interface{}             `json:"content" binding:"required" swaggerignore:"true"`
	ContentType     int32                              `json:"contentType" binding:"required"`
	SessionType     int32                              `json:"sessionType" binding:"required"`
	IsOnlineOnly    bool                               `json:"isOnlineOnly"`
	NotOfflinePush  bool                               `json:"notOfflinePush"`
	OfflinePushInfo *server_api_params.OfflinePushInfo `json:"offlinePushInfo"`
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
	IsSendAll  bool     `json:"isSendAll"`
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

type GetUsersReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserName    string `json:"userName"`
	UserID      string `json:"userID"`
	Content     string `json:"content"`
	PageNumber  int32  `json:"pageNumber" binding:"required"`
	ShowNumber  int32  `json:"showNumber" binding:"required"`
}

type CMSUser struct {
	UserID           string `json:"userID"`
	Nickname         string `json:"nickname"`
	FaceURL          string `json:"faceURL"`
	Gender           int32  `json:"gender"`
	PhoneNumber      string `json:"phoneNumber"`
	Birth            uint32 `json:"birth"`
	Email            string `json:"email"`
	Ex               string `json:"ex"`
	CreateIp         string `json:"createIp"`
	CreateTime       uint32 `json:"createTime"`
	LastLoginIp      string `json:"LastLoginIp"`
	LastLoginTime    uint32 `json:"LastLoginTime"`
	AppMangerLevel   int32  `json:"appMangerLevel"`
	GlobalRecvMsgOpt int32  `json:"globalRecvMsgOpt"`
	IsBlock          bool   `json:"isBlock"`
}

type GetUsersResp struct {
	CommResp
	Data struct {
		UserList    []*CMSUser `json:"userList"`
		TotalNum    int32      `json:"totalNum"`
		CurrentPage int32      `json:"currentPage"`
		ShowNumber  int32      `json:"showNumber"`
	} `json:"data"`
}
