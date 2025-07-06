package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"api-iras/internal/services"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

// CorpPass Authentication endpoints

// @Summary CorpPass Authentication
// @Description Initiate CorpPass authentication flow
// @Tags CorpPass
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param scope query string false "Scope"
// @Param callback_url query string false "Callback URL"
// @Param state query string false "State"
// @Param tax_agent query bool false "Tax Agent"
// @Success 200 {object} models.CorpPassAuthResponse
// @Router /iras/sb/Authentication/CorpPassAuth [get]
func (ctrl *GSTController) CorpPassAuth(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, models.CorpPassAuthResponse{
			ReturnCode: 40,
			Info: &models.CorpPassAuthInfo{
				MessageCode: "40003",
				Message:     "Missing required headers",
				FieldInfoList: []models.CorpPassFieldError{
					{
						Field:   "headers",
						Message: "X-IBM-Client-Id and X-IBM-Client-Secret are required",
					},
				},
			},
		})
		return
	}

	// Get query parameters
	scope := c.Query("scope")
	callbackURL := c.Query("callback_url")
	state := c.Query("state")
	taxAgent := c.Query("tax_agent") == "true"

	// Simulate validation - callback_url must be registered
	if callbackURL != "" && !ctrl.isValidCallbackURL(callbackURL) {
		c.JSON(http.StatusBadRequest, models.CorpPassAuthResponse{
			ReturnCode: 40,
			Info: &models.CorpPassAuthInfo{
				MessageCode: "850301",
				Message:     "Arguments Error",
				FieldInfoList: []models.CorpPassFieldError{
					{
						Field:   "callback_url",
						Message: "The callback_url specified is not registered",
					},
				},
			},
		})
		return
	}

	// Generate mock CorpPass URL
	corpPassURL := ctrl.generateCorpPassURL(scope, callbackURL, state, taxAgent)

	c.JSON(http.StatusOK, models.CorpPassAuthResponse{
		ReturnCode: 10,
		Data: &models.CorpPassAuthData{
			URL: corpPassURL,
		},
	})
}

// @Summary CorpPass Token
// @Description Exchange authorization code for access token
// @Tags CorpPass
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param body body models.CorpPassTokenRequest true "Token Request"
// @Success 200 {object} models.CorpPassTokenResponse
// @Router /iras/sb/Authentication/CorpPassToken [post]
func (ctrl *GSTController) CorpPassToken(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, models.CorpPassTokenResponse{
			ReturnCode: 40,
			Info: &models.CorpPassAuthInfo{
				MessageCode: "40003",
				Message:     "Missing required headers",
				FieldInfoList: []models.CorpPassFieldError{
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
	var req models.CorpPassTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.CorpPassTokenResponse{
			ReturnCode: 40,
			Info: &models.CorpPassAuthInfo{
				MessageCode: "40004",
				Message:     "Invalid request format",
				FieldInfoList: []models.CorpPassFieldError{
					{
						Field:   "body",
						Message: "Invalid JSON format",
					},
				},
			},
		})
		return
	}

	// Validate ID
	if req.ID <= 0 {
		c.JSON(http.StatusBadRequest, models.CorpPassTokenResponse{
			ReturnCode: 40,
			Info: &models.CorpPassAuthInfo{
				MessageCode: "40005",
				Message:     "Invalid ID",
				FieldInfoList: []models.CorpPassFieldError{
					{
						Field:   "id",
						Message: "ID must be a positive number",
					},
				},
			},
		})
		return
	}

	// Generate mock access token
	accessToken := ctrl.generateAccessToken(req.ID)

	c.JSON(http.StatusOK, models.CorpPassTokenResponse{
		ReturnCode: 10,
		Data: &models.CorpPassTokenData{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    3600, // 1 hour
			RefreshToken: ctrl.generateRefreshToken(req.ID),
		},
	})
}

// Helper methods for CorpPass simulation
func (ctrl *GSTController) isValidCallbackURL(callbackURL string) bool {
	// Simulate registered callback URLs
	validURLs := []string{
		"http://localhost:3000/callback",
		"https://abcpayroll.com/callback",
		"http://po.ec/vefocuf",
		"https://demo.example.com/callback",
	}

	for _, validURL := range validURLs {
		if callbackURL == validURL {
			return true
		}
	}
	return false
}

func (ctrl *GSTController) generateCorpPassURL(scope, callbackURL, state string, taxAgent bool) string {
	baseURL := "https://stg-saml.corppass.gov.sg/FIM/sps/CorpIDPFed/saml20/logininitial"

	// Default parameters for simulation
	if scope == "" {
		scope = "EmpIncomeSub"
	}
	if callbackURL == "" {
		callbackURL = "https://demo.example.com/callback"
	}
	if state == "" {
		state = "1234"
	}

	// Simulate CorpPass URL generation
	return baseURL + "?RequestBinding=HTTPArtifact&ResponseBinding=HTTPArtifact" +
		"&PartnerId=https%3A%2F%2Fstg-home.corppass.gov.sg%2Fconsent%2Firas-cp" +
		"&Target=https://stg-home.corppass.gov.sg/consent/oauth2/authorize" +
		"?realm=/consent/iras-cp&response_type=code&appName=IRASDemo" +
		"&state=" + state + "&client_id=iras&scope=" + scope +
		"&redirect_uri=" + callbackURL
}

func (ctrl *GSTController) generateAccessToken(id int64) string {
	// Generate mock access token
	return "corppass_access_token_" + strconv.FormatInt(id, 10) + "_demo_12345"
}

func (ctrl *GSTController) generateRefreshToken(id int64) string {
	// Generate mock refresh token
	return "corppass_refresh_token_" + strconv.FormatInt(id, 10) + "_demo_67890"
}

// eStamp endpoints

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
func (ctrl *GSTController) StampTenancyAgreement(c *gin.Context) {
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
func (ctrl *GSTController) ShareTransfer(c *gin.Context) {
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
func (ctrl *GSTController) processStampTenancyAgreement(req *models.StampTenancyAgreementRequest) models.EStampResponse {
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

func (ctrl *GSTController) processShareTransfer(req *models.ShareTransferRequest) models.EStampResponse {
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

func (ctrl *GSTController) generateDocumentReference() string {
	// Generate mock document reference number
	return fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(1000))
}

func (ctrl *GSTController) generateMockPDF(docType, docRefNo string) string {
	// Generate mock base64 PDF content
	mockPDFContent := fmt.Sprintf("PDF-1.4 Mock %s Document - Ref: %s - Generated: %s",
		docType, docRefNo, time.Now().Format("2006-01-02 15:04:05"))
	return base64.StdEncoding.EncodeToString([]byte(mockPDFContent))
}
