package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"api-iras/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SingPassController struct {
	singpassService *services.SingPassService
}

func NewSingPassController(singpassService *services.SingPassService) *SingPassController {
	return &SingPassController{singpassService: singpassService}
}

// @Summary SingPass Service Authentication
// @Description Initiates SingPass authentication and returns authentication URL
// @Tags SingPass
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param scope query string false "Scope for authentication"
// @Param callback_url query string false "Callback URL"
// @Param state query string false "State parameter"
// @Param body body models.SingPassServiceAuthRequest false "SingPass Auth Request (for POST)"
// @Success 200 {object} models.SingPassServiceAuthResponse
// @Router /iras/prod/Authentication/SingPassServiceAuth [post]
func (ctrl *SingPassController) SingPassServiceAuth(c *gin.Context) {
	// Validate headers
	clientID := c.GetHeader("X-IBM-Client-Id")
	clientSecret := c.GetHeader("X-IBM-Client-Secret")

	// For development, accept demo credentials
	if config.AppConfig.Env == "development" {
		if clientID == "" {
			clientID = config.AppConfig.IBMClientID
		}
		if clientSecret == "" {
			clientSecret = config.AppConfig.IBMClientSecret
		}
	}

	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusUnauthorized, models.SingPassServiceAuthResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 40003,
				Message:     "Missing required headers",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "headers",
						Message: "X-IBM-Client-Id and X-IBM-Client-Secret are required",
					},
				},
			},
		})
		return
	}

	// Parse request parameters - support both query parameters and JSON body
	var req models.SingPassServiceAuthRequest
	
	// Try to bind JSON body first (for POST requests)
	if c.Request.Method == "POST" && c.GetHeader("Content-Type") == "application/json" {
		if err := c.ShouldBindJSON(&req); err != nil {
			// If JSON binding fails, try query parameters
			if err := c.ShouldBindQuery(&req); err != nil {
				c.JSON(http.StatusBadRequest, models.SingPassServiceAuthResponse{
					ReturnCode: 40,
					Info: &models.SingPassServiceAuthInfo{
						MessageCode: 40004,
						Message:     "Invalid request parameters",
						FieldInfoList: []models.SingPassServiceAuthFieldError{
							{
								Field:   "request",
								Message: "Invalid JSON body or query parameter format",
							},
						},
					},
				})
				return
			}
		}
	} else {
		// For GET requests or non-JSON POST, use query parameters
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.SingPassServiceAuthResponse{
				ReturnCode: 40,
				Info: &models.SingPassServiceAuthInfo{
					MessageCode: 40004,
					Message:     "Invalid query parameters",
					FieldInfoList: []models.SingPassServiceAuthFieldError{
						{
							Field:   "query",
							Message: "Invalid query parameter format",
						},
					},
				},
			})
			return
		}
	}

	// Call service to initiate SingPass auth
	response, err := ctrl.singpassService.SingPassServiceAuth(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.SingPassServiceAuthResponse{
			ReturnCode: 50,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 50001,
				Message:     "Internal server error",
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

// @Summary SingPass Service Authentication Token
// @Description Exchange authorization code for access token
// @Tags SingPass
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param body body models.SingPassServiceAuthTokenRequest true "SingPass Auth Token Request"
// @Success 200 {object} models.SingPassServiceAuthTokenResponse
// @Router /iras/prod/Authentication/SingPassServiceAuthToken [post]
func (ctrl *SingPassController) SingPassServiceAuthToken(c *gin.Context) {
	// Validate headers
	clientID := c.GetHeader("X-IBM-Client-Id")
	clientSecret := c.GetHeader("X-IBM-Client-Secret")

	// For development, accept demo credentials
	if config.AppConfig.Env == "development" {
		if clientID == "" {
			clientID = config.AppConfig.IBMClientID
		}
		if clientSecret == "" {
			clientSecret = config.AppConfig.IBMClientSecret
		}
	}

	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusUnauthorized, models.SingPassServiceAuthTokenResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 40003,
				Message:     "Missing required headers",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "headers",
						Message: "X-IBM-Client-Id and X-IBM-Client-Secret are required",
					},
				},
			},
		})
		return
	}

	// Parse request body
	var req models.SingPassServiceAuthTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.SingPassServiceAuthTokenResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 40004,
				Message:     "Invalid request format",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "body",
						Message: "Invalid JSON format",
					},
				},
			},
		})
		return
	}

	// Call service to exchange code for token
	response, err := ctrl.singpassService.SingPassServiceAuthToken(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.SingPassServiceAuthTokenResponse{
			ReturnCode: 50,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 50001,
				Message:     "Internal server error",
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

// Admin endpoints for managing SingPass auth records

// CreateSingPassAuthRecord creates a new SingPass auth record
func (ctrl *SingPassController) CreateSingPassAuthRecord(c *gin.Context) {
	var record models.SingPassAuthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.singpassService.CreateSingPassAuthRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create SingPass auth record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "SingPass auth record created successfully",
		Data:    record,
	})
}

