package services

import (
	"api-iras/internal/models"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PropertyService struct {
	db *gorm.DB
}

func NewPropertyService(db *gorm.DB) *PropertyService {
	return &PropertyService{db: db}
}

// RetrieveConsolidatedStatement retrieves property consolidated statement
func (s *PropertyService) RetrieveConsolidatedStatement(req *models.PropertyConsolidatedStatementRequest) (*models.PropertyConsolidatedStatementResponse, error) {
	// Validate reference number
	if strings.TrimSpace(req.RefNo) == "" {
		return &models.PropertyConsolidatedStatementResponse{
			ReturnCode: 40,
			Info: &models.PropertyConsolidatedStatementInfo{
				Message:     "Invalid reference number",
				MessageCode: 40001,
				FieldInfoList: []models.PropertyConsolidatedFieldError{
					{
						Field:   "refNo",
						Message: "Reference number is required",
					},
				},
			},
		}, nil
	}

	// Validate property tax reference
	if strings.TrimSpace(req.PropertyTaxRef) == "" {
		return &models.PropertyConsolidatedStatementResponse{
			ReturnCode: 40,
			Info: &models.PropertyConsolidatedStatementInfo{
				Message:     "Invalid property tax reference",
				MessageCode: 40002,
				FieldInfoList: []models.PropertyConsolidatedFieldError{
					{
						Field:   "propertyTaxRef",
						Message: "Property tax reference is required",
					},
				},
			},
		}, nil
	}

	// Search for property consolidated statement in database
	var propertyRecord models.PropertyConsolidatedStatementRecord
	err := s.db.Where("ref_no = ? AND property_tax_ref = ?", req.RefNo, req.PropertyTaxRef).First(&propertyRecord).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// For demo purposes, return a mock response if not found in database
			return s.generateMockConsolidatedStatement(req), nil
		}
		// Return server error
		return &models.PropertyConsolidatedStatementResponse{
			ReturnCode: 50,
			Info: &models.PropertyConsolidatedStatementInfo{
				Message:     "Internal server error",
				MessageCode: 50001,
			},
		}, err
	}

	// Parse consolidated data from JSON
	var consolidatedStatement models.ConsolidatedStatement
	if propertyRecord.ConsolidatedData != "" {
		if err := json.Unmarshal([]byte(propertyRecord.ConsolidatedData), &consolidatedStatement); err != nil {
			return &models.PropertyConsolidatedStatementResponse{
				ReturnCode: 50,
				Info: &models.PropertyConsolidatedStatementInfo{
					Message:     "Error parsing consolidated data",
					MessageCode: 50002,
				},
			}, err
		}
	}

	// Return successful response
	return &models.PropertyConsolidatedStatementResponse{
		ReturnCode: 0,
		Data: &models.PropertyConsolidatedStatementData{
			RefNo:                 propertyRecord.RefNo,
			PropertyTaxRef:        propertyRecord.PropertyTaxRef,
			ConsolidatedStatement: &consolidatedStatement,
		},
	}, nil
}

// generateMockConsolidatedStatement generates mock data for demo purposes
func (s *PropertyService) generateMockConsolidatedStatement(req *models.PropertyConsolidatedStatementRequest) *models.PropertyConsolidatedStatementResponse {
	// Generate mock data based on the reference numbers provided
	currentDate := time.Now().Format("2006-01-02")

	return &models.PropertyConsolidatedStatementResponse{
		ReturnCode: 0,
		Data: &models.PropertyConsolidatedStatementData{
			RefNo:          req.RefNo,
			PropertyTaxRef: req.PropertyTaxRef,
			ConsolidatedStatement: &models.ConsolidatedStatement{
				StatementDate: currentDate,
				TotalAmount:   "2,500.00",
				PropertyDetails: []models.PropertyDetail{
					{
						PropertyID:   "PROP001",
						Address:      "123 Orchard Road, Singapore 238858",
						PropertyType: "Residential",
						TaxAmount:    "1,200.00",
						DueDate:      "2025-03-31",
						Status:       "Outstanding",
					},
					{
						PropertyID:   "PROP002",
						Address:      "456 Marina Bay, Singapore 018956",
						PropertyType: "Commercial",
						TaxAmount:    "1,300.00",
						DueDate:      "2025-03-31",
						Status:       "Outstanding",
					},
				},
				PaymentHistory: []models.PaymentHistoryItem{
					{
						PaymentDate:    "2024-12-15",
						Amount:         "2,400.00",
						PaymentMethod:  "Online Banking",
						TransactionRef: "TXN202412150001",
					},
					{
						PaymentDate:    "2024-06-15",
						Amount:         "2,350.00",
						PaymentMethod:  "Credit Card",
						TransactionRef: "TXN202406150001",
					},
				},
			},
		},
	}
}

