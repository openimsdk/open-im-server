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

package user

import (
	"fmt"

	gettoken "github.com/openimsdk/open-im-server/v3/test/e2e/api/token"
	"github.com/openimsdk/open-im-server/v3/test/e2e/framework/config"
)

// UserInfoRequest represents a request to get or update user information.
type UserInfoRequest struct {
	UserIDs  []string       `json:"userIDs,omitempty"`
	UserInfo *gettoken.User `json:"userInfo,omitempty"`
}

// GetUsersOnlineStatusRequest represents a request to get users' online status.
type GetUsersOnlineStatusRequest struct {
	UserIDs []string `json:"userIDs"`
}

// GetUsersInfo retrieves detailed information for a list of user IDs.
func GetUsersInfo(token string, userIDs []string) error {

	url := fmt.Sprintf("http://%s:%s/user/get_users_info", config.LoadConfig().APIHost, config.LoadConfig().APIPort)

	requestBody := UserInfoRequest{
		UserIDs: userIDs,
	}
	return sendPostRequestWithToken(url, token, requestBody)
}

// UpdateUserInfo updates the information for a user.
func UpdateUserInfo(token, userID, nickname, faceURL string) error {

	url := fmt.Sprintf("http://%s:%s/user/update_user_info", config.LoadConfig().APIHost, config.LoadConfig().APIPort)

	requestBody := UserInfoRequest{
		UserInfo: &gettoken.User{
			UserID:   userID,
			Nickname: nickname,
			FaceURL:  faceURL,
		},
	}
	return sendPostRequestWithToken(url, token, requestBody)
}

// GetUsersOnlineStatus retrieves the online status for a list of user IDs.
func GetUsersOnlineStatus(token string, userIDs []string) error {

	url := fmt.Sprintf("http://%s:%s/user/get_users_online_status", config.LoadConfig().APIHost, config.LoadConfig().APIPort)

	requestBody := GetUsersOnlineStatusRequest{
		UserIDs: userIDs,
	}

	return sendPostRequestWithToken(url, token, requestBody)
}
