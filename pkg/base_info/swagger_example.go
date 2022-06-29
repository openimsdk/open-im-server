package base_info

type Swagger400Resp struct {
	ErrCode int32  `json:"errCode" example:"400"`
	ErrMsg  string `json:"errMsg" example:"err msg"`
}

type Swagger500Resp struct {
	ErrCode int32  `json:"errCode" example:"500"`
	ErrMsg  string `json:"errMsg" example:"err msg"`
}
