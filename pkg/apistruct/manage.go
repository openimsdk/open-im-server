package apistruct

import (
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type DeleteUsersReq struct {
	OperationID      string   `json:"operationID" binding:"required"`
	DeleteUserIDList []string `json:"deleteUserIDList" binding:"required"`
}
type DeleteUsersResp struct {
	FailedUserIDList []string `json:"data"`
}
type GetAllUsersUidReq struct {
	OperationID string `json:"operationID" binding:"required"`
}
type GetAllUsersUidResp struct {
	UserIDList []string `json:"data"`
}
type GetUsersOnlineStatusReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
}
type GetUsersOnlineStatusResp struct {

	//SuccessResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult `json:"data"`
}
type AccountCheckReq struct {
	OperationID     string   `json:"operationID" binding:"required"`
	CheckUserIDList []string `json:"checkUserIDList" binding:"required,lte=100"`
}
type AccountCheckResp struct {
}

type ManagementSendMsg struct {
	SendID           string                 `json:"sendID" binding:"required"`
	GroupID          string                 `json:"groupID" binding:"required_if=SessionType 2|required_if=SessionType 3"`
	SenderNickname   string                 `json:"senderNickname" `
	SenderFaceURL    string                 `json:"senderFaceURL" `
	SenderPlatformID int32                  `json:"senderPlatformID"`
	Content          map[string]interface{} `json:"content" binding:"required" swaggerignore:"true"`
	ContentType      int32                  `json:"contentType" binding:"required"`
	SessionType      int32                  `json:"sessionType" binding:"required"`
	IsOnlineOnly     bool                   `json:"isOnlineOnly"`
	NotOfflinePush   bool                   `json:"notOfflinePush"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}

type ManagementSendMsgReq struct {
	SendID           string                 `json:"sendID" binding:"required"`
	RecvID           string                 `json:"recvID" binding:"required_if" message:"recvID is required if sessionType is SingleChatType or NotificationChatType"`
	GroupID          string                 `json:"groupID" binding:"required_if" message:"groupID is required if sessionType is GroupChatType or SuperGroupChatType"`
	SenderNickname   string                 `json:"senderNickname" `
	SenderFaceURL    string                 `json:"senderFaceURL" `
	SenderPlatformID int32                  `json:"senderPlatformID"`
	Content          map[string]interface{} `json:"content" binding:"required" swaggerignore:"true"`
	ContentType      int32                  `json:"contentType" binding:"required"`
	SessionType      int32                  `json:"sessionType" binding:"required"`
	IsOnlineOnly     bool                   `json:"isOnlineOnly"`
	NotOfflinePush   bool                   `json:"notOfflinePush"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
}

type ManagementSendMsgResp struct {
	ResultList sdkws.UserSendMsgResp `json:"data"`
}

type ManagementBatchSendMsgReq struct {
	ManagementSendMsg
	IsSendAll  bool     `json:"isSendAll"`
	RecvIDList []string `json:"recvIDList"`
}

type ManagementBatchSendMsgResp struct {
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
	Status int32 `json:"status"`
}
