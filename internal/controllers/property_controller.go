package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"api-iras/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PropertyController struct {
	propertyService *services.PropertyService
}

func NewPropertyController(propertyService *services.PropertyService) *PropertyController {
	return &PropertyController{propertyService: propertyService}
}

// @Summary Retrieve Property Consolidated Statement
// @Description Retrieve consolidated property tax statement based on reference number and property tax reference
// @Tags Property
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param body body models.PropertyConsolidatedStatementRequest true "Property Consolidated Statement Request"
// @Success 200 {object} models.PropertyConsolidatedStatementResponse
// @Router /iras/sb/PropertyConsolidatedStatement/retrieve [post]
func (ctrl *PropertyController) RetrieveConsolidatedStatement(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, models.PropertyConsolidatedStatementResponse{
			ReturnCode: 40,
			Info: &models.PropertyConsolidatedStatementInfo{
				Message:     "Missing required headers",
				MessageCode: 40003,
				FieldInfoList: []models.PropertyConsolidatedFieldError{
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
	var req models.PropertyConsolidatedStatementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PropertyConsolidatedStatementResponse{
			ReturnCode: 40,
			Info: &models.PropertyConsolidatedStatementInfo{
				Message:     "Invalid request format",
				MessageCode: 40004,
				FieldInfoList: []models.PropertyConsolidatedFieldError{
					{
						Field:   "body",
						Message: "Invalid JSON format",
					},
				},
			},
		})
		return
	}

	// Validate request
	if req.RefNo == "" {
		c.JSON(http.StatusBadRequest, models.PropertyConsolidatedStatementResponse{
			ReturnCode: 40,
			Info: &models.PropertyConsolidatedStatementInfo{
				Message:     "Missing reference number",
				MessageCode: 40001,
				FieldInfoList: []models.PropertyConsolidatedFieldError{
					{
						Field:   "refNo",
						Message: "Reference number is required",
					},
				},
			},
		})
		return
	}

	if req.PropertyTaxRef == "" {
		c.JSON(http.StatusBadRequest, models.PropertyConsolidatedStatementResponse{
			ReturnCode: 40,
			Info: &models.PropertyConsolidatedStatementInfo{
				Message:     "Missing property tax reference",
				MessageCode: 40002,
				FieldInfoList: []models.PropertyConsolidatedFieldError{
					{
						Field:   "propertyTaxRef",
						Message: "Property tax reference is required",
					},
				},
			},
		})
		return
	}

	// Call service to retrieve consolidated statement
	response, err := ctrl.propertyService.RetrieveConsolidatedStatement(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PropertyConsolidatedStatementResponse{
			ReturnCode: 50,
			Info: &models.PropertyConsolidatedStatementInfo{
				Message:     "Internal server error",
				MessageCode: 50001,
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

// Admin endpoints for managing property consolidated statement records

// CreateConsolidatedStatementRecord creates a new property consolidated statement record
func (ctrl *PropertyController) CreateConsolidatedStatementRecord(c *gin.Context) {
	var record models.PropertyConsolidatedStatementRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.propertyService.CreateConsolidatedStatementRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create property consolidated statement record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Property consolidated statement record created successfully",
		Data:    record,
	})
}

// GetConsolidatedStatementRecords retrieves all property consolidated statement records with pagination
func (ctrl *PropertyController) GetConsolidatedStatementRecords(c *gin.Context) {
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

	records, total, err := ctrl.propertyService.GetConsolidatedStatementRecords(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch property consolidated statement records",
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
		Message: "Property consolidated statement records retrieved successfully",
		Data:    response,
	})
}

// GetConsolidatedStatementRecord retrieves a specific property consolidated statement record by ID
func (ctrl *PropertyController) GetConsolidatedStatementRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	record, err := ctrl.propertyService.GetConsolidatedStatementRecordByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Property consolidated statement record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Property consolidated statement record retrieved successfully",
		Data:    record,
	})
}

// UpdateConsolidatedStatementRecord updates a property consolidated statement record
func (ctrl *PropertyController) UpdateConsolidatedStatementRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var record models.PropertyConsolidatedStatementRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.propertyService.UpdateConsolidatedStatementRecord(uint(id), &record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update property consolidated statement record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Property consolidated statement record updated successfully",
	})
}

// DeleteConsolidatedStatementRecord deletes a property consolidated statement record
func (ctrl *PropertyController) DeleteConsolidatedStatementRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	if err := ctrl.propertyService.DeleteConsolidatedStatementRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete property consolidated statement record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Property consolidated statement record deleted successfully",
	})
}
