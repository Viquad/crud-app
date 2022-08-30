package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ConnectionInfo struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	DBName     string `mapstructure:"db_name"`
	SSLMode    string `mapstructure:"ssl_mode"`
	Password   string `mapstructure:"password"`
	Connection struct {
		Attempts int           `mapstructure:"attempts"`
		Wait     time.Duration `mapstructure:"wait"`
	} `mapstructure:"connection"`
}

var ErrZeroAttempts = errors.New("zero attempts to connection")

func NewPostgresConnection(info ConnectionInfo) (db *sql.DB, err error) {
	err = ErrZeroAttempts
	for i := 1; i <= info.Connection.Attempts; i++ {
		db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password))
		if err != nil {
			logrus.Infof("Attempt #%d Postgres is unavailable - sleeping", i)
			time.Sleep(info.Connection.Wait)
			continue
		}

		if err := db.Ping(); err != nil {
			logrus.Infof("Attempt #%d Postgres is unavailable - sleeping", i)
			time.Sleep(info.Connection.Wait)
			continue
		}

		logrus.Infof("Attempt #%d Postgres is ready - continue", i)
		break
	}

	return db, err
}
