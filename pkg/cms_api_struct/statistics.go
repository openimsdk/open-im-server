package cms_api_struct

type StatisticsRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// 单聊
type MessageStatisticsResponse struct {
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

// 用户统计
type UserStatisticsResponse struct {
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
		TotalUserNum string `json:"total_user_num"`
	} `json:"total_user_num_list"`
}

// 群聊统计
type GroupMessageStatisticsResponse struct {
	IncreaseGroupNum     int `json:"increase_group_num"`
	TotalGroupNum        int `json:"total_group_num"`
	IncreaseGroupNumList []struct {
		Date             string `json:"date"`
		IncreaseGroupNum int    `json:"increase_group_num"`
	} `json:"increase_group_num_list"`
	TotalGroupNumList []struct {
		Date          string `json:"date"`
		TotalGroupNum string `json:"total_group_num"`
	} `json:"total_group_num_list"`
}

type ActiveUserStatisticsResponse struct {
	ActiveUserList []struct {
		NickName   string `json:"nick_name"`
		Id         int    `json:"id"`
		MessageNum int    `json:"message_num"`
	} `json:"active_user_list"`
}

type ActiveGroupStatisticsResponse struct {
	ActiveGroupList []struct {
		GroupNickName string `json:"group_nick_name"`
		GroupId       int    `json:"group_id"`
		MessageNum    int    `json:"message_num"`
	} `json:"active_group_list"`
}
