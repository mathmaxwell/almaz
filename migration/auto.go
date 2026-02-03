package main

import (
	"demo/almaz/internal/admincart"
	"demo/almaz/internal/announcements"
	"demo/almaz/internal/auth"
	"demo/almaz/internal/buy"
	"demo/almaz/internal/games"
	"demo/almaz/internal/offers"
	"demo/almaz/internal/payment"
	"demo/almaz/internal/promocode"
	"demo/almaz/internal/transactions"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&auth.User{})
	db.AutoMigrate(&games.Games{})
	db.AutoMigrate(&offers.Offers{})
	db.AutoMigrate(&announcements.Announcements{})
	db.AutoMigrate(&admincart.Admincart{})
	db.AutoMigrate(&payment.Payment{})
	db.AutoMigrate(&transactions.Transaction{})
	db.AutoMigrate(&promocode.PromoCode{})
	db.AutoMigrate(&buy.Buy{})
}
