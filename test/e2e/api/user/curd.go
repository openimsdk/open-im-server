package user

import (
	gettoken "github.com/openimsdk/open-im-server/v3/test/e2e/api/token"
)

// UserInfoRequest represents a request to get or update user information
type UserInfoRequest struct {
	UserIDs  []string       `json:"userIDs,omitempty"`
	UserInfo *gettoken.User `json:"userInfo,omitempty"`
}

// GetUsersOnlineStatusRequest represents a request to get users' online status
type GetUsersOnlineStatusRequest struct {
	UserIDs []string `json:"userIDs"`
}

// GetUsersInfo retrieves detailed information for a list of user IDs
func GetUsersInfo(token string, userIDs []string) error {
	requestBody := UserInfoRequest{
		UserIDs: userIDs,
	}
	return sendPostRequestWithToken("http://your-api-host:port/user/get_users_info", token, requestBody)
}

// UpdateUserInfo updates the information for a user
func UpdateUserInfo(token, userID, nickname, faceURL string) error {
	requestBody := UserInfoRequest{
		UserInfo: &gettoken.User{
			UserID:   userID,
			Nickname: nickname,
			FaceURL:  faceURL,
		},
	}
	return sendPostRequestWithToken("http://your-api-host:port/user/update_user_info", token, requestBody)
}

// GetUsersOnlineStatus retrieves the online status for a list of user IDs
func GetUsersOnlineStatus(token string, userIDs []string) error {
	requestBody := GetUsersOnlineStatusRequest{
		UserIDs: userIDs,
	}
	return sendPostRequestWithToken("http://your-api-host:port/user/get_users_online_status", token, requestBody)
}
