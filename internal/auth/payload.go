package auth

import (
	"demo/almaz/configs"
	"demo/almaz/pkg/db"
)

type User struct {
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Balance  int    `json:"balance"`
	UserRole string `json:"userRole"`
}
type UpdateUserRequest struct {
	Token    string `json:"token"`
	UserId   string `json:"userId"`
	UserRole string `json:"userRole"`
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
type GetBalanceRequest struct {
	UserId string `json:"userId"`
}
type DeleteRequest struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}
type AuthRepository struct {
	DataBase *db.Db
}

type AuthRepositoryDeps struct {
	DataBase *db.Db
}
type GetUsersRequest struct {
	Page         int     `json:"page"`
	Count        int     `json:"count"`
	Login        *string `json:"login"`
	Token        *string `json:"token"`
	StartBalance *int    `json:"startBalance"`
	UserRole     *string `json:"userRole"`
}
