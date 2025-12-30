package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort            string
	DBUrl               string
	RedisAddr           string
	RedisPassword       string
	RedisDB             int
	LeaderboardCacheTTL time.Duration
	TokensAutoCount     int
	TokensThreshold     int
	TokensBatchSize     int
	BonusRequiredPts    int
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
	ttlSeconds := getEnvInt("LEADERBOARD_CACHE_TTL_SECONDS", 60)

	cfg := Config{
		HTTPPort:            os.Getenv("PORT"),
		DBUrl:               dbURL,
		RedisAddr:           os.Getenv("REDIS_ADDR"),
		RedisPassword:       os.Getenv("REDIS_PASSWORD"),
		RedisDB:             getEnvInt("REDIS_DB", 0),
		LeaderboardCacheTTL: time.Second * time.Duration(ttlSeconds),
		TokensAutoCount:     getEnvInt("TOKENS_AUTO_COUNT", 10),
		TokensThreshold:     getEnvInt("TOKENS_THRESHOLD", 5),
		TokensBatchSize:     getEnvInt("TOKENS_BATCH_SIZE", 10),
		BonusRequiredPts:    getEnvInt("BONUS_REQUIRED_POINTS", 50),
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
