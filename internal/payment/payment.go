package payment

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"net/http"
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
	router.HandleFunc("/payment/createTelegram", handler.createTelegram())
	router.HandleFunc("/payment/getAllPayment", handler.getAllPayment())
	router.HandleFunc("/payment/createPayment", handler.createPayment())
	router.HandleFunc("/payment/updatePayment", handler.updatePayment())
	router.HandleFunc("/payment/deletePayment", handler.deletePayment())
	return handler
}
func (handler *PaymentHandler) createPayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createPaymentRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		payment := Payment{
			Id:        token.CreateId(),
			Year:      body.Year,
			Month:     body.Month,
			Day:       body.Day,
			Hour:      body.Hour,
			Minute:    body.Minute,
			UserId:    body.UserId,
			Price:     body.Price,
			IsWorking: body.IsWorking,
		}
		handler.PaymentRepository.DataBase.Create(&payment)
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
		res.Json(w, payment, 200)
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

		res.Json(w, payments, 200)
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

		res.Json(w, payments, 200)
	}
}
func (handler *PaymentHandler) createTelegram() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createPaymentRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		payment := Payment{
			Id:        token.CreateId(),
			Year:      body.Year,
			Month:     body.Month,
			Day:       body.Day,
			Hour:      body.Hour,
			Minute:    body.Minute,
			UserId:    body.UserId,
			Price:     body.Price,
			IsWorking: body.IsWorking,
		}
		handler.PaymentRepository.DataBase.Create(&payment)
		res.Json(w, payment, 200)
	}
} //ishlamiydi hoizr
