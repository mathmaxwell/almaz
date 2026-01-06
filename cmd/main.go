package main

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/internal/games"
	"demo/almaz/internal/offers"
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
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
