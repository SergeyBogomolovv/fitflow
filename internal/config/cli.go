package config

import (
	"github.com/SergeyBogomolovv/fitflow/pkg/env"
	"github.com/joho/godotenv"
)

type CLIConfig struct {
	PostgresURL string
}

func NewCLIConfig() *CLIConfig {
	godotenv.Load()
	return &CLIConfig{
		PostgresURL: env.MustLoad("POSTGRES_URL"),
	}
}
