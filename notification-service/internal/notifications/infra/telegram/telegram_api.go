package telegram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"os"
	"time"
)

type TGApi struct {
	token  string
	client *http.Client
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type sendMessageResponse struct {
	Ok          bool            `json:"ok"`
	Description string          `json:"description,omitempty"`
	Result      json.RawMessage `json:"result,omitempty"`
}

func NewTGApiClient() *TGApi {
	return &TGApi{
		token: os.Getenv("TELEGRAM_TOKEN"),
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (t *TGApi) Send(ctx context.Context, to string, message string, parseMode string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	reqBody, err := json.Marshal(sendMessageRequest{
		ChatID:    to,
		Text:      message,
		ParseMode: parseMode,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram api request: %w", err)
	}
	defer resp.Body.Close()

	var res sendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if !res.Ok {
		return fmt.Errorf("telegram api error: %s, %v", res.Description, res)
	}

	return nil
}
