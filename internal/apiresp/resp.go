package apiresp

type ApiResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data"`
}

func Success(data any) *ApiResponse {
	return &ApiResponse{
		Data: data,
	}
}

func Error(err error) *ApiResponse {
	return &ApiResponse{}
}
