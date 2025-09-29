package main

import (
	"log"

	_ "github.com/kulikovroman08/reviewlink-backend/docs"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/app"
)

// @title           Reviewlink API
// @version         1.0
// @description     API для пользователей, мест, отзывов и токенов
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := configs.LoadConfig()

	reviewLinkApp := app.InitApp(&cfg)

	log.Println("Server running on :" + cfg.HTTPPort)

	if err := reviewLinkApp.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
