package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Config interface {
	GetToken() string
	GetUrl() string
}

type ScenarioProcessor interface {
	Process(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error)
}

type ScenarioProcessorV2 interface {
	Process(ctx context.Context, msg *tgbotapi.Message) ([]*tgbotapi.MessageConfig, error)
}
