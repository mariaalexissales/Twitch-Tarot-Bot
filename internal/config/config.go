package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	TwitchToken   string `env:"TWITCH_TOKEN,required"`
	TwitchUser    string `env:"TWITCH_USER,required"`
	TwitchChannel string `env:"TWITCH_CHANNEL,required"`
	DatabaseURL   string `env:"DATABASE_URL"`
	RedisURL      string `env:"REDIS_URL"`
	Environment   string `env:"ENV" envDefault:"development"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
