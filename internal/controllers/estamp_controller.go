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

// @Summary Stamp Mortgage
// @Description Submit mortgage document for stamping
// @Tags eStamp
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param access_token header string true "CorpPass Access Token"
// @Param body body models.StampMortgageRequest true "Mortgage Stamp Request"
// @Success 200 {object} models.EStampResponse
// @Router /iras/sb/eStamp/StampMortgage [post]
func (ctrl *EStampController) StampMortgage(c *gin.Context) {
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
	var req models.StampMortgageRequest
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
	if req.AssignID == "" || req.DocumentDescription == "" || req.AmountOfLoan <= 0 {
		c.JSON(http.StatusBadRequest, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40006,
				Message:     "Missing required fields",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "assignId,documentDescription,amountOfLoan",
						Message: "AssignID, DocumentDescription, and AmountOfLoan are required",
					},
				},
			},
		})
		return
	}

	// Process stamp calculation (simulation)
	response := ctrl.processStampMortgage(&req)

	c.JSON(http.StatusOK, response)
}

// @Summary Sale Purchase Buyers Stamping
// @Description Submit sale purchase buyers document for stamping
// @Tags eStamp
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param access_token header string true "CorpPass Access Token"
// @Param body body models.SalePurchaseBuyersRequest true "Sale Purchase Buyers Stamp Request"
// @Success 200 {object} models.EStampResponse
// @Router /iras/sb/eStamp/SalePurchaseBuyers [post]
func (ctrl *EStampController) SalePurchaseBuyers(c *gin.Context) {
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
	var req models.SalePurchaseBuyersRequest
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
	if req.AssignID == "" || req.DocumentDescription == "" || req.PurchasePrice <= 0 {
		c.JSON(http.StatusBadRequest, models.EStampResponse{
			ReturnCode: 40,
			Info: &models.EStampInfo{
				MessageCode: 40006,
				Message:     "Missing required fields",
				FieldInfoList: []models.EStampFieldError{
					{
						Field:   "assignId,documentDescription,purchasePrice",
						Message: "AssignID, DocumentDescription, and PurchasePrice are required",
					},
				},
			},
		})
		return
	}

	// Process stamp calculation (simulation)
	response := ctrl.processSalePurchaseBuyers(&req)

	c.JSON(http.StatusOK, response)
}

