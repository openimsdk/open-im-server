package call_back_struct

import (
	commonPb "Open_IM/pkg/proto/sdk_ws"
)

type CallbackBeforeCreateGroupReq struct {
	CallbackCommand string `json:"callbackCommand"`
	commonPb.GroupInfo
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
