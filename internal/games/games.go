package games

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

func NewGamesRepository(dataBase *db.Db) *GamesRepository {
	return &GamesRepository{
		DataBase: dataBase,
	}
}
func NewGamesHandler(router *http.ServeMux, deps GameshandlerDeps) *GamesHandler {
	handler := &GamesHandler{
		Config:          deps.Config,
		GamesRepository: *deps.GamesRepository,
		AuthHandler:     deps.AuthHandler,
	}
	router.HandleFunc("/games/create", handler.create())
	router.HandleFunc("/games/getGames", handler.getGames())
	router.HandleFunc("/games/getGameById", handler.getGameById())
	router.HandleFunc("/games/updateGame", handler.updateGame())
	router.HandleFunc("/games/deleteGame", handler.deleteGame())
	return handler
}

func (handler *GamesHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.FormValue("token")
		if userToken != handler.Config.Token.AdminToken {
			res.Json(w, "you are not admin", 401)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
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
		videoFile, videoHeader, videoErr := r.FormFile("video")
		if videoErr != nil {
			res.Json(w, "video is not found", 400)
			return
		}
		videoPath, err := files.SaveFile(videoFile, videoHeader)
		if err != nil {
			res.Json(w, "failed to save video", http.StatusInternalServerError)
			return
		}
		fileHelper, headerHelper, err := r.FormFile("helpImage")
		if err != nil {
			res.Json(w, "image is not found", 400)
			return
		}
		photoPathHelper, err := files.SaveFile(fileHelper, headerHelper)
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
		newGame := Games{
			Id:          gameId,
			Video:       videoPath,
			Name:        r.FormValue("name"),
			Place:       r.FormValue("place"),
			HowToUseRu:  r.FormValue("howToUseRu"),
			HowToUseUz:  r.FormValue("howToUseUz"),
			Description: r.FormValue("description"),
			Image:       photoPath,
			HelpImage:   photoPathHelper,
		}
		if err := handler.GamesRepository.DataBase.Create(&newGame).Error; err != nil {
			res.Json(w, "db error", 500)
			return
		}
		res.Json(w, newGame, 200)
	}
}
func (handler *GamesHandler) getGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetGamesRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		var games []Games
		err = handler.GamesRepository.DataBase.Model(&Games{}).Find(&games).Error
		if err != nil {
			res.Json(w, "failed to get games", 500)
			return
		}
		res.Json(w, games, 200)
	}
}
func (handler *GamesHandler) getGameById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DeleteGameRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 400)
			return
		}
		var game Games
		err = handler.GamesRepository.DataBase.Where("id = ?", body.Id).First(&game).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "game not found", 404)
				return
			}
			res.Json(w, "failed to get game", 500)
			return
		}
		res.Json(w, game, 200)
	}
}
func (handler *GamesHandler) updateGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(200 << 20); err != nil {
			res.Json(w, "не удалось разобрать форму", http.StatusBadRequest)
			return
		}
		userToken := r.FormValue("token")
		if userToken != handler.Config.Token.AdminToken {
			res.Json(w, "доступ запрещён", http.StatusForbidden)
			return
		}
		gameId := r.FormValue("id")
		if gameId == "" {
			res.Json(w, "ID игры обязателен", http.StatusBadRequest)
			return
		}

		var game Games
		if err := handler.GamesRepository.DataBase.
			Where("id = ?", gameId).
			First(&game).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "игра не найдена", http.StatusNotFound)
			} else {
				res.Json(w, "ошибка базы данных", http.StatusInternalServerError)
			}
			return
		}

		updates := map[string]interface{}{}
		if file, header, err := r.FormFile("image"); err == nil {
			defer file.Close()

			path, err := files.SaveFile(file, header)
			if err != nil {
				res.Json(w, "не удалось сохранить изображение", http.StatusInternalServerError)
				return
			}
			updates["image"] = path
		}

		if file, header, err := r.FormFile("helpImage"); err == nil {
			defer file.Close()

			path, err := files.SaveFile(file, header)
			if err != nil {
				res.Json(w, "не удалось сохранить изображение помощи", http.StatusInternalServerError)
				return
			}
			updates["help_image"] = path
		}

		file, header, err := r.FormFile("video")
		if err != nil && err != http.ErrMissingFile {
			res.Json(w, "ошибка загрузки видео", http.StatusBadRequest)
			return
		}

		if err == nil {
			defer file.Close()

			path, err := files.SaveFile(file, header)
			if err != nil {
				res.Json(w, "не удалось сохранить видео", http.StatusInternalServerError)
				return
			}
			updates["video"] = path
		}

		if v := r.FormValue("name"); v != "" {
			updates["name"] = v
		}
		if v := r.FormValue("howToUseRu"); v != "" {
			updates["how_to_use_ru"] = v
		}
		if v := r.FormValue("howToUseUz"); v != "" {
			updates["how_to_use_uz"] = v
		}
		if v := r.FormValue("place"); v != "" {
			updates["place"] = v
		}
		if v := r.FormValue("description"); v != "" {
			updates["description"] = v
		}
		if len(updates) == 0 {
			res.Json(w, "нет данных для обновления", http.StatusBadRequest)
			return
		}

		if err := handler.GamesRepository.DataBase.
			Model(&game).
			Updates(updates).Error; err != nil {
			res.Json(w, "не удалось обновить игру", http.StatusInternalServerError)
			return
		}

		res.Json(w, map[string]interface{}{
			"message": "игра обновлена",
			"game":    updates,
		}, http.StatusOK)
	}
}

func (handler *GamesHandler) deleteGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DeleteGameRequest](&w, r)
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
		db := handler.GamesRepository.DataBase
		result := db.Delete(&Games{}, "id = ?", body.Id)
		if result.Error != nil {
			res.Json(w, result.Error.Error(), 500)
			return
		}
		if result.RowsAffected == 0 {
			res.Json(w, "game is not found", 404)
			return
		}
		res.Json(w, "game deleted", 200)
	}
}
