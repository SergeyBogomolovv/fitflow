package config

import (
	"time"

	"github.com/SergeyBogomolovv/fitflow/pkg/env"
	"github.com/joho/godotenv"
)

type ApiConfig struct {
	Addr string

	AllowedOrigins []string

	JwtSecret   []byte
	JwtTTL      time.Duration
	PostgresURL string
}

func NewApiConfig() *ApiConfig {
	godotenv.Load()
	return &ApiConfig{
		Addr:           env.MustLoad("ADDR"),
		AllowedOrigins: []string{"http://localhost:3000"}, // TODO: load from env
		JwtSecret:      []byte(env.MustLoad("JWT_SECRET")),
		JwtTTL:         env.MustLoadDuration("JWT_TTL"),
		PostgresURL:    env.MustLoad("POSTGRES_URL"),
	}
}
