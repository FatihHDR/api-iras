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

// eStamp models based on IRAS API spec
type EStampResponse struct {
	ReturnCode int         `json:"returnCode"`
	Data       *EStampData `json:"data,omitempty"`
	Info       *EStampInfo `json:"info,omitempty"`
}

type EStampData struct {
	DocRefNo        string `json:"docRefNo"`
	SDAmount        string `json:"sdAmount"`
	SDPenalty       string `json:"sdPenalty"`
	TotalAmtPayable string `json:"totalAmtPayable"`
	PaymentDueDate  string `json:"paymentDueDate"`
	PDFBase64       string `json:"pdfBase64"`
}

type EStampInfo struct {
	Message       string             `json:"message"`
	MessageCode   int                `json:"messageCode"`
	FieldInfoList []EStampFieldError `json:"fieldInfoList,omitempty"`
}

type EStampFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Simplified eStamp request models (for demo purposes)
type StampTenancyAgreementRequest struct {
	AssignID                string                    `json:"assignId" validate:"required"`
	FTReferenceNo           string                    `json:"ftReferenceNo"`
	DocumentDescription     string                    `json:"documentDescription" validate:"required"`
	DocumentReferenceNo     string                    `json:"documentReferenceNo"`
	FormatOfDocument        int                       `json:"formatOfDocument"`
	ModeOfOffer             int                       `json:"modeOfOffer"`
	ModeOfAcceptance        int                       `json:"modeOfAcceptance"`
	IsSignedInSingapore     bool                      `json:"isSignedInSingapore"`
	DateOfDocument          string                    `json:"dateOfDocument"`
	ReceivingDateOfDocument string                    `json:"receivingDateOfDocument"`
	Submission              SubmissionData            `json:"submission"`
	Assets                  AssetsData                `json:"assets"`
	AssessmentRental        []AssessmentRentalData    `json:"assessmentRental"`
	LandlordLessor          []PartyData               `json:"landlordLessor"`
	TenantLessee            []PartyData               `json:"tenantLessee"`
	AssessmentRemissions    []AssessmentRemissionData `json:"assessmentRemissions"`
}

type ShareTransferRequest struct {
	AssignID                                        string                    `json:"assignId" validate:"required"`
	FTReferenceNo                                   string                    `json:"ftReferenceNo"`
	DocumentDescription                             string                    `json:"documentDescription" validate:"required"`
	DocumentReferenceNo                             string                    `json:"documentReferenceNo"`
	FormatOfDocument                                int                       `json:"formatOfDocument"`
	ModeOfOffer                                     int                       `json:"modeOfOffer"`
	ModeOfAcceptance                                int                       `json:"modeOfAcceptance"`
	IsSignedInSingapore                             bool                      `json:"isSignedInSingapore"`
	DateOfDocument                                  string                    `json:"dateOfDocument"`
	ReceivingDateOfDocument                         string                    `json:"receivingDateOfDocument"`
	Submission                                      SubmissionData            `json:"submission"`
	Assets                                          AssetsData                `json:"assets"`
	ConsiderationAmount                             float64                   `json:"considerationAmount"`
	HasIntentToHoldTheSharesTrustForBeneficialOwner bool                      `json:"hasIntentToHoldTheSharesTrustForBeneficialOwner"`
	Transferor                                      []PartyData               `json:"transferor"`
	Transferee                                      []PartyData               `json:"transferee"`
	TargetCompany                                   TargetCompanyData         `json:"targetCompany"`
	AssessmentRemissions                            []AssessmentRemissionData `json:"assessmentRemissions"`
}

// Supporting data structures (simplified for demo)
type SubmissionData struct {
	Declaration bool `json:"declaration"`
}

type AssetsData struct {
	Properties   []PropertyData   `json:"properties"`
	Lands        []LandData       `json:"lands"`
	StocksShares []StockShareData `json:"stocksShares"`
	Security     []SecurityData   `json:"security"`
}

type PropertyData struct {
	Sequence                             int     `json:"sequence"`
	PostalCode                           string  `json:"postalCode"`
	StreetName                           string  `json:"streetName"`
	BlockNo                              string  `json:"blockNo"`
	PropertyType                         int     `json:"propertyType"`
	BuyingPriceMarketValueResidential    float64 `json:"buyingPriceMarketValueResidential"`
	BuyingPriceMarketValueNonResidential float64 `json:"buyingPriceMarketValueNonResidential"`
	TotalFloorArea                       float64 `json:"totalFloorArea"`
	ValuationType                        int     `json:"valuationType"`
	ValuationValue                       float64 `json:"valuationValue"`
}

