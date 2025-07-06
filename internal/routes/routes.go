package routes

import (
	"api-iras/internal/controllers"
	"api-iras/internal/middleware"
	"api-iras/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "IRAS GST API is running",
			"version": "1.0.7",
		})
	})

	// Initialize services
	gstService := services.NewGSTService(db)
	authService := services.NewAuthService(db)

	// Initialize controllers
	gstController := controllers.NewGSTController(gstService)
	authController := controllers.NewAuthController(authService)

	// IRAS GST API routes (following the swagger spec basePath)
	irasGroup := router.Group("/iras/prod/GSTListing")
	{
		// Main GST search endpoint as per IRAS API spec
		irasGroup.POST("/SearchGSTRegistered", gstController.SearchGSTRegistered)
	}

	// IRAS CorpPass Authentication routes
	corpPassGroup := router.Group("/iras/sb/Authentication")
	{
		corpPassGroup.GET("/CorpPassAuth", gstController.CorpPassAuth)
		corpPassGroup.POST("/CorpPassToken", gstController.CorpPassToken)
	}

	// Authentication routes (public)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
		authGroup.GET("/demo-token", authController.GenerateDemoToken) // Development only

		// Protected auth routes
		authGroup.Use(middleware.AuthRequired())
		authGroup.GET("/profile", authController.GetProfile)
		authGroup.PUT("/profile", authController.UpdateProfile)
	}

	// Admin routes for managing GST registrations (protected with auth)
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthRequired()) // Add authentication middleware
	{
		// GST Registration management endpoints
		adminGroup.POST("/gst-registrations", gstController.CreateGSTRegistration)
		adminGroup.GET("/gst-registrations", gstController.GetGSTRegistrations)
		adminGroup.GET("/gst-registrations/:id", gstController.GetGSTRegistration)
		adminGroup.PUT("/gst-registrations/:id", gstController.UpdateGSTRegistration)
		adminGroup.DELETE("/gst-registrations/:id", gstController.DeleteGSTRegistration)

		// User management endpoints (admin only)
		adminGroup.GET("/users", authController.GetAllUsers)
		adminGroup.PUT("/users/:id/deactivate", authController.DeactivateUser)
	}

	// API info endpoint
	router.GET("/api/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"title":       "Check GST Register",
			"description": "The Check GST Register API enables you to check whether businesses are GST-registered based on their GST registration number, UEN or NRIC.",
			"version":     "1.0.7",
			"basePath":    "/iras/prod/GSTListing",
			"schemes":     []string{"https"},
			"host":        "apiservices.iras.gov.sg",
			"consumes":    []string{"application/json"},
			"produces":    []string{"application/json"},
			"endpoints": gin.H{
				"main": "/iras/prod/GSTListing/SearchGSTRegistered",
				"admin": gin.H{
					"create": "/admin/gst-registrations",
					"list":   "/admin/gst-registrations",
					"get":    "/admin/gst-registrations/{id}",
					"update": "/admin/gst-registrations/{id}",
					"delete": "/admin/gst-registrations/{id}",
				},
			},
		})
	})

	return router
}
