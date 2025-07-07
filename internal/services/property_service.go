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
