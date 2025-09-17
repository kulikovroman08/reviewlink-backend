package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string
	DBUrl    string
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

	fmt.Println("APP_ENV:", os.Getenv("APP_ENV"))
	fmt.Println("DB_URL:", dbURL)

	return Config{
		HTTPPort: os.Getenv("PORT"),
		DBUrl:    dbURL,
	}
}
