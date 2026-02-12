package transactions

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/pkg/db"
)

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
type User struct {
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Balance  int    `json:"balance"`
}
type TransactionhandlerDeps struct {
	*configs.Config
	TransactionRepository *TransactionRepository
	AuthHandler           *auth.AuthHandler
}
type TransactionRepository struct {
	DataBase *db.Db
}
type getRequest struct {
	UserId string `json:"userId"`
}
type createRequest struct {
	UserId    string `json:"userId"`
	Price     int    `json:"price"`
	GameName  string `json:"gameName"`
	DonatName string `json:"donatName"`
	CreatedBy string `json:"createdBy"`
}
type TransactionHandler struct {
	*configs.Config
	TransactionRepository TransactionRepository
	AuthHandler           *auth.AuthHandler
}
type deleteRequest struct {
	Id string `json:"id"`
}
type getByPeriodRequest struct {
	Token      string `json:"token"`
	StartDay   int    `json:"startDay"`
	StartMonth int    `json:"startMonth"`
	StartYear  int    `json:"startYear"`
	EndDay     int    `json:"endDay"`
	EndMonth   int    `json:"endMonth"`
	EndYear    int    `json:"endYear"`
}
