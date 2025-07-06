package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"api-iras/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GSTController struct {
	gstService *services.GSTService
}

func NewGSTController(gstService *services.GSTService) *GSTController {
	return &GSTController{gstService: gstService}
}

// @Summary Search GST Registered
// @Description Check whether businesses are GST-registered based on their GST registration number, UEN or NRIC
// @Tags GST
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param body body models.GSTRequest true "GST Search Request"
// @Success 200 {object} models.GSTResponse
// @Router /iras/prod/GSTListing/SearchGSTRegistered [post]
func (ctrl *GSTController) SearchGSTRegistered(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, models.GSTResponse{
			ReturnCode: 40,
			Info: &models.GSTInfo{
				Message:     "Missing required headers",
				MessageCode: 40003,
				FieldInfoList: []models.GSTFieldError{
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
	var req models.GSTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.GSTResponse{
			ReturnCode: 40,
			Info: &models.GSTInfo{
				Message:     "Invalid request format",
				MessageCode: 40004,
				FieldInfoList: []models.GSTFieldError{
					{
						Field:   "body",
						Message: "Invalid JSON format",
					},
				},
			},
		})
		return
	}

	// Override clientID from header if provided
	if req.ClientID == "" {
		req.ClientID = clientID
	}

	// Perform GST search
	response, err := ctrl.gstService.SearchGSTRegistered(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.GSTResponse{
			ReturnCode: 50,
			Info: &models.GSTInfo{
				Message:     "Internal server error",
				MessageCode: 50001,
			},
		})
		return
	}

	// Return response based on return code
	switch response.ReturnCode {
	case 10:
		c.JSON(http.StatusOK, response)
	case 20:
		c.JSON(http.StatusNotFound, response)
	case 40:
		c.JSON(http.StatusBadRequest, response)
	default:
		c.JSON(http.StatusInternalServerError, response)
	}
}

// Admin endpoints for managing GST registrations (for setup/maintenance)

// @Summary Create GST Registration
// @Description Create a new GST registration record (Admin only)
// @Tags GST Admin
// @Accept json
// @Produce json
// @Param gst body models.GSTRegistration true "GST Registration data"
// @Success 201 {object} models.GSTRegistration
// @Router /admin/gst-registrations [post]
func (ctrl *GSTController) CreateGSTRegistration(c *gin.Context) {
	var gstReg models.GSTRegistration
	if err := c.ShouldBindJSON(&gstReg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	if err := ctrl.gstService.CreateGSTRegistration(&gstReg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create GST registration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gstReg)
}

// @Summary Get GST Registrations
// @Description Get all GST registrations with pagination (Admin only)
// @Tags GST Admin
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} models.GSTRegistration
// @Router /admin/gst-registrations [get]
func (ctrl *GSTController) GetGSTRegistrations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	gstRegs, total, err := ctrl.gstService.GetAllGSTRegistrations(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get GST registrations",
			"details": err.Error(),
		})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data":        gstRegs,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// @Summary Get GST Registration by ID
// @Description Get a single GST registration by ID (Admin only)
// @Tags GST Admin
// @Accept json
// @Produce json
// @Param id path int true "GST Registration ID"
// @Success 200 {object} models.GSTRegistration
// @Router /admin/gst-registrations/{id} [get]
func (ctrl *GSTController) GetGSTRegistration(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid GST registration ID",
		})
		return
	}

	gstReg, err := ctrl.gstService.GetGSTRegistrationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "GST registration not found",
		})
		return
	}

	c.JSON(http.StatusOK, gstReg)
}

// @Summary Update GST Registration
// @Description Update GST registration by ID (Admin only)
// @Tags GST Admin
// @Accept json
// @Produce json
// @Param id path int true "GST Registration ID"
// @Param gst body object true "GST Registration updates"
// @Success 200 {object} gin.H
// @Router /admin/gst-registrations/{id} [put]
func (ctrl *GSTController) UpdateGSTRegistration(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid GST registration ID",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	if err := ctrl.gstService.UpdateGSTRegistration(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update GST registration",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "GST registration updated successfully",
	})
}

// @Summary Delete GST Registration
// @Description Delete GST registration by ID (Admin only)
// @Tags GST Admin
// @Accept json
// @Produce json
// @Param id path int true "GST Registration ID"
// @Success 200 {object} gin.H
// @Router /admin/gst-registrations/{id} [delete]
func (ctrl *GSTController) DeleteGSTRegistration(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid GST registration ID",
		})
		return
	}

	if err := ctrl.gstService.DeleteGSTRegistration(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete GST registration",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "GST registration deleted successfully",
	})
}
