package app

import (
	"github.com/google/logger"
	"github.com/mrdniwe/pasatyje/pkg/intf/app"
	"github.com/spf13/viper"
)

type appCfg struct {
	env     app.Environment
	logType string // LOG_TYPE -  TEXT* | JSON
}

func (a *appCfg) GetEnvironment() app.Environment {
	if a.env == "" {
		logger.Fatal("no environment config specified")
	}
	return a.env
}

func (a *appCfg) GetLogType() string {
	if a.logType == "" {
		logger.Fatal("log type is empty")
	}
	return a.logType
}

// New creates App config object
//
// environment vars:
// APP_LOG_TYPE - TEXT* | JSON
// APP_ENV - development* | test | staging | production
func New(v *viper.Viper) (app.Config, error) {
	v.SetDefault("log_type", "TEXT")
	v.BindEnv("log_type", "APP_LOG_TYPE")
	v.SetDefault("env", "development")
	v.BindEnv("env", "APP_ENV")
	return &appCfg{
		env:     app.Environment(v.GetString("env")),
		logType: v.GetString("log_type"),
	}, nil
}
