package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type EStampController struct{}

func NewEStampController() *EStampController {
	return &EStampController{}
}

// @Summary Stamp Tenancy Agreement
// @Description Submit tenancy agreement for stamping
// @Tags eStamp
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param access_token header string true "CorpPass Access Token"
// @Param body body models.StampTenancyAgreementRequest true "Tenancy Agreement Stamp Request"
// @Success 200 {object} models.EStampResponse
// @Router /iras/sb/eStamp/StampTenancyAgreement [post]
func (ctrl *EStampController) StampTenancyAgreement(c *gin.Context) {
	// Validate headers
	clientID := c.GetHeader("X-IBM-Client-Id")
	clientSecret := c.GetHeader("X-IBM-Client-Secret")
	accessToken := c.GetHeader("access_token")

	// For development, accept demo credentials
	if config.AppConfig.Env == "development" {
		if clientID == "" {
			clientID = config.AppConfig.IBMClientID
		}
		if clientSecret == "" {
			clientSecret = config.AppConfig.IBMClientSecret
		}
		if accessToken == "" {
			accessToken = "demo_access_token_123456"
		}
	}

	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusUnauthorized, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40003,
				Message:     "Missing required headers",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "headers",
						Message: "X-IBM-Client-Id and X-IBM-Client-Secret are required",
					},
				},
			},
		})
		return
	}

	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40004,
				Message:     "Missing access token",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "access_token",
						Message: "CorpPass access token is required",
					},
				},
			},
		})
		return
	}

	// Parse request body
	var req models.StampTenancyAgreementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40005,
				Message:     "Invalid request format",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "body",
						Message: "Invalid JSON format or missing required fields",
					},
				},
			},
		})
		return
	}

	// Validate required fields
	if req.AssignID == "" || req.DocumentDescription == "" {
		c.JSON(http.StatusBadRequest, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40006,
				Message:     "Missing required fields",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "assignId,documentDescription",
						Message: "AssignID and DocumentDescription are required",
					},
				},
			},
		})
		return
	}

	// Process stamp calculation (simulation)
	response := ctrl.processStampTenancyAgreement(&req)

	c.JSON(http.StatusOK, response)
}

// @Summary Share Transfer Stamping
// @Description Submit share transfer document for stamping
// @Tags eStamp
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param access_token header string true "CorpPass Access Token"
// @Param body body models.ShareTransferRequest true "Share Transfer Stamp Request"
// @Success 200 {object} models.EStampResponse
// @Router /iras/sb/eStamp/ShareTransfer [post]
func (ctrl *EStampController) ShareTransfer(c *gin.Context) {
	// Validate headers
	clientID := c.GetHeader("X-IBM-Client-Id")
	clientSecret := c.GetHeader("X-IBM-Client-Secret")
	accessToken := c.GetHeader("access_token")

	// For development, accept demo credentials
	if config.AppConfig.Env == "development" {
		if clientID == "" {
			clientID = config.AppConfig.IBMClientID
		}
		if clientSecret == "" {
			clientSecret = config.AppConfig.IBMClientSecret
		}
		if accessToken == "" {
			accessToken = "demo_access_token_123456"
		}
	}

	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusUnauthorized, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40003,
				Message:     "Missing required headers",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "headers",
						Message: "X-IBM-Client-Id and X-IBM-Client-Secret are required",
					},
				},
			},
		})
		return
	}

	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40004,
				Message:     "Missing access token",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "access_token",
						Message: "CorpPass access token is required",
					},
				},
			},
		})
		return
	}

	// Parse request body
	var req models.ShareTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40005,
				Message:     "Invalid request format",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "body",
						Message: "Invalid JSON format or missing required fields",
					},
				},
			},
		})
		return
	}

	// Validate required fields
	if req.AssignID == "" || req.DocumentDescription == "" {
		c.JSON(http.StatusBadRequest, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40006,
				Message:     "Missing required fields",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "assignId,documentDescription",
						Message: "AssignID and DocumentDescription are required",
					},
				},
			},
		})
		return
	}

	// Process stamp calculation (simulation)
	response := ctrl.processShareTransfer(&req)

	c.JSON(http.StatusOK, response)
}

// Helper methods for eStamp processing
func (ctrl *EStampController) processStampTenancyAgreement(req *models.StampTenancyAgreementRequest) models.EStampResponse {
	// Simulate stamp duty calculation for tenancy agreement
	var totalRent float64
	for _, rental := range req.AssessmentRental {
		totalRent += rental.TotalGrossRentAmount
	}

	// Simple calculation: 0.4% of total rent or $1 whichever is higher
	stampDuty := totalRent * 0.004
	if stampDuty < 1.0 {
		stampDuty = 1.0
	}

	// Generate mock document reference
	docRefNo := "TA" + ctrl.generateDocumentReference()

	return models.EStampResponse{
		ReturnCode: 10,
		Data: &models.EStampData{
			DocRefNo:        docRefNo,
			SDAmount:        fmt.Sprintf("%.2f", stampDuty),
			SDPenalty:       "0.00",
			TotalAmtPayable: fmt.Sprintf("%.2f", stampDuty),
			PaymentDueDate:  "2025-08-06", // 30 days from now
			PDFBase64:       ctrl.generateMockPDF("Tenancy Agreement", docRefNo),
		},
	}
}

func (ctrl *EStampController) processShareTransfer(req *models.ShareTransferRequest) models.EStampResponse {
	// Simulate stamp duty calculation for share transfer
	considerationAmount := req.ConsiderationAmount
	if considerationAmount == 0 {
		considerationAmount = req.TargetCompany.TotalMarketPrice
	}

	// Simple calculation: 0.2% of consideration amount or $1 whichever is higher
	stampDuty := considerationAmount * 0.002
	if stampDuty < 1.0 {
		stampDuty = 1.0
	}

	// Generate mock document reference
	docRefNo := "ST" + ctrl.generateDocumentReference()

	return models.EStampResponse{
		ReturnCode: 10,
		Data: &models.EStampData{
			DocRefNo:        docRefNo,
			SDAmount:        fmt.Sprintf("%.2f", stampDuty),
			SDPenalty:       "0.00",
			TotalAmtPayable: fmt.Sprintf("%.2f", stampDuty),
			PaymentDueDate:  "2025-08-06", // 30 days from now
			PDFBase64:       ctrl.generateMockPDF("Share Transfer", docRefNo),
		},
	}
}

func (ctrl *EStampController) generateDocumentReference() string {
	// Generate mock document reference number
	return fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(1000))
}

func (ctrl *EStampController) generateMockPDF(docType, docRefNo string) string {
	// Generate mock base64 PDF content
	mockPDFContent := fmt.Sprintf("PDF-1.4 Mock %s Document - Ref: %s - Generated: %s",
		docType, docRefNo, time.Now().Format("2006-01-02 15:04:05"))
	return base64.StdEncoding.EncodeToString([]byte(mockPDFContent))
}
