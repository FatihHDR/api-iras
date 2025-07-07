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
	Sequence                              int                     `json:"sequence"`
	PostalCode                            string                  `json:"postalCode"`
	StreetName                            string                  `json:"streetName"`
	BlockNo                               string                  `json:"blockNo"`
	LevelUnits                            []LevelUnitData         `json:"levelUnits,omitempty"`
	ShareOfPropertyTransferred            int                     `json:"shareOfPropertyTransferred,omitempty"`
	MannerOfHolding                       int                     `json:"mannerOfHolding,omitempty"`
	FractionNumerator                     int                     `json:"fractionNumerator,omitempty"`
	FractionDenominator                   int                     `json:"fractionDenominator,omitempty"`
	PropertyOwnerships                    []PropertyOwnershipData `json:"propertyOwnerships,omitempty"`
	PropertyType                          int                     `json:"propertyType"`
	BuyingPriceMarketValueResidential     float64                 `json:"buyingPriceMarketValueResidential"`
	BuyingPriceMarketValueNonResidential  float64                 `json:"buyingPriceMarketValueNonResidential"`
	SellingPriceMarketValueResidential    float64                 `json:"sellingPriceMarketValueResidential,omitempty"`
	SellingPriceMarketValueNonResidential float64                 `json:"sellingPriceMarketValueNonResidential,omitempty"`
	IsWhollyRented                        bool                    `json:"isWhollyRented,omitempty"`
	TotalFloorArea                        float64                 `json:"totalFloorArea"`
	FloorAreaMeasurementType              int                     `json:"floorAreaMeasurementType,omitempty"`
	TotalFloorAreaNotAvailable            bool                    `json:"totalFloorAreaNotAvailable,omitempty"`
	ValuationType                         int                     `json:"valuationType"`
	ValuationValue                        float64                 `json:"valuationValue"`
	AbsdRate                              float64                 `json:"absdRate,omitempty"`
}

