package config

import "github.com/SergeyBogomolovv/fitflow/pkg/env"

type BotConfig struct {
	Token       string
	PostgresURL string
}

func NewBotConfig() *BotConfig {
	return &BotConfig{
		Token:       env.MustLoad("BOT_TOKEN"),
		PostgresURL: env.MustLoad("POSTGRES_URL"),
	}
}
