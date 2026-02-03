package main

import (
	"demo/almaz/configs"
	"demo/almaz/internal/admincart"
	"demo/almaz/internal/announcements"
	"demo/almaz/internal/auth"
	"demo/almaz/internal/buy"
	"demo/almaz/internal/games"
	"demo/almaz/internal/offers"
	"demo/almaz/internal/payment"
	"demo/almaz/internal/promocode"
	"demo/almaz/internal/transactions"
	"demo/almaz/pkg/cors"
	"demo/almaz/pkg/db"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDB(conf)
	router := http.NewServeMux()
	router.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	mux := cors.Cors(router)
	//register
	authRepo := auth.NewUserRepository(database)
	database.AutoMigrate(&auth.User{})
	authHandler := auth.NewAuthHandler(router, auth.AuthhandlerDeps{
		Config:         conf,
		AuthRepository: authRepo,
	})
	//games
	gamesRepo := games.NewGamesRepository(database)
	database.AutoMigrate(&games.Games{})
	games.NewGamesHandler(router, games.GameshandlerDeps{
		Config:          conf,
		GamesRepository: gamesRepo,
		AuthHandler:     authHandler,
	})
	//offers
	offersRepo := offers.NewOffersRepository(database)
	database.AutoMigrate(&offers.Offers{})
	offers.NewOffersHandler(router, offers.OffersshandlerDeps{
		Config:           conf,
		OffersRepository: offersRepo,
		AuthHandler:      authHandler,
	})

	//announcements elonlar
	announcementsRepo := announcements.NewOffersRepository(database)
	database.AutoMigrate(&announcements.Announcements{})
	announcements.NewOffersHandler(router, announcements.AnnouncementshandlerDeps{
		Config:                  conf,
		AnnouncementsRepository: announcementsRepo,
		AuthHandler:             authHandler,
	})
	//admin cards
	admincardRepo := admincart.NewAdmincartRepository(database)
	database.AutoMigrate(&admincart.Admincart{})
	admincart.NewAdmincartHandler(router, admincart.AdmincarthandlerDeps{
		Config:              conf,
		AdmincartRepository: admincardRepo,
		AuthHandler:         authHandler,
	})
	//Payment
	paymentRepo := payment.NewPaymentRepository(database)
	database.AutoMigrate(&payment.Payment{})
	payment.NewPaymentHandler(router, payment.PaymenthandlerDeps{
		Config:            conf,
		PaymentRepository: paymentRepo,
		AuthHandler:       authHandler,
	})
	//transactions
	transactionsRepo := transactions.NewTransactionRepository(database)
	database.AutoMigrate(&transactions.Transaction{})
	transactions.NewTranactionHandler(router, &transactions.TransactionhandlerDeps{
		Config:                conf,
		TransactionRepository: transactionsRepo,
		AuthHandler:           authHandler,
	})
	//promocodes
	promocodesRepo := promocode.NewPromocodesRepository(database)
	database.AutoMigrate(&promocode.PromoCode{})
	promocode.NewPromocodeHandler(router, &promocode.PromocodeshandlerDeps{
		Config:              conf,
		PromocodeRepository: promocodesRepo,
		AuthHandler:         authHandler,
	})
	//buy
	buyRepo := buy.NewBuyRepository(database)
	database.AutoMigrate(&buy.Buy{})
	buy.NewGamesHandler(router, &buy.BuyhandlerDeps{
		Config:        conf,
		BuyRepository: buyRepo,
		AuthHandler:   authHandler,
	})
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
