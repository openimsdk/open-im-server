package apistruct

import (
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type KickGroupMemberReq struct {
	GroupID          string   `json:"groupID" binding:"required"`
	KickedUserIDList []string `json:"kickedUserIDList" binding:"required"`
	Reason           string   `json:"reason"`
	OperationID      string   `json:"operationID" binding:"required"`
}
type KickGroupMemberResp struct {

	//UserIDResultList []*UserIDResult `json:"data"`
}

type GetGroupMembersInfoReq struct {
	GroupID     string   `json:"groupID" binding:"required"`
	MemberList  []string `json:"memberList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}
type GetGroupMembersInfoResp struct {
	MemberList []*sdkws.GroupMemberFullInfo `json:"-"`
	Data       []map[string]interface{}     `json:"data" swaggerignore:"true"`
}

type InviteUserToGroupReq struct {
	GroupID           string   `json:"groupID" binding:"required"`
	InvitedUserIDList []string `json:"invitedUserIDList" binding:"required"`
	Reason            string   `json:"reason"`
	OperationID       string   `json:"operationID" binding:"required"`
}
type InviteUserToGroupResp struct {

	//UserIDResultList []*UserIDResult `json:"data"`
}

type GetJoinedGroupListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetJoinedGroupListResp struct {
	GroupInfoList []*sdkws.GroupInfo       `json:"-"`
	Data          []map[string]interface{} `json:"data" swaggerignore:"true"`
}

type GetGroupMemberListReq struct {
	GroupID     string `json:"groupID"`
	Filter      int32  `json:"filter"`
	NextSeq     int32  `json:"nextSeq"`
	OperationID string `json:"operationID"`
}
type GetGroupMemberListResp struct {
	NextSeq    int32                        `json:"nextSeq"`
	MemberList []*sdkws.GroupMemberFullInfo `json:"-"`
	Data       []map[string]interface{}     `json:"data" swaggerignore:"true"`
}

type GetGroupAllMemberReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
	Offset      int32  `json:"offset"`
	Count       int32  `json:"count"`
}
type GetGroupAllMemberResp struct {
	MemberList []*sdkws.GroupMemberFullInfo `json:"-"`
	Data       []map[string]interface{}     `json:"data" swaggerignore:"true"`
}

//
//type GetGroupAllMemberListBySplitReq struct {
//	GroupID     string `json:"groupID" binding:"required"`
//	OperationID string `json:"operationID" binding:"required"`
//	Offset      int32  `json:"offset" binding:"required"`
//	Count       int32  `json:"count" binding:"required"`
//}
//type GetGroupAllMemberListBySplitResp struct {
//
//	MemberList []*sdkws.GroupMemberFullInfo `json:"-"`
//	Map       []map[string]interface{}           `json:"data" swaggerignore:"true"`
//}

type CreateGroupReq struct {
	MemberList   []*GroupAddMemberInfo `json:"memberList"`
	OwnerUserID  string                `json:"ownerUserID"`
	GroupType    int32                 `json:"groupType"`
	GroupName    string                `json:"groupName"`
	Notification string                `json:"notification"`
	Introduction string                `json:"introduction"`
	FaceURL      string                `json:"faceURL"`
	Ex           string                `json:"ex"`
	OperationID  string                `json:"operationID" binding:"required"`
	GroupID      string                `json:"groupID"`
}
type CreateGroupResp struct {
	GroupInfo sdkws.GroupInfo        `json:"-"`
	Data      map[string]interface{} `json:"data" swaggerignore:"true"`
}

type GetGroupApplicationListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"` //作为管理员或群主收到的 进群申请
}
type GetGroupApplicationListResp struct {
	GroupRequestList []*sdkws.GroupRequest    `json:"-"`
	Data             []map[string]interface{} `json:"data" swaggerignore:"true"`
}

type GetUserReqGroupApplicationListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

type GetUserRespGroupApplicationResp struct {
	GroupRequestList []*sdkws.GroupRequest `json:"data"`
}

type GetGroupInfoReq struct {
	GroupIDList []string `json:"groupIDList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}
type GetGroupInfoResp struct {
	GroupInfoList []*sdkws.GroupInfo       `json:"-"`
	Data          []map[string]interface{} `json:"data" swaggerignore:"true"`
}

//type GroupInfoAlias struct {
//	sdkws.GroupInfo
//	NeedVerification int32 `protobuf:"bytes,13,opt,name=needVerification" json:"needVerification,omitempty"`
//}

//type GroupInfoAlias struct {
//	GroupID          string `protobuf:"bytes,1,opt,name=groupID" json:"groupID,omitempty"`
//	GroupName        string `protobuf:"bytes,2,opt,name=groupName" json:"groupName,omitempty"`
//	Notification     string `protobuf:"bytes,3,opt,name=notification" json:"notification,omitempty"`
//	Introduction     string `protobuf:"bytes,4,opt,name=introduction" json:"introduction,omitempty"`
//	FaceURL          string `protobuf:"bytes,5,opt,name=faceURL" json:"faceURL,omitempty"`
//	OwnerUserID      string `protobuf:"bytes,6,opt,name=ownerUserID" json:"ownerUserID,omitempty"`
//	CreateTime       uint32 `protobuf:"varint,7,opt,name=createTime" json:"createTime,omitempty"`
//	MemberCount      uint32 `protobuf:"varint,8,opt,name=memberCount" json:"memberCount,omitempty"`
//	Ex               string `protobuf:"bytes,9,opt,name=ex" json:"ex,omitempty"`
//	Status           int32  `protobuf:"varint,10,opt,name=status" json:"status,omitempty"`
//	CreatorUserID    string `protobuf:"bytes,11,opt,name=creatorUserID" json:"creatorUserID,omitempty"`
//	GroupType        int32  `protobuf:"varint,12,opt,name=groupType" json:"groupType,omitempty"`
//	NeedVerification int32  `protobuf:"bytes,13,opt,name=needVerification" json:"needVerification,omitempty"`
//}

type ApplicationGroupResponseReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	GroupID      string `json:"groupID" binding:"required"`
	FromUserID   string `json:"fromUserID" binding:"required"` //application from FromUserID
	HandledMsg   string `json:"handledMsg"`
	HandleResult int32  `json:"handleResult" binding:"required,oneof=-1 1"`
}
type ApplicationGroupResponseResp struct {
}

type JoinGroupReq struct {
	GroupID       string `json:"groupID" binding:"required"`
	ReqMessage    string `json:"reqMessage"`
	OperationID   string `json:"operationID" binding:"required"`
	JoinSource    int32  `json:"joinSource"`
	InviterUserID string `json:"inviterUserID"`
}

type JoinGroupResp struct {
}

type QuitGroupReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type QuitGroupResp struct {
}

type SetGroupInfoReq struct {
	GroupID           string `json:"groupID" binding:"required"`
	GroupName         string `json:"groupName"`
	Notification      string `json:"notification"`
	Introduction      string `json:"introduction"`
	FaceURL           string `json:"faceURL"`
	Ex                string `json:"ex"`
	OperationID       string `json:"operationID" binding:"required"`
	NeedVerification  *int32 `json:"needVerification"`
	LookMemberInfo    *int32 `json:"lookMemberInfo"`
	ApplyMemberFriend *int32 `json:"applyMemberFriend"`
}

type SetGroupInfoResp struct {
}

type TransferGroupOwnerReq struct {
	GroupID        string `json:"groupID" binding:"required"`
	OldOwnerUserID string `json:"oldOwnerUserID" binding:"required"`
	NewOwnerUserID string `json:"newOwnerUserID" binding:"required"`
	OperationID    string `json:"operationID" binding:"required"`
}
type TransferGroupOwnerResp struct {
}

type DismissGroupReq struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type DismissGroupResp struct {
}

type MuteGroupMemberReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	GroupID      string `json:"groupID" binding:"required"`
	UserID       string `json:"userID" binding:"required"`
	MutedSeconds uint32 `json:"mutedSeconds" binding:"required"`
}
type MuteGroupMemberResp struct {
}

type CancelMuteGroupMemberReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}
type CancelMuteGroupMemberResp struct {
}

type MuteGroupReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
}
type MuteGroupResp struct {
}

type CancelMuteGroupReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
}
type CancelMuteGroupResp struct {
}

type SetGroupMemberNicknameReq struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
	Nickname    string `json:"nickname"`
}

type SetGroupMemberNicknameResp struct {
}

type SetGroupMemberInfoReq struct {
	OperationID string  `json:"operationID" binding:"required"`
	GroupID     string  `json:"groupID" binding:"required"`
	UserID      string  `json:"userID" binding:"required"`
	Nickname    *string `json:"nickname"`
	FaceURL     *string `json:"userGroupFaceUrl"`
	RoleLevel   *int32  `json:"roleLevel" validate:"gte=1,lte=3"`
	Ex          *string `json:"ex"`
}

type SetGroupMemberInfoResp struct {
}

type GetGroupAbstractInfoReq struct {
	OperationID string `json:"operationID"`
	GroupID     string `json:"groupID"`
}

type GetGroupAbstractInfoResp struct {
	GroupMemberNumber   int32  `json:"groupMemberNumber"`
	GroupMemberListHash uint64 `json:"groupMemberListHash"`
}
