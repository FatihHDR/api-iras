package middleware

import (
	"api-iras/internal/config"
	"api-iras/pkg/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log details
		timeStamp := time.Now()
		latency := timeStamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// You can customize this log format
		gin.DefaultWriter.Write([]byte(
			fmt.Sprintf("[GIN] %s | %d | %s | %s | %s | %s\n",
				timeStamp.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency.String(),
				clientIP,
				method,
				path),
		))
	}
}

// ErrorHandler middleware
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Invalid request format",
					"error":   err.Error(),
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal server error",
					"error":   err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal server error",
				})
			}
		}
	}
}

// Validator middleware
func ValidateJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Validate struct
		if err := validate.Struct(obj); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Error())
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validation failed",
				"errors":  errors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// JWT Auth middleware
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token format",
			})
			c.Abort()
			return
		}

		// For development: accept simple demo tokens
		if strings.HasPrefix(tokenString, "demo-token") && len(tokenString) >= 10 {
			c.Set("user_id", "999")
			c.Set("username", "demo")
			c.Set("email", "demo@example.com")
			c.Set("role", "admin")
			c.Next()
			return
		}

		// Validate JWT token
		claims, err := utils.ValidateJWT(tokenString, config.AppConfig.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}
