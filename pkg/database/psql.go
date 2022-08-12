package database

import (
	"database/sql"
	"fmt"
)

type ConnectionInfo struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	Password string `mapstructure:"password"`
}

func NewPostgresConnection(info ConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
