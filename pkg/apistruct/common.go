package apistruct

type RequestPagination struct {
	PageNumber int `json:"pageNumber"  binding:"required"`
	ShowNumber int `json:"showNumber"  binding:"required"`
}
