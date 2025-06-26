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

	// Initialize controllers
	gstController := controllers.NewGSTController(gstService)

	// IRAS GST API routes (following the swagger spec basePath)
	irasGroup := router.Group("/iras/prod/GSTListing")
	{
		// Main GST search endpoint as per IRAS API spec
		irasGroup.POST("/SearchGSTRegistered", gstController.SearchGSTRegistered)
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

			// Products (public read access)
			public.GET("/products", productController.GetProducts)
			public.GET("/products/:id", productController.GetProduct)

			// Users (public registration)
			public.POST("/users/register", userController.CreateUser)

			// Auth endpoints (for future JWT implementation)
			public.POST("/auth/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Login endpoint - JWT implementation needed",
				})
			})
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// User management
			users := protected.Group("/users")
			{
				users.GET("/", userController.GetUsers)
				users.GET("/:id", userController.GetUser)
				users.PUT("/:id", userController.UpdateUser)
				users.DELETE("/:id", userController.DeleteUser)
			}

			// Product management
			products := protected.Group("/products")
			{
				products.POST("/", productController.CreateProduct)
				products.PUT("/:id", productController.UpdateProduct)
				products.DELETE("/:id", productController.DeleteProduct)
			}

			// Category management
			categories := protected.Group("/categories")
			{
				categories.POST("/", categoryController.CreateCategory)
				categories.PUT("/:id", categoryController.UpdateCategory)
				categories.DELETE("/:id", categoryController.DeleteCategory)
			}
		}

		// Admin routes (for future role-based access)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired())
		// admin.Use(middleware.AdminRequired()) // Implement this middleware
		{
			admin.GET("/stats", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Admin stats endpoint",
				})
			})
		}
	}

	return router
}
