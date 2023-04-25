package schemax

type PageRequest struct {
	PageNum  int `json:"page_num" query:"page_num" validate:"required,min=1"`
	PageSize int `json:"page_size" query:"page_size" validate:"required,min=1,max=100"`
}
