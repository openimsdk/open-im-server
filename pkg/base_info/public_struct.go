package base_info

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserInfo struct {
	UserID      string `json:"userID" binding:"required,min=1,max=64"`
	Nickname    string `json:"nickname" binding:"required,min=1,max=64"`
	FaceUrl     string `json:"faceUrl" binding:"omitempty,max=1024"`
	Gender      int32  `json:"gender" binding:"omitempty,oneof=0 1 2"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=32"`
	Birth       string `json:"birth" binding:"omitempty,max=16"`
	Email       string `json:"email" binding:"omitempty,max=64"`
	Ex          string `json:"ex" binding:"omitempty,max=1024"`
}

func SetErrCodeMsg(c *gin.Context, status int) *CommResp {
	resp := CommResp{ErrCode: int32(status), ErrMsg: http.StatusText(status)}
	c.JSON(status, resp)
	return &resp
}

//
//type GroupInfo struct {
//	GroupID       string `json:"groupID"`
//	GroupName     string `json:"groupName"`
//	Notification  string `json:"notification"`
//	Introduction  string `json:"introduction"`
//	FaceUrl       string `json:"faceUrl"`
//	OperationID   string `json:"operationID"`
//	OwnerUserID   string `json:"ownerUserID"`
//	CreateTime    int64  `json:"createTime"`
//	MemberCount   uint32 `json:"memberCount"`
//	Ex            string `json:"ex"`
//	Status        int32  `json:"status"`
//	CreatorUserID string `json:"creatorUserID"`
//	GroupType     int32  `json:"groupType"`
//}
//
//type GroupMemberFullInfo struct {
//	GroupID        string `json:"groupID"`
//	UserID         string `json:"userID"`
//	RoleLevel      int32  `json:"roleLevel"`
//	JoinTime       uint64 `json:"joinTime"`
//	Nickname       string `json:"nickname"`
//	FaceUrl        string `json:"faceUrl"`
//	FriendRemark   string `json:"friendRemark"`
//	AppMangerLevel int32  `json:"appMangerLevel"`
//	JoinSource     int32  `json:"joinSource"`
//	OperatorUserID string `json:"operatorUserID"`
//	Ex             string `json:"ex"`
//}
//
//type PublicUserInfo struct {
//	UserID   string `json:"userID"`
//	Nickname string `json:"nickname"`
//	FaceUrl  string `json:"faceUrl"`
//	Gender   int32  `json:"gender"`
//}
//
//type UserInfo struct {
//	UserID   string `json:"userID"`
//	Nickname string `json:"nickname"`
//	FaceUrl  string `json:"faceUrl"`
//	Gender   int32  `json:"gender"`
//	Mobile   string `json:"mobile"`
//	Birth    string `json:"birth"`
//	Email    string `json:"email"`
//	Ex       string `json:"ex"`
//}
//
//type FriendInfo struct {
//	OwnerUserID    string   `json:"ownerUserID"`
//	Remark         string   `json:"remark"`
//	CreateTime     int64    `json:"createTime"`
//	FriendUser     UserInfo `json:"friendUser"`
//	AddSource      int32    `json:"addSource"`
//	OperatorUserID string   `json:"operatorUserID"`
//	Ex             string   `json:"ex"`
//}
//
//type BlackInfo struct {
//	OwnerUserID    string         `json:"ownerUserID"`
//	CreateTime     int64          `json:"createTime"`
//	BlackUser      PublicUserInfo `json:"friendUser"`
//	AddSource      int32          `json:"addSource"`
//	OperatorUserID string         `json:"operatorUserID"`
//	Ex             string         `json:"ex"`
//}
//
//type GroupRequest struct {
//	UserID       string `json:"userID"`
//	GroupID      string `json:"groupID"`
//	HandleResult string `json:"handleResult"`
//	ReqMsg       string `json:"reqMsg"`
//	HandleMsg    string `json:"handleMsg"`
//	ReqTime      int64  `json:"reqTime"`
//	HandleUserID string `json:"handleUserID"`
//	HandleTime   int64  `json:"handleTime"`
//	Ex           string `json:"ex"`
//}
//
//type FriendRequest struct {
//	FromUserID    string `json:"fromUserID"`
//	ToUserID      string `json:"toUserID"`
//	HandleResult  int32  `json:"handleResult"`
//	ReqMessage    string `json:"reqMessage"`
//	CreateTime    int64  `json:"createTime"`
//	HandlerUserID string `json:"handlerUserID"`
//	HandleMsg     string `json:"handleMsg"`
//	HandleTime    int64  `json:"handleTime"`
//	Ex            string `json:"ex"`
//}
//
//
//