type LandData struct {
	Sequence                             int     `json:"sequence"`
	LandIDType                           int     `json:"landIdType"`
	MKOrTSNo                             string  `json:"mkOrTSNo"`
	StreetName                           string  `json:"streetName"`
	BuyingPriceMarketValueResidential    float64 `json:"buyingPriceMarketValueResidential"`
	BuyingPriceMarketValueNonResidential float64 `json:"buyingPriceMarketValueNonResidential"`
	ValuationType                        int     `json:"valuationType"`
	ValuationValue                       float64 `json:"valuationValue"`
}

type StockShareData struct {
	Sequence      int    `json:"sequence"`
	EntityType    int    `json:"entityType"`
	TaxEntityID   string `json:"taxEntityId"`
	FTCompanyName string `json:"ftCompanyName"`
	NoStockShares int    `json:"noStockShares"`
}

type SecurityData struct {
	Sequence      int    `json:"sequence"`
	FTDescription string `json:"ftDescription"`
}

type AssessmentRentalData struct {
	IsPremiumConsiderationMade    bool               `json:"isPremiumConsiderationMade"`
	PremiumConsiderationAmount    float64            `json:"premiumConsiderationAmount"`
	ResidentialComponentAmount    float64            `json:"residentialComponentAmount"`
	NonResidentialComponentAmount float64            `json:"nonResidentialComponentAmount"`
	IsMonthlyRentPayable          bool               `json:"isMonthlyRentPayable"`
	RentalDetails                 []RentalDetailData `json:"rentalDetails"`
	TotalGrossRentAmount          float64            `json:"totalGrossRentAmount"`
	AverageRentAmount             float64            `json:"averageRentAmount"`
}

type RentalDetailData struct {
	StartPeriodOfLease string  `json:"startPeriodOfLease"`
	EndPeriodOfLease   string  `json:"endPeriodOfLease"`
	RentAmount         float64 `json:"rentAmount"`
	MarketRentalValue  float64 `json:"marketRentalValue"`
}

type PartyData struct {
	Sequence             int                `json:"sequence"`
	TypeOfProfile        int                `json:"typeOfProfile"`
	TaxEntityIDType      int                `json:"taxEntityIdType"`
	TaxEntityIDNo        string             `json:"taxEntityIdNo"`
	FTTaxEntityName      string             `json:"ftTaxEntityName"`
	Gender               int                `json:"gender"`
	DateOfBirth          string             `json:"dateOfBirth"`
	MailingAddress       MailingAddressData `json:"mailingAddress"`
	PartyType            string             `json:"partyType"`
	CountryOfNationality string             `json:"countryOfNationality"`
	IsVerifiedEntity     bool               `json:"isVerifiedEntity"`
	IsLiableParty        bool               `json:"isLiableParty"`
}

type MailingAddressData struct {
	MailingAddressType int    `json:"mailingAddressType"`
	Country            string `json:"country"`
	PostalCode         string `json:"postalCode"`
	BlockNo            string `json:"blockNo"`
	StreetName         string `json:"streetName"`
	FloorNo            string `json:"floorNo"`
	UnitNo             string `json:"unitNo"`
}

type TargetCompanyData struct {
	EntityType            int     `json:"entityType"`
	TaxEntityIDNo         string  `json:"taxEntityIdNo"`
	FTCompanyName         string  `json:"ftCompanyName"`
	DateOfIncorporation   string  `json:"dateOfIncorporation"`
	CompanyType           string  `json:"companyType"`
	MarketPricePerShare   float64 `json:"marketPricePerShare"`
	NoOfSharesTransferred int     `json:"noOfSharesTransferred"`
	TotalMarketPrice      float64 `json:"totalMarketPrice"`
	NetAssetValue         float64 `json:"netAssetValue"`
	HasOneClassOfShares   bool    `json:"hasOneClassOfShares"`
}

type AssessmentRemissionData struct {
	Sequence             int    `json:"sequence"`
	RemissionType        string `json:"remissionType"`
	RemissionOptionText1 string `json:"remissionOptionText1"`
}
