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
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/wrapperspb"
)

type CallbackBeforeUpdateUserInfoReq struct {
	CallbackCommand `json:"callbackCommand"`
	UserID          string  `json:"userID"`
	Nickname        *string `json:"nickName"`
	FaceURL         *string `json:"faceURL"`
	Ex              *string `json:"ex"`
}

type CallbackBeforeUpdateUserInfoResp struct {
	CommonCallbackResp
	Nickname *string `json:"nickName"`
	FaceURL  *string `json:"faceURL"`
	Ex       *string `json:"ex"`
}

type CallbackAfterUpdateUserInfoReq struct {
	CallbackCommand `json:"callbackCommand"`
	UserID          string `json:"userID"`
	Nickname        string `json:"nickName"`
	FaceURL         string `json:"faceURL"`
	Ex              string `json:"ex"`
}
type CallbackAfterUpdateUserInfoResp struct {
	CommonCallbackResp
}

type CallbackBeforeUpdateUserInfoExReq struct {
	CallbackCommand `json:"callbackCommand"`
	UserID          string                  `json:"userID"`
	Nickname        *wrapperspb.StringValue `json:"nickName"`
	FaceURL         *wrapperspb.StringValue `json:"faceURL"`
	Ex              *wrapperspb.StringValue `json:"ex"`
}
type CallbackBeforeUpdateUserInfoExResp struct {
	CommonCallbackResp
	Nickname *wrapperspb.StringValue `json:"nickName"`
	FaceURL  *wrapperspb.StringValue `json:"faceURL"`
	Ex       *wrapperspb.StringValue `json:"ex"`
}

type CallbackAfterUpdateUserInfoExReq struct {
	CallbackCommand `json:"callbackCommand"`
	UserID          string                  `json:"userID"`
	Nickname        *wrapperspb.StringValue `json:"nickName"`
	FaceURL         *wrapperspb.StringValue `json:"faceURL"`
	Ex              *wrapperspb.StringValue `json:"ex"`
}
type CallbackAfterUpdateUserInfoExResp struct {
	CommonCallbackResp
}

type CallbackBeforeUserRegisterReq struct {
	CallbackCommand `json:"callbackCommand"`
	Users           []*sdkws.UserInfo `json:"users"`
}

type CallbackBeforeUserRegisterResp struct {
	CommonCallbackResp
	Users []*sdkws.UserInfo `json:"users"`
}

type CallbackAfterUserRegisterReq struct {
	CallbackCommand `json:"callbackCommand"`
	Users           []*sdkws.UserInfo `json:"users"`
}

type CallbackAfterUserRegisterResp struct {
	CommonCallbackResp
}
