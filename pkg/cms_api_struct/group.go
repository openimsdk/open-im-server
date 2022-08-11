package cms_api_struct

import (
	server_api_params "Open_IM/pkg/proto/sdk_ws"
)

type GroupResponse struct {
	GroupOwnerName string `json:"GroupOwnerName"`
	GroupOwnerID   string `json:"GroupOwnerID"`
	//*server_api_params.GroupInfo
	GroupID                string `json:"groupID"`
	GroupName              string `json:"groupName"`
	Notification           string `json:"notification"`
	Introduction           string `json:"introduction"`
	FaceURL                string `json:"faceURL"`
	OwnerUserID            string `json:"ownerUserID"`
	CreateTime             uint32 `json:"createTime"`
	MemberCount            uint32 `json:"memberCount"`
	Ex                     string `json:"ex"`
	Status                 int32  `json:"status"`
	CreatorUserID          string `json:"creatorUserID"`
	GroupType              int32  `json:"groupType"`
	NeedVerification       int32  `json:"needVerification"`
	LookMemberInfo         int32  `json:"lookMemberInfo"`
	ApplyMemberFriend      int32  `json:"applyMemberFriend"`
	NotificationUpdateTime uint32 `json:"notificationUpdateTime"`
	NotificationUserID     string `json:"notificationUserID"`
}

type GetGroupByIDRequest struct {
	GroupID string `form:"groupID" binding:"required"`
}

type GetGroupByIDResponse struct {
	GroupResponse
}

type GetGroupRequest struct {
	GroupName string `form:"groupName" binding:"required"`
	RequestPagination
}

type GetGroupResponse struct {
	Groups    []GroupResponse `json:"groups"`
	GroupNums int             `json:"groupNums"`
	ResponsePagination
}

type GetGroupsRequest struct {
	RequestPagination
}

type GetGroupsResponse struct {
	Groups    []GroupResponse `json:"groups"`
	GroupNums int             `json:"groupNums"`
	ResponsePagination
}

type CreateGroupRequest struct {
	GroupName     string   `json:"groupName" binding:"required"`
	GroupMasterId string   `json:"groupOwnerID" binding:"required"`
	GroupMembers  []string `json:"groupMembers" binding:"required"`
}

type CreateGroupResponse struct {
}

type SetGroupMasterRequest struct {
	GroupId string `json:"groupID" binding:"required"`
	UserId  string `json:"userID" binding:"required"`
}

type SetGroupMasterResponse struct {
}

type SetGroupMemberRequest struct {
	GroupId string `json:"groupID" binding:"required"`
	UserId  string `json:"userID" binding:"required"`
}

type SetGroupMemberRespones struct {
}

type BanGroupChatRequest struct {
	GroupId string `json:"groupID" binding:"required"`
}

type BanGroupChatResponse struct {
}

type BanPrivateChatRequest struct {
	GroupId string `json:"groupID" binding:"required"`
}

type BanPrivateChatResponse struct {
}

type DeleteGroupRequest struct {
	GroupId string `json:"groupID" binding:"required"`
}

type DeleteGroupResponse struct {
}

type GetGroupMembersRequest struct {
	GroupID  string `form:"groupID" binding:"required"`
	UserName string `form:"userName"`
	RequestPagination
}

type GetGroupMembersResponse struct {
	GroupMembers []server_api_params.GroupMemberFullInfo `json:"groupMembers"`
	ResponsePagination
	MemberNums int `json:"memberNums"`
}

type GroupMemberRequest struct {
	GroupId string   `json:"groupID" binding:"required"`
	Members []string `json:"members" binding:"required"`
}

type GroupMemberOperateResponse struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
}

type AddGroupMembersRequest struct {
	GroupMemberRequest
}

type AddGroupMembersResponse struct {
	GroupMemberOperateResponse
}

type RemoveGroupMembersRequest struct {
	GroupMemberRequest
}

type RemoveGroupMembersResponse struct {
	GroupMemberOperateResponse
}

type AlterGroupInfoRequest struct {
	GroupID      string `json:"groupID"`
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	ProfilePhoto string `json:"profilePhoto"`
	GroupType    int    `json:"groupType"`
}

type AlterGroupInfoResponse struct {
}
