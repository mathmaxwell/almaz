package buy

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/pkg/db"
)

type BuyRepository struct {
	DataBase *db.Db
}

type BuyRepositoryDeps struct {
	DataBase *db.Db
}
type BuyHandler struct {
	*configs.Config
	BuyRepository BuyRepository
	AuthHandler   *auth.AuthHandler
}
type BuyhandlerDeps struct {
	*configs.Config
	BuyRepository *BuyRepository
	AuthHandler   *auth.AuthHandler
}
type Games struct {
	Id         string `json:"id"`
	Video      string `json:"video"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	HowToUseUz string `json:"howToUseUz"`
	HowToUseRu string `json:"howToUseRu"`
	HelpImage  string `json:"helpImage"`
	Place      string `json:"place"`
}
type Buy struct {
	Id       string `json:"id"`
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Day      int    `json:"day"`
	Hour     int    `json:"hour"`
	Minute   int    `json:"minute"`
	UserId   string `json:"userId"` //token
	Price    int    `json:"price"`
	Status   string `json:"status"` //sucses wait
	GameId   string `json:"gameId"`
	BotId    string `json:"botId"`
	ServerId string `json:"serverId"`
	PlayerId string `json:"playerId"`
	Order    string `json:"order"`
}
type createBuyRequest struct {
	Token    string `json:"token"`
	GameId   string `json:"gameId"`
	PlayerId string `json:"playerId"`
	ServerId string `json:"serverId"`
	BotId    string `json:"botId"`
	OfferId  string `json:"offerId"`
}
type Transaction struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Price     int    `json:"price"`
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
	GameName  string `json:"gameName"`
	DonatName string `json:"donatName"`
	CreatedBy string `json:"createdBy"` // system | admin | gateway
	Order     string `json:"order"`
}
type Offers struct {
	Id     string `json:"id"`
	GameId string `json:"gameId"`
	BotId  string `json:"botId"`
	Image  string `json:"image"`
	UzName string `json:"uzName"`
	RuName string `json:"ruName"`
	Price  string `json:"price"`
	RuDesc string `json:"ruDesc"`
	UzDesc string `json:"uzDesc"`
	Status string `json:"status"`
}
