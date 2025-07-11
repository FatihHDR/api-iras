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

	// Drop and recreate users table if it has schema issues
	log.Println("Checking users table schema...")

	// Check if table exists and try to fix it
	if db.Migrator().HasTable(&models.User{}) {
		// Try to drop the table to recreate with correct schema
		log.Println("Dropping existing users table to fix schema...")
		if err := db.Migrator().DropTable(&models.User{}); err != nil {
			log.Printf("Warning: Failed to drop users table: %v", err)
		}
	}

	// Auto-migrate all models (this will create tables with correct schema)
	log.Println("Running auto-migration...")
	err := db.AutoMigrate(
		&models.User{},
		&models.GSTRegistration{},
		&models.PropertyConsolidatedStatementRecord{},
		&models.PropertyTaxBalanceRecord{},
		&models.RentalSubmissionRecord{},
		&models.CITConversionRecord{},
		&models.SingPassAuthRecord{},
		&models.SingPassTokenRecord{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
