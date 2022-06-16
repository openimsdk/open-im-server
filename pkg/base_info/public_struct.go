package base_info

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiUserInfo struct {
	UserID      string `json:"userID" binding:"required,min=1,max=64"`
	Nickname    string `json:"nickname" binding:"omitempty,min=1,max=64"`
	FaceURL     string `json:"faceURL" binding:"omitempty,max=1024"`
	Gender      int32  `json:"gender" binding:"omitempty,oneof=0 1 2"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=32"`
	Birth       uint32 `json:"birth" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,max=64"`
	Ex          string `json:"ex" binding:"omitempty,max=1024"`
}

//type Conversation struct {
//	OwnerUserID      string `gorm:"column:owner_user_id;primary_key;type:char(128)" json:"OwnerUserID"`
//	ConversationID   string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
//	ConversationType int32  `gorm:"column:conversation_type" json:"conversationType"`
//	UserID           string `gorm:"column:user_id;type:char(64)" json:"userID"`
//	GroupID          string `gorm:"column:group_id;type:char(128)" json:"groupID"`
//	RecvMsgOpt       int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
//	UnreadCount      int32  `gorm:"column:unread_count" json:"unreadCount"`
//	DraftTextTime    int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
//	IsPinned         bool   `gorm:"column:is_pinned" json:"isPinned"`
//	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
//	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
//}

type GroupAddMemberInfo struct {
	UserID    string `json:"userID" binding:"required"`
	RoleLevel int32  `json:"roleLevel" binding:"required"`
}

func SetErrCodeMsg(c *gin.Context, status int) *CommResp {
	resp := CommResp{ErrCode: int32(status), ErrMsg: http.StatusText(status)}
	c.JSON(status, resp)
	return &resp
}

//GroupName    string                `json:"groupName"`
//	Introduction string                `json:"introduction"`
//	Notification string                `json:"notification"`
//	FaceUrl      string                `json:"faceUrl"`
//	OperationID  string                `json:"operationID" binding:"required"`
//	GroupType    int32                 `json:"groupType"`
//	Ex           string                `json:"ex"`

//type GroupInfo struct {
//	GroupID       string `json:"groupID"`
//	GroupName     string `json:"groupName"`
//	Notification  string `json:"notification"`
//	Introduction  string `json:"introduction"`
//	FaceUrl       string `json:"faceUrl"`
//	OwnerUserID   string `json:"ownerUserID"`
//	Ex            string `json:"ex"`
//	GroupType     int32  `json:"groupType"`
//}

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
