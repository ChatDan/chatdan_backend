package api

type LoginRequest struct {
	username *string
	email    *string
	password string
}

type UserResponse struct {
	id       string
	username string
}
