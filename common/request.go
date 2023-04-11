package common

type PageRequest struct {
	PageNum  int `json:"page_num" query:"page_num"`
	PageSize int `json:"page_size" query:"page_size"`
}
