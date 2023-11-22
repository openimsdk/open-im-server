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

type CallbackBeforeAddFriendReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
	OperationID     string `json:"operationID"`
}

type CallbackBeforeAddFriendResp struct {
	CommonCallbackResp
}
type CallbackAfterAddFriendReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
	OperationID     string `json:"operationID"`
}

type CallbackAfterAddFriendResp struct {
	CommonCallbackResp
}
type CallbackBeforeAddBlackReq struct {
	CallbackCommand `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID" `
	BlackUserID     string `json:"blackUserID"`
	OperationID     string `json:"operationID"`
}

type CallbackBeforeAddBlackResp struct {
	CommonCallbackResp
}

type CallbackBeforeAddFriendAgreeReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"blackUserID"`
	HandleResult    int32  `json:"HandleResult"'`
	HandleMsg       string `json:"HandleMsg"'`
	OperationID     string `json:"operationID"`
}

type CallbackBeforeAddFriendAgreeResp struct {
	CommonCallbackResp
}
