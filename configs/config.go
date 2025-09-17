package configs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppEnv    string `env:"APP_ENV" env-default:"dev"`
	HTTPPort  string `env:"PORT" env-default:"8080"`
	DBUrl     string `env:"DB_URL"`
	DBUrlTest string `env:"DB_URL_TEST"`
}

func LoadConfig() *Config {
	var cfg Config

	// Читаем APP_ENV напрямую
	appEnv := os.Getenv("APP_ENV")

	// Определяем абсолютный путь к .env / .env.test
	wd, _ := os.Getwd()
	envFile := filepath.Join(wd, ".env")
	if appEnv == "test" {
		envFile = filepath.Join(wd, ".env.test")
	}

	// Читаем конфиг
	if err := cleanenv.ReadConfig(envFile, &cfg); err != nil {
		log.Fatalf("failed to read config from %s: %v", envFile, err)
	}

	// Подменяем строку подключения в тестах
	if cfg.AppEnv == "test" && cfg.DBUrlTest != "" {
		cfg.DBUrl = cfg.DBUrlTest
	}

	log.Printf("APP_ENV: %s", cfg.AppEnv)
	log.Printf("DB_URL: %s", cfg.DBUrl)

	return &cfg
}