// @Summary Sale Purchase Sellers Stamping
// @Description Submit sale purchase sellers document for stamping
// @Tags eStamp
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param access_token header string true "CorpPass Access Token"
// @Param body body models.SalePurchaseSellersRequest true "Sale Purchase Sellers Stamp Request"
// @Success 200 {object} models.EStampResponse
// @Router /iras/sb/eStamp/SalePurchaseSellers [post]
func (ctrl *EStampController) SalePurchaseSellers(c *gin.Context) {
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
	var req models.SalePurchaseSellersRequest
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
	response := ctrl.processSalePurchaseSellers(&req)

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

func (ctrl *EStampController) processStampMortgage(req *models.StampMortgageRequest) models.EStampResponse {
	// Simulate stamp duty calculation for mortgage
	loanAmount := req.AmountOfLoan

	// Simple calculation: 0.4% of loan amount or $5 whichever is higher
	stampDuty := loanAmount * 0.004
	if stampDuty < 5.0 {
		stampDuty = 5.0
	}

	// Generate mock document reference
	docRefNo := "MG" + ctrl.generateDocumentReference()

	return models.EStampResponse{
		ReturnCode: 10,
		Data: &models.EStampData{
			DocRefNo:        docRefNo,
			SDAmount:        fmt.Sprintf("%.2f", stampDuty),
			SDPenalty:       "0.00",
			TotalAmtPayable: fmt.Sprintf("%.2f", stampDuty),
			PaymentDueDate:  "2025-08-06", // 30 days from now
			PDFBase64:       ctrl.generateMockPDF("Mortgage", docRefNo),
		},
	}
}

func (ctrl *EStampController) processSalePurchaseBuyers(req *models.SalePurchaseBuyersRequest) models.EStampResponse {
	// Simulate stamp duty calculation for sale purchase
	purchasePrice := req.PurchasePrice
	considerationAmount := req.ConsiderationAmount
	if considerationAmount == 0 {
		considerationAmount = purchasePrice
	}

	// Simple calculation: 3% of purchase price or consideration amount, whichever is higher, or $5 whichever is higher
	var baseAmount float64
	if considerationAmount > purchasePrice {
		baseAmount = considerationAmount
	} else {
		baseAmount = purchasePrice
	}

	stampDuty := baseAmount * 0.03
	if stampDuty < 5.0 {
		stampDuty = 5.0
	}

	// Apply ABSD if applicable
	if req.IntentToClaimAbsdRefund == 1 {
		absdAmount := baseAmount * 0.20 // 20% ABSD
		stampDuty += absdAmount
	}

	// Generate mock document reference
	docRefNo := "SP" + ctrl.generateDocumentReference()

	return models.EStampResponse{
		ReturnCode: 10,
		Data: &models.EStampData{
			DocRefNo:        docRefNo,
			SDAmount:        fmt.Sprintf("%.2f", stampDuty),
			SDPenalty:       "0.00",
			TotalAmtPayable: fmt.Sprintf("%.2f", stampDuty),
			PaymentDueDate:  "2025-08-06", // 30 days from now
			PDFBase64:       ctrl.generateMockPDF("Sale Purchase Buyers", docRefNo),
		},
	}
}

func (ctrl *EStampController) processSalePurchaseSellers(req *models.SalePurchaseSellersRequest) models.EStampResponse {
	// Simulate stamp duty calculation for sale purchase sellers
	purchasePrice := req.PurchasePrice
	sellingPrice := req.SellingPrice
	considerationAmount := req.ConsiderationAmount

	// Use the highest value for calculation
	var baseAmount float64
	if sellingPrice > 0 && sellingPrice > purchasePrice && sellingPrice > considerationAmount {
		baseAmount = sellingPrice
	} else if considerationAmount > 0 && considerationAmount > purchasePrice {
		baseAmount = considerationAmount
	} else {
		baseAmount = purchasePrice
	}

	// Simple calculation: 3% of base amount or $5 whichever is higher
	stampDuty := baseAmount * 0.03
	if stampDuty < 5.0 {
		stampDuty = 5.0
	}

	// Apply ABSD if applicable
	if req.IntentToClaimAbsdRefund == 1 {
		absdAmount := baseAmount * 0.20 // 20% ABSD
		stampDuty += absdAmount
	}

	// Generate mock document reference
	docRefNo := "SPS" + ctrl.generateDocumentReference()

	return models.EStampResponse{
		ReturnCode: 10,
		Data: &models.EStampData{
			DocRefNo:        docRefNo,
			SDAmount:        fmt.Sprintf("%.2f", stampDuty),
			SDPenalty:       "0.00",
			TotalAmtPayable: fmt.Sprintf("%.2f", stampDuty),
			PaymentDueDate:  "2025-08-06", // 30 days from now
			PDFBase64:       ctrl.generateMockPDF("Sale Purchase Sellers", docRefNo),
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

// @Summary Stamp Certificate Authenticity Check
// @Description Check the authenticity of a stamp certificate
// @Tags Stamp
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param body body models.SCAuthenticityRequest true "Stamp Certificate Authenticity Request"
// @Success 200 {object} models.SCAuthenticityResponse
// @Router /iras/prod/SD/SCAuthenticity [post]
func (ctrl *EStampController) SCAuthenticity(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, models.SCAuthenticityResponse{
			ReturnCode: 400,
			Info: &models.SCAuthenticityInfo{
				Message:     "Missing required headers",
				MessageCode: 400,
				FieldInfoList: []models.SCAuthenticityFieldInfo{
					{Field: "X-IBM-Client-Id", Message: "Client ID is required"},
					{Field: "X-IBM-Client-Secret", Message: "Client Secret is required"},
				},
			},
		})
		return
	}

	// Parse request body
	var req models.SCAuthenticityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.SCAuthenticityResponse{
			ReturnCode: 400,
			Info: &models.SCAuthenticityInfo{
				Message:     "Invalid request body",
				MessageCode: 400,
				FieldInfoList: []models.SCAuthenticityFieldInfo{
					{Field: "request", Message: err.Error()},
				},
			},
		})
		return
	}

	// Validate required fields
	if req.DocRefNo == 0 {
		c.JSON(http.StatusBadRequest, models.SCAuthenticityResponse{
			ReturnCode: 400,
			Info: &models.SCAuthenticityInfo{
				Message:     "Validation failed",
				MessageCode: 400,
				FieldInfoList: []models.SCAuthenticityFieldInfo{
					{Field: "docRefNo", Message: "Document reference number is required"},
				},
			},
		})
		return
	}

	if req.StampCertRef == "" {
		c.JSON(http.StatusBadRequest, models.SCAuthenticityResponse{
			ReturnCode: 400,
			Info: &models.SCAuthenticityInfo{
				Message:     "Validation failed",
				MessageCode: 400,
				FieldInfoList: []models.SCAuthenticityFieldInfo{
					{Field: "stampCertRef", Message: "Stamp certificate reference is required"},
				},
			},
		})
		return
	}

	// Generate mock response data based on the API specification
	mockData := &models.SCAuthenticityData{
		AddBuyerSD:         54.30191592,
		AdjudicationFee:    54.64431772,
		AppRefNo:           "kobejcedfa",
		AssmtType:          "nihulejac",
		BuyerSD:            61.67180885,
		CertType:           "mofalpihimhunipk",
		DateOfDoc:          "2/10/2109",
		DocDescription:     "Tenancy Agreement for residential property",
		DocRefNo:           float64(req.DocRefNo),
		DocVerNo:           37.22779704,
		Duplicate:          15.08565901,
		Fines:              84.80979919,
		Penalty:            44.79689412,
		SDAmount:           47.94581197,
		Securities:         []string{"onre", "holapu", "vegha"},
		StampCertIssueDate: time.Now().AddDate(0, -1, 0).Format("2/1/2006"),
		StampCertRef:       req.StampCertRef,
		TotalAmtPayable:    90.0944656,
		ValuationFee:       84.73702907,
		PropertyList: []models.SCAuthPropertyData{
			{
				BlkHseNo:   "123A",
				PostalCode: "S123456",
				Street:     "Orchard Road",
				UnitLevel:  "#12-34",
			},
			{
				BlkHseNo:   "456B",
				PostalCode: "S654321",
				Street:     "Marina Bay Drive",
				UnitLevel:  "#05-67",
			},
		},
		StockSharesList: []models.SCAuthStockSharesData{
			{
				EntityID:       "698679549755392",
				EntityType:     "Private Limited Company",
				NameOfCompany:  "ABC Holdings Pte Ltd",
				NoStocksShares: 13.94657641,
			},
			{
				EntityID:       "5811371709038592",
				EntityType:     "Public Limited Company",
				NameOfCompany:  "XYZ Corporation Ltd",
				NoStocksShares: 19.96927301,
			},
		},
		VacantLandList: []models.SCAuthVacantLandData{
			{
				LotNo:        "Lot 123",
				MkTSNo:       "MK789",
				PlPTParcelNo: 23.16430078,
				StreetName:   "Sentosa Cove",
			},
			{
				LotNo:        "Lot 456",
				MkTSNo:       "MK012",
				PlPTParcelNo: 25.35278033,
				StreetName:   "Jurong Island",
			},
		},
	}

	// Return successful response
	response := models.SCAuthenticityResponse{
		ReturnCode: 200,
		Data:       mockData,
		Info: &models.SCAuthenticityInfo{
			Message:       "Stamp certificate authenticity check completed successfully",
			MessageCode:   200,
			FieldInfoList: []models.SCAuthenticityFieldInfo{},
		},
	}

	c.JSON(http.StatusOK, response)
}
