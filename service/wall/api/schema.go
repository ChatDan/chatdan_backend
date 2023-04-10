package api

type post struct {
	id         string `json:"id"`
	posterID   string `json:"poster_id"`
	content    string `json:"content"`
	visibility string `json:"visibility"`
}

type WallRequest struct {
	pageNum  int `json:"page_num"`
	pageSize int `json:"page_size"`
}

type WallResponse struct {
	posts []post `json:"posts"`
}
