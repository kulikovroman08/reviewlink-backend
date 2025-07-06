package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBUrl string
}

func LoadConfig() Config {
	_ = godotenv.Load() // не паникуем, если .env нет

	return Config{
		Port:  getEnv("PORT", "8080"),
		DBUrl: getEnv("DB_URL", "postgres://user:pass@localhost:5432/reviewlink"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
