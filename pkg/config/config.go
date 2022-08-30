package config

import (
	"os"
	"time"

	"github.com/Viquad/crud-app/pkg/database"
	"github.com/spf13/viper"
)

type Config struct {
	DB    database.ConnectionInfo `mapstructure:"db"`
	Cache struct {
		TTL time.Duration `mapstructure:"ttl"`
	} `mapstructure:"cache"`
	Auth struct {
		AccessTokenTTL  time.Duration `mapstructure:"access_ttl"`
		RefreshTokenTTL time.Duration `mapstructure:"refresh_ttl"`
	} `mapstructure:"auth"`
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
