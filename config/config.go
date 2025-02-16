package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP HTTP `yaml:"http"`
		JWT  JWT  `yaml:"jwt"`
		Log  Log  `yaml:"logger"`
		TG   TG   `yaml:"telegram"`
		AI   AI   `yaml:"ai"`
		S3   S3   `yaml:"s3"`
		PG   PG
	}

	HTTP struct {
		Port    int      `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		Host    string   `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Origins []string `env-required:"true" yaml:"allowed_origins" env:"HTTP_ALLOWED_ORIGINS"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		URL string `env-required:"true" env:"POSTGRES_URL"`
	}

	TG struct {
		Token         string `env-required:"true" env:"BOT_TOKEN"`
		BroadcastSpec string `env-required:"true" yaml:"broadcast_spec" env:"BOT_BROADCAST_SPEC"`
		LevelSpec     string `env-required:"true" yaml:"level_spec" env:"BOT_LEVEL_SPEC"`
	}

	JWT struct {
		Secret []byte        `env-required:"true" env:"JWT_SECRET"`
		TTL    time.Duration `env-required:"true" yaml:"ttl" env:"JWT_TTL"`
	}

	AI struct {
		Key           string `env-required:"true" env:"AI_KEY"`
		Model         string `env-required:"true" yaml:"model" env:"AI_MODEL"`
		DefaultPrompt string `env-required:"true" yaml:"default_prompt" env:"AI_DEFAULT_PROMPT"`
	}

	S3 struct {
		AccessKey string `env-required:"true" env:"S3_ACCESS_KEY"`
		SecretKey string `env-required:"true" env:"S3_SECRET_KEY"`
		Region    string `env-required:"true" env:"S3_REGION" yaml:"region"`
		Endpoint  string `env-required:"true" env:"S3_ENDPOINT" yaml:"endpoint"`
		Bucket    string `env-required:"true" env:"S3_BUCKET" yaml:"bucket"`
	}
)

func MustNewConfig(path string) *Config {
	cfg := new(Config)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Fatalf("config error: %s", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("config env error: %s", err)
	}

	return cfg
}
