package main

import (
	"log"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/app"
)

func main() {
	cfg := configs.LoadConfig()

	reviewLinkApp := app.InitApp(&cfg)

	log.Println("Server running on :" + cfg.HTTPPort)

	if err := reviewLinkApp.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
