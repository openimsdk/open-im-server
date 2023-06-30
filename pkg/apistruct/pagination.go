package apistruct

type Pagination struct {
	PageNumber int32 `json:"pageNumber" binding:"required"`
	ShowNumber int32 `json:"showNumber" binding:"required"`
}
