// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package callbackstruct

import (
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
)

const (
	Next = 1
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
	SendTime         int64    `json:"sendTime"`
	CreateTime       int64    `json:"createTime"`
	Content          string   `json:"content"`
	Seq              uint32   `json:"seq"`
	AtUserIDList     []string `json:"atUserList"`
	SenderFaceURL    string   `json:"faceURL"`
	Ex               string   `json:"ex"`
}

func (c *CommonCallbackReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

type CallbackReq interface {
	GetCallbackCommand() string
}

type CallbackResp interface {
	Parse() (err error)
}

type CommonCallbackResp struct {
	ActionCode int32  `json:"actionCode"`
	ErrCode    int32  `json:"errCode"`
	ErrMsg     string `json:"errMsg"`
	ErrDlt     string `json:"errDlt"`
	NextCode   int32  `json:"nextCode"`
}

func (c CommonCallbackResp) Parse() error {
	if c.ActionCode == servererrs.NoError && c.NextCode == Next {
		return errs.NewCodeError(int(c.ErrCode), c.ErrMsg).WithDetail(c.ErrDlt)
	}
	return nil
}

type UserStatusBaseCallback struct {
	CallbackCommand string `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	PlatformID      int    `json:"platformID"`
	Platform        string `json:"platform"`
}

func (c UserStatusBaseCallback) GetCallbackCommand() string {
	return c.CallbackCommand
}

type UserStatusCallbackReq struct {
	UserStatusBaseCallback
	UserID string `json:"userID"`
}

type UserStatusBatchCallbackReq struct {
	UserStatusBaseCallback
	UserIDList []string `json:"userIDList"`
}
