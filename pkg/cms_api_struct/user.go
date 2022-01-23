package cms_api_struct

type UserResponse struct {
	ProfilePhoto string `json:"profile_photo"`
	Nickname     string `json:"nick_name"`
	UserId       string `json:"user_id"`
	CreateTime   string `json:"create_time"`
}

type GetUserRequest struct {
	UserId string `form:"user_id"`
}

type GetUserResponse struct {
	UserResponse
}

type GetUsersRequest struct {
	RequestPagination
}

type GetUsersResponse struct {
	Users []*UserResponse `json:"users"`
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
}

type AddUserResponse struct {
}

type BlockUserRequest struct {
	UserId string `json:"user_id"`
}

type BlockUserResponse struct {
}

type UnblockUserRequest struct {
	UserId string `json:"user_id"`
}

type UnBlockUserResponse struct {
}

type GetBlockUsersRequest struct {
	RequestPagination
}

type GetBlockUsersResponse struct {
}
