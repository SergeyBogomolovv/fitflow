package config

import (
	"github.com/SergeyBogomolovv/fitflow/pkg/env"
	"github.com/joho/godotenv"
)

type BotConfig struct {
	Token       string
	PostgresURL string
}

func NewBotConfig() *BotConfig {
	godotenv.Load()
	return &BotConfig{
		Token:       env.MustLoad("BOT_TOKEN"),
		PostgresURL: env.MustLoad("POSTGRES_URL"),
	}
}
