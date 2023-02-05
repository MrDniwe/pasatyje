package tg

import (
	"context"
	"fmt"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrdniwe/pasatyje/pkg/intf/app"
	"github.com/mrdniwe/pasatyje/pkg/intf/tg"
	"github.com/mrdniwe/pasatyje/pkg/lib/uerror"
	"github.com/sirupsen/logrus"
)

const (
	MESSAGE_WORKERS       = 3
	CALLBACK_WORKERS      = 3
	MSG_CHAN_BUFSIZE      = 3
	CALLBACK_CHAN_BUFSIZE = 3
	UPD_TIMEOUT           = 60
)

type TgBotServer interface {
	Run()
}

type bot struct {
	logger  *logrus.Logger
	botAPI  *api.BotAPI
	sp      tg.ScenarioProcessorV2
	msgChan chan *api.Message
	cbChan  chan *api.CallbackQuery
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
		msgChan: make(chan *api.Message, MSG_CHAN_BUFSIZE),
		cbChan:  make(chan *api.CallbackQuery, CALLBACK_CHAN_BUFSIZE),
		botCfg:  botCfg,
		appCfg:  appCfg,
	}
}

func (b *bot) Run() {
	go b.listen()
	for i := 0; i < MESSAGE_WORKERS; i++ {
		go b.runMsgProcWorker(i)
	}
	for i := 0; i < CALLBACK_WORKERS; i++ {
		go b.runCallbackWorker(i)
	}
}

func (b *bot) listen() {

	u := api.NewUpdate(0)
	u.Timeout = UPD_TIMEOUT

	updates := b.botAPI.GetUpdatesChan(u)
	b.logger.Info("Telegram bot listener started")
	for update := range updates {
		// TODO: handle all types of updates
		if update.Message != nil { // ignore any non-Message Updates
			b.logger.Debugf("Recieved update from '%s', message: %s", update.Message.From.UserName, update.Message.Text)
			b.msgChan <- update.Message
			continue
		}
		if update.CallbackQuery != nil {
			b.logger.Debugf("Recieved update from '%s', callback query: %+v", update.Message.From.UserName, update.CallbackQuery)
			b.cbChan <- update.CallbackQuery
		}
	}
}

func (b *bot) runCallbackWorker(index int) {
	b.logger.Infof("Telegram bot callback worker #%d started", index)
	for upd := range b.cbChan {
		if b.sp == nil {
			if upd.Message == nil {
				continue
			}
			b.handleMsgError(upd.Message, uerror.ServerError.New("нет обработчика для калбэков"))
			continue
		}
		// TODO: real context
		resp, err := b.sp.HandleCallback(context.Background(), upd)
		if err != nil {
			b.handleMsgError(upd.Message, err)
			continue
		}
		for _, r := range resp {
			if r == nil {
				continue
			}
			b.botAPI.Send(r)
		}
	}
}

func (b *bot) runMsgProcWorker(index int) {
	b.logger.Infof("Telegram bot message processing worker #%d started", index)
	for msg := range b.msgChan {

		// TODO: проверки IsChannel/IsBot/IsPrivate/IsGroup итд в привязке к конфигу
		if msg.From.IsBot {
			b.handleMsgError(msg, uerror.Forbidden.New("бот не общается с ботами"))
			continue
		}
		if !msg.Chat.IsPrivate() {
			continue
		}
		if b.sp == nil {
			b.handleMsgError(msg, uerror.ServerError.New("нет обработчика для сценариев"))
			continue
		}
		// TODO: real context
		resp, err := b.sp.Process(context.Background(), msg)
		if err != nil {
			b.handleMsgError(msg, err)
			continue
		}
		for _, r := range resp {
			if r == nil {
				continue
			}
			b.botAPI.Send(r)
		}
	}
}

func (b *bot) handleMsgError(msg *api.Message, err error) {
	if msg == nil {
		b.handleChatError(0, err)
		return
	}
	b.handleChatError(msg.Chat.ID, err)
}

func (b *bot) handleChatError(chatID int64, err error) {
	var errPattern string
	switch uerror.GetType(err) {
	case uerror.BadRequest:
		b.logger.Warn(err.Error())
		errPattern = "Ошибка в запросе"
	case uerror.ServerError:
		b.logger.Error(err.Error())
		errPattern = "Ошибка на стороне сервера"
	case uerror.NotFound:
		b.logger.Warn(err.Error())
		errPattern = "Не удалось найти"
	case uerror.NotAuthorized:
		b.logger.Warn(err.Error())
		errPattern = "Не авторизован"
	case uerror.Forbidden:
		b.logger.Warn(err.Error())
		errPattern = "Данная функциональность для вас под запретом"
	default:
		b.logger.Error(err.Error())
		errPattern = "Неизвестная ошибка"
	}
	if chatID != 0 {
		b.botAPI.Send(api.NewMessage(chatID, fmt.Sprintf(errPattern+": %s", uerror.GetLastError(err))))
	}
}
