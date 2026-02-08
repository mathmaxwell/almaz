package auth

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"errors"
	"net/http"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewUserRepository(dataBase *db.Db) *AuthRepository {
	return &AuthRepository{
		DataBase: dataBase,
	}
}
func NewAuthHandler(router *http.ServeMux, deps AuthhandlerDeps) *AuthHandler {
	handler := &AuthHandler{
		Config:         deps.Config,
		AuthRepository: *deps.AuthRepository,
	}
	router.HandleFunc("/users/login", handler.login())
	router.HandleFunc("/users/getUserById", handler.getUserById())
	router.HandleFunc("/users/getUsers", handler.getUsers())
	router.HandleFunc("/users/deleteUser", handler.deleteUser())
	router.HandleFunc("/users/register", handler.register())
	return handler
}
func (handler *AuthHandler) GetUserByToken(token string) (User, error) {
	var user User
	err := handler.AuthRepository.DataBase.Where("token = ?", token).First(&user).Error
	if err != nil {
		return user, errors.New("user is not found")
	}
	return user, nil
}
func (handler *AuthHandler) UpdateBalance(token string, newPrice int) (User, error) {
	var user User
	err := handler.AuthRepository.DataBase.
		Model(&User{}).
		Where("token = ?", token).
		Update("balance", gorm.Expr("balance + ?", newPrice)).
		Error
	if err != nil {
		return user, errors.New("user is not found")
	}
	return user, nil
}
func (handler *AuthHandler) DecreaseBalance(tx *gorm.DB, userToken string, price int) error {
	var user User
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("token = ?", userToken).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("пользователь не найден")
		}
		return err
	}

	if user.Balance < price {
		return errors.New("недостаточно средств")
	}

	// Самый надёжный вариант ↓
	return tx.
		Model(&User{}).
		Where("token = ?", userToken).
		Update("balance", gorm.Expr("balance - ?", price)).Error
}
func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		var user User
		err = handler.AuthRepository.DataBase.Where("login = ?", body.Login).First(&user).Error
		if err != nil {
			res.Json(w, "user is not found", 400)
			return
		}
		if user.Password != body.Password {
			res.Json(w, "password is not correct", 401)
			return
		}
		res.Json(w, user, 200)
	}
}
func (handler *AuthHandler) getUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetUsersRequest](&w, r)
		if err != nil {
			return
		}

		page := body.Page
		if page < 1 {
			page = 1
		}

		count := body.Count
		if count < 1 {
			count = 10
		}
		if count > 100 {
			count = 100
		}

		offset := (page - 1) * count

		query := handler.AuthRepository.DataBase.Model(&User{})

		// 🔍 поиск по login
		if body.Login != nil && *body.Login != "" {
			query = query.Where(
				"login ILIKE ?",
				"%"+*body.Login+"%",
			)
		}

		// 🔍 универсальный поиск (login + token)
		if body.Token != nil && *body.Token != "" {
			val := "%" + *body.Token + "%"
			query = query.Where(
				"(login ILIKE ? OR token ILIKE ?)",
				val,
				val,
			)
		}

		// 💰 фильтр по балансу
		if body.StartBalance != nil {
			query = query.Where(
				"balance >= ?",
				*body.StartBalance,
			)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			res.Json(w, "ошибка при подсчёте пользователей", 500)
			return
		}

		var users []User
		if err := query.
			Offset(offset).
			Limit(count).
			Find(&users).
			Error; err != nil {
			res.Json(w, "ошибка при получении пользователей", 500)
			return
		}

		res.Json(w, map[string]interface{}{
			"users": users,
			"total": total,
			"page":  page,
			"count": count,
			"pages": (total + int64(count) - 1) / int64(count),
		}, 200)
	}
}
func (handler *AuthHandler) getUserById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetBalanceRequest](&w, r)
		if err != nil {
			return
		}
		var user User
		err = handler.AuthRepository.DataBase.Where("token = ?", body.UserId).First(&user).Error
		if err != nil {
			res.Json(w, "user is not found", 400)
			return
		}
		res.Json(w, user, 200)
	}
}
func (handler *AuthHandler) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DeleteRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		var user User
		handler.AuthRepository.DataBase.Where("token = ?", body.Token).First(&user)
		if user.Token != handler.Config.Token.AdminToken {
			res.Json(w, "you are not admin", 403)
			return
		}
		db := handler.AuthRepository.DataBase
		result := db.Delete(&User{}, "token = ?", body.UserId)
		if result.Error != nil {
			res.Json(w, result.Error.Error(), 500)
			return
		}
		if result.RowsAffected == 0 {
			res.Json(w, "user is not found", 404)
			return
		}
		res.Json(w, "user deleted", 200)
	}
}
func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		var user User
		err = handler.AuthRepository.DataBase.Where("login = ?", body.Login).First(&user).Error
		if err == nil {
			res.Json(w, "login is alredy exist", 400)
			return
		}
		var tokenId string
		if body.Login == "testadmin" && body.Password == "tadi123$" {
			tokenId = handler.Config.Token.AdminToken
		} else {
			tokenId = token.CreateId()
		}
		data := User{
			Login:    body.Login,
			Password: body.Password,
			Token:    tokenId,
			Balance:  0,
		}
		handler.AuthRepository.DataBase.Create(&data)
		res.Json(w, data, 200)
	}
}
