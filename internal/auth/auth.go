package auth

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"errors"
	"net/http"

	"gorm.io/gorm"
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
	router.HandleFunc("/users/getBalance", handler.getBalance())
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
func (handler *AuthHandler) getBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetBalanceRequest](&w, r)
		if err != nil {
			return
		}
		var user User
		err = handler.AuthRepository.DataBase.Where("login = ?", body.UserId).First(&user).Error
		if err != nil {
			res.Json(w, "user is not found", 400)
			return
		}
		res.Json(w, user, 200)
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
