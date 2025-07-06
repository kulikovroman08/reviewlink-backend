package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/configs"
)

func main() {
	cfg := configs.LoadConfig()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
