package config

import (
	"os"

	"github.com/Viquad/crud-app/pkg/database"
	"github.com/spf13/viper"
)

type Config struct {
	DB database.ConnectionInfo `mapstructure:"db"`
}

func New(path, name string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(name)

	var cfg Config

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.DB.Password = os.Getenv("POSTGRES_PASSWORD")

	return &cfg, nil
}
