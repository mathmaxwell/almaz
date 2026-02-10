package offers

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/pkg/db"
)

type OffersRepository struct {
	DataBase *db.Db
}

type OffersRepositoryDeps struct {
	DataBase *db.Db
}

type Offers struct {
	Id         string `json:"id"`
	GameId     string `json:"gameId"`
	BotId      string `json:"botId"`
	Image      string `json:"image"`
	UzName     string `json:"uzName"`
	RuName     string `json:"ruName"`
	Price      string `json:"price"`
	SuperPrice string `json:"superPrice"`
	RuDesc     string `json:"ruDesc"`
	UzDesc     string `json:"uzDesc"`
	Status     string `json:"status"`
}
type User struct {
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Balance  int    `json:"balance"`
}

type OffersHandler struct {
	*configs.Config
	OffersRepository OffersRepository
	AuthHandler      *auth.AuthHandler
}

type OffersshandlerDeps struct {
	*configs.Config
	OffersRepository *OffersRepository
	AuthHandler      *auth.AuthHandler
}
type GetOffersRequest struct {
	Token  string `json:"token"`
	GameId string `json:"gameId"`
}
type DeleteOffersRequest struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}
