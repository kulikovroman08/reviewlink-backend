package main

import (
	"log"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/app"
)

func main() {
	cfg := configs.LoadConfig()

	r := app.InitApp(&cfg)

	log.Println("Server running on :" + cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
