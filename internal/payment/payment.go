package payment

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"fmt"
	"net/http"
	"time"
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
				err := handler.PaymentRepository.DataBase.
					Model(&Payment{}).
					Where("id = ?", p.Id).
					Update("is_working", false).Error
				if err != nil {
					res.Json(w, err, 500)
					return
				}
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
				err := handler.PaymentRepository.DataBase.
					Model(&Payment{}).
					Where("id = ?", p.Id).
					Update("is_working", false).Error
				if err != nil {
					res.Json(w, err, 500)
					return
				}
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
				err := handler.PaymentRepository.DataBase.
					Model(&Payment{}).
					Where("id = ?", p.Id).
					Update("is_working", false).Error
				if err != nil {
					res.Json(w, err, 500)
					return
				}
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
				handler.PaymentRepository.DataBase.
					Model(&Payment{}).
					Where("id = ?", p.Id).
					Update("is_working", false)
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
		_, err = handler.AuthHandler.UpdateBalance(body.UserId, payment.Price)

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
} //works transactions
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

		telegramTime := time.Date(
			body.Year,
			time.Month(body.Month),
			body.Day,
			body.Hour,
			body.Minute,
			0, 0, time.Local,
		)

		fromTime := telegramTime.Add(-48 * time.Hour)
		toTime := telegramTime.Add(5 * time.Minute)

		fromY, fromM, fromD := fromTime.Year(), int(fromTime.Month()), fromTime.Day()
		toY, toM, toD := toTime.Year(), int(toTime.Month()), toTime.Day()
		fromH, fromMin := fromTime.Hour(), fromTime.Minute()
		toH, toMin := toTime.Hour(), toTime.Minute()

		var payment Payment
		err = handler.PaymentRepository.DataBase.
			Where("price = ?", body.Amount).
			Where(`(
				(year > ?) OR 
				(year = ? AND month > ?) OR 
				(year = ? AND month = ? AND day > ?) OR 
				(year = ? AND month = ? AND day = ? AND hour > ?) OR 
				(year = ? AND month = ? AND day = ? AND hour = ? AND minute >= ?)
			) AND (
				(year < ?) OR 
				(year = ? AND month < ?) OR 
				(year = ? AND month = ? AND day < ?) OR 
				(year = ? AND month = ? AND day = ? AND hour < ?) OR 
				(year = ? AND month = ? AND day = ? AND hour = ? AND minute <= ?)
			)`,
				fromY, fromY, fromM, fromY, fromM, fromD, fromY, fromM, fromD, fromH, fromY, fromM, fromD, fromH, fromMin,
				toY, toY, toM, toY, toM, toD, toY, toM, toD, toH, toY, toM, toD, toH, toMin,
			).
			Order("year DESC, month DESC, day DESC, hour DESC, minute DESC").
			First(&payment).Error

		if err != nil {
			res.Json(w, map[string]string{"error": "no matching reservation found"}, 404)
			return
		}

		// Проверяем, не обработана ли уже бронь
		if !payment.IsWorking {
			// Можно дополнительно проверить, есть ли уже транзакция по этой брони
			var count int64
			handler.PaymentRepository.DataBase.Model(&Transaction{}).
				Where("payment_id = ?", payment.Id).
				Count(&count)
			if count > 0 {
				res.Json(w, map[string]string{"error": "reservation already processed"}, 409)
				return
			}
			// Если is_working = false, но транзакции нет — возможно админ вручную закрыл, но не начислил → продолжаем
		}

		// Всё в транзакции БД
		txDb := handler.PaymentRepository.DataBase.Begin()
		if txDb.Error != nil {
			res.Json(w, txDb.Error, 500)
			return
		}

		// Создаём транзакцию
		transaction := Transaction{
			Id:        token.CreateId(),
			UserId:    payment.UserId,
			Price:     body.Amount,
			Year:      body.Year,
			Month:     body.Month,
			Day:       body.Day,
			Hour:      body.Hour,
			Minute:    body.Minute,
			GameName:  "-",
			DonatName: "-",
			CreatedBy: body.Sender,
		}

		if err := txDb.Create(&transaction).Error; err != nil {
			txDb.Rollback()
			res.Json(w, err, 500)
			return
		}

		// Начисляем баланс
		_, err = handler.AuthHandler.UpdateBalance(payment.UserId, payment.Price)
		if err != nil {
			txDb.Rollback()
			res.Json(w, err, 500)
			return
		}

		// Закрываем бронь
		if err := txDb.Model(&Payment{}).
			Where("id = ?", payment.Id).
			Update("is_working", false).Error; err != nil {
			txDb.Rollback()
			res.Json(w, err, 500)
			return
		}

		txDb.Commit()

		fmt.Println("Processed payment:", body.Amount, "for user:", payment.UserId, "at", telegramTime)
		res.Json(w, map[string]string{
			"status":         "success",
			"user_id":        payment.UserId,
			"amount":         fmt.Sprintf("%d", body.Amount),
			"transaction_id": transaction.Id,
		}, 200)
	}
} //works transactions. not finished
