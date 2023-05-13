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

type GetStatisticsRequest struct {
	From        string `json:"from" binding:"required"`
	To          string `json:"to" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type GetMessageStatisticsRequest struct {
	GetStatisticsRequest
}

type GetMessageStatisticsResponse struct {
	PrivateMessageNum     int `json:"privateMessageNum"`
	GroupMessageNum       int `json:"groupMessageNum"`
	PrivateMessageNumList []struct {
		Date       string `json:"date"`
		MessageNum int    `json:"messageNum"`
	} `json:"privateMessageNumList"`
	GroupMessageNumList []struct {
		Date       string `json:"date"`
		MessageNum int    `json:"messageNum"`
	} `json:"groupMessageNumList"`
}

type GetUserStatisticsRequest struct {
	GetStatisticsRequest
}

type GetUserStatisticsResponse struct {
	IncreaseUserNum     int `json:"increaseUserNum"`
	ActiveUserNum       int `json:"activeUserNum"`
	TotalUserNum        int `json:"totalUserNum"`
	IncreaseUserNumList []struct {
		Date            string `json:"date"`
		IncreaseUserNum int    `json:"increaseUserNum"`
	} `json:"increaseUserNumList"`
	ActiveUserNumList []struct {
		Date          string `json:"date"`
		ActiveUserNum int    `json:"activeUserNum"`
	} `json:"activeUserNumList"`
	TotalUserNumList []struct {
		Date         string `json:"date"`
		TotalUserNum int    `json:"totalUserNum"`
	} `json:"totalUserNumList"`
}

type GetGroupStatisticsRequest struct {
	GetStatisticsRequest
}

// 群聊统计
type GetGroupStatisticsResponse struct {
	IncreaseGroupNum     int `json:"increaseGroupNum"`
	TotalGroupNum        int `json:"totalGroupNum"`
	IncreaseGroupNumList []struct {
		Date             string `json:"date"`
		IncreaseGroupNum int    `json:"increaseGroupNum"`
	} `json:"increaseGroupNumList"`
	TotalGroupNumList []struct {
		Date          string `json:"date"`
		TotalGroupNum int    `json:"totalGroupNum"`
	} `json:"totalGroupNumList"`
}

type GetActiveUserRequest struct {
	GetStatisticsRequest
	// RequestPagination
}

type GetActiveUserResponse struct {
	ActiveUserList []struct {
		NickName   string `json:"nickName"`
		UserId     string `json:"userID"`
		MessageNum int    `json:"messageNum"`
	} `json:"activeUserList"`
}

type GetActiveGroupRequest struct {
	GetStatisticsRequest
	// RequestPagination
}

type GetActiveGroupResponse struct {
	ActiveGroupList []struct {
		GroupName  string `json:"groupName"`
		GroupId    string `json:"groupID"`
		MessageNum int    `json:"messageNum"`
	} `json:"activeGroupList"`
}
