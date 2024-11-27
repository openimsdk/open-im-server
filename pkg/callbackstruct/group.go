// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package callbackstruct

import (
	common "github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/wrapperspb"

	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
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

type CallbackAfterCreateGroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	*common.GroupInfo
	InitMemberList []*apistruct.GroupAddMemberInfo `json:"initMemberList"`
}

type CallbackAfterCreateGroupResp struct {
	CommonCallbackResp
}

type CallbackGroupMember struct {
	UserID string `json:"userID"`
	Ex     string `json:"ex"`
}

type CallbackBeforeMembersJoinGroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	GroupID         string                 `json:"groupID"`
	MembersList     []*CallbackGroupMember `json:"memberList"`
	GroupEx         string                 `json:"groupEx"`
}

type MemberJoinGroupCallBack struct {
	UserID      *string `json:"userID"`
	Nickname    *string `json:"nickname"`
	FaceURL     *string `json:"faceURL"`
	RoleLevel   *int32  `json:"roleLevel"`
	MuteEndTime *int64  `json:"muteEndTime"`
	Ex          *string `json:"ex"`
}

type CallbackBeforeMembersJoinGroupResp struct {
	CommonCallbackResp
	MemberCallbackList []*MemberJoinGroupCallBack `json:"memberCallbackList"`
}

type CallbackBeforeSetGroupMemberInfoReq struct {
	CallbackCommand `json:"callbackCommand"`
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

type CallbackAfterSetGroupMemberInfoReq struct {
	CallbackCommand `json:"callbackCommand"`
	GroupID         string  `json:"groupID"`
	UserID          string  `json:"userID"`
	Nickname        *string `json:"nickName"`
	FaceURL         *string `json:"faceURL"`
	RoleLevel       *int32  `json:"roleLevel"`
	Ex              *string `json:"ex"`
}

type CallbackAfterSetGroupMemberInfoResp struct {
	CommonCallbackResp
}

type CallbackQuitGroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	UserID          string `json:"userID"`
}

type CallbackQuitGroupResp struct {
	CommonCallbackResp
}

type CallbackKillGroupMemberReq struct {
	CallbackCommand `json:"callbackCommand"`
	GroupID         string   `json:"groupID"`
	KickedUserIDs   []string `json:"kickedUserIDs"`
	Reason          string   `json:"reason"`
}

type CallbackKillGroupMemberResp struct {
	CommonCallbackResp
}

type CallbackDisMissGroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	GroupID         string   `json:"groupID"`
	OwnerID         string   `json:"ownerID"`
	GroupType       string   `json:"groupType"`
	MembersID       []string `json:"membersID"`
}

type CallbackDisMissGroupResp struct {
	CommonCallbackResp
}

type CallbackJoinGroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	GroupType       string `json:"groupType"`
	ApplyID         string `json:"applyID"`
	ReqMessage      string `json:"reqMessage"`
	Ex              string `json:"ex"`
}

type CallbackJoinGroupResp struct {
	CommonCallbackResp
}

type CallbackTransferGroupOwnerReq struct {
	CallbackCommand `json:"callbackCommand"`
	GroupID         string `json:"groupID"`
	OldOwnerUserID  string `json:"oldOwnerUserID"`
	NewOwnerUserID  string `json:"newOwnerUserID"`
}

type CallbackTransferGroupOwnerResp struct {
	CommonCallbackResp
}

type CallbackBeforeInviteUserToGroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string   `json:"operationID"`
	GroupID         string   `json:"groupID"`
	Reason          string   `json:"reason"`
	InvitedUserIDs  []string `json:"invitedUserIDs"`
}
type CallbackBeforeInviteUserToGroupResp struct {
	CommonCallbackResp
	RefusedMembersAccount []string `json:"refusedMembersAccount,omitempty"` // Optional field to list members whose invitation is refused.
}

type CallbackAfterJoinGroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	GroupID         string `json:"groupID"`
	ReqMessage      string `json:"reqMessage"`
	JoinSource      int32  `json:"joinSource"`
	InviterUserID   string `json:"inviterUserID"`
}
type CallbackAfterJoinGroupResp struct {
	CommonCallbackResp
}

