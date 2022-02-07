package cms_api_struct

type GetStatisticsRequest struct {
	From string `form:"from" binding:"required"`
	To   string `form:"to" binding:"required"`
}

type GetMessageStatisticsRequest struct {
	GetStatisticsRequest
}

type GetMessageStatisticsResponse struct {
	PrivateMessageNum     int `json:"private_message_num"`
	GroupMessageNum       int `json:"group_message_num"`
	PrivateMessageNumList []struct {
		Date       string `json:"date"`
		MessageNum int    `json:"message_num"`
	} `json:"private_message_num_list"`
	GroupMessageNumList []struct {
		Date       string `json:"date"`
		MessageNum int    `json:"message_num"`
	} `json:"group_message_num_list"`
}

type GetUserStatisticsRequest struct {
	GetStatisticsRequest
}

type GetUserStatisticsResponse struct {
	IncreaseUserNum     int `json:"increase_user_num"`
	ActiveUserNum       int `json:"active_user_num"`
	TotalUserNum        int `json:"total_user_num"`
	IncreaseUserNumList []struct {
		Date            string `json:"date"`
		IncreaseUserNum int    `json:"increase_user_num"`
	} `json:"increase_user_num_list"`
	ActiveUserNumList []struct {
		Date          string `json:"date"`
		ActiveUserNum int    `json:"active_user_num"`
	} `json:"active_user_num_list"`
	TotalUserNumList []struct {
		Date         string `json:"date"`
		TotalUserNum int    `json:"total_user_num"`
	} `json:"total_user_num_list"`
}

type GetGroupStatisticsRequest struct {
	GetStatisticsRequest
}

// 群聊统计
type GetGroupStatisticsResponse struct {
	IncreaseGroupNum     int `json:"increase_group_num"`
	TotalGroupNum        int `json:"total_group_num"`
	IncreaseGroupNumList []struct {
		Date             string `json:"date"`
		IncreaseGroupNum int    `json:"increase_group_num"`
	} `json:"increase_group_num_list"`
	TotalGroupNumList []struct {
		Date          string `json:"date"`
		TotalGroupNum int    `json:"total_group_num"`
	} `json:"total_group_num_list"`
}

type GetActiveUserRequest struct {
	GetStatisticsRequest
	// RequestPagination
}

type GetActiveUserResponse struct {
	ActiveUserList []struct {
		NickName   string `json:"nick_name"`
		UserId     string `json:"user_id"`
		MessageNum int    `json:"message_num"`
	} `json:"active_user_list"`
}

type GetActiveGroupRequest struct {
	GetStatisticsRequest
	// RequestPagination
}

type GetActiveGroupResponse struct {
	ActiveGroupList []struct {
		GroupName  string `json:"group_name"`
		GroupId    string `json:"group_id"`
		MessageNum int    `json:"message_num"`
	} `json:"active_group_list"`
}
