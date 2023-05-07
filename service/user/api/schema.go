package api

type LoginRequest struct {
	Username string `json:"username" validate:"min=8"`
	Password string `json:"password" validate:"min=8"`
}

type UserResponse struct {
	Username string `json:"username"`
}

type ResetRequest struct {
	OldPassword string `json:"old_password" validate:"min=8"`
	NewPassword string `json:"new_password" validate:"min=8"`
}
