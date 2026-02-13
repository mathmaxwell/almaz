package buy

import (
	"demo/almaz/pkg/db"
	"demo/almaz/pkg/req"
	"demo/almaz/pkg/res"
	"demo/almaz/pkg/token"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const USD_TO_SUM = 12500

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
	router.HandleFunc("/buy/orderStatus", handler.orderStatus())
	return handler
}

func (handler *BuyHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createBuyRequest](&w, r)
		if err != nil {
			res.Json(w, "bad request", 400)
			return
		}

		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "пользователь не найден", 404)
			return
		}

		var offer Offers
		if err := handler.BuyRepository.DataBase.
			Where("id = ?", body.OfferId).
			First(&offer).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "offer не найдено", 404)
				return
			}
			res.Json(w, "ошибка базы данных", 500)
			return
		}

		var offerPrice int
		if user.UserRole == "superUser" {
			offerPrice, err = strconv.Atoi(offer.SuperPrice)
		} else {
			offerPrice, err = strconv.Atoi(offer.Price)
		}
		if err != nil || offerPrice <= 0 {
			res.Json(w, "некорректная цена", 400)
			return
		}

		var game Games
		if err := handler.BuyRepository.DataBase.
			Where("id = ?", body.GameId).
			First(&game).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "игра не найдена", 404)
				return
			}
			res.Json(w, "ошибка при получении игры", 500)
			return
		}

		botIdNumber, err := strconv.Atoi(body.BotId)
		if err != nil {
			res.Json(w, "некорректный bot id", 400)
			return
		}

		provider := &BulkProvider{
			ApiURL: os.Getenv("BULKAPI"),
			ApiKey: os.Getenv("BULKKEY"),
		}

		order := "empty"

		if game.Name == "Mobile Legends Global" || game.Name == "Mobile Legends" || game.Name == "PUBG Mobile" || game.Name == "Freefire Global" {
			link := body.PlayerId
			if game.Description == "two" {
				if body.ServerId == "" {
					res.Json(w, "не указан server id", 400)
					return
				}
				link = body.PlayerId + "|" + body.ServerId
			}

			// 🔹 Проверка баланса провайдера в сум
			balanceStr, _, err := provider.GetBalance()
			if err != nil {
				res.Json(w, "не удалось получить баланс провайдера", 500)
				return
			}

			providerBalanceUSD, err := strconv.ParseFloat(balanceStr, 64)
			if err != nil {
				res.Json(w, "ошибка обработки баланса провайдера", 500)
				return
			}

			providerBalanceSom := providerBalanceUSD * USD_TO_SUM
			offerPriceFloat := float64(offerPrice)

			if providerBalanceSom < offerPriceFloat {
				res.Json(w, "недостаточно средств у провайдера", 400)
				return
			}

			if providerBalanceSom-offerPriceFloat < 100000 {
				res.Json(w, "баланс провайдера ниже допустимого порога", 400)
				return
			}

			// 🔹 Транзакция: списание + создание заказа + сохранение
			now := time.Now()
			var txId string
			err = handler.BuyRepository.DataBase.Transaction(func(tx *gorm.DB) error {
				if err := handler.AuthHandler.DecreaseBalance(tx, body.Token, offerPrice); err != nil {
					return err
				}

				orderId, err := provider.CreateOrder(botIdNumber, link)
				if err != nil {
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
					DonatName: offer.UzName,
					CreatedBy: body.BotId,
					Order:     orderId,
				}

				if err := tx.Create(&transaction).Error; err != nil {
					return err
				}

				order = orderId
				return nil
			})

			if err != nil {
				if err.Error() == "недостаточно средств" {
					res.Json(w, map[string]string{
						"error": "Недостаточно средств на балансе",
					}, 400)
					return
				}
				if err.Error() == "пользователь не найден" {
					res.Json(w, map[string]string{
						"error": "Пользователь не найден",
					}, 404)
					return
				}
				res.Json(w, "ошибка провайдера или базы данных", 500)
				return
			}
		}

		res.Json(w, map[string]string{
			"order":   order,
			"message": "Покупка успешно выполнена",
		}, 200)
	}
}

func (handler *BuyHandler) orderStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[OrderStatusRequest](&w, r)
		if err != nil {
			res.Json(w, "bad request", http.StatusBadRequest)
			return
		}

		var game Games
		if err := handler.BuyRepository.DataBase.
			Where("id = ?", body.GameId).
			First(&game).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "игра не найдена", http.StatusNotFound)
				return
			}

			res.Json(w, "ошибка при получении игры", http.StatusInternalServerError)
			return
		}

		switch game.Name {
		case "Mobile Legends Global", "Mobile Legends":
			provider := &BulkProvider{
				ApiURL: os.Getenv("BULKAPI"),
				ApiKey: os.Getenv("BULKKEY"),
			}

			status, err := provider.OrderStatus(body.Order)
			if err != nil {
				res.Json(w, "ошибка получения статуса заказа", http.StatusBadGateway)
				return
			}

			res.Json(w, status, http.StatusOK)
			return

		default:
			res.Json(w, "игра не поддерживает проверку статуса", http.StatusBadRequest)
			return
		}
	}
}