type LandData struct {
	Sequence                              int                 `json:"sequence"`
	LandIDType                            int                 `json:"landIdType"`
	MKOrTSNo                              string              `json:"mkOrTSNo"`
	StreetName                            string              `json:"streetName"`
	FTLotNo                               string              `json:"ftLotNo,omitempty"`
	FTPlOrPtParcelNo                      string              `json:"ftPlOrPtParcelNo,omitempty"`
	ShareOfLandTransferred                int                 `json:"shareOfLandTransferred,omitempty"`
	FractionNumerator                     int                 `json:"fractionNumerator,omitempty"`
	FractionDenominator                   int                 `json:"fractionDenominator,omitempty"`
	MannerOfHolding                       int                 `json:"mannerOfHolding,omitempty"`
	MasterPlanZoning                      int                 `json:"masterPlanZoning,omitempty"`
	BuyingPriceMarketValueResidential     float64             `json:"buyingPriceMarketValueResidential"`
	BuyingPriceMarketValueNonResidential  float64             `json:"buyingPriceMarketValueNonResidential"`
	SellingPriceMarketValueResidential    float64             `json:"sellingPriceMarketValueResidential,omitempty"`
	SellingPriceMarketValueNonResidential float64             `json:"sellingPriceMarketValueNonResidential,omitempty"`
	LandOwnership                         []LandOwnershipData `json:"landOwnership,omitempty"`
	ValuationType                         int                 `json:"valuationType"`
	ValuationValue                        float64             `json:"valuationValue"`
	AbsdRate                              float64             `json:"absdRate,omitempty"`
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

// Stamp Mortgage Models
type StampMortgageRequest struct {
	AssignID                  string              `json:"assignId" binding:"required"`
	FTReferenceNo             string              `json:"ftReferenceNo"`
	DocumentDescription       string              `json:"documentDescription" binding:"required"`
	DocumentReferenceNo       string              `json:"documentReferenceNo"`
	FormatOfDocument          int                 `json:"formatOfDocument"`
	ModeOfOffer               int                 `json:"modeOfOffer"`
	ModeOfAcceptance          int                 `json:"modeOfAcceptance"`
	IsSignedInSingapore       bool                `json:"isSignedInSingapore"`
	DateOfDocument            string              `json:"dateOfDocument"`
	ReceivingDateOfDocument   string              `json:"receivingDateOfDocument"`
	IsDocumentDateUnavailable bool                `json:"isDocumentDateUnavailable"`
	Submission                SubmissionData      `json:"submission"`
	Assets                    MortgageAssetsData  `json:"assets"`
	TypeOfMortgage            int                 `json:"typeOfMortgage"`
	AmountOfLoan              float64             `json:"amountOfLoan" binding:"required"`
	Mortgagors                []MortgagePartyData `json:"mortgagors"`
	Mortgagees                []MortgagePartyData `json:"mortgagees"`
}

type MortgageAssetsData struct {
	Properties   []PropertyData   `json:"properties"`
	Lands        []LandData       `json:"lands"`
	StocksShares []StockShareData `json:"stocksShares"`
	Security     []SecurityData   `json:"security"`
}

type MortgagePartyData struct {
	Sequence         int                `json:"sequence"`
	TypeOfProfile    int                `json:"typeOfProfile"`
	TaxEntityIDType  int                `json:"taxEntityIdType"`
	TaxEntityIDNo    string             `json:"taxEntityIdNo"`
	FTTaxEntityName  string             `json:"ftTaxEntityName"`
	Gender           int                `json:"gender"`
	DateOfBirth      string             `json:"dateOfBirth"`
	MailingAddress   MailingAddressData `json:"mailingAddress"`
	PartyType        string             `json:"partyType"`
	IsVerifiedEntity bool               `json:"isVerifiedEntity"`
	IsLiableParty    bool               `json:"isLiableParty"`
}

// Sale Purchase Buyers Models
type SalePurchaseBuyersRequest struct {
	AssignID                                        string                    `json:"assignId" binding:"required"`
	FTReferenceNo                                   string                    `json:"ftReferenceNo"`
	DocumentDescription                             string                    `json:"documentDescription" binding:"required"`
	DocumentReferenceNo                             string                    `json:"documentReferenceNo"`
	FormatOfDocument                                int                       `json:"formatOfDocument"`
	ModeOfOffer                                     int                       `json:"modeOfOffer"`
	ModeOfAcceptance                                int                       `json:"modeOfAcceptance"`
	IsSignedInSingapore                             bool                      `json:"isSignedInSingapore"`
	DateOfDocument                                  string                    `json:"dateOfDocument"`
	ReceivingDateOfDocument                         string                    `json:"receivingDateOfDocument"`
	IsDocumentDateUnavailable                       bool                      `json:"isDocumentDateUnavailable"`
	Submission                                      SubmissionData            `json:"submission"`
	Assets                                          SalePurchaseAssetsData    `json:"assets"`
	RetrievedStampingRecord                         string                    `json:"retrievedStampingRecord"`
	HasIntentToHoldThePropInTrustForBeneficialOwner bool                      `json:"hasIntentToHoldThePropInTrustForBeneficialOwner"`
	HasIntentToTransferPropViaConveyanceDirection   bool                      `json:"hasIntentToTransferPropViaConveyanceDirection"`
	PurchasePrice                                   float64                   `json:"purchasePrice" binding:"required"`
	AnyConsiderationPaid                            float64                   `json:"anyConsiderationPaid"`
	ConsiderationAmount                             float64                   `json:"considerationAmount"`
	ConservancyCharges                              float64                   `json:"conservancyCharges"`
	ConservancyChargesUnitOfMeasure                 int                       `json:"conservancyChargesUnitOfMeasure"`
	TrusteeBeneficiaryDetails                       TrusteeBeneficiaryData    `json:"trusteeBeneficiaryDetails"`
	SellerTransferor                                []SalePurchasePartyData   `json:"sellerTransferor"`
	BuyerTransferee                                 []SalePurchasePartyData   `json:"buyerTransferee"`
	Trustee                                         []SalePurchasePartyData   `json:"trustee"`
	Beneficiary                                     []SalePurchasePartyData   `json:"beneficiary"`
	AssessmentRemissions                            []AssessmentRemissionData `json:"assessmentRemissions"`
	PropertyBuyers                                  []PropertyBuyerData       `json:"propertyBuyers"`
	IntentToClaimAbsdRefund                         int                       `json:"intentToClaimAbsdRefund"`
	AbsdRefundAssets                                AbsdRefundAssetsData      `json:"absdRefundAssets"`
	MaritalStatus                                   MaritalStatusData         `json:"maritalStatus"`
}

type SalePurchaseAssetsData struct {
	Properties   []PropertyData   `json:"properties"`
	Lands        []LandData       `json:"lands"`
	StocksShares []StockShareData `json:"stocksShares"`
	Security     []SecurityData   `json:"security"`
}

type LevelUnitData struct {
	FloorNo                           string  `json:"floorNo"`
	UnitNo                            string  `json:"unitNo"`
	AbsdRate                          float64 `json:"absdRate"`
	BuyingPriceMarketValueResidential float64 `json:"buyingPriceMarketValueResidential"`
}

type PropertyOwnershipData struct {
	Sequence         int `json:"sequence"`
	ShareNumerator   int `json:"shareNumerator"`
	ShareDenominator int `json:"shareDenominator"`
}

type LandOwnershipData struct {
	Sequence         int `json:"sequence"`
	ShareNumerator   int `json:"shareNumerator"`
	ShareDenominator int `json:"shareDenominator"`
}

type SalePurchasePartyData struct {
	Sequence                    int                `json:"sequence"`
	TypeOfProfile               int                `json:"typeOfProfile"`
	TaxEntityIDType             int                `json:"taxEntityIdType"`
	TaxEntityIDNo               string             `json:"taxEntityIdNo"`
	FTTaxEntityName             string             `json:"ftTaxEntityName"`
	Gender                      int                `json:"gender"`
	DateOfBirth                 string             `json:"dateOfBirth"`
	IsSubFundOwnerLessee        bool               `json:"isSubFundOwnerLessee"`
	MailingAddress              MailingAddressData `json:"mailingAddress"`
	PartyType                   string             `json:"partyType"`
	CountryOfNationality        string             `json:"countryOfNationality"`
	IsVerifiedEntity            bool               `json:"isVerifiedEntity"`
	BeneficiaryIsUnidentifiable bool               `json:"beneficiaryIsUnidentifiable"`
	FTDescription               string             `json:"ftDescription"`
	IsLiableParty               bool               `json:"isLiableParty"`
}

type TrusteeBeneficiaryData struct {
	RelationshipBetweenTrusteeAndBeneficiary int                `json:"relationshipBetweenTrusteeAndBeneficiary"`
	FTRelationshipOthersValue                string             `json:"ftRelationshipOthersValue"`
	Reasons                                  TrusteeReasonsData `json:"reasons"`
}

type TrusteeReasonsData struct {
	IsReasonBeneficiaryIsMinor           bool   `json:"isReasonBeneficiaryIsMinor"`
	IsReasonBeneficiaryIsNotLegalEntity  bool   `json:"isReasonBeneficiaryIsNotLegalEntity"`
	IsReasonEstatePlanning               bool   `json:"isReasonEstatePlanning"`
	IsReasonPursuantToNomineeArrangement bool   `json:"isReasonPursuantToNomineeArrangement"`
	IsReasonOthers                       bool   `json:"isReasonOthers"`
	FTReasonsOthersValue                 string `json:"ftReasonsOthersValue"`
}

type PropertyBuyerData struct {
	TypeOfProfile             int    `json:"typeOfProfile"`
	TaxEntityIDType           int    `json:"taxEntityIdType"`
	TaxEntityIDNo             string `json:"taxEntityIdNo"`
	PartyType                 string `json:"partyType"`
	NoOfProperties            int    `json:"noOfProperties"`
	BuyerName                 string `json:"buyerName"`
	Sequence                  int    `json:"sequence"`
	AssessmentPartiesSequence int    `json:"assessmentPartiesSequence"`
}

type AbsdRefundAssetsData struct {
	Assets                AbsdRefundAssetsDetail `json:"assets"`
	CompletedProperty     bool                   `json:"completedProperty"`
	PaymentMode           int                    `json:"paymentMode"`
	FTCpfHdbRefno         string                 `json:"ftCpfHdbRefno"`
	CashRecipientSequence int                    `json:"cashRecipientSequence"`
	CpfRecipient          []CpfRecipientData     `json:"cpfRecipient"`
	Declaration           bool                   `json:"declaration"`
}

type AbsdRefundAssetsDetail struct {
	Properties     []PropertyData   `json:"properties"`
	Lands          []LandData       `json:"lands"`
	StocksShares   []StockShareData `json:"stocksShares"`
	Security       []SecurityData   `json:"security"`
	AbsdProperties []PropertyData   `json:"absdProperties"`
	AbsdLands      []LandData       `json:"absdLands"`
}

type CpfRecipientData struct {
	CpfAmountUsed float64 `json:"cpfAmountUsed"`
	Sequence      int     `json:"sequence"`
}

type MaritalStatusData struct {
	MarriedCouple int `json:"marriedCouple"`
}

// Sale Purchase Sellers Models
type SalePurchaseSellersRequest struct {
	AssignID                        string                        `json:"assignId" binding:"required"`
	FTReferenceNo                   string                        `json:"ftReferenceNo"`
	DocumentDescription             string                        `json:"documentDescription" binding:"required"`
	DocumentReferenceNo             string                        `json:"documentReferenceNo"`
	FormatOfDocument                int                           `json:"formatOfDocument"`
	ModeOfOffer                     int                           `json:"modeOfOffer"`
	FTModeOfOfferOthersValue        string                        `json:"ftModeOfOfferOthersValue"`
	ModeOfAcceptance                int                           `json:"modeOfAcceptance"`
	FTModeOfAcceptanceOthersValue   string                        `json:"ftModeOfAcceptanceOthersValue"`
	IsSignedInSingapore             bool                          `json:"isSignedInSingapore"`
	DateOfDocument                  string                        `json:"dateOfDocument"`
	ReceivingDateOfDocument         string                        `json:"receivingDateOfDocument"`
	Submission                      SubmissionData                `json:"submission"`
	IsDocumentDateUnavailable       bool                          `json:"isDocumentDateUnavailable"`
	Assets                          SalePurchaseAssetsData        `json:"assets"`
	RetrievedStampingRecord         string                        `json:"retrievedStampingRecord"`
	PurchasePrice                   float64                       `json:"purchasePrice"`
	AnyConsiderationPaid            bool                          `json:"anyConsiderationPaid"`
	ConsiderationAmount             float64                       `json:"considerationAmount"`
	ConservancyCharge               float64                       `json:"conservancyCharge"`
	SellingPrice                    float64                       `json:"sellingPrice"`
	ConservancyChargesUnitOfMeasure int                           `json:"conservancyChargesUnitOfMeasure"`
	TrusteeBeneficiaryDetails       TrusteeBeneficiaryData        `json:"trusteeBeneficiaryDetails"`
	SellerTransferor                []SalePurchaseAdvancedParty   `json:"sellerTransferor"`
	BuyerTransferee                 []SalePurchaseAdvancedParty   `json:"buyerTransferee"`
	Trustee                         []SalePurchaseAdvancedParty   `json:"trustee"`
	Beneficiary                     []SalePurchaseAdvancedParty   `json:"beneficiary"`
	TransferorInitialPurchaser      []SalePurchaseAdvancedParty   `json:"transferorInitialPurchaser"`
	Transferee                      []SalePurchaseAdvancedParty   `json:"transferee"`
	AssessmentRemissions            []AdvancedAssessmentRemission `json:"assessmentRemissions"`
	PropertyBuyers                  []PropertyBuyerData           `json:"propertyBuyers"`
	IntentToClaimAbsdRefund         int                           `json:"intentToClaimAbsdRefund"`
	DateOfAcquisition               string                        `json:"dateOfAcquisition"`
}

// Advanced party structure for SalePurchaseSellers with additional fields
type SalePurchaseAdvancedParty struct {
	Sequence                    int                        `json:"sequence"`
	TypeOfProfile               int                        `json:"typeOfProfile"`
	TaxEntityIDType             int                        `json:"taxEntityIdType"`
	TaxEntityIDNo               string                     `json:"taxEntityIdNo"`
	FTTaxEntityName             string                     `json:"ftTaxEntityName"`
	Gender                      int                        `json:"gender"`
	DateOfBirth                 string                     `json:"dateOfBirth"`
	IsSubFundOwnerLessee        bool                       `json:"isSubFundOwnerLessee"`
	SubFunds                    []SubFundData              `json:"subFunds"`
	MailingAddress              AdvancedMailingAddressData `json:"mailingAddress"`
	PartyType                   string                     `json:"partyType"`
	CountryOfNationality        string                     `json:"countryOfNationality"`
	IsVerifiedEntity            bool                       `json:"isVerifiedEntity"`
	PartyRelationshipDetails    []PartyRelationshipData    `json:"partyRelationshipDetails"`
	BeneficiaryIsUnidentifiable bool                       `json:"beneficiaryIsUnidentifiable"`
	FTDescription               string                     `json:"ftDescription"`
	IsLiableParty               bool                       `json:"isLiableParty"`
	Lawyer                      []LawyerData               `json:"lawyer"`
}

// SubFund data structure
type SubFundData struct {
	SubFundNo string `json:"subFundNo"`
	FTName    string `json:"ftName"`
	IsOthers  bool   `json:"isOthers"`
}

// Advanced mailing address with additional fields
type AdvancedMailingAddressData struct {
	MailingAddressType int    `json:"mailingAddressType"`
	Country            string `json:"country"`
	PostalCode         string `json:"postalCode"`
	BlockNo            string `json:"blockNo"`
	StreetName         string `json:"streetName"`
	FloorNo            string `json:"floorNo"`
	UnitNo             string `json:"unitNo"`
	FTAddressLine1     string `json:"ftAddressLine1"`
	FTAddressLine2     string `json:"ftAddressLine2"`
	FTAddressLine3     string `json:"ftAddressLine3"`
	FTAddressLine4     string `json:"ftAddressLine4"`
}

// Party relationship details
type PartyRelationshipData struct {
	TypeOfProfile    int    `json:"typeOfProfile"`
	IdentityType     int    `json:"identityType"`
	IdentityNo       string `json:"identityNo"`
	RelationshipType int    `json:"relationshipType"`
	Sequence         int    `json:"sequence"`
}

// Lawyer data structure
type LawyerData struct {
	TaxEntityIDType int    `json:"taxEntityIdType"`
	TaxEntityIDNo   string `json:"taxEntityIdNo"`
	FTName          string `json:"ftName"`
}

// Advanced assessment remission structure
type AdvancedAssessmentRemission struct {
	Sequence               int                 `json:"sequence"`
	RemissionType          string              `json:"remissionType"`
	RemissionOption1       RemissionOptionData `json:"remissionOption1"`
	RemissionOption2       RemissionOptionData `json:"remissionOption2"`
	RemissionOption3       RemissionOptionData `json:"remissionOption3"`
	RemissionOption4       RemissionOptionData `json:"remissionOption4"`
	RemissionOption5       RemissionOptionData `json:"remissionOption5"`
	RemissionOptionText1   string              `json:"remissionOptionText1"`
	FTRemissionOptionText2 string              `json:"ftRemissionOptionText2"`
}

// Remission option data structure
type RemissionOptionData struct {
	SelectedOption    int `json:"selectedOption"`
	SelectedSubOption int `json:"selectedSubOption"`
}

// AIS Organization Search models based on IRAS API spec
type AISorgSearchRequest struct {
	ClientID       string `json:"clientID" validate:"required"`
	OrganizationID string `json:"organizationID" validate:"required"`
	BasisYear      int    `json:"basisYear" validate:"required"`
}

type AISorgSearchResponse struct {
	ReturnCode int         `json:"returnCode"`
	Data       *AISorgData `json:"data,omitempty"`
	Info       *AISorgInfo `json:"info,omitempty"`
}

type AISorgData struct {
	OrganizationInAIS string `json:"organizationInAIS"`
}

type AISorgInfo struct {
	MessageCode   string             `json:"messageCode"`
	Message       string             `json:"message"`
	FieldInfoList []AISorgFieldError `json:"fieldInfoList,omitempty"`
}

type AISorgFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Property Consolidated Statement models based on IRAS API spec
type PropertyConsolidatedStatementRequest struct {
	RefNo          string `json:"refNo" validate:"required"`
	PropertyTaxRef string `json:"propertyTaxRef" validate:"required"`
}

type PropertyConsolidatedStatementResponse struct {
	ReturnCode int                                `json:"returnCode"`
	Data       *PropertyConsolidatedStatementData `json:"data,omitempty"`
	Info       *PropertyConsolidatedStatementInfo `json:"info,omitempty"`
}

type PropertyConsolidatedStatementData struct {
	RefNo                 string                 `json:"refNo"`
	PropertyTaxRef        string                 `json:"propertyTaxRef"`
	ConsolidatedStatement *ConsolidatedStatement `json:"consolidatedStatement,omitempty"`
}

type ConsolidatedStatement struct {
	StatementDate   string               `json:"statementDate"`
	TotalAmount     string               `json:"totalAmount"`
	PropertyDetails []PropertyDetail     `json:"propertyDetails"`
	PaymentHistory  []PaymentHistoryItem `json:"paymentHistory"`
}

type PropertyDetail struct {
	PropertyID   string `json:"propertyId"`
	Address      string `json:"address"`
	PropertyType string `json:"propertyType"`
	TaxAmount    string `json:"taxAmount"`
	DueDate      string `json:"dueDate"`
	Status       string `json:"status"`
}

type PaymentHistoryItem struct {
	PaymentDate    string `json:"paymentDate"`
	Amount         string `json:"amount"`
	PaymentMethod  string `json:"paymentMethod"`
	TransactionRef string `json:"transactionRef"`
}

type PropertyConsolidatedStatementInfo struct {
	Message       string                           `json:"message"`
	MessageCode   int                              `json:"messageCode"`
	FieldInfoList []PropertyConsolidatedFieldError `json:"fieldInfoList,omitempty"`
}

type PropertyConsolidatedFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Property Consolidated Statement storage model for database
type PropertyConsolidatedStatementRecord struct {
	BaseModel
	RefNo            string `json:"ref_no" gorm:"not null;index" validate:"required"`
	PropertyTaxRef   string `json:"property_tax_ref" gorm:"not null;index" validate:"required"`
	StatementDate    string `json:"statement_date"`
	TotalAmount      string `json:"total_amount"`
	ConsolidatedData string `json:"consolidated_data" gorm:"type:text"` // JSON serialized data
	Status           string `json:"status" gorm:"default:active"`
}

// Property Tax Balance Search models based on IRAS API spec
type PropertyTaxBalanceSearchRequest struct {
	ClientID      string `json:"clientID" validate:"required"`
	Criteria      string `json:"criteria,omitempty"`
	BlkHouseNo    string `json:"blkHouseNo,omitempty"`
	StreetName    string `json:"streetName,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	StoreyNo      string `json:"storeyNo,omitempty"`
	UnitNo        string `json:"unitNo,omitempty"`
	OwnerTaxRefID string `json:"ownerTaxRefID,omitempty"`
	PptyTaxRefNo  string `json:"pptyTaxRefNo,omitempty"`
}

type PropertyTaxBalanceSearchResponse struct {
	ReturnCode int                           `json:"returnCode"`
	Data       *PropertyTaxBalanceSearchData `json:"data,omitempty"`
	Info       *PropertyTaxBalanceSearchInfo `json:"info,omitempty"`
}

type PropertyTaxBalanceSearchData struct {
	DateOfSearch           string  `json:"dateOfSearch"`
	PropertyDescription    string  `json:"propertyDescription"`
	PropertyTaxReferenceNo string  `json:"propertyTaxReferenceNo"`
	OutstandingBalance     float64 `json:"outstandingBalance"`
	PaymentByGiro          string  `json:"paymentByGiro"`
}

type PropertyTaxBalanceSearchInfo struct {
	Message       string                               `json:"message"`
	MessageCode   int                                  `json:"messageCode"`
	FieldInfoList []PropertyTaxBalanceSearchFieldError `json:"fieldInfoList,omitempty"`
}

type PropertyTaxBalanceSearchFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Property Tax Balance Search storage model for database
type PropertyTaxBalanceRecord struct {
	BaseModel
	ClientID               string  `json:"client_id" gorm:"not null;index" validate:"required"`
	PropertyTaxReferenceNo string  `json:"property_tax_reference_no" gorm:"not null;index" validate:"required"`
	PropertyDescription    string  `json:"property_description"`
	OutstandingBalance     float64 `json:"outstanding_balance"`
	PaymentByGiro          string  `json:"payment_by_giro"`
	BlkHouseNo             string  `json:"blk_house_no"`
	StreetName             string  `json:"street_name"`
	PostalCode             string  `json:"postal_code"`
	StoreyNo               string  `json:"storey_no"`
	UnitNo                 string  `json:"unit_no"`
	OwnerTaxRefID          string  `json:"owner_tax_ref_id"`
	Status                 string  `json:"status" gorm:"default:active"`
}

// Rental Submission models based on IRAS API spec
type RentalSubmissionRequest struct {
	OrgAndSubmissionInfo OrgAndSubmissionInfo       `json:"orgAndSubmissionInfo" validate:"required"`
	PropertyDtl          []PropertyDetailSubmission `json:"propertyDtl" validate:"required,min=1"`
}

type OrgAndSubmissionInfo struct {
	AssmtYear             float64 `json:"assmtYear" validate:"required"`
	AuthorisedPersonEmail string  `json:"authorisedPersonEmail" validate:"required,email"`
	AuthorisedPersonName  string  `json:"authorisedPersonName" validate:"required"`
	DevelopmentName       string  `json:"developmentName" validate:"required"`
}

type PropertyDetailSubmission struct {
	AdvPromotionAmt float64 `json:"advPromotionAmt"`
	DateGTOEnd      float64 `json:"dateGTOEnd"`
	DateGTOStart    float64 `json:"dateGTOStart"`
	DateLeaseEnd    string  `json:"dateLeaseEnd"`
	DateLeaseStart  string  `json:"dateLeaseStart"`
	GTOAmt          float64 `json:"GTOAmt"`
	GTOInfo         string  `json:"GTOInfo"`
	InfoRemarks     string  `json:"infoRemarks"`
	LetArea         string  `json:"letArea"`
	NetRentAmt      float64 `json:"netRentAmt"`
	PropertyTaxRef  string  `json:"propertyTaxRef" validate:"required"`
	RecordID        float64 `json:"recordID"`
	SvcChargeAmt    float64 `json:"svcChargeAmt"`
	TenantName      string  `json:"tenantName"`
	UnitNo          string  `json:"unitNo"`
	VacantInd       string  `json:"vacantInd"`
}

type RentalSubmissionResponse struct {
	ReturnCode int                   `json:"returnCode"`
	Data       *RentalSubmissionData `json:"data,omitempty"`
	Info       *RentalSubmissionInfo `json:"info,omitempty"`
}

type RentalSubmissionData struct {
	RefNo string `json:"refNo"`
}

type RentalSubmissionInfo struct {
	Message       string                     `json:"message"`
	MessageCode   int                        `json:"messageCode"`
	FieldInfoList *RentalSubmissionFieldInfo `json:"fieldInfoList,omitempty"`
}

type RentalSubmissionFieldInfo struct {
	FieldInfo []RentalSubmissionFieldError `json:"fieldInfo"`
}

type RentalSubmissionFieldError struct {
	Field     string `json:"field"`
	Message   string `json:"message"`
	RecordID  string `json:"recordID"`
	RecordIDs string `json:"recordIDs"`
}

// Rental Submission storage model for database
type RentalSubmissionRecord struct {
	BaseModel
	RefNo                 string  `json:"ref_no" gorm:"not null;index" validate:"required"`
	AssmtYear             float64 `json:"assmt_year" validate:"required"`
	AuthorisedPersonEmail string  `json:"authorised_person_email" validate:"required"`
	AuthorisedPersonName  string  `json:"authorised_person_name" validate:"required"`
	DevelopmentName       string  `json:"development_name" validate:"required"`
	SubmissionData        string  `json:"submission_data" gorm:"type:text"` // JSON serialized property details
	TotalProperties       int     `json:"total_properties"`
	Status                string  `json:"status" gorm:"default:submitted"`
}
