package main

import (
	"fmt"
	"log"
	"messageboard/config"
	"messageboard/model"
	"messageboard/router"
	"messageboard/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	// Set Gin mode
	gin.SetMode(config.AppConfig.Server.Mode)

	// Initialize database
	if err := model.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("Database initialized successfully")

	// Ensure upload directory exists
	if err := utils.EnsureDir(config.AppConfig.Upload.UploadDir); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	if err := utils.EnsureDir(config.AppConfig.Upload.UploadDir + "/avatars"); err != nil {
		log.Fatalf("Failed to create avatars directory: %v", err)
	}
	if err := utils.EnsureDir(config.AppConfig.Upload.UploadDir + "/posts"); err != nil {
		log.Fatalf("Failed to create posts directory: %v", err)
	}
	fmt.Println("Upload directories created successfully")

	// Setup router
	r := router.SetupRouter()

	// Start server
	port := config.AppConfig.Server.Port
	fmt.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}