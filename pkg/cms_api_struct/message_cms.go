package cms_api_struct

import (
	pbCommon "Open_IM/pkg/proto/sdk_ws"
)

type GetChatLogsReq struct {
	SessionType int    `json:"sessionType"`
	ContentType int    `json:"contentType"`
	Content     string `json:"content"`
	SendID      string `json:"userID"`
	RecvID      string `json:"recvID"`
	GroupID     string `json:"groupID"`
	SendTime    string `json:"sendTime"`
	RequestPagination
	OperationID string `json:"operationID"`
}

type GetChatLogsResp struct {
	ChatLogs    []*pbCommon.MsgData `json:"chatLogs"`
	ChatLogsNum int                 `json:"logNums"`
	ResponsePagination
}
