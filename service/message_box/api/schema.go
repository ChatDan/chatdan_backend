package api

type BoxPostRequest struct {
	title string `json:"title"`
}

type BoxPostResponse struct {
	id      string `json:"id"`
	ownerID string `json:"owner_id"`
	title   string `json:"title"`
}

type BoxesGetRequest struct {
	pageNum  int    `json:"page_num"`
	pageSize int    `json:"page_size"`
	title    string `json:"title"`
	owner    int    `json:"owner"`
}

type BoxesGetResponse struct {
	id      string `json:"id"`
	ownerID string `json:"owner_id"`
	title   string `json:"title"`
}

type BoxGetRequest struct {
	id string `json:"id"`
}

type BoxGetResponse struct {
	id      string    `json:"id"`
	ownerID string    `json:"owner_id"`
	title   string    `json:"title"`
	posts   *[]string `json:"posts"`
}

type BoxModifyRequest struct {
	title string `json:"title"`
}

type BoxPutRequest struct {
	id   string           `json:"id"`
	body BoxModifyRequest `json:"body"`
}

type BoxPutResponse struct {
	id      string `json:"id"`
	ownerID string `json:"owner_id"`
	title   string `json:"title"`
}

type BoxDeleteRequest struct {
	id string `json:"id"`
}

type BoxDeleteResponse struct {
	message string `json:"message"`
}
