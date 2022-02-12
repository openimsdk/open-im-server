package cms_api_struct

type AdminLoginRequest struct {
	AdminName string `json:"admin_name" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}