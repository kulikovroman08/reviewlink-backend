package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит переменные окружения
type Config struct {
	Port  string
	DBUrl string
}

// LoadConfig читает переменные из .env и возвращает структуру Config
func LoadConfig() *Config {
	// Загружаем .env, если он существует
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // значение по умолчанию
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not set in environment")
	}

	return &Config{
		Port:  port,
		DBUrl: dbUrl,
	}
}
