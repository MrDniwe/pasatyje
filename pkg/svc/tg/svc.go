package tg

import (
	"context"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mrdniwe/pasatyje/pkg/intf/app"
	"github.com/mrdniwe/pasatyje/pkg/intf/tg"
	"github.com/sirupsen/logrus"
)

const (
	messageWorkers     = 3
	messageChanBufsize = 3
	updateTimeout      = 60
)

type TgBotServer interface {
	Run()
}

type bot struct {
	logger  *logrus.Logger
	botAPI  *api.BotAPI
	sp      tg.ScenarioProcessorV2
	msgChan chan *api.Message
	botCfg  tg.Config
	appCfg  app.Config
}

func New(
	botCfg tg.Config,
	appCfg app.Config,
	logger *logrus.Logger,
	sp tg.ScenarioProcessorV2,
) TgBotServer {
	b, err := api.NewBotAPI(botCfg.GetToken())
	if err != nil {
		log.Fatalf("TgBot unable to start: %v", err)
	}
	if appCfg.GetEnvironment() != "production" {
		b.Debug = true
	}
	log.Printf("Authorized on account %s", b.Self.UserName)
	return &bot{
		botAPI:  b,
		logger:  logger,
		sp:      sp,
		msgChan: make(chan *api.Message, messageChanBufsize),
		botCfg:  botCfg,
		appCfg:  appCfg,
	}
}

func (b *bot) Run() {
	go b.listen()
	for i := 0; i < messageWorkers; i++ {
		go b.runMsgProcWorker(i)
	}
}

func (b *bot) listen() {

	u := api.NewUpdate(0)
	u.Timeout = updateTimeout

	updates, err := b.botAPI.GetUpdatesChan(u)
	if err != nil {
		b.logger.Fatal(err)
	}
	b.logger.Info("Telegram bot listener started")
	for update := range updates {
		// TODO: handle all types of updates
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		b.logger.Debugf("Recieved update from '%s', message: %s", update.Message.From.UserName, update.Message.Text)
		b.msgChan <- update.Message
	}
}

func (b *bot) runMsgProcWorker(index int) {
	b.logger.Infof("Telegram bot message processing worker #%d started", index)
	for msg := range b.msgChan {

		// TODO: проверки IsChannel/IsBot/IsPrivate/IsGroup итд в привязке к конфигу
		if msg.From.IsBot {
			b.botAPI.Send(api.NewMessage(
				msg.Chat.ID,
				"Бот не общается с ботами",
			))
			continue
		}
		if !msg.Chat.IsPrivate() {
			continue
		}
		if b.sp == nil {
			b.logger.Errorf("Unable to handle TG command: no scenario procesor")
			b.botAPI.Send(api.NewMessage(msg.Chat.ID, "Ошибка на стороне сервера: нет обработчика для сценариев"))
			continue
		}
		resp, err := b.sp.Process(context.Background(), msg)
		if err != nil {
			b.logger.Errorf("Error processing message: %v", err)
			b.botAPI.Send(api.NewMessage(msg.Chat.ID, err.Error()))
			continue
		}
		for _, r := range resp {
			if resp == nil {
				continue
			}
			b.botAPI.Send(r)
		}
	}
}
