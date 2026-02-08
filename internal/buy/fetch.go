package buy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type BulkProvider struct {
	ApiURL string
	ApiKey string
}
type Provider interface {
	CreateOrder(service int, playerId string) (order string, err error)
}

func (b *BulkProvider) CreateOrder(service int, link string) (string, error) {

	payload := map[string]interface{}{
		"key":      b.ApiKey,
		"action":   "add",
		"service":  service,
		"link":     link,
		"quantity": 1,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("err", err)
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		b.ApiURL,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		fmt.Println("err", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("err", err)
		return "", fmt.Errorf("bulk bad status: %d, body: %s", resp.StatusCode, body)
	}

	var result CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("err", err)
		return "", err
	}

	// ❌ Ошибка от провайдера
	if result.Error != "" {
		fmt.Println("Ошибка от провайдера", result.Error)
		return "", errors.New(result.Error)
	}

	// ❌ Провайдер не вернул order
	if result.Order == "" {
		fmt.Println("Провайдер не вернул order", err)
		return "", errors.New("bulk did not return order id")
	}
	return result.Order.String(), nil
}

func (b *BulkProvider) OrderStatus(order string) (*OrderStatusResponse, error) {

	orderInt, err := strconv.Atoi(order)
	if err != nil {
		return nil, fmt.Errorf("invalid order id: %s", order)
	}

	payload := map[string]interface{}{
		"key":    b.ApiKey,
		"action": "status",
		"order":  orderInt, // ⚠️ ВАЖНО: int
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		b.ApiURL,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bulk bad status: %d, body: %s", resp.StatusCode, body)
	}

	var result OrderStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status == "" {
		return nil, fmt.Errorf("bulk status empty, response: %+v", result)
	}

	return &result, nil
}
