package apiresp

type apiResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data"`
}

func apiSuccess(data any) *apiResponse {
	return &apiResponse{
		Data: data,
	}
}

func apiError(err error) *apiResponse {
	return &apiResponse{ErrCode: 10000, ErrMsg: err.Error()}
}
