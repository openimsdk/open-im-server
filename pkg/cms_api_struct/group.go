package cms_api_struct

type GroupResponse struct {
	GroupName        string `json:"group_name"`
	GroupID          string `json:"group_id"`
	GroupMasterName  string `json:"group_master_name"`
	GroupMasterId    string `json:"group_master_id"`
	CreateTime       string `json:"create_time"`
	IsBanChat        bool   `json:"is_ban_chat"`
	IsBanPrivateChat bool   `json:"is_ban_private_chat"`
	ProfilePhoto string `json:"profile_photo"`
}

type GetGroupByIdRequest struct {
	GroupId string `form:"group_id" binding:"required"`
}

type GetGroupByIdResponse struct {
	GroupResponse
}

type GetGroupRequest struct {
	GroupName string `form:"group_name" binding:"required"`
	RequestPagination
}

type GetGroupResponse struct {
	Groups    []GroupResponse `json:"groups"`
	GroupNums int             `json:"group_nums"`
	ResponsePagination
}

type GetGroupsRequest struct {
	RequestPagination
}

type GetGroupsResponse struct {
	Groups    []GroupResponse `json:"groups"`
	GroupNums int             `json:"group_nums"`
	ResponsePagination
}

type CreateGroupRequest struct {
	GroupName     string   `json:"group_name" binding:"required"`
	GroupMasterId string   `json:"group_master_id" binding:"required"`
	GroupMembers  []string `json:"group_members" binding:"required"`
}

type CreateGroupResponse struct {
}

type SetGroupMasterRequest struct {
	GroupId string `json:"group_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

type SetGroupMasterResponse struct {
}

type SetGroupMemberRequest struct {
	GroupId string `json:"group_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

type SetGroupMemberRespones struct {

}

type BanGroupChatRequest struct {
	GroupId string `json:"group_id" binding:"required"`
}

type BanGroupChatResponse struct {
}

type BanPrivateChatRequest struct {
	GroupId string `json:"group_id" binding:"required"`
}

type BanPrivateChatResponse struct {
}

type DeleteGroupRequest struct {
	GroupId string `json:"group_id" binding:"required"`
}

type DeleteGroupResponse struct {
}

type GetGroupMembersRequest struct {
	GroupId string `form:"group_id" binding:"required"`
	UserName string `form:"user_name"`
	RequestPagination
}

type GroupMemberResponse struct {
	MemberPosition int    `json:"member_position"`
	MemberNickName string `json:"member_nick_name"`
	MemberId       string    `json:"member_id"`
	JoinTime       string `json:"join_time"`
}

type GetGroupMembersResponse struct {
	GroupMembers []GroupMemberResponse    `json:"group_members"`
	ResponsePagination
	MemberNums int `json:"member_nums"`
}

type GroupMemberRequest struct {
	GroupId string `json:"group_id" binding:"required"`
	Members []string `json:"members" binding:"required"`
}

type GroupMemberOperateResponse struct {
	Success []string `json:"success"`
	Failed []string `json:"failed"`
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

type RemoveGroupMembersResponse struct{
	GroupMemberOperateResponse
}

type AlterGroupInfoRequest struct {
	GroupID       string `json:"group_id"`
	GroupName     string `json:"group_name"`
	Notification  string `json:"notification"`
		Introduction  string `json:"introduction"`
	ProfilePhoto  string `json:"profile_photo"`
	GroupType     int `json:"group_type"`
}

type AlterGroupInfoResponse struct {

}
