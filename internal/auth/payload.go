package auth

import "demo/purpleSchool/configs"

type LoginResponse struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type AuthHandler struct {
	*configs.Config
}

type AuthhandlerDeps struct {
	*configs.Config
}
type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
