package controllers

import (
	"api-iras/internal/config"
	"api-iras/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CorpPassController struct{}

func NewCorpPassController() *CorpPassController {
	return &CorpPassController{}
}

// @Summary CorpPass Authentication
// @Description Initiate CorpPass authentication flow
// @Tags CorpPass
// @Accept json
// @Produce json
// @Param X-IBM-Client-Id header string true "Client ID"
// @Param X-IBM-Client-Secret header string true "Client Secret"
// @Param scope query string false "OAuth scope"
// @Param callback_url query string false "Callback URL"
// @Param state query string false "State parameter"
// @Param tax_agent query bool false "Tax agent flag"
// @Success 200 {object} models.CorpPassAuthResponse
// @Router /iras/sb/Authentication/CorpPassAuth [get]
func (ctrl *CorpPassController) CorpPassAuth(c *gin.Context) {
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
func (ctrl *CorpPassController) CorpPassToken(c *gin.Context) {
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
func (ctrl *CorpPassController) isValidCallbackURL(callbackURL string) bool {
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

func (ctrl *CorpPassController) generateCorpPassURL(scope, callbackURL, state string, taxAgent bool) string {
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

	// Add tax agent parameter to scope if applicable
	if taxAgent {
		scope += ",TaxAgent"
	}

	// Simulate CorpPass URL generation
	return baseURL + "?RequestBinding=HTTPArtifact&ResponseBinding=HTTPArtifact" +
		"&PartnerId=https%3A%2F%2Fstg-home.corppass.gov.sg%2Fconsent%2Firas-cp" +
		"&Target=https://stg-home.corppass.gov.sg/consent/oauth2/authorize" +
		"?realm=/consent/iras-cp&response_type=code&appName=IRASDemo" +
		"&state=" + state + "&client_id=iras&scope=" + scope +
		"&redirect_uri=" + callbackURL
}

func (ctrl *CorpPassController) generateAccessToken(id int64) string {
	// Generate mock access token
	return "corppass_access_token_" + strconv.FormatInt(id, 10) + "_demo_12345"
}

func (ctrl *CorpPassController) generateRefreshToken(id int64) string {
	// Generate mock refresh token
	return "corppass_refresh_token_" + strconv.FormatInt(id, 10) + "_demo_67890"
}