type CallbackBeforeSetGroupInfoReq struct {
	CallbackCommand   `json:"callbackCommand"`
	OperationID       string `json:"operationID"`
	GroupID           string `json:"groupID"`
	GroupName         string `json:"groupName"`
	Notification      string `json:"notification"`
	Introduction      string `json:"introduction"`
	FaceURL           string `json:"faceURL"`
	Ex                string `json:"ex"`
	NeedVerification  int32  `json:"needVerification"`
	LookMemberInfo    int32  `json:"lookMemberInfo"`
	ApplyMemberFriend int32  `json:"applyMemberFriend"`
}

type CallbackBeforeSetGroupInfoResp struct {
	CommonCallbackResp
	GroupID           string  ` json:"groupID"`
	GroupName         string  `json:"groupName"`
	Notification      string  `json:"notification"`
	Introduction      string  `json:"introduction"`
	FaceURL           string  `json:"faceURL"`
	Ex                *string `json:"ex"`
	NeedVerification  *int32  `json:"needVerification"`
	LookMemberInfo    *int32  `json:"lookMemberInfo"`
	ApplyMemberFriend *int32  `json:"applyMemberFriend"`
}

type CallbackAfterSetGroupInfoReq struct {
	CallbackCommand   `json:"callbackCommand"`
	OperationID       string  `json:"operationID"`
	GroupID           string  `json:"groupID"`
	GroupName         string  `json:"groupName"`
	Notification      string  `json:"notification"`
	Introduction      string  `json:"introduction"`
	FaceURL           string  `json:"faceURL"`
	Ex                *string `json:"ex"`
	NeedVerification  *int32  `json:"needVerification"`
	LookMemberInfo    *int32  `json:"lookMemberInfo"`
	ApplyMemberFriend *int32  `json:"applyMemberFriend"`
}

type CallbackAfterSetGroupInfoResp struct {
	CommonCallbackResp
}

type CallbackBeforeSetGroupInfoExReq struct {
	CallbackCommand   `json:"callbackCommand"`
	OperationID       string                  `json:"operationID"`
	GroupID           string                  `json:"groupID"`
	GroupName         *wrapperspb.StringValue `json:"groupName"`
	Notification      *wrapperspb.StringValue `json:"notification"`
	Introduction      *wrapperspb.StringValue `json:"introduction"`
	FaceURL           *wrapperspb.StringValue `json:"faceURL"`
	Ex                *wrapperspb.StringValue `json:"ex"`
	NeedVerification  *wrapperspb.Int32Value  `json:"needVerification"`
	LookMemberInfo    *wrapperspb.Int32Value  `json:"lookMemberInfo"`
	ApplyMemberFriend *wrapperspb.Int32Value  `json:"applyMemberFriend"`
}

type CallbackBeforeSetGroupInfoExResp struct {
	CommonCallbackResp
	GroupID           string                  `json:"groupID"`
	GroupName         *wrapperspb.StringValue `json:"groupName"`
	Notification      *wrapperspb.StringValue `json:"notification"`
	Introduction      *wrapperspb.StringValue `json:"introduction"`
	FaceURL           *wrapperspb.StringValue `json:"faceURL"`
	Ex                *wrapperspb.StringValue `json:"ex"`
	NeedVerification  *wrapperspb.Int32Value  `json:"needVerification"`
	LookMemberInfo    *wrapperspb.Int32Value  `json:"lookMemberInfo"`
	ApplyMemberFriend *wrapperspb.Int32Value  `json:"applyMemberFriend"`
}

type CallbackAfterSetGroupInfoExReq struct {
	CallbackCommand   `json:"callbackCommand"`
	OperationID       string                  `json:"operationID"`
	GroupID           string                  `json:"groupID"`
	GroupName         *wrapperspb.StringValue `json:"groupName"`
	Notification      *wrapperspb.StringValue `json:"notification"`
	Introduction      *wrapperspb.StringValue `json:"introduction"`
	FaceURL           *wrapperspb.StringValue `json:"faceURL"`
	Ex                *wrapperspb.StringValue `json:"ex"`
	NeedVerification  *wrapperspb.Int32Value  `json:"needVerification"`
	LookMemberInfo    *wrapperspb.Int32Value  `json:"lookMemberInfo"`
	ApplyMemberFriend *wrapperspb.Int32Value  `json:"applyMemberFriend"`
}

type CallbackAfterSetGroupInfoExResp struct {
	CommonCallbackResp
}
