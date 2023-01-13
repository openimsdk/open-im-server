package callback_struct

import (
	"Open_IM/pkg/common/constant"
	"fmt"
)

type CommonCallbackReq struct {
	SendID           string   `json:"sendID"`
	CallbackCommand  string   `json:"callbackCommand"`
	ServerMsgID      string   `json:"serverMsgID"`
	ClientMsgID      string   `json:"clientMsgID"`
	OperationID      string   `json:"operationID"`
	SenderPlatformID int32    `json:"senderPlatformID"`
	SenderNickname   string   `json:"senderNickname"`
	SessionType      int32    `json:"sessionType"`
	MsgFrom          int32    `json:"msgFrom"`
	ContentType      int32    `json:"contentType"`
	Status           int32    `json:"status"`
	CreateTime       int64    `json:"createTime"`
	Content          string   `json:"content"`
	Seq              uint32   `json:"seq"`
	AtUserIDList     []string `json:"atUserList"`
	SenderFaceURL    string   `json:"faceURL"`
	Ex               string   `json:"ex"`
}

type CallbackResp interface {
	Parse() (err error)
}

type CommonCallbackResp struct {
	ActionCode  int    `json:"actionCode"`
	ErrCode     int32  `json:"errCode"`
	ErrMsg      string `json:"errMsg"`
	OperationID string `json:"operationID"`
}

func (c *CommonCallbackResp) Parse() (err error) {
	if c.ActionCode != constant.NoError || c.ErrCode != constant.NoError {
		newErr := constant.ErrCallback
		newErr.ErrCode = c.ErrCode
		newErr.DetailErrMsg = fmt.Sprintf("callback response error actionCode is %d, errCode is %d, errMsg is %s", c.ActionCode, c.ErrCode, c.ErrMsg)
		err = newErr
		return
	}
	return
}

type UserStatusBaseCallback struct {
	CallbackCommand string `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	PlatformID      int32  `json:"platformID"`
	Platform        string `json:"platform"`
}

type UserStatusCallbackReq struct {
	UserStatusBaseCallback
	UserID string `json:"userID"`
}

type UserStatusBatchCallbackReq struct {
	UserStatusBaseCallback
	UserIDList []string `json:"userIDList"`
}
