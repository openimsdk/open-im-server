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

package apistruct

type UserRegisterReq struct {
	Secret   string `json:"secret"      binding:"required,max=32"`
	Platform int32  `json:"platform"    binding:"required,min=1,max=12"`
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenInfo struct {
	UserID      string `json:"userID"`
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}
type UserRegisterResp struct {
	UserToken UserTokenInfo `json:"data"`
}

type UserTokenReq struct {
	Secret      string `json:"secret"      binding:"required,max=32"`
	Platform    int32  `json:"platform"    binding:"required,min=1,max=12"`
	UserID      string `json:"userID"      binding:"required,min=1,max=64"`
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenResp struct {
	UserToken UserTokenInfo `json:"data"`
}

type ForceLogoutReq struct {
	Platform    int32  `json:"platform"    binding:"required,min=1,max=12"`
	FromUserID  string `json:"fromUserID"  binding:"required,min=1,max=64"`
	OperationID string `json:"operationID" binding:"required"`
}

type ForceLogoutResp struct{}

type ParseTokenReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

//type ParseTokenResp struct {
//
//	ExpireTime int64 `json:"expireTime" binding:"required"`
//}

type ExpireTime struct {
	ExpireTimeSeconds uint32 `json:"expireTimeSeconds"`
}

type ParseTokenResp struct {
	Data       map[string]interface{} `json:"data" swaggerignore:"true"`
	ExpireTime ExpireTime             `json:"-"`
}
