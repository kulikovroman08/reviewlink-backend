package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort        string
	DBUrl           string
	TokensAutoCount int
	TokensThreshold int
	TokensBatchSize int
}

func LoadConfig() Config {
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.test")

	dbURL := os.Getenv("DB_URL")
	if os.Getenv("APP_ENV") == "test" {
		if val := os.Getenv("DB_URL_TEST"); val != "" {
			dbURL = val
		}
	}

	cfg := Config{
		HTTPPort:        os.Getenv("PORT"),
		DBUrl:           dbURL,
		TokensAutoCount: getEnvInt("TOKENS_AUTO_COUNT", 10),
		TokensThreshold: getEnvInt("TOKENS_THRESHOLD", 5),
		TokensBatchSize: getEnvInt("TOKENS_BATCH_SIZE", 10),
	}

	fmt.Println("APP_ENV:", os.Getenv("APP_ENV"))
	fmt.Println("DB_URL:", dbURL)
	fmt.Println("TOKENS_AUTO_COUNT:", cfg.TokensAutoCount)
	fmt.Println("TOKENS_THRESHOLD:", cfg.TokensThreshold)
	fmt.Println("TOKENS_BATCH_SIZE:", cfg.TokensBatchSize)

	return cfg
}

func getEnvInt(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return def
}
