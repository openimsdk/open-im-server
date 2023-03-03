package apiresp

type ApiResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data"`
}

func apiSuccess(data any) *ApiResponse {
	return &ApiResponse{
		Data: data,
	}
}

func apiError(err error) *ApiResponse {
	return &ApiResponse{ErrCode: 10000, ErrMsg: err.Error()}
}
