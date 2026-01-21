package admincart

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/pkg/db"
)

type AdmincartRepository struct {
	DataBase *db.Db
}

type AdmincartRepositoryDeps struct {
	DataBase *db.Db
}

type Admincart struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
	Type   string `json:"type"`
}
type User struct {
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Balance  int    `json:"balance"`
}

type AdmincartHandler struct {
	*configs.Config
	AdmincartRepository AdmincartRepository
	AuthHandler         *auth.AuthHandler
}

type AdmincarthandlerDeps struct {
	*configs.Config
	AdmincartRepository *AdmincartRepository
	AuthHandler         *auth.AuthHandler
}
type GetAdmincartRequest struct {
	Token string `json:"token"`
}
type CreateAdmincartRequest struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Number string `json:"number"`
	Type   string `json:"type"`
}
type UpdateAdmincartRequest struct {
	Id     string `json:"id"`
	Token  string `json:"token"`
	Name   string `json:"name"`
	Number string `json:"number"`
	Type   string `json:"type"`
}
type DeleteAdmincartRequest struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}
