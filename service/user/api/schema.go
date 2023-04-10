package api

type LoginRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
}

type UserResponse struct {
	Username string `json:"username"`
}
