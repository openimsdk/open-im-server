package callbackstruct

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apistruct"
	common "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type CallbackCommand string

func (c CallbackCommand) GetCallbackCommand() string {
	return string(c)
}

type CallbackBeforeCreateGroupReq struct {
	OperationID     string `json:"operationID"`
	CallbackCommand `json:"callbackCommand"`
	*common.GroupInfo
	InitMemberList []*apistruct.GroupAddMemberInfo `json:"initMemberList"`
}

type CallbackBeforeCreateGroupResp struct {
	CommonCallbackResp
	GroupID           *string `json:"groupID"`
	GroupName         *string `json:"groupName"`
	Notification      *string `json:"notification"`
	Introduction      *string `json:"introduction"`
	FaceURL           *string `json:"faceURL"`
	OwnerUserID       *string `json:"ownerUserID"`
	Ex                *string `json:"ex"`
	Status            *int32  `json:"status"`
	CreatorUserID     *string `json:"creatorUserID"`
	GroupType         *int32  `json:"groupType"`
	NeedVerification  *int32  `json:"needVerification"`
	LookMemberInfo    *int32  `json:"lookMemberInfo"`
	ApplyMemberFriend *int32  `json:"applyMemberFriend"`
}

type CallbackBeforeMemberJoinGroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	GroupID         string `json:"groupID"`
	UserID          string `json:"userID"`
	Ex              string `json:"ex"`
	GroupEx         string `json:"groupEx"`
}

type CallbackBeforeMemberJoinGroupResp struct {
	CommonCallbackResp
	Nickname    *string `json:"nickname"`
	FaceURL     *string `json:"faceURL"`
	RoleLevel   *int32  `json:"roleLevel"`
	MuteEndTime *int64  `json:"muteEndTime"`
	Ex          *string `json:"ex"`
}

type CallbackBeforeSetGroupMemberInfoReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string  `json:"operationID"`
	GroupID         string  `json:"groupID"`
	UserID          string  `json:"userID"`
	Nickname        *string `json:"nickName"`
	FaceURL         *string `json:"faceURL"`
	RoleLevel       *int32  `json:"roleLevel"`
	Ex              *string `json:"ex"`
}

type CallbackBeforeSetGroupMemberInfoResp struct {
	CommonCallbackResp
	Ex        *string `json:"ex"`
	Nickname  *string `json:"nickName"`
	FaceURL   *string `json:"faceURL"`
	RoleLevel *int32  `json:"roleLevel"`
}
