package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

func MustLoad(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable %s is required", key)
	}
	return val
}

func MustLoadInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable %s is required", key)
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("environment variable %s is not a number", key)
	}
	return i
}

func MustLoadDuration(key string) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable %s is required", key)
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("environment variable %s is not a duration", key)
	}
	return d
}
