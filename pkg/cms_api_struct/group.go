package cms_api_struct

type SearchGroupsResponse struct {
	GroupList []struct {
		GroupNickName string `json:"group_nick_name"`
		GroupID       int    `json:"group_id"`
		MasterName    string `json:"master_name"`
		MasterId      int    `json:"master_id"`
		CreatTime     string `json:"creat_time"`
	} `json:"group_list"`
}

type SearchGroupMemberResponse struct {
	GroupMemberList []struct {
		MemberPosition int    `json:"member_position"`
		MemberNickName string `json:"member_nick_name"`
		MemberId       int    `json:"member_id"`
		JoinTime       string `json:"join_time"`
	} `json:"group_member_list"`
}
