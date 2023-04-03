package mysql

type Config interface {
	GetDSN() string
}
