package cms_api_struct

import (
	pbCommon "Open_IM/pkg/proto/sdk_ws"
)

type GetChatLogsReq struct {
	SessionType int    `json:"sessionType"`
	ContentType int    `json:"contentType"`
	Content     string `json:"content"`
	SendID      string `json:"sendID"`
	RecvID      string `json:"recvID"`
	GroupID     string `json:"groupID"`
	SendTime    string `json:"sendTime"`
	RequestPagination
	OperationID string `json:"operationID"`
}

type ChatLog struct {
	SendID           string                    `json:"sendID,omitempty"`
	RecvID           string                    `json:"recvID,omitempty"`
	GroupID          string                    `json:"groupID,omitempty"`
	ClientMsgID      string                    `json:"clientMsgID,omitempty"`
	ServerMsgID      string                    `json:"serverMsgID,omitempty"`
	SenderPlatformID int32                     `json:"senderPlatformID,omitempty"`
	SenderNickname   string                    `json:"senderNickname,omitempty"`
	SenderFaceURL    string                    `json:"senderFaceURL,omitempty"`
	SessionType      int32                     `json:"sessionType,omitempty"`
	MsgFrom          int32                     `json:"msgFrom,omitempty"`
	ContentType      int32                     `json:"contentType,omitempty"`
	Content          string                    `json:"content,omitempty"`
	Seq              uint32                    `json:"seq,omitempty"`
	SendTime         int64                     `json:"sendTime,omitempty"`
	CreateTime       int64                     `json:"createTime,omitempty"`
	Status           int32                     `json:"status,omitempty"`
	Options          map[string]bool           `json:"options,omitempty"`
	OfflinePushInfo  *pbCommon.OfflinePushInfo `json:"offlinePushInfo,omitempty"`
	AtUserIDList     []string                  `json:"atUserIDList,omitempty"`
	MsgDataList      []byte                    `json:"msgDataList,omitempty"`
	AttachedInfo     string                    `json:"attachedInfo,omitempty"`
	Ex               string                    `json:"ex,omitempty"`
}

type GetChatLogsResp struct {
	ChatLogs    []*ChatLog `json:"chatLogs"`
	ChatLogsNum int        `json:"logNums"`
	ResponsePagination
}
