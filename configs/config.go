package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string
	DBUrl    string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	return Config{
		HTTPPort: getEnv("PORT", "8080"),
		DBUrl:    getEnv("DB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
