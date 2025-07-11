package services

import (
	"api-iras/internal/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SingPassService struct {
	db *gorm.DB
}

func NewSingPassService(db *gorm.DB) *SingPassService {
	return &SingPassService{db: db}
}

// SingPassServiceAuth handles the GET /SingPassServiceAuth endpoint
func (s *SingPassService) SingPassServiceAuth(req *models.SingPassServiceAuthRequest) (*models.SingPassServiceAuthResponse, error) {
	// Validate callback URL if provided
	if req.CallbackURL != "" && !s.isValidCallbackURL(req.CallbackURL) {
		return &models.SingPassServiceAuthResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 850301,
				Message:     "Arguments Error",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "callback_url",
						Message: "The callback_url specified is not registered",
					},
				},
			},
		}, nil
	}

	// Generate state if not provided
	state := req.State
	if state == "" {
		state = uuid.New().String()
	}

	// Set default scope if not provided
	scope := req.Scope
	if scope == "" {
		scope = "GSTReturnsSub+GSTTransListSub"
	}

	// Set default callback URL if not provided
	callbackURL := req.CallbackURL
	if callbackURL == "" {
		callbackURL = "http://www.iras.gov.sg/callback"
	}

	// Generate SingPass authentication URL
	authURL := s.generateSingPassAuthURL(scope, callbackURL, state)

	// Store auth record in database
	authRecord := &models.SingPassAuthRecord{
		AuthURL:     authURL,
		State:       state,
		Scope:       scope,
		CallbackURL: callbackURL,
		ClientID:    "singpass-client",
		Status:      "pending",
	}

	if err := s.db.Create(authRecord).Error; err != nil {
		return &models.SingPassServiceAuthResponse{
			ReturnCode: 50,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 50001,
				Message:     "Internal server error",
			},
		}, fmt.Errorf("failed to save auth record: %w", err)
	}

	// Return successful response
	return &models.SingPassServiceAuthResponse{
		ReturnCode: 10,
		Data: &models.SingPassServiceAuthData{
			URL:   authURL,
			State: state,
		},
	}, nil
}

// SingPassServiceAuthToken handles the POST /SingPassServiceAuthToken endpoint
func (s *SingPassService) SingPassServiceAuthToken(req *models.SingPassServiceAuthTokenRequest) (*models.SingPassServiceAuthTokenResponse, error) {
	// Validate required fields
	if strings.TrimSpace(req.Code) == "" {
		return &models.SingPassServiceAuthTokenResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 40001,
				Message:     "Missing required field",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "code",
						Message: "Code is required",
					},
				},
			},
		}, nil
	}

	if strings.TrimSpace(req.State) == "" {
		return &models.SingPassServiceAuthTokenResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 40002,
				Message:     "Missing required field",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "state",
						Message: "State is required",
					},
				},
			},
		}, nil
	}

	if strings.TrimSpace(req.CallbackURL) == "" {
		return &models.SingPassServiceAuthTokenResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 40003,
				Message:     "Missing required field",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "callback_url",
						Message: "Callback URL is required",
					},
				},
			},
		}, nil
	}

	if strings.TrimSpace(req.Scope) == "" {
		return &models.SingPassServiceAuthTokenResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 40004,
				Message:     "Missing required field",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "scope",
						Message: "Scope is required",
					},
				},
			},
		}, nil
	}

	// Validate callback URL
	if !s.isValidCallbackURL(req.CallbackURL) {
		return &models.SingPassServiceAuthTokenResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 850301,
				Message:     "Arguments Error",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "callback_url",
						Message: "The callback_url specified is not registered",
					},
				},
			},
		}, nil
	}

	// Check if auth record exists for this state
	var authRecord models.SingPassAuthRecord
	err := s.db.Where("state = ?", req.State).First(&authRecord).Error
	if err != nil {
		return &models.SingPassServiceAuthTokenResponse{
			ReturnCode: 40,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 40005,
				Message:     "Invalid state",
				FieldInfoList: []models.SingPassServiceAuthFieldError{
					{
						Field:   "state",
						Message: "State not found or expired",
					},
				},
			},
		}, nil
	}

	// Generate access token
	accessToken := s.generateAccessToken()
	refreshToken := s.generateRefreshToken()

	// Create token record
	tokenRecord := &models.SingPassTokenRecord{
		Code:         req.Code,
		State:        req.State,
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		RefreshToken: refreshToken,
		Scope:        req.Scope,
		CallbackURL:  req.CallbackURL,
		ClientID:     "singpass-client",
		Status:       "active",
	}

	if err := s.db.Create(tokenRecord).Error; err != nil {
		return &models.SingPassServiceAuthTokenResponse{
			ReturnCode: 50,
			Info: &models.SingPassServiceAuthInfo{
				MessageCode: 50001,
				Message:     "Internal server error",
			},
		}, fmt.Errorf("failed to save token record: %w", err)
	}

	// Update auth record status
	s.db.Model(&authRecord).Update("status", "completed")

	// Return successful response
	return &models.SingPassServiceAuthTokenResponse{
		ReturnCode: 10,
		Data: &models.SingPassServiceAuthTokenData{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			RefreshToken: refreshToken,
			Scope:        req.Scope,
		},
	}, nil
}

