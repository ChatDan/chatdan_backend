package api

// TODO:
type BoxRequest struct {
	title string `json:"title"`
}

type BoxResponse struct {
	id      string    `json:"id"`
	ownerID string    `json:"owner_id"`
	title   string    `json:"title"`
	posts   *[]string `json:"posts"`
}