// CreateConsolidatedStatementRecord creates a new property consolidated statement record
func (s *PropertyService) CreateConsolidatedStatementRecord(record *models.PropertyConsolidatedStatementRecord) error {
	return s.db.Create(record).Error
}

// GetConsolidatedStatementRecords retrieves all property consolidated statement records with pagination
func (s *PropertyService) GetConsolidatedStatementRecords(offset, limit int) ([]models.PropertyConsolidatedStatementRecord, int64, error) {
	var records []models.PropertyConsolidatedStatementRecord
	var total int64

	// Get total count
	if err := s.db.Model(&models.PropertyConsolidatedStatementRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get records with pagination
	err := s.db.Offset(offset).Limit(limit).Find(&records).Error
	return records, total, err
}

// GetConsolidatedStatementRecordByID retrieves a property consolidated statement record by ID
func (s *PropertyService) GetConsolidatedStatementRecordByID(id uint) (*models.PropertyConsolidatedStatementRecord, error) {
	var record models.PropertyConsolidatedStatementRecord
	err := s.db.First(&record, id).Error
	return &record, err
}

// UpdateConsolidatedStatementRecord updates a property consolidated statement record
func (s *PropertyService) UpdateConsolidatedStatementRecord(id uint, record *models.PropertyConsolidatedStatementRecord) error {
	return s.db.Model(&models.PropertyConsolidatedStatementRecord{}).Where("id = ?", id).Updates(record).Error
}

// DeleteConsolidatedStatementRecord soft deletes a property consolidated statement record
func (s *PropertyService) DeleteConsolidatedStatementRecord(id uint) error {
	return s.db.Delete(&models.PropertyConsolidatedStatementRecord{}, id).Error
}

// Property Tax Balance Search methods

// SearchPropertyTaxBalance searches property tax balance based on criteria
func (s *PropertyService) SearchPropertyTaxBalance(req *models.PropertyTaxBalanceSearchRequest) (*models.PropertyTaxBalanceSearchResponse, error) {
	// Validate client ID
	if strings.TrimSpace(req.ClientID) == "" {
		return &models.PropertyTaxBalanceSearchResponse{
			ReturnCode: 40,
			Info: &models.PropertyTaxBalanceSearchInfo{
				Message:     "Invalid client ID",
				MessageCode: 40001,
				FieldInfoList: []models.PropertyTaxBalanceSearchFieldError{
					{
						Field:   "clientID",
						Message: "Client ID is required",
					},
				},
			},
		}, nil
	}

	// Build query based on provided criteria
	query := s.db.Where("client_id = ?", req.ClientID)

	if req.PptyTaxRefNo != "" {
		query = query.Where("property_tax_reference_no = ?", req.PptyTaxRefNo)
	}
	if req.PostalCode != "" {
		query = query.Where("postal_code = ?", req.PostalCode)
	}
	if req.BlkHouseNo != "" {
		query = query.Where("blk_house_no = ?", req.BlkHouseNo)
	}
	if req.StreetName != "" {
		query = query.Where("street_name ILIKE ?", "%"+req.StreetName+"%")
	}
	if req.StoreyNo != "" {
		query = query.Where("storey_no = ?", req.StoreyNo)
	}
	if req.UnitNo != "" {
		query = query.Where("unit_no = ?", req.UnitNo)
	}
	if req.OwnerTaxRefID != "" {
		query = query.Where("owner_tax_ref_id = ?", req.OwnerTaxRefID)
	}

	// Search for property tax balance in database
	var propertyRecord models.PropertyTaxBalanceRecord
	err := query.First(&propertyRecord).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// For demo purposes, return a mock response if not found in database
			return s.generateMockPropertyTaxBalance(req), nil
		}
		// Return server error
		return &models.PropertyTaxBalanceSearchResponse{
			ReturnCode: 50,
			Info: &models.PropertyTaxBalanceSearchInfo{
				Message:     "Internal server error",
				MessageCode: 50001,
			},
		}, err
	}

	// Return successful response
	return &models.PropertyTaxBalanceSearchResponse{
		ReturnCode: 0,
		Data: &models.PropertyTaxBalanceSearchData{
			DateOfSearch:           time.Now().Format("02-Jan-2006"),
			PropertyDescription:    propertyRecord.PropertyDescription,
			PropertyTaxReferenceNo: propertyRecord.PropertyTaxReferenceNo,
			OutstandingBalance:     propertyRecord.OutstandingBalance,
			PaymentByGiro:          propertyRecord.PaymentByGiro,
		},
	}, nil
}

