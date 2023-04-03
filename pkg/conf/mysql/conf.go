package db

import (
	"fmt"

	"github.com/mrdniwe/pasatyje/pkg/intf/conf/mysql"
	"github.com/spf13/viper"
)

type mysqlCfg struct {
	username string
	password string
	hostname string
	port     int
	database string
}

func (m *mysqlCfg) GetDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.username, m.password, m.hostname, m.port, m.database)
	return dsn
}

// New creates a MySQL database config object
//
// environment vars:
// MYSQL_USERNAME - string representing the username to connect to the MySQL database
// MYSQL_PASSWORD - string representing the password to connect to the MySQL database
// MYSQL_HOSTNAME - string representing the hostname of the MySQL server
// MYSQL_PORT - int representing the port number of the MySQL server
// MYSQL_DATABASE - string representing the name of the MySQL database
func New(v *viper.Viper, prefix string) (mysql.Config, error) {
	v.SetDefault(prefix+"username", "")
	v.BindEnv(prefix+"username", "MYSQL_USERNAME")
	v.SetDefault(prefix+"password", "")
	v.BindEnv(prefix+"password", "MYSQL_PASSWORD")
	v.SetDefault(prefix+"hostname", "")
	v.BindEnv(prefix+"hostname", "MYSQL_HOSTNAME")
	v.SetDefault(prefix+"port", 3306)
	v.BindEnv(prefix+"port", "MYSQL_PORT")
	v.SetDefault(prefix+"database", "")
	v.BindEnv(prefix+"database", "MYSQL_DATABASE")

	cfg := &mysqlCfg{
		username: v.GetString(prefix + "username"),
		password: v.GetString(prefix + "password"),
		hostname: v.GetString(prefix + "hostname"),
		port:     v.GetInt(prefix + "port"),
		database: v.GetString(prefix + "database"),
	}

	return cfg, nil
}
