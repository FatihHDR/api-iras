package routes

import (
	"api-iras/internal/controllers"
	"api-iras/internal/middleware"
	"api-iras/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	// Initialize services
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	categoryService := services.NewCategoryService(db)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	categoryController := controllers.NewCategoryController(categoryService)

	// Setup Gin router
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "API is running",
		})
	})

	// API version 1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := v1.Group("/")
		{
			// Categories (public read access)
			public.GET("/categories", categoryController.GetCategories)
			public.GET("/categories/:id", categoryController.GetCategory)

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
