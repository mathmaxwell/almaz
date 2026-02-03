package buy

import (
	"bytes"
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func NewBuyRepository(dataBase *db.Db) *BuyRepository {
	return &BuyRepository{
		DataBase: dataBase,
	}
}
func NewGamesHandler(router *http.ServeMux, deps *BuyhandlerDeps) *BuyHandler {
	handler := &BuyHandler{
		Config:        deps.Config,
		BuyRepository: *deps.BuyRepository,
		AuthHandler:   deps.AuthHandler,
	}
	router.HandleFunc("/buy/create", handler.create())
	return handler
}

func (handler *BuyHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createBuyRequest](&w, r)
		if err != nil {
			res.Json(w, "bad request", 400)
			return
		}
		var offers Offers
		if err := handler.BuyRepository.DataBase.
			Where("bot_id = ?", body.BotId).
			First(&offers).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "offer not found", 404)
				return
			}
			res.Json(w, "db error", 500)
			return
		}

		offerPrice, err := strconv.Atoi(offers.Price)
		if err != nil {
			res.Json(w, "price error", 500)
			return
		}
		var game Games
		if err := handler.BuyRepository.DataBase.
			Where("id = ?", body.GameId).
			First(&game).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "game not found", 404)
				return
			}
			res.Json(w, "failed to get game", 500)
			return
		}
		now := time.Now()
		var txId string
		err = handler.BuyRepository.DataBase.Transaction(func(tx *gorm.DB) error {
			if err := handler.AuthHandler.DecreaseBalance(body.Token, offerPrice); err != nil {
				return err
			}
			txId = token.CreateId()
			transaction := Transaction{
				Id:        txId,
				UserId:    body.Token,
				Price:     -offerPrice,
				Year:      now.Year(),
				Month:     int(now.Month()),
				Day:       now.Day(),
				Hour:      now.Hour(),
				Minute:    now.Minute(),
				GameName:  game.Name,
				DonatName: offers.UzName,
				CreatedBy: body.BotId,
				Order:     "gameOrder",
			}
			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}
			newBuy := Buy{
				Id:       token.CreateId(),
				Year:     now.Year(),
				Month:    int(now.Month()),
				Day:      now.Day(),
				Hour:     now.Hour(),
				Minute:   now.Minute(),
				UserId:   body.Token,
				Status:   "wait",
				BotId:    body.BotId,
				ServerId: body.ServerId,
				PlayerId: body.PlayerId,
				GameId:   body.GameId,
				Price:    offerPrice,
			}
			if err := tx.Create(&newBuy).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			if err.Error() == "not enough balance" {
				res.Json(w, "not enough balance", 400)
				return
			}
			res.Json(w, "transaction failed", 500)
			return
		}
		botIdNumber, err := strconv.Atoi(body.BotId)
		if err != nil {
			res.Json(w, "bot id error", 500)
			return
		}
		payload := map[string]interface{}{
			"key":      os.Getenv("BULKKEY"),
			"action":   "add",
			"service":  botIdNumber,
			"link":     body.PlayerId,
			"quantity": 1,
		}
		jsonBody, _ := json.Marshal(payload)
		reqBulk, err := http.NewRequest(
			http.MethodPost,
			os.Getenv("BULKAPI"),
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			res.Json(w, "bulk request error", 500)
			return
		}
		reqBulk.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(reqBulk)
		if err != nil {
			res.Json(w, "bulk api error", 500)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			res.Json(w, "bulk api bad status", 500)
			return
		}
		var raw map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			res.Json(w, "bulk response error", 500)
			return
		}
		order := fmt.Sprintf("%v", raw["order"])
		handler.BuyRepository.DataBase.
			Model(&Transaction{}).
			Where("id = ?", txId).
			Update("order", order)
		res.Json(w, map[string]string{
			"status": "ok",
			"order":  order,
		}, 200)
	}
}
