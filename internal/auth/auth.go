package auth

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
)

func NewUserRepository(dataBase *db.Db) *AuthRepository {
	return &AuthRepository{
		DataBase: dataBase,
	}
}

func NewAuthHandler(router *http.ServeMux, deps AuthhandlerDeps) {
	handler := &AuthHandler{
		AuthRepository: *deps.AuthRepository,
		Config:         deps.Config,
	}
	router.HandleFunc("/users/login", handler.login())
	router.HandleFunc("/users/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		//ищет по логину, если пароль совполает, то дает токен
		email := body.Login
		password := body.Password
		token := token.CreateId()
		data := LoginResponse{
			Login:    email,
			Password: password,
			Token:    token,
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		//создает новый запись
		email := body.Login
		password := body.Password
		token := token.CreateId()
		data := LoginResponse{
			Login:    email,
			Password: password,
			Token:    token,
		}
		res.Json(w, data, 200)
	}
}
