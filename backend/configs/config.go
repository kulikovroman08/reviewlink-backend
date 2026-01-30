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
	OSMCacheTTL         time.Duration
	OSMCacheHardTTL     time.Duration
	PublicNewsRSSURL    string
	PublicNewsCacheTTL  time.Duration
	TokensAutoCount     int
	TokensThreshold     int
	TokensBatchSize     int
	BonusRequiredPts    int
}

func LoadConfig() Config {
	env := os.Getenv("APP_ENV")
	_ = godotenv.Load(".env")

	if env == "test" {
		_ = godotenv.Overload(".env.test")
	}

	dbURL := os.Getenv("DB_URL")
	if env == "test" {
		if val := os.Getenv("DB_URL_TEST"); val != "" {
			dbURL = val
		}
	}

	cfg := Config{
		HTTPPort:            os.Getenv("PORT"),
		DBUrl:               dbURL,
		RedisAddr:           os.Getenv("REDIS_ADDR"),
		RedisPassword:       os.Getenv("REDIS_PASSWORD"),
		RedisDB:             getEnvInt("REDIS_DB", 0),
		LeaderboardCacheTTL: getEnvDuration("LEADERBOARD_CACHE_TTL", 1*time.Hour),
		OSMCacheTTL:         getEnvDuration("OSM_CACHE_TTL", 6*time.Hour),
		OSMCacheHardTTL:     getEnvDuration("OSM_CACHE_HARD_TTL", 24*time.Hour),
		PublicNewsRSSURL:    os.Getenv("PUBLIC_NEWS_RSS_URL"),
		PublicNewsCacheTTL:  getEnvDuration("PUBLIC_NEWS_CACHE_TTL", 3*time.Hour),
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

	if cfg.PublicNewsRSSURL == "" {
		panic("PUBLIC_NEWS_RSS_URL is required")
	}

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

func getEnvDuration(key string, def time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return def
}
