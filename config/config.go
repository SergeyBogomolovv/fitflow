package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP HTTP `yaml:"http"`
		JWT  JWT  `yaml:"jwt"`
		Log  Log  `yaml:"logger"`
		PG   PG
		TG   TG
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
		Token string `env-required:"true" env:"BOT_TOKEN"`
	}

	JWT struct {
		Secret string        `env-required:"true" env:"JWT_SECRET"`
		TTL    time.Duration `env-required:"true" yaml:"ttl" env:"JWT_TTL"`
	}
)

func MustNewConfig(path string) *Config {
	godotenv.Load()
	cfg := new(Config)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Fatalf("config error: %s", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("config env error: %s", err)
	}

	return cfg
}
