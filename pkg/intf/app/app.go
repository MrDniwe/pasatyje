package app

type Environment string

const (
	DEV  Environment = "development"
	TEST Environment = "test"
	SGT  Environment = "staging"
	PROD Environment = "production"
)

type Config interface {
	GetEnvironment() Environment
	GetLogType() string
}
