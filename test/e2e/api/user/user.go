package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	gettoken "github.com/openimsdk/open-im-server/v3/test/e2e/api/token"
)

// ForceLogoutRequest represents a request to force a user logout
type ForceLogoutRequest struct {
	PlatformID int    `json:"platformID"`
	UserID     string `json:"userID"`
}

// CheckUserAccountRequest represents a request to check a user account
type CheckUserAccountRequest struct {
	CheckUserIDs []string `json:"checkUserIDs"`
}

// GetUsersRequest represents a request to get a list of users
type GetUsersRequest struct {
	Pagination Pagination `json:"pagination"`
}

// Pagination specifies the page number and number of items per page
type Pagination struct {
	PageNumber int `json:"pageNumber"`
	ShowNumber int `json:"showNumber"`
}

// ForceLogout forces a user to log out
func ForceLogout(token, userID string, platformID int) error {
	requestBody := ForceLogoutRequest{
		PlatformID: platformID,
		UserID:     userID,
	}
	return sendPostRequestWithToken("http://your-api-host:port/auth/force_logout", token, requestBody)
}

// CheckUserAccount checks if the user accounts exist
func CheckUserAccount(token string, userIDs []string) error {
	requestBody := CheckUserAccountRequest{
		CheckUserIDs: userIDs,
	}
	return sendPostRequestWithToken("http://your-api-host:port/user/account_check", token, requestBody)
}

// GetUsers retrieves a list of users with pagination
func GetUsers(token string, pageNumber, showNumber int) error {
	requestBody := GetUsersRequest{
		Pagination: Pagination{
			PageNumber: pageNumber,
			ShowNumber: showNumber,
		},
	}
	return sendPostRequestWithToken("http://your-api-host:port/user/get_users", token, requestBody)
}

// sendPostRequestWithToken sends a POST request with a token in the header
func sendPostRequestWithToken(url, token string, body interface{}) error {
	reqBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("operationID", gettoken.OperationID)
	req.Header.Add("token", token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal(respBody, &respData); err != nil {
		return err
	}

	if errCode, ok := respData["errCode"].(float64); ok && errCode != 0 {
		return fmt.Errorf("error in response: %v", respData)
	}

	return nil
}
