package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/configs"
	"log"
)

func main() {
	cfg := configs.LoadConfig()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server running on :" + cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
