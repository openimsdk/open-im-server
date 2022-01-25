package cms_api_struct

type UserResponse struct {
	ProfilePhoto string `json:"profile_photo"`
	Nickname     string `json:"nick_name"`
	UserId       string `json:"user_id"`
	CreateTime   string `json:"create_time"`
	IsBlock bool `json:"is_block"`
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
	UserNum int `json:"user_num"`
	ResponsePagination
}

type ResignUserRequest struct {
	UserId string `json:"user_id"`
}

type ResignUserResponse struct {
}

type AlterUserRequest struct {
	UserId string `json:"user_id"`
}

type AlterUserResponse struct {
}

type AddUserRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type AddUserResponse struct {
}

type BlockUserRequest struct {
	UserId string `json:"user_id" binding:"required"`
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
	BlockUsers []UserResponse `json:"block_users"`
	BlockUserNum int `json:"block_user_num"`
	ResponsePagination
}
