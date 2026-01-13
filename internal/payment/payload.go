package payment

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/pkg/db"
)

type PaymentRepository struct {
	DataBase *db.Db
}

type PaymentRepositoryDeps struct {
	DataBase *db.Db
}

type User struct {
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Balance  int    `json:"balance"`
}

type PaymentHandler struct {
	*configs.Config
	PaymentRepository PaymentRepository
	AuthHandler       *auth.AuthHandler
}

type PaymenthandlerDeps struct {
	*configs.Config
	PaymentRepository *PaymentRepository
	AuthHandler       *auth.AuthHandler
}
type Payment struct {
	Id        string `json:"id"`
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
	UserId    string `json:"userId"`
	Price     int    `json:"price"`
	IsWorking bool   `json:"isWorking"`
}
type createPaymentRequest struct {
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
	UserId    string `json:"userId"`
	Price     int    `json:"price"`
	IsWorking bool   `json:"isWorking"`
}
type updatePaymentRequest struct {
	Token     string `json:"token"`
	Id        string `json:"id"`
	IsWorking bool   `json:"isWorking"`
}
type deletePaymentRequest struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}
type getPaymentRequest struct {
	Token string `json:"token"`
}
