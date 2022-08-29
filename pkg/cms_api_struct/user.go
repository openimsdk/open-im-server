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

type GetUserRequest struct {
	UserId string `form:"user_id" binding:"required"`
}

type GetUserResponse struct {
	UserResponse
}

type GetUsersRequest struct {
	RequestPagination
}

type GetUsersResponse struct {
	Users []*UserResponse `json:"users"`
	ResponsePagination
	UserNums int32 `json:"user_nums"`
}

type GetUsersByNameRequest struct {
	UserName string `form:"user_name" binding:"required"`
	RequestPagination
}

type GetUsersByNameResponse struct {
	Users []*UserResponse `json:"users"`
	ResponsePagination
	UserNums int32 `json:"user_nums"`
}

type ResignUserRequest struct {
	UserId string `json:"user_id"`
}

type ResignUserResponse struct {
}

type AlterUserRequest struct {
	UserId      string `json:"user_id" binding:"required"`
	Nickname    string `json:"nickname"`
	PhoneNumber string `json:"phone_number" validate:"len=11"`
	Email       string `json:"email"`
	Birth       string `json:"birth"`
	Gender      string `json:"gender"`
	Photo       string `json:"photo"`
}

type AlterUserResponse struct {
}

type AddUserRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	UserId      string `json:"user_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email"`
	Birth       string `json:"birth"`
	Gender      string `json:"gender"`
	Photo       string `json:"photo"`
}

type AddUserResponse struct {
}

type BlockUser struct {
	UserResponse
	BeginDisableTime string `json:"begin_disable_time"`
	EndDisableTime   string `json:"end_disable_time"`
}

type BlockUserRequest struct {
	UserId         string `json:"user_id" binding:"required"`
	EndDisableTime string `json:"end_disable_time" binding:"required"`
}

type BlockUserResponse struct {
}

type UnblockUserRequest struct {
	UserId string `json:"user_id" binding:"required"`
}

type UnBlockUserResponse struct {
}

type GetBlockUsersRequest struct {
	RequestPagination
}

type GetBlockUsersResponse struct {
	BlockUsers []BlockUser `json:"block_users"`
	ResponsePagination
	UserNums int32 `json:"user_nums"`
}

type GetBlockUserRequest struct {
	UserId string `form:"user_id" binding:"required"`
}

type GetBlockUserResponse struct {
	BlockUser
}

type DeleteUserRequest struct {
	UserId string `json:"user_id" binding:"required"`
}

type DeleteUserResponse struct {
}
