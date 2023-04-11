package api

type post struct {
	ID         string `json:"id"`
	PosterID   string `json:"poster_id"`
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
}

type WallRequest struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type WallResponse struct {
	Posts []post `json:"posts"`
}
