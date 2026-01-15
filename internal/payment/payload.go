package payment

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/pkg/db"
	"time"
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
type createPaymentTelegram struct {
	Text string `json:"text"`
}

func isExpired(p Payment) bool {
	created := time.Date(
		p.Year,
		time.Month(p.Month),
		p.Day,
		p.Hour,
		p.Minute,
		0,
		0,
		time.Local,
	)

	return time.Since(created) >= 6*time.Minute
}
