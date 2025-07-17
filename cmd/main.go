package main

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/router"
	"log"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	_ "github.com/kulikovroman08/reviewlink-backend/internal/router"
)

func main() {
	cfg := configs.LoadConfig()

	r := router.SetupRouter()

	log.Println("Server running on :" + cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
