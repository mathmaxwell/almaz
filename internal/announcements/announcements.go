package announcements

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

func NewOffersRepository(dataBase *db.Db) *AnnouncementsRepository {
	return &AnnouncementsRepository{
		DataBase: dataBase,
	}
}
func NewOffersHandler(router *http.ServeMux, deps AnnouncementshandlerDeps) *AnnouncementsHandler {
	handler := &AnnouncementsHandler{
		Config:                  deps.Config,
		AnnouncementsRepository: *deps.AnnouncementsRepository,
		AuthHandler:             deps.AuthHandler,
	}
	router.HandleFunc("/announcements/create", handler.create())
	router.HandleFunc("/announcements/getAnnouncements", handler.get())
	router.HandleFunc("/announcements/update", handler.update())
	router.HandleFunc("/announcements/delete", handler.delete())
	return handler
}

func (handler *AnnouncementsHandler) create() http.HandlerFunc {
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

		newAnnouncements := Announcements{
			Id:     token.CreateId(),
			Image:  photoPath,
			Uz:     r.FormValue("uz"),
			Ru:     r.FormValue("ru"),
			RuText: r.FormValue("ruText"),
			UzText: r.FormValue("uzText"),
		}
		if err := handler.AnnouncementsRepository.DataBase.Create(&newAnnouncements).Error; err != nil {
			res.Json(w, "db error", 500)
			return
		}
		res.Json(w, newAnnouncements, 200)
	}
}
func (handler *AnnouncementsHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetAnnouncementsRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 400)
			return
		}
		var offers []Announcements
		err = handler.AnnouncementsRepository.DataBase.Model(&Announcements{}).Find(&offers).Error
		if err != nil {
			res.Json(w, "failed to get offers", 500)
			return
		}
		res.Json(w, offers, 200)
	}
}
func (handler *AnnouncementsHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(200); err != nil {
			res.Json(w, "не удалось разобрать форму", http.StatusBadRequest)
			return
		}
		userToken := r.FormValue("token")
		if userToken != handler.Config.Token.AdminToken {
			res.Json(w, "доступ запрещён: требуется роль администратора", http.StatusForbidden)
			return
		}
		id := r.FormValue("id")
		if id == "" {
			res.Json(w, "ID игры обязателен", http.StatusBadRequest)
			return
		}
		var announcement Announcements
		if err := handler.AnnouncementsRepository.DataBase.Where("id = ?", id).First(&announcement).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "offer is not found", http.StatusNotFound)
			} else {
				res.Json(w, "ошибка базы данных", http.StatusInternalServerError)
			}
			return
		}
		photoPath := announcement.Image
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
		updateOffer := Announcements{
			Image:  photoPath,
			Id:     r.FormValue("id"),
			Uz:     r.FormValue("uz"),
			Ru:     r.FormValue("ru"),
			RuText: r.FormValue("ruText"),
			UzText: r.FormValue("uzText"),
		}
		if err := handler.AnnouncementsRepository.DataBase.Model(&announcement).Updates(updateOffer).Error; err != nil {
			res.Json(w, "не удалось обновить предложение", http.StatusInternalServerError)
			return
		}
		res.Json(w, map[string]interface{}{
			"message": "Предложения обновлена",
			"game":    updateOffer,
		}, http.StatusOK)
	}
}
func (handler *AnnouncementsHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DeleteAnnouncementsRequest](&w, r)
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
		db := handler.AnnouncementsRepository.DataBase
		result := db.Delete(&Announcements{}, "id = ?", body.Id)
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
