package cms_api_struct

type CommonMessage struct {
	ChatType       int    `json:"chat_type"`
	MessageType    int    `json:"message_type"`
	SenderNickName string `json:"sender_nick_name"`
	SenderId       int    `json:"sender_id"`
	SearchContent  string `json:"search_content"`
	WholeContent   string `json:"whole_content"`
}

type SearchMessageByUserResponse struct {
	MessageList []struct {
		CommonMessage
		ReceiverNickName string `json:"receiver_nick_name"`
		ReceiverID       int    `json:"receiver_id"`
		Date             string `json:"date"`
	} `json:"massage_list"`
}

type SearchMessageByGroupResponse struct {
	MessageList []struct {
		CommonMessage
		Date           string `json:"date"`
	} `json:"massage_list"`
}
