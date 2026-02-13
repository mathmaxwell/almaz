package transactions

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"fmt"
	"net/http"
	"time"
)

func NewTransactionRepository(dataBase *db.Db) *TransactionRepository {
	return &TransactionRepository{
		DataBase: dataBase,
	}
}
func NewTranactionHandler(router *http.ServeMux, deps *TransactionhandlerDeps) *TransactionHandler {
	handler := &TransactionHandler{
		Config:                deps.Config,
		TransactionRepository: *deps.TransactionRepository,
		AuthHandler:           deps.AuthHandler,
	}
	router.HandleFunc("/transactions/transactionsCreate", handler.create())
	router.HandleFunc("/transactions/transactionsGet", handler.getAll())
	router.HandleFunc("/transactions/transactionsDelete", handler.delete())
	router.HandleFunc("/transactions/getByPeriod", handler.getByPeriod())
	return handler
}

func (handler *TransactionHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		now := time.Now()
		tx := Transaction{
			Id:        token.CreateId(),
			UserId:    body.UserId,
			Price:     body.Price,
			Year:      now.Year(),
			Month:     int(now.Month()),
			Day:       now.Day(),
			Hour:      now.Hour(),
			Minute:    now.Minute(),
			GameName:  body.GameName,
			DonatName: body.DonatName,
			CreatedBy: body.CreatedBy,
			Order:     "-",
		}
		if err := handler.TransactionRepository.DataBase.Create(&tx).Error; err != nil {
			res.Json(w, err, 500)
			return
		}
		_, err = handler.AuthHandler.UpdateBalance(body.UserId, body.Price)
		if err != nil {
			res.Json(w, err, 500)
			return
		}
		res.Json(w, tx, 200)
	}
}
func (handler *TransactionHandler) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}

		txs := make([]Transaction, 0)
		if err := handler.TransactionRepository.DataBase.
			Where("user_id = ?", body.UserId).
			Find(&txs).Error; err != nil {
			res.Json(w, err, 500)
			return
		}
		res.Json(w, txs, 200)
	}
}
func (handler *TransactionHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[deleteRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		if err := handler.TransactionRepository.DataBase.
			Where("id = ?", body.Id).
			Delete(&Transaction{}).Error; err != nil {
			res.Json(w, err, 500)
			return
		}
		res.Json(w, map[string]string{"status": "deleted"}, 200)
	}
}

func (handler *TransactionHandler) getByPeriod() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getByPeriodRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		if body.Token != handler.Config.Token.AdminToken {
			res.Json(w, "you are not admin", 403)
			return
		}
		if body.StartYear < 2000 || body.EndYear < 2000 ||
			body.StartMonth < 1 || body.StartMonth > 12 ||
			body.EndMonth < 1 || body.EndMonth > 12 ||
			body.StartDay < 1 || body.StartDay > 31 ||
			body.EndDay < 1 || body.EndDay > 31 {
			res.Json(w, map[string]string{"error": "некорректные значения дат"}, 400)
			return
		}

		startYMD := fmt.Sprintf("%04d%02d%02d", body.StartYear, body.StartMonth, body.StartDay)
		endYMD := fmt.Sprintf("%04d%02d%02d", body.EndYear, body.EndMonth, body.EndDay)

		if startYMD > endYMD {
			res.Json(w, map[string]string{"error": "start дата позже end даты"}, 400)
			return
		}

		var txs []Transaction

		err = handler.TransactionRepository.DataBase.Where("CONCAT(year::text, LPAD(month::text, 2, '0'), LPAD(day::text, 2, '0')) >= ?", startYMD).
			Where("CONCAT(year::text, LPAD(month::text, 2, '0'), LPAD(day::text, 2, '0')) <= ?", endYMD).
			Order("year DESC, month DESC, day DESC, hour DESC, minute DESC").
			Find(&txs).Error

		if err != nil {
			res.Json(w, err, 500)
			return
		}

		res.Json(w, txs, 200)
	}
}
