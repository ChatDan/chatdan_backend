package api

type BoxPostRequest struct {
	Title string `json:"title"`
}

type BoxPostResponse struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Title   string `json:"title"`
}

type BoxesGetRequest struct {
	PageNum  int    `json:"page_num"`
	PageSize int    `json:"page_size"`
	Title    string `json:"title"`
	Owner    int    `json:"owner"`
}

type BoxesGetResponse struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Title   string `json:"title"`
}

type BoxGetRequest struct {
	ID string `json:"id"`
}

type BoxGetResponse struct {
	ID      string    `json:"id"`
	OwnerID string    `json:"owner_id"`
	Title   string    `json:"title"`
	Posts   *[]string `json:"posts"`
}

type BoxModifyRequest struct {
	Title string `json:"title"`
}

type BoxPutRequest struct {
	ID   string           `json:"id"`
	Body BoxModifyRequest `json:"body"`
}

type BoxPutResponse struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Title   string `json:"title"`
}

type BoxDeleteRequest struct {
	ID string `json:"id"`
}

type BoxDeleteResponse struct {
	Message string `json:"message"`
}

type PostPostRequest struct {
	MessageBoxID string `json:"message_box_id"`
	Content      string `json:"content"`
	Visibility   string `json:"visibility"`
}

type PostPostResponse struct {
	MessageBoxID string `json:"message_box_id"`
	Content      string `json:"content"`
	Visibility   string `json:"visibility"`
}

type PostsGetRequest struct {
	PageNum      int    `json:"page_num"`
	PageSize     int    `json:"page_size"`
	MessageBoxID string `json:"message_box_id"`
}

type PostsGetResponse struct {
	ID         string `json:"id"`
	PosterID   string `json:"poster_id"`
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
}

type PostGetRequest struct {
	ID string `json:"id"`
}

type PostGetResponse struct {
	ID         string   `json:"id"`
	PosterID   string   `json:"poster_id"`
	Content    string   `json:"content"`
	Visibility string   `json:"visibility"`
	Channels   []string `json:"channels"`
}

type PostDeleteRequest struct {
	ID string `json:"id"`
}

type PostDeleteResponse struct {
	Message string `json:"message"`
}
