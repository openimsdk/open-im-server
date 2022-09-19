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
