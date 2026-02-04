package buy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type BulkProvider struct {
	ApiURL string
	ApiKey string
}
type Provider interface {
	CreateOrder(service int, playerId string) (order string, err error)
}

func (b *BulkProvider) CreateOrder(service int, playerId string) (string, error) {
	payload := map[string]interface{}{
		"key":      b.ApiKey,
		"action":   "add",
		"service":  service,
		"link":     playerId,
		"quantity": 1,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		b.ApiURL,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("bulk bad status: %d", resp.StatusCode)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return "", err
	}

	order := fmt.Sprintf("%v", raw["order"])
	if order == "" {
		return "", errors.New("bulk order empty")
	}

	return order, nil
}
