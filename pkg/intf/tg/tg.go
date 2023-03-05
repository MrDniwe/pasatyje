package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config interface {
	GetToken() string
	GetUrl() string
}

type ScenarioProcessor interface {
	Process(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error)
}

type ScenarioProcessorV2 interface {
	Process(ctx context.Context, msg *tgbotapi.Message) ([]tgbotapi.Chattable, error)
	HandleCallback(ctx context.Context, query *tgbotapi.CallbackQuery) (*tgbotapi.CallbackConfig, []tgbotapi.Chattable, error)
}
