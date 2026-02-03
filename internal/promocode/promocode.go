package promocode

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"net/http"
	"strings"
	"time"
)

func NewPromocodesRepository(dataBase *db.Db) *PromocodeRepository {
	return &PromocodeRepository{
		DataBase: dataBase,
	}
}

func NewPromocodeHandler(router *http.ServeMux, deps *PromocodeshandlerDeps) *PromocodeHandler {
	handler := &PromocodeHandler{
		Config:              deps.Config,
		PromocodeRepository: deps.PromocodeRepository,
		AuthHandler:         deps.AuthHandler,
	}

	router.HandleFunc("/promocode/create", handler.create())
	router.HandleFunc("/promocode/get", handler.get())       // все или по id через query ?id=...
	router.HandleFunc("/promocode/update", handler.update()) // PATCH, ожидает id в query ?id=...
	router.HandleFunc("/promocode/delete", handler.delete()) // DELETE, ожидает id в query ?id=...

	return handler
}

// ─── create (твой оригинал, чуть подправил только мелкие вещи) ───

func (handler *PromocodeHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CreatePromocodeRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}

		if body.Token != handler.Config.Token.AdminToken {
			res.Json(w, "you are not admin", 401)
			return
		}

		if body.DiscountType != "percent" && body.DiscountType != "fixed" {
			res.Json(w, "invalid discount type", 400)
			return
		}

		promo := PromoCode{
			Id:           token.CreateId(),
			Code:         body.Code,
			ExpiresAt:    body.ExpiresAt,
			UsageLimit:   body.UsageLimit,
			UsagePerUser: body.UsagePerUser,
			UsedCount:    0,
			DiscountType: body.DiscountType,
			Discount:     body.Discount,
			MaxDiscount:  body.MaxDiscount,
			MinPrice:     body.MinPrice,
			IsActive:     true,
		}

		result := handler.PromocodeRepository.DataBase.Create(&promo)
		if result.Error != nil {
			res.Json(w, result.Error.Error(), 500)
			return
		}

		res.Json(w, promo, 200)
	}
}

// ─── get (все или один по ?id=xxx) ───

func (handler *PromocodeHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !handler.isAdmin(r) {
			res.Json(w, "you are not admin", 401)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			// возвращаем все
			var promos []PromoCode
			if err := handler.PromocodeRepository.DataBase.Find(&promos); err != nil {
				res.Json(w, "db error", 500)
				return
			}
			res.Json(w, promos, 200)
			return
		}

		// возвращаем один
		var promo PromoCode
		if err := handler.PromocodeRepository.DataBase.First(&promo, "id = ?", id); err != nil {
			res.Json(w, "promocode not found", 404)
			return
		}

		res.Json(w, promo, 200)
	}
}

// ─── update (PATCH, ожидает ?id=xxx и тело с полями) ───

type UpdatePromocodeRequest struct {
	ExpiresAt    *time.Time `json:"expiresAt"`
	UsageLimit   *int       `json:"usageLimit"`
	UsagePerUser *int       `json:"usagePerUser"`
	IsActive     *bool      `json:"isActive"`
	// можно добавить остальные поля, если хочешь
}

func (handler *PromocodeHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !handler.isAdmin(r) {
			res.Json(w, "you are not admin", 401)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			res.Json(w, "id required", 400)
			return
		}

		body, err := req.HandleBody[UpdatePromocodeRequest](&w, r)
		if err != nil {
			return
		}

		var promo PromoCode
		if err := handler.PromocodeRepository.DataBase.First(&promo, "id = ?", id); err != nil {
			res.Json(w, "promocode not found", 404)
			return
		}

		// обновляем только то, что пришло
		if body.ExpiresAt != nil {
			promo.ExpiresAt = *body.ExpiresAt
		}
		if body.UsageLimit != nil {
			promo.UsageLimit = *body.UsageLimit
		}
		if body.UsagePerUser != nil {
			promo.UsagePerUser = *body.UsagePerUser
		}
		if body.IsActive != nil {
			promo.IsActive = *body.IsActive
		}

		if err := handler.PromocodeRepository.DataBase.Save(&promo); err != nil {
			res.Json(w, "db error", 500)
			return
		}

		res.Json(w, promo, 200)
	}
}

// ─── delete (?id=xxx) ───

func (handler *PromocodeHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !handler.isAdmin(r) {
			res.Json(w, "you are not admin", 401)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			res.Json(w, "id required", 400)
			return
		}

		result := handler.PromocodeRepository.DataBase.Delete(&PromoCode{}, "id = ?", id)
		if result.Error != nil {
			res.Json(w, "db error", 500)
			return
		}
		if result.RowsAffected == 0 {
			res.Json(w, "promocode not found", 404)
			return
		}

		res.Json(w, "promocode deleted", 200)
	}
}

// ─── вспомогательная проверка админа (чтобы не дублировать) ───

func (handler *PromocodeHandler) isAdmin(r *http.Request) bool {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		tokenStr = r.URL.Query().Get("token") // оставил совместимость с твоим create
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	// здесь твоя логика проверки токена
	// пока оставляю как было у тебя в create
	return tokenStr == handler.Config.Token.AdminToken
}
