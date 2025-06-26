package main

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"api-iras/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	// Auto-migrate database schemas
	if err := autoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Setup Gin mode
	if config.AppConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup routes
	router := routes.SetupRoutes(config.AppConfig.DB)

	// Start server
	log.Printf("Server starting on port %s", config.AppConfig.Port)
	log.Printf("Environment: %s", config.AppConfig.Env)

	if err := router.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func autoMigrate() error {
	db := config.AppConfig.DB

	// Auto-migrate GST Registration model
	err := db.AutoMigrate(
		&models.GSTRegistration{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
