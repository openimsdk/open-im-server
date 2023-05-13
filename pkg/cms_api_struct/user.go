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

type UserResponse struct {
	FaceURL       string `json:"faceURL"`
	Nickname      string `json:"nickName"`
	UserID        string `json:"userID"`
	CreateTime    string `json:"createTime,omitempty"`
	CreateIp      string `json:"createIp,omitempty"`
	LastLoginTime string `json:"lastLoginTime,omitempty"`
	LastLoginIp   string `json:"lastLoginIP,omitempty"`
	LoginTimes    int32  `json:"loginTimes"`
	LoginLimit    int32  `json:"loginLimit"`
	IsBlock       bool   `json:"isBlock"`
	PhoneNumber   string `json:"phoneNumber"`
	Email         string `json:"email"`
	Birth         string `json:"birth"`
	Gender        int    `json:"gender"`
}

type AddUserRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	UserId      string `json:"userID" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email"`
	Birth       string `json:"birth"`
	Gender      string `json:"gender"`
	FaceURL     string `json:"faceURL"`
}

type AddUserResponse struct {
}

type BlockUser struct {
	UserResponse
	BeginDisableTime string `json:"beginDisableTime"`
	EndDisableTime   string `json:"endDisableTime"`
}

type BlockUserRequest struct {
	OperationID    string `json:"operationID" binding:"required"`
	UserID         string `json:"userID" binding:"required"`
	EndDisableTime string `json:"endDisableTime" binding:"required"`
}

type BlockUserResponse struct {
}

type UnblockUserRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

type UnBlockUserResponse struct {
}

type GetBlockUsersRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	RequestPagination
}

type GetBlockUsersResponse struct {
	BlockUsers []BlockUser `json:"blockUsers"`
	ResponsePagination
	UserNums int32 `json:"userNums"`
}

type GetUserIDByEmailAndPhoneNumberRequest struct {
	OperationID string `json:"operationID" binding:"required"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type GetUserIDByEmailAndPhoneNumberResponse struct {
	UserIDList []string `json:"userIDList"`
}
