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
	common "github.com/OpenIMSDK/protocol/sdkws"

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
	OperationID     string `json:"operationID"`
	CallbackCommand `json:"callbackCommand"`
	*common.GroupInfo
	InitMemberList []*apistruct.GroupAddMemberInfo `json:"initMemberList"`
}

type CallbackAfterCreateGroupResp struct {
	CommonCallbackResp
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

type CallbackAfterGroupMemberExitReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	GroupID         string `json:"groupID"`
	UserID          string `json:"userID"`
	GroupType       *int32 `json:"groupType"`
	ExitType        string `json:"exitType"`
	MuteEndTime     *int64 `json:"muteEndTime"`
}

type CallbackAfterGroupMemberExitResp struct {
	CommonCallbackResp
}

type CallbackAfterUngroupReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string   `json:"operationID"`
	GroupID         string   `json:"groupID"`
	GroupType       *int32   `json:"groupType"`
	OwnerID         string   `json:"ownerID"`
	MemberList      []string `json:"memberList"`
	MuteEndTime     *int64   `json:"muteEndTime"`
}

type CallbackAfterUngroupResp struct {
	CommonCallbackResp
}

type CallbackAfterSetGroupInfoReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	GroupID         string `json:"groupID"`
	GroupType       *int32 `json:"groupType"`
	UserID          string `json:"userID"`
	Name            string `json:"name"`
	Notification    string `json:"notification"`
	GroupUrl        string `json:"groupUrl"`
	MuteEndTime     *int64 `json:"muteEndTime"`
}

type CallbackAfterSetGroupInfoResp struct {
	CommonCallbackResp
}

type CallbackAfterRevokeMsgReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	GroupID         string `json:"groupID"`
	GroupType       *int32 `json:"groupType"`
	UserID          string `json:"userID"`
	Content         string `json:"content"`
	MuteEndTime     *int64 `json:"muteEndTime"`
}

type CallbackAfterRevokeMsgResp struct {
	CommonCallbackResp
}

type CallbackGroupMsgReadReq struct {
	CallbackCommand `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	SendID          string `json:"sendID"`
	ReceiveID       string `json:"receiveID"`
	UnreadMsgNum    int64  `json:"UnreadMsgNum"`
	MuteEndTime     *int64 `json:"muteEndTime"`
}

type CallbackGroupMsgReadResp struct {
	CommonCallbackResp
}
