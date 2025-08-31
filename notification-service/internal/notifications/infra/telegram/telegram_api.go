package telegram

import (
	"context"
	"fmt"
	"os"
)

type TGApi struct {
	token string
}

func NewTGApiClient() *TGApi {
	return &TGApi{
		token: os.Getenv("TELEGRAM_TOKEN"),
	}
}

func (t *TGApi) Send(ctx context.Context, payload []byte) error {
	fmt.Println("telegram sended")
	return nil
}
