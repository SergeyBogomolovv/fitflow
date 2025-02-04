package config

import "github.com/SergeyBogomolovv/fitflow/pkg/env"

type CLIConfig struct {
	PostgresURL string
}

func NewCLIConfig() *CLIConfig {
	return &CLIConfig{
		PostgresURL: env.MustLoad("POSTGRES_URL"),
	}
}
