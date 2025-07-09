package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"api-iras/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CITController struct {
	citService *services.CITService
}

func NewCITController(citService *services.CITService) *CITController {
	return &CITController{citService: citService}
}

// @Summary Convert Form CS
// @Description Convert Form CS based on the provided ID
// @Tags CIT
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param body body models.CITConversionRequest true "CIT Conversion Request"
// @Success 200 {object} models.CITConversionResponse
// @Router /iras/prod/ct/convertformcs [post]
func (ctrl *CITController) ConvertFormCS(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, models.CITConversionResponse{
			ReturnCode: 40,
			Info: &models.CITConversionInfo{
				Message:     "Missing required headers",
				MessageCode: 40003,
				FieldInfoList: []models.CITConversionFieldError{
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
	var req models.CITConversionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.CITConversionResponse{
			ReturnCode: 40,
			Info: &models.CITConversionInfo{
				Message:     "Invalid request format",
				MessageCode: 40004,
				FieldInfoList: []models.CITConversionFieldError{
					{
						Field:   "body",
						Message: "Invalid JSON format",
					},
				},
			},
		})
		return
	}

	// Call service to convert form CS
	response, err := ctrl.citService.ConvertFormCS(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CITConversionResponse{
			ReturnCode: 50,
			Info: &models.CITConversionInfo{
				Message:     "Internal server error",
				MessageCode: 50001,
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

// Admin endpoints for managing CIT conversion records

// CreateCITConversionRecord creates a new CIT conversion record
func (ctrl *CITController) CreateCITConversionRecord(c *gin.Context) {
	var record models.CITConversionRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.citService.CreateCITConversionRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create CIT conversion record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "CIT conversion record created successfully",
		Data:    record,
	})
}

// GetCITConversionRecords retrieves all CIT conversion records with pagination
func (ctrl *CITController) GetCITConversionRecords(c *gin.Context) {
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

	records, total, err := ctrl.citService.GetCITConversionRecords(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch CIT conversion records",
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
		Message: "CIT conversion records retrieved successfully",
		Data:    response,
	})
}

// GetCITConversionRecord retrieves a specific CIT conversion record by ID
func (ctrl *CITController) GetCITConversionRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	record, err := ctrl.citService.GetCITConversionRecordByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "CIT conversion record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "CIT conversion record retrieved successfully",
		Data:    record,
	})
}

// GetCITConversionRecordByConversionID retrieves a specific CIT conversion record by conversion ID
func (ctrl *CITController) GetCITConversionRecordByConversionID(c *gin.Context) {
	conversionID := c.Param("conversionId")
	if conversionID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Conversion ID is required",
		})
		return
	}

	record, err := ctrl.citService.GetCITConversionRecordByConversionID(conversionID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "CIT conversion record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "CIT conversion record retrieved successfully",
		Data:    record,
	})
}

// GetCITConversionRecordByRequestID retrieves a specific CIT conversion record by request ID
func (ctrl *CITController) GetCITConversionRecordByRequestID(c *gin.Context) {
	requestIDStr := c.Param("requestId")
	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request ID format",
		})
		return
	}

	record, err := ctrl.citService.GetCITConversionRecordByRequestID(requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "CIT conversion record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "CIT conversion record retrieved successfully",
		Data:    record,
	})
}

// UpdateCITConversionRecord updates a CIT conversion record
func (ctrl *CITController) UpdateCITConversionRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var record models.CITConversionRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.citService.UpdateCITConversionRecord(uint(id), &record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update CIT conversion record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "CIT conversion record updated successfully",
	})
}

// DeleteCITConversionRecord deletes a CIT conversion record
func (ctrl *CITController) DeleteCITConversionRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	if err := ctrl.citService.DeleteCITConversionRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete CIT conversion record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "CIT conversion record deleted successfully",
	})
}
