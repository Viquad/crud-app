package config

import (
	"fmt"
	"os"
	"time"

	"github.com/Viquad/crud-app/pkg/database"
	"github.com/spf13/viper"
)

const _SECRET = "SECRET"
const _POSTGRES_PASSWORD = "POSTGRES_PASSWORD"

type Config struct {
	DB     database.ConnectionInfo `mapstructure:"db"`
	Secret string
	Cache  struct {
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

	cfg.Secret = os.Getenv(_SECRET)
	cfg.DB.Password = os.Getenv(_POSTGRES_PASSWORD)
	if len(cfg.Secret) == 0 {
		return &cfg, fmt.Errorf("missed $(%s) environment variable", _SECRET)
	}
	if len(cfg.DB.Password) == 0 {
		return &cfg, fmt.Errorf("missed $(%s) environment variable", _POSTGRES_PASSWORD)
	}

	return &cfg, nil
}
