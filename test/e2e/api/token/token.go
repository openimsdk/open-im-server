package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// API endpoints and other constants
const (
	APIHost         = "http://127.0.0.1:10002"
	UserTokenURL    = APIHost + "/auth/user_token"
	UserRegisterURL = APIHost + "/user/user_register"
	SecretKey       = "openIM123"
	OperationID     = "1646445464564"
)

// UserTokenRequest represents a request to get a user token
type UserTokenRequest struct {
	Secret     string `json:"secret"`
	PlatformID int    `json:"platformID"`
	UserID     string `json:"userID"`
}

// UserTokenResponse represents a response containing a user token
type UserTokenResponse struct {
	Token   string `json:"token"`
	ErrCode int    `json:"errCode"`
}

// User represents user data for registration
type User struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
	FaceURL  string `json:"faceURL"`
}

// UserRegisterRequest represents a request to register a user
type UserRegisterRequest struct {
	Secret string `json:"secret"`
	Users  []User `json:"users"`
}

func main() {
	// Example usage of functions
	token, err := GetUserToken("openIM123456")
	if err != nil {
		log.Fatalf("Error getting user token: %v", err)
	}
	fmt.Println("Token:", token)

	err = RegisterUser(token, "testUserID", "TestNickname", "https://example.com/image.jpg")
	if err != nil {
		log.Fatalf("Error registering user: %v", err)
	}
}

// GetUserToken requests a user token from the API
func GetUserToken(userID string) (string, error) {
	reqBody := UserTokenRequest{
		Secret:     SecretKey,
		PlatformID: 1,
		UserID:     userID,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(UserTokenURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp UserTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("error in token response: %v", tokenResp.ErrCode)
	}

	return tokenResp.Token, nil
}

// RegisterUser registers a new user using the API
func RegisterUser(token, userID, nickname, faceURL string) error {
	user := User{
		UserID:   userID,
		Nickname: nickname,
		FaceURL:  faceURL,
	}
	reqBody := UserRegisterRequest{
		Secret: SecretKey,
		Users:  []User{user},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", UserRegisterURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("operationID", OperationID)
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
		return fmt.Errorf("error in user registration response: %v", respData)
	}

	return nil
}
