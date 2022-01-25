package cms_api_struct

type RequestPagination struct {
	PageNumber int `form:"page_number" binding:"required"`
	ShowNumber int `form:"show_number" binding:"required"`
}

type ResponsePagination struct {
	CurrentPage int `json:"current_number" binding:"required"`
	ShowNumber int `json:"show_number" binding:"required"`
}