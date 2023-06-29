package apistruct

type UserRegisterReq struct {
	Secret   string `json:"secret" binding:"required,max=32"`
	Platform int32  `json:"platform" binding:"required,min=1,max=12"`
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenInfo struct {
	UserID      string `json:"userID"`
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}
type UserRegisterResp struct {
	UserToken UserTokenInfo `json:"data"`
}

type UserTokenReq struct {
	Secret      string `json:"secret" binding:"required,max=32"`
	Platform    int32  `json:"platform" binding:"required,min=1,max=12"`
	UserID      string `json:"userID" binding:"required,min=1,max=64"`
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenResp struct {
	UserToken UserTokenInfo `json:"data"`
}

type ForceLogoutReq struct {
	Platform    int32  `json:"platform" binding:"required,min=1,max=12"`
	FromUserID  string `json:"fromUserID" binding:"required,min=1,max=64"`
	OperationID string `json:"operationID" binding:"required"`
}

type ForceLogoutResp struct {
}

type ParseTokenReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

//type ParseTokenResp struct {
//
//	ExpireTime int64 `json:"expireTime" binding:"required"`
//}

type ExpireTime struct {
	ExpireTimeSeconds uint32 `json:"expireTimeSeconds" `
}

type ParseTokenResp struct {
	Data       map[string]interface{} `json:"data" swaggerignore:"true"`
	ExpireTime ExpireTime             `json:"-"`
}
