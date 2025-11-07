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
