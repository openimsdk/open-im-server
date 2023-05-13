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

package cms_api_struct

import (
	"Open_IM/pkg/base_info"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
)

type AdminLoginRequest struct {
	AdminName   string `json:"adminID" binding:"required"`
	Secret      string `json:"secret" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type AdminLoginResponse struct {
	Token    string `json:"token"`
	UserName string `json:"userName"`
	FaceURL  string `json:"faceURL"`
}

type GetUserTokenRequest struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
	PlatFormID  int32  `json:"platformID" binding:"required"`
}

type GetUserTokenResponse struct {
	Token   string `json:"token"`
	ExpTime int64  `json:"expTime"`
}

type AddUserRegisterAddFriendIDListRequest struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}

type AddUserRegisterAddFriendIDListResponse struct {
}

type ReduceUserRegisterAddFriendIDListRequest struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
	Operation   int32    `json:"operation"`
}

type ReduceUserRegisterAddFriendIDListResponse struct {
}

type GetUserRegisterAddFriendIDListRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	base_info.RequestPagination
}

type GetUserRegisterAddFriendIDListResponse struct {
	Users []*server_api_params.UserInfo `json:"users"`
	base_info.ResponsePagination
}
