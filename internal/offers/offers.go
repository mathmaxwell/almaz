package offers

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/files"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

func NewOffersRepository(dataBase *db.Db) *OffersRepository {
	return &OffersRepository{
		DataBase: dataBase,
	}
}
func NewOffersHandler(router *http.ServeMux, deps OffersshandlerDeps) *OffersHandler {
	handler := &OffersHandler{
		Config:           deps.Config,
		OffersRepository: *deps.OffersRepository,
		AuthHandler:      deps.AuthHandler,
	}
	router.HandleFunc("/offers/create", handler.create())
	router.HandleFunc("/offers/getOffers", handler.getOffers())
	router.HandleFunc("/offers/getOffersById", handler.getOffersById())
	router.HandleFunc("/offers/updateOffer", handler.updateOffer())
	router.HandleFunc("/offers/deleteOffer", handler.deleteOffer())
	return handler
}

func (handler *OffersHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.FormValue("token")
		if userToken != handler.Config.Token.AdminToken {
			res.Json(w, "you are not admin", 401)
			return
		}

		if err := r.ParseMultipartForm(200 << 20); err != nil {
			res.Json(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			res.Json(w, "image is not found", 400)
			return
		}
		photoPath, err := files.SaveFile(file, header)
		if err != nil {
			res.Json(w, "failed to save image", http.StatusInternalServerError)
			return
		}

		var gameId string
		if r.FormValue("id") == "" {
			gameId = token.CreateId()
		} else {
			gameId = r.FormValue("id")
		}
		newGame := Offers{
			Id:     gameId,
			Image:  photoPath,
			Status: r.FormValue("status"),
			GameId: r.FormValue("gameId"),
			UzName: r.FormValue("uzName"),
			RuName: r.FormValue("ruName"),
			Price:  r.FormValue("price"),
			RuDesc: r.FormValue("ruDesc"),
			UzDesc: r.FormValue("uzDesc"),
			BotId:  r.FormValue("botId"),
		}
		if err := handler.OffersRepository.DataBase.Create(&newGame).Error; err != nil {
			res.Json(w, "db error", 500)
			return
		}
		res.Json(w, newGame, 200)
	}
}
func (handler *OffersHandler) getOffers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetOffersRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 400)
			return
		}
		var offers []Offers
		err = handler.OffersRepository.DataBase.Model(&Offers{}).Where("game_id = ?", body.GameId).Find(&offers).Error
		if err != nil {
			res.Json(w, "failed to get offers", 500)
			return
		}
		res.Json(w, offers, 200)
	}
}
func (handler *OffersHandler) getOffersById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DeleteOffersRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 400)
			return
		}
		var offers Offers
		err = handler.OffersRepository.DataBase.Where("id = ?", body.Id).First(&offers).Error

		if err != nil {
			res.Json(w, "failed to get offers", 500)
			return
		}
		res.Json(w, offers, 200)
	}
}
func (handler *OffersHandler) updateOffer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(200 << 20); err != nil {
			res.Json(w, "не удалось разобрать форму", http.StatusBadRequest)
			return
		}
		userToken := r.FormValue("token")
		if userToken != handler.Config.Token.AdminToken {
			res.Json(w, "доступ запрещён: требуется роль администратора", http.StatusForbidden)
			return
		}
		gameId := r.FormValue("id")
		if gameId == "" {
			res.Json(w, "ID игры обязателен", http.StatusBadRequest)
			return
		}
		var offer Offers
		if err := handler.OffersRepository.DataBase.Where("id = ?", gameId).First(&offer).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "offer is not found", http.StatusNotFound)
			} else {
				res.Json(w, "ошибка базы данных", http.StatusInternalServerError)
			}
			return
		}
		photoPath := offer.Image
		file, header, err := r.FormFile("image")
		if err == nil && file != nil {
			defer file.Close()
			var saveErr error
			photoPath, saveErr = files.SaveFile(file, header)
			if saveErr != nil {
				res.Json(w, "не удалось сохранить изображение", http.StatusInternalServerError)
				return
			}
		}
		updateOffer := Offers{
			Image:  photoPath,
			Status: r.FormValue("status"),
			Id:     r.FormValue("id"),
			GameId: r.FormValue("gameId"),
			UzName: r.FormValue("uzName"),
			RuName: r.FormValue("ruName"),
			Price:  r.FormValue("price"),
			RuDesc: r.FormValue("ruDesc"),
			UzDesc: r.FormValue("uzDesc"),
			BotId:  r.FormValue("botId"),
		}
		if err := handler.OffersRepository.DataBase.Model(&offer).Updates(updateOffer).Error; err != nil {
			res.Json(w, "не удалось обновить предложение", http.StatusInternalServerError)
			return
		}
		res.Json(w, map[string]interface{}{
			"message": "Предложения обновлена",
			"game":    updateOffer,
		}, http.StatusOK)
	}
}
func (handler *OffersHandler) deleteOffer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DeleteOffersRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.Token != handler.Config.Token.AdminToken {
			res.Json(w, "you are not admin", 403)
			return
		}
		db := handler.OffersRepository.DataBase
		result := db.Delete(&Offers{}, "id = ?", body.Id)
		if result.Error != nil {
			res.Json(w, result.Error.Error(), 500)
			return
		}
		if result.RowsAffected == 0 {
			res.Json(w, "offer is not found", 404)
			return
		}
		res.Json(w, "offer deleted", 200)
	}
}
