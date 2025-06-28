package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"api-iras/internal/services"
	"api-iras/pkg/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService *services.AuthService
	validator   *validator.Validate
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validator:   validator.New(),
	}
}

// @Summary User Registration
// @Description Register a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.APIResponse
// @Router /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format", err))
		return
	}

	// Log the request for debugging
	fmt.Printf("Register request: %+v\n", req)

	// Validate request
	if err := ctrl.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validation failed", err))
		return
	}

	// Register user
	user, err := ctrl.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusConflict, utils.ErrorResponse("Registration failed", err))
		return
	}

	// Prepare response (exclude password)
	userInfo := &models.UserInfo{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("User registered successfully", userInfo))
}

// @Summary User Login
// @Description Authenticate user and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format", err))
		return
	}

	// Validate request
	if err := ctrl.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validation failed", err))
		return
	}

	// Authenticate user
	user, err := ctrl.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Authentication failed", err))
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(
		fmt.Sprintf("%d", user.ID),
		user.Username,
		user.Email,
		user.Role,
		config.AppConfig.JWTSecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate token", err))
		return
	}

	// Prepare response
	response := &models.LoginResponse{
		Token:     token,
		ExpiresIn: 86400,
		UserInfo: &models.UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", response))
}

// @Summary Get Current User Profile
// @Description Get current authenticated user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Router /auth/profile [get]
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	// Get user ID from JWT token (set by middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated", nil))
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID", err))
		return
	}

	// Get user from database
	user, err := ctrl.authService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found", err))
		return
	}

	// Prepare response
	userInfo := &models.UserInfo{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Profile retrieved successfully", userInfo))
}

// @Summary Update User Profile
// @Description Update current authenticated user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body map[string]interface{} true "Update data"
// @Success 200 {object} models.APIResponse
// @Router /auth/profile [put]
func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	// Get user ID from JWT token
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated", nil))
		return
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID", err))
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request data", err))
		return
	}

	// Remove sensitive fields that shouldn't be updated via this endpoint
	delete(updates, "role")      // Role should be updated by admin only
	delete(updates, "is_active") // Should be updated via separate endpoint

	// Update user
	if err := ctrl.authService.UpdateUser(uint(userID), updates); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update profile", err))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Profile updated successfully", nil))
}

// @Summary Generate Demo Token
// @Description Generate a demo token for testing (development only)
// @Tags Auth
// @Produce json
// @Success 200 {object} models.LoginResponse
// @Router /auth/demo-token [get]
func (ctrl *AuthController) GenerateDemoToken(c *gin.Context) {
	if config.AppConfig.Env != "development" {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Demo tokens only available in development", nil))
		return
	}

	// Generate demo token
	token, err := utils.GenerateJWT(
		"demo-user",
		"demo",
		"demo@example.com",
		"admin",
		config.AppConfig.JWTSecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate token", err))
		return
	}

	response := &models.LoginResponse{
		Token:     token,
		ExpiresIn: 86400,
		UserInfo: &models.UserInfo{
			ID:       999,
			Name:     "Demo User",
			Username: "demo",
			Email:    "demo@example.com",
			Role:     "admin",
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Demo token generated", response))
}

// Admin endpoints
// @Summary Get All Users (Admin Only)
// @Description Get list of all users with pagination
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.APIResponse
// @Router /admin/users [get]
func (ctrl *AuthController) GetAllUsers(c *gin.Context) {
	// Check if user is admin
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Admin access required", nil))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := ctrl.authService.GetAllUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get users", err))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Users retrieved successfully", result))
}

// @Summary Deactivate User (Admin Only)
// @Description Deactivate a user account
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.APIResponse
// @Router /admin/users/{id}/deactivate [put]
func (ctrl *AuthController) DeactivateUser(c *gin.Context) {
	// Check if user is admin
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Admin access required", nil))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID", err))
		return
	}

	if err := ctrl.authService.DeactivateUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to deactivate user", err))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("User deactivated successfully", nil))
}
