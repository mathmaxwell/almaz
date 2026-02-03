package transactions

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
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
			Order: "-",
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
