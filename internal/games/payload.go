package games

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/pkg/db"
)

type GamesRepository struct {
	DataBase *db.Db
}

type GamesRepositoryDeps struct {
	DataBase *db.Db
}

type Games struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
type User struct {
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Balance  int    `json:"balance"`
}

type GamesHandler struct {
	*configs.Config
	GamesRepository GamesRepository
	AuthHandler     *auth.AuthHandler
}

type GameshandlerDeps struct {
	*configs.Config
	GamesRepository *GamesRepository
	AuthHandler     *auth.AuthHandler
}
type GetGamesRequest struct {
	Token string `json:"token"`
}
type DeleteGameRequest struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}
