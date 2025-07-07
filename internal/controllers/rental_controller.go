package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"api-iras/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RentalController struct {
	rentalService *services.RentalService
}

func NewRentalController(rentalService *services.RentalService) *RentalController {
	return &RentalController{rentalService: rentalService}
}

// @Summary Submit Rental Information
// @Description Submit rental information for property assessment
// @Tags Rental
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param body body models.RentalSubmissionRequest true "Rental Submission Request"
// @Success 200 {object} models.RentalSubmissionResponse
// @Router /iras/sb/rental/Submission [post]
func (ctrl *RentalController) SubmitRental(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, models.RentalSubmissionResponse{
			ReturnCode: 40,
			Info: &models.RentalSubmissionInfo{
				Message:     "Missing required headers",
				MessageCode: 40003,
				FieldInfoList: &models.RentalSubmissionFieldInfo{
					FieldInfo: []models.RentalSubmissionFieldError{
						{
							Field:   "headers",
							Message: "X-IBM-Client-Id and X-IBM-Client-Secret are required",
						},
					},
				},
			},
		})
		return
	}

	// Parse request body
	var req models.RentalSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.RentalSubmissionResponse{
			ReturnCode: 40,
			Info: &models.RentalSubmissionInfo{
				Message:     "Invalid request format",
				MessageCode: 40004,
				FieldInfoList: &models.RentalSubmissionFieldInfo{
					FieldInfo: []models.RentalSubmissionFieldError{
						{
							Field:   "body",
							Message: "Invalid JSON format",
						},
					},
				},
			},
		})
		return
	}

	// Call service to submit rental
	response, err := ctrl.rentalService.SubmitRental(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RentalSubmissionResponse{
			ReturnCode: 50,
			Info: &models.RentalSubmissionInfo{
				Message:     "Internal server error",
				MessageCode: 50001,
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

// Admin endpoints for managing rental submission records

// CreateRentalSubmissionRecord creates a new rental submission record
func (ctrl *RentalController) CreateRentalSubmissionRecord(c *gin.Context) {
	var record models.RentalSubmissionRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.rentalService.CreateRentalSubmissionRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create rental submission record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Rental submission record created successfully",
		Data:    record,
	})
}

// GetRentalSubmissionRecords retrieves all rental submission records with pagination
func (ctrl *RentalController) GetRentalSubmissionRecords(c *gin.Context) {
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

	records, total, err := ctrl.rentalService.GetRentalSubmissionRecords(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch rental submission records",
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
		Message: "Rental submission records retrieved successfully",
		Data:    response,
	})
}

// GetRentalSubmissionRecord retrieves a specific rental submission record by ID
func (ctrl *RentalController) GetRentalSubmissionRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	record, err := ctrl.rentalService.GetRentalSubmissionRecordByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Rental submission record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Rental submission record retrieved successfully",
		Data:    record,
	})
}

// GetRentalSubmissionRecordByRefNo retrieves a specific rental submission record by reference number
func (ctrl *RentalController) GetRentalSubmissionRecordByRefNo(c *gin.Context) {
	refNo := c.Param("refNo")
	if refNo == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Reference number is required",
		})
		return
	}

	record, err := ctrl.rentalService.GetRentalSubmissionRecordByRefNo(refNo)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Rental submission record not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Rental submission record retrieved successfully",
		Data:    record,
	})
}

// UpdateRentalSubmissionRecord updates a rental submission record
func (ctrl *RentalController) UpdateRentalSubmissionRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	var record models.RentalSubmissionRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if err := ctrl.rentalService.UpdateRentalSubmissionRecord(uint(id), &record); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update rental submission record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Rental submission record updated successfully",
	})
}

// DeleteRentalSubmissionRecord deletes a rental submission record
func (ctrl *RentalController) DeleteRentalSubmissionRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ID format",
		})
		return
	}

	if err := ctrl.rentalService.DeleteRentalSubmissionRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete rental submission record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Rental submission record deleted successfully",
	})
}