// generateSingPassAuthURL generates the SingPass authentication URL
func (s *SingPassService) generateSingPassAuthURL(scope, callbackURL, state string) string {
	baseURL := "https://stg-saml.singpass.gov.sg/FIM/sps/SingpassIDPFed/saml20/logininitial"
	clientID := "a1234b5c-1234-abcd-efgh-a1234b5cdef"
	
	return fmt.Sprintf("%s?client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		baseURL, clientID, scope, callbackURL, state)
}

// generateAccessToken generates an access token
func (s *SingPassService) generateAccessToken() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("SP_AT_%s_%s", timestamp, uuid.New().String()[:8])
}

// generateRefreshToken generates a refresh token
func (s *SingPassService) generateRefreshToken() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("SP_RT_%s_%s", timestamp, uuid.New().String()[:8])
}

// isValidCallbackURL validates if the callback URL is registered
func (s *SingPassService) isValidCallbackURL(callbackURL string) bool {
	// List of valid callback URLs (in real implementation, this would come from database)
	validURLs := []string{
		"http://www.iras.gov.sg/callback",
		"https://www.iras.gov.sg/callback",
		"http://localhost:8090/callback",
		"https://localhost:8090/callback",
		"http://dirtor.mv/ma", // From the example in the spec
	}

	for _, validURL := range validURLs {
		if callbackURL == validURL {
			return true
		}
	}
	return false
}

// Admin CRUD operations for SingPass auth records

// CreateSingPassAuthRecord creates a new SingPass auth record
func (s *SingPassService) CreateSingPassAuthRecord(record *models.SingPassAuthRecord) error {
	return s.db.Create(record).Error
}

// GetSingPassAuthRecords retrieves SingPass auth records with pagination
func (s *SingPassService) GetSingPassAuthRecords(offset, limit int) ([]models.SingPassAuthRecord, int64, error) {
	var records []models.SingPassAuthRecord
	var total int64

	// Get total count
	if err := s.db.Model(&models.SingPassAuthRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get records with pagination
	if err := s.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetSingPassAuthRecordByID retrieves a SingPass auth record by ID
func (s *SingPassService) GetSingPassAuthRecordByID(id uint) (*models.SingPassAuthRecord, error) {
	var record models.SingPassAuthRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GetSingPassAuthRecordByState retrieves a SingPass auth record by state
func (s *SingPassService) GetSingPassAuthRecordByState(state string) (*models.SingPassAuthRecord, error) {
	var record models.SingPassAuthRecord
	if err := s.db.Where("state = ?", state).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// UpdateSingPassAuthRecord updates a SingPass auth record
func (s *SingPassService) UpdateSingPassAuthRecord(id uint, record *models.SingPassAuthRecord) error {
	return s.db.Model(&models.SingPassAuthRecord{}).Where("id = ?", id).Updates(record).Error
}

// DeleteSingPassAuthRecord deletes a SingPass auth record
func (s *SingPassService) DeleteSingPassAuthRecord(id uint) error {
	return s.db.Delete(&models.SingPassAuthRecord{}, id).Error
}

// Admin CRUD operations for SingPass token records

// CreateSingPassTokenRecord creates a new SingPass token record
func (s *SingPassService) CreateSingPassTokenRecord(record *models.SingPassTokenRecord) error {
	return s.db.Create(record).Error
}

// GetSingPassTokenRecords retrieves SingPass token records with pagination
func (s *SingPassService) GetSingPassTokenRecords(offset, limit int) ([]models.SingPassTokenRecord, int64, error) {
	var records []models.SingPassTokenRecord
	var total int64

	// Get total count
	if err := s.db.Model(&models.SingPassTokenRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get records with pagination
	if err := s.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetSingPassTokenRecordByID retrieves a SingPass token record by ID
func (s *SingPassService) GetSingPassTokenRecordByID(id uint) (*models.SingPassTokenRecord, error) {
	var record models.SingPassTokenRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GetSingPassTokenRecordByState retrieves a SingPass token record by state
func (s *SingPassService) GetSingPassTokenRecordByState(state string) (*models.SingPassTokenRecord, error) {
	var record models.SingPassTokenRecord
	if err := s.db.Where("state = ?", state).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// UpdateSingPassTokenRecord updates a SingPass token record
func (s *SingPassService) UpdateSingPassTokenRecord(id uint, record *models.SingPassTokenRecord) error {
	return s.db.Model(&models.SingPassTokenRecord{}).Where("id = ?", id).Updates(record).Error
}

// DeleteSingPassTokenRecord deletes a SingPass token record
func (s *SingPassService) DeleteSingPassTokenRecord(id uint) error {
	return s.db.Delete(&models.SingPassTokenRecord{}, id).Error
}