// generateMockPropertyTaxBalance generates mock data for demo purposes
func (s *PropertyService) generateMockPropertyTaxBalance(req *models.PropertyTaxBalanceSearchRequest) *models.PropertyTaxBalanceSearchResponse {
	// Generate mock data based on the criteria provided
	propertyRef := req.PptyTaxRefNo
	if propertyRef == "" {
		propertyRef = "3004250U"
	}

	propertyDesc := "Test Property"
	if req.StreetName != "" {
		propertyDesc = "Property at " + req.StreetName
	}
	if req.PostalCode != "" {
		propertyDesc += " (" + req.PostalCode + ")"
	}

	return &models.PropertyTaxBalanceSearchResponse{
		ReturnCode: 0,
		Data: &models.PropertyTaxBalanceSearchData{
			DateOfSearch:           time.Now().Format("02-Jan-2006"),
			PropertyDescription:    propertyDesc,
			PropertyTaxReferenceNo: propertyRef,
			OutstandingBalance:     1800.00,
			PaymentByGiro:          "Yes",
		},
	}
}

// Property Tax Balance Record CRUD methods

// CreatePropertyTaxBalanceRecord creates a new property tax balance record
func (s *PropertyService) CreatePropertyTaxBalanceRecord(record *models.PropertyTaxBalanceRecord) error {
	return s.db.Create(record).Error
}

// GetPropertyTaxBalanceRecords retrieves all property tax balance records with pagination
func (s *PropertyService) GetPropertyTaxBalanceRecords(offset, limit int) ([]models.PropertyTaxBalanceRecord, int64, error) {
	var records []models.PropertyTaxBalanceRecord
	var total int64

	// Get total count
	if err := s.db.Model(&models.PropertyTaxBalanceRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get records with pagination
	err := s.db.Offset(offset).Limit(limit).Find(&records).Error
	return records, total, err
}

// GetPropertyTaxBalanceRecordByID retrieves a property tax balance record by ID
func (s *PropertyService) GetPropertyTaxBalanceRecordByID(id uint) (*models.PropertyTaxBalanceRecord, error) {
	var record models.PropertyTaxBalanceRecord
	err := s.db.First(&record, id).Error
	return &record, err
}

// UpdatePropertyTaxBalanceRecord updates a property tax balance record
func (s *PropertyService) UpdatePropertyTaxBalanceRecord(id uint, record *models.PropertyTaxBalanceRecord) error {
	return s.db.Model(&models.PropertyTaxBalanceRecord{}).Where("id = ?", id).Updates(record).Error
}

// DeletePropertyTaxBalanceRecord soft deletes a property tax balance record
func (s *PropertyService) DeletePropertyTaxBalanceRecord(id uint) error {
	return s.db.Delete(&models.PropertyTaxBalanceRecord{}, id).Error
}
