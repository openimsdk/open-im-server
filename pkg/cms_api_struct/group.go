package cms_api_struct

type GroupResponse struct {
	GroupName        string `json:"group_name"`
	GroupID          string `json:"group_id"`
	GroupMasterName  string `json:"group_master_name"`
	GroupMasterId    string `json:"group_master_id"`
	CreateTime       string `json:"create_time"`
	IsBanChat        bool   `json:"is_ban_chat"`
	IsBanPrivateChat bool   `json:"is_ban_private_chat"`
}

type GetGroupRequest struct {
	GroupName string `form:"group_name"`
}

type GetGroupResponse struct {
	GroupResponse
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
	GroupName     string   `json:"group_name"`
	GroupMasterId string   `json:"group_master_id"`
	GroupMembers  []string `json:"group_members"`
}

type CreateGroupResponse struct {
}

type SetGroupMasterRequest struct {
	GroupId string `json:"group_id"`
	UserId  string `json:"user_id"`
}

type SetGroupMasterResponse struct {
}

type BanGroupChatRequest struct {
	GroupId string `json:"group_id"`
}

type BanGroupChatResponse struct {
}

type BanPrivateChatRequest struct {
	GroupId string `json:"group_id"`
}

type BanPrivateChatResponse struct {
}

type DeleteGroupRequest struct {
	GroupId string `json:"group_id"`
}

type DeleteGroupResponse struct {
}

type GetGroupMemberRequest struct {
	GroupId string `json:"group_id"`
}

type GroupMemberResponse struct {
	MemberPosition int    `json:"member_position"`
	MemberNickName string `json:"member_nick_name"`
	MemberId       int    `json:"member_id"`
	JoinTime       string `json:"join_time"`
}

type GetGroupMemberResponse struct {
	GroupMemberList []GroupMemberResponse `json:"group_member_list"`
	GroupMemberNums int                   `json:"group_member_nums"`
	ResponsePagination
}
