package auth

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/pkg/db"
)

type LoginResponse struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type AuthHandler struct {
	*configs.Config
	AuthRepository AuthRepository
}

type AuthhandlerDeps struct {
	*configs.Config
	AuthRepository *AuthRepository
}
type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type AuthRepository struct {
	DataBase *db.Db
}

type AuthRepositoryDeps struct {
	DataBase *db.Db
}
