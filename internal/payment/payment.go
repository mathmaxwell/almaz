package payment

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func NewPaymentRepository(dataBase *db.Db) *PaymentRepository {
	return &PaymentRepository{
		DataBase: dataBase,
	}
}
func NewPaymentHandler(router *http.ServeMux, deps PaymenthandlerDeps) *PaymentHandler {
	handler := &PaymentHandler{
		Config:            deps.Config,
		PaymentRepository: *deps.PaymentRepository,
		AuthHandler:       deps.AuthHandler,
	}
	router.HandleFunc("/payment/getPayment", handler.getPayment())
	router.HandleFunc("/payment/getPaymentByUser", handler.getPaymentByUser())
	router.HandleFunc("/payment/getAllPayment", handler.getAllPayment())
	router.HandleFunc("/payment/createPayment", handler.createPayment())
	router.HandleFunc("/payment/updatePayment", handler.updatePayment())
	router.HandleFunc("/payment/deletePayment", handler.deletePayment())
	router.HandleFunc("/payment/createTelegram", handler.createTelegram())
	return handler
}
func (handler *PaymentHandler) getPayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getPaymentRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, err, 401)
			return
		}
		var payments []Payment
		err = handler.PaymentRepository.DataBase.
			Where("is_working = ?", true).
			Find(&payments).Error
		if err != nil {
			res.Json(w, err, 500)
			return
		}
		var result []Payment
		for _, p := range payments {
			if isExpired(p) {
				handler.PaymentRepository.DataBase.
					Where("id = ?", p.Id).
					Delete(&Payment{})
				continue
			}
			result = append(result, p)
		}

		res.Json(w, result, 200)

	}
}
func (handler *PaymentHandler) getPaymentByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[deletePaymentRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, err, 401)
			return
		}
		var payments []Payment
		err = handler.PaymentRepository.DataBase.
			Where("user_id = ?", body.Id).
			Find(&payments).Error
		if err != nil {
			res.Json(w, err, 500)
			return
		}
		var result []Payment
		for _, p := range payments {
			if isExpired(p) {
				handler.PaymentRepository.DataBase.
					Where("id = ?", p.Id).
					Delete(&Payment{})
				continue
			}
			result = append(result, p)
		}

		res.Json(w, result, 200)

	}
}
func (handler *PaymentHandler) getAllPayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getPaymentRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, err, 401)
			return
		}
		var payments []Payment
		err = handler.PaymentRepository.DataBase.Find(&payments).Error
		if err != nil {
			res.Json(w, err, 500)
			return
		}
		var result []Payment
		for _, p := range payments {
			if isExpired(p) {
				handler.PaymentRepository.DataBase.
					Where("id = ?", p.Id).
					Delete(&Payment{})
				continue
			}
			result = append(result, p)
		}

		res.Json(w, result, 200)

	}
}
func (handler *PaymentHandler) createPayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createPaymentRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		now := time.Now()
		tx := handler.PaymentRepository.DataBase.Begin()
		if tx.Error != nil {
			res.Json(w, tx.Error, 500)
			return
		}
		var payments []Payment
		if err := tx.
			Set("gorm:query_option", "FOR UPDATE").
			Where("is_working = ?", true).
			Find(&payments).Error; err != nil {
			tx.Rollback()
			res.Json(w, err, 500)
			return
		}
		for _, p := range payments {
			if p.UserId == body.UserId {
				tx.Rollback()
				res.Json(w, p, 409)
				return
			}
			if isExpired(p) {
				tx.Where("id = ?", p.Id).Delete(&Payment{})
				continue
			}
			if p.Price == body.Price {
				tx.Rollback()
				res.Json(w, map[string]string{
					"error": "payment is busy",
				}, 409)
				return
			}
		}
		payment := Payment{
			Id:        token.CreateId(),
			Year:      now.Year(),
			Month:     int(now.Month()),
			Day:       now.Day(),
			Hour:      now.Hour(),
			Minute:    now.Minute(),
			UserId:    body.UserId,
			Price:     body.Price,
			IsWorking: true,
		}
		if err := tx.Create(&payment).Error; err != nil {
			tx.Rollback()
			res.Json(w, err, 500)
			return
		}
		tx.Commit()
		res.Json(w, payment, 200)
	}
}
func (handler *PaymentHandler) updatePayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[updatePaymentRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, err, 401)
			return
		}
		var payment Payment
		err = handler.PaymentRepository.DataBase.
			Where("id = ?", body.Id).
			First(&payment).Error
		if err != nil {
			res.Json(w, err, 404)
			return
		}
		payment.IsWorking = body.IsWorking
		err = handler.PaymentRepository.DataBase.Save(&payment).Error
		if err != nil {
			res.Json(w, err, 500)
			return
		}
		now := time.Now()
		tx := Transaction{
			Id:        token.CreateId(),
			UserId:    body.UserId,
			Price:     payment.Price,
			Year:      now.Year(),
			Month:     int(now.Month()),
			Day:       now.Day(),
			Hour:      now.Hour(),
			Minute:    now.Minute(),
			GameName:  "-",
			DonatName: "-",
			CreatedBy: "admin",
		}
		var user User
		err = handler.PaymentRepository.DataBase.
			Where("token = ?", body.UserId).
			First(&user).Error
		if err != nil {
			res.Json(w, err, 404)
			return
		}

		err = handler.PaymentRepository.DataBase.
			Model(&User{}).
			Where("token = ?", body.UserId).
			Update("balance", gorm.Expr("balance + ?", payment.Price)).
			Error

		if err != nil {
			res.Json(w, err, 500)
			return
		}
		if err := handler.PaymentRepository.DataBase.Create(&tx).Error; err != nil {
			res.Json(w, err, 500)
			return
		}
		res.Json(w, tx, 200)
	}
}
func (handler *PaymentHandler) deletePayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[deletePaymentRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		_, err = handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, err, 401)
			return
		}
		err = handler.PaymentRepository.DataBase.
			Where("id = ?", body.Id).
			Delete(&Payment{}).Error
		if err != nil {
			res.Json(w, err, 500)
			return
		}
		res.Json(w, map[string]string{
			"status": "deleted",
		}, 200)
	}
}
func (handler *PaymentHandler) createTelegram() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createPaymentTelegram](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		fmt.Println("body", body)
		res.Json(w, body.Text, 200)
	}
} //ishlamiydi hoizr