// GetSingPassAuthRecords retrieves all SingPass auth records with pagination
func (ctrl *SingPassController) GetSingPassAuthRecords(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	records, total, err := ctrl.singpassService.GetSingPassAuthRecords(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch SingPass auth records",
			Error:   err.Error(),
		})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := models.PaginationResponse{
		Data:       records,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass auth records retrieved successfully",
		Data:    response,
	})
}

// GetSingPassAuthRecord retrieves a specific SingPass auth record by ID
func (ctrl *SingPassController) GetSingPassAuthRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	record, err := ctrl.singpassService.GetSingPassAuthRecordByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "SingPass auth record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass auth record retrieved successfully",
		Data:    record,
	})
}

// GetSingPassAuthRecordByState retrieves a specific SingPass auth record by state
func (ctrl *SingPassController) GetSingPassAuthRecordByState(c *gin.Context) {
	state := c.Param("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "State is required",
		})
		return
	}

	record, err := ctrl.singpassService.GetSingPassAuthRecordByState(state)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "SingPass auth record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass auth record retrieved successfully",
		Data:    record,
	})
}

// UpdateSingPassAuthRecord updates a SingPass auth record
func (ctrl *SingPassController) UpdateSingPassAuthRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var record models.SingPassAuthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.singpassService.UpdateSingPassAuthRecord(uint(id), &record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update SingPass auth record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass auth record updated successfully",
	})
}

// DeleteSingPassAuthRecord deletes a SingPass auth record
func (ctrl *SingPassController) DeleteSingPassAuthRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	if err := ctrl.singpassService.DeleteSingPassAuthRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete SingPass auth record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass auth record deleted successfully",
	})
}

// Token management endpoints

// CreateSingPassTokenRecord creates a new SingPass token record
func (ctrl *SingPassController) CreateSingPassTokenRecord(c *gin.Context) {
	var record models.SingPassTokenRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.singpassService.CreateSingPassTokenRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create SingPass token record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "SingPass token record created successfully",
		Data:    record,
	})
}

// GetSingPassTokenRecords retrieves all SingPass token records with pagination
func (ctrl *SingPassController) GetSingPassTokenRecords(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	records, total, err := ctrl.singpassService.GetSingPassTokenRecords(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch SingPass token records",
			Error:   err.Error(),
		})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := models.PaginationResponse{
		Data:       records,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass token records retrieved successfully",
		Data:    response,
	})
}

// GetSingPassTokenRecord retrieves a specific SingPass token record by ID
func (ctrl *SingPassController) GetSingPassTokenRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	record, err := ctrl.singpassService.GetSingPassTokenRecordByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "SingPass token record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass token record retrieved successfully",
		Data:    record,
	})
}

// GetSingPassTokenRecordByState retrieves a specific SingPass token record by state
func (ctrl *SingPassController) GetSingPassTokenRecordByState(c *gin.Context) {
	state := c.Param("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "State is required",
		})
		return
	}

	record, err := ctrl.singpassService.GetSingPassTokenRecordByState(state)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "SingPass token record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass token record retrieved successfully",
		Data:    record,
	})
}

// UpdateSingPassTokenRecord updates a SingPass token record
func (ctrl *SingPassController) UpdateSingPassTokenRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var record models.SingPassTokenRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.singpassService.UpdateSingPassTokenRecord(uint(id), &record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update SingPass token record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass token record updated successfully",
	})
}

// DeleteSingPassTokenRecord deletes a SingPass token record
func (ctrl *SingPassController) DeleteSingPassTokenRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	if err := ctrl.singpassService.DeleteSingPassTokenRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete SingPass token record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SingPass token record deleted successfully",
	})
}
