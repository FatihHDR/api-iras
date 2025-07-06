package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"api-iras/internal/models"
	"api-iras/internal/services"
)

type AISController struct {
	aisService *services.AISService
	validator  *validator.Validate
}

func NewAISController(aisService *services.AISService) *AISController {
	return &AISController{
		aisService: aisService,
		validator:  validator.New(),
	}
}

// AISOrgSearch handles organization search in AIS
func (ac *AISController) AISOrgSearch(c *gin.Context) {
	var request models.AISorgSearchRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// Validate request
	if err := ac.validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Call service to search organization in AIS
	response, err := ac.aisService.SearchOrganization(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to search organization in AIS",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Organization search completed successfully",
		Data:    response,
	})
}
