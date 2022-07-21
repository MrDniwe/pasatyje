package tg

import (
	"github.com/google/logger"
	"github.com/mrdniwe/pasatyje/pkg/intf/tg"
	"github.com/spf13/viper"
)

type botCfg struct {
	token string
	url   string
}

func (b *botCfg) GetToken() string {
	if b.token == "" {
		logger.Fatal("token is empty")
	}
	return b.token
}

func (b *botCfg) GetUrl() string {
	if b.token == "" {
		logger.Fatal("url is empty")
	}
	return b.token
}

// New creates Telegram bot config object
//
// environment vars:
// TOKEN - string representing the tg bot token
// URL - tg bot URL
func New(v *viper.Viper, prefix string) (tg.Config, error) {
	v.SetDefault(prefix+"token", "")
	v.BindEnv(prefix+"token", prefix+"TOKEN")
	v.SetDefault(prefix+"url", "")
	v.BindEnv(prefix+"url", prefix+"URL")
	return &botCfg{
		token: v.GetString(prefix + "token"),
		url:   v.GetString(prefix + "url"),
	}, nil
}
