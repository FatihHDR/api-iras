package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common columns for all tables
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// GST Registration model for storing GST registration data
type GSTRegistration struct {
	BaseModel
	ClientID              string `json:"client_id" gorm:"not null" validate:"required"`
	RegistrationID        string `json:"registration_id" gorm:"uniqueIndex;not null" validate:"required"`
	GSTRegistrationNumber string `json:"gst_registration_number" gorm:"index"`
	Name                  string `json:"name" gorm:"type:text"`
	RegisteredFrom        string `json:"registered_from"`
	RegisteredTo          string `json:"registered_to"`
	Status                string `json:"status"`
	Remarks               string `json:"remarks" gorm:"type:text"`
}

// GST Response models based on IRAS API spec
type GSTResponse struct {
	ReturnCode int      `json:"returnCode"`
	Data       *GSTData `json:"data,omitempty"`
	Info       *GSTInfo `json:"info,omitempty"`
}

type GSTData struct {
	Name                  string `json:"name"`
	GSTRegistrationNumber string `json:"gstRegistrationNumber"`
	RegistrationID        string `json:"registrationId"`
	RegisteredFrom        string `json:"RegisteredFrom"`
	RegisteredTo          string `json:"RegisteredTo"`
	Remarks               string `json:"Remarks"`
	Status                string `json:"Status"`
}

type GSTInfo struct {
	Message       string          `json:"message"`
	MessageCode   int             `json:"messageCode"`
	FieldInfoList []GSTFieldError `json:"fieldInfoList,omitempty"`
}

type GSTFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// GST Request models based on IRAS API spec
type GSTRequest struct {
	ClientID string `json:"clientID" validate:"required"`
	RegID    string `json:"regID" validate:"required"`
}

// User model for authentication
type User struct {
	BaseModel
	Name     string `json:"name" gorm:"type:varchar(255)" validate:"required,min=2,max=100"`
	Username string `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50"`
	Email    string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password string `json:"-" gorm:"not null" validate:"required,min=6"`
	Role     string `json:"role" gorm:"default:user" validate:"oneof=admin user"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

// Auth Request models
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Auth Response models
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresIn int       `json:"expires_in"`
	UserInfo  *UserInfo `json:"user_info"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// API Response models
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// CorpPass Authentication models based on IRAS API spec
type CorpPassAuthResponse struct {
	ReturnCode int               `json:"returnCode"`
	Data       *CorpPassAuthData `json:"data,omitempty"`
	Info       *CorpPassAuthInfo `json:"info,omitempty"`
}

type CorpPassAuthData struct {
	URL string `json:"url"`
}

type CorpPassAuthInfo struct {
	MessageCode   string               `json:"messageCode"`
	Message       string               `json:"message"`
	FieldInfoList []CorpPassFieldError `json:"fieldInfoList,omitempty"`
}

type CorpPassFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// CorpPass Token Request/Response models
type CorpPassTokenRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type CorpPassTokenResponse struct {
	ReturnCode int                `json:"returnCode"`
	Data       *CorpPassTokenData `json:"data,omitempty"`
	Info       *CorpPassAuthInfo  `json:"info,omitempty"`
}

type CorpPassTokenData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
