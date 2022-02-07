package cms_api_struct

type BroadcastRequest struct {
	Message string `json:"message"`
}

type BroadcastResponse struct {
}

type CommonMessage struct {
	SessionType    int    `json:"session_type"`
	ContentType    int    `json:"content_type"`
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
		Date string `json:"date"`
	} `json:"massage_list"`
}
