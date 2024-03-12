package model

type AuthRequest struct {
	Username string `form:"username" binding:"required,min=6"`
	Password string `form:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
