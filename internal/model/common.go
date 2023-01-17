package model

type CommonResPrefix struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PaginationInput struct {
	PageNum  int `json:"page_num" v:"required|integer|min:0" description:"必填|min:0"`
	PageSize int `json:"page_size" v:"required|integer|min:5|max:25" description:"必填|min:5|max:25"`
}
