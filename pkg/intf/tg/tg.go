package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Config interface {
	GetToken() string
	GetUrl() string
}

type MsgProc interface {
	Process(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error)
}
