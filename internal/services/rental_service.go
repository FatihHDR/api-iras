package services

import (
	"api-iras/internal/models"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type RentalService struct {
	db *gorm.DB
}

func NewRentalService(db *gorm.DB) *RentalService {
	return &RentalService{db: db}
}

// SubmitRental submits rental information
func (s *RentalService) SubmitRental(req *models.RentalSubmissionRequest) (*models.RentalSubmissionResponse, error) {
	// Validate organization and submission info
	if req.OrgAndSubmissionInfo.AssmtYear <= 0 {
		return &models.RentalSubmissionResponse{
			ReturnCode: 40,
			Info: &models.RentalSubmissionInfo{
				Message:     "Invalid assessment year",
				MessageCode: 40001,
				FieldInfoList: &models.RentalSubmissionFieldInfo{
					FieldInfo: []models.RentalSubmissionFieldError{
						{
							Field:   "assmtYear",
							Message: "Assessment year must be greater than 0",
						},
					},
				},
			},
		}, nil
	}

	if strings.TrimSpace(req.OrgAndSubmissionInfo.AuthorisedPersonEmail) == "" {
		return &models.RentalSubmissionResponse{
			ReturnCode: 40,
			Info: &models.RentalSubmissionInfo{
				Message:     "Invalid authorised person email",
				MessageCode: 40002,
				FieldInfoList: &models.RentalSubmissionFieldInfo{
					FieldInfo: []models.RentalSubmissionFieldError{
						{
							Field:   "authorisedPersonEmail",
							Message: "Authorised person email is required",
						},
					},
				},
			},
		}, nil
	}

	if strings.TrimSpace(req.OrgAndSubmissionInfo.AuthorisedPersonName) == "" {
		return &models.RentalSubmissionResponse{
			ReturnCode: 40,
			Info: &models.RentalSubmissionInfo{
				Message:     "Invalid authorised person name",
				MessageCode: 40003,
				FieldInfoList: &models.RentalSubmissionFieldInfo{
					FieldInfo: []models.RentalSubmissionFieldError{
						{
							Field:   "authorisedPersonName",
							Message: "Authorised person name is required",
						},
					},
				},
			},
		}, nil
	}

	if strings.TrimSpace(req.OrgAndSubmissionInfo.DevelopmentName) == "" {
		return &models.RentalSubmissionResponse{
			ReturnCode: 40,
			Info: &models.RentalSubmissionInfo{
				Message:     "Invalid development name",
				MessageCode: 40004,
				FieldInfoList: &models.RentalSubmissionFieldInfo{
					FieldInfo: []models.RentalSubmissionFieldError{
						{
							Field:   "developmentName",
							Message: "Development name is required",
						},
					},
				},
			},
		}, nil
	}

	// Validate property details
	if len(req.PropertyDtl) == 0 {
		return &models.RentalSubmissionResponse{
			ReturnCode: 40,
			Info: &models.RentalSubmissionInfo{
				Message:     "No property details provided",
				MessageCode: 40005,
				FieldInfoList: &models.RentalSubmissionFieldInfo{
					FieldInfo: []models.RentalSubmissionFieldError{
						{
							Field:   "propertyDtl",
							Message: "At least one property detail is required",
						},
					},
				},
			},
		}, nil
	}

	// Validate each property detail
	var fieldErrors []models.RentalSubmissionFieldError
	for _, property := range req.PropertyDtl {
		if strings.TrimSpace(property.PropertyTaxRef) == "" {
			fieldErrors = append(fieldErrors, models.RentalSubmissionFieldError{
				Field:    "propertyTaxRef",
				Message:  "Property tax reference is required",
				RecordID: fmt.Sprintf("%.0f", property.RecordID),
			})
		}
	}

	if len(fieldErrors) > 0 {
		return &models.RentalSubmissionResponse{
			ReturnCode: 40,
			Info: &models.RentalSubmissionInfo{
				Message:     "Validation errors in property details",
				MessageCode: 40006,
				FieldInfoList: &models.RentalSubmissionFieldInfo{
					FieldInfo: fieldErrors,
				},
			},
		}, nil
	}

	// Generate reference number
	refNo := s.generateRefNo()

	// Serialize property details to JSON
	propertyDataJSON, err := json.Marshal(req.PropertyDtl)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize property data: %w", err)
	}

	// Store submission record in database
	record := &models.RentalSubmissionRecord{
		RefNo:                 refNo,
		AssmtYear:             req.OrgAndSubmissionInfo.AssmtYear,
		AuthorisedPersonEmail: req.OrgAndSubmissionInfo.AuthorisedPersonEmail,
		AuthorisedPersonName:  req.OrgAndSubmissionInfo.AuthorisedPersonName,
		DevelopmentName:       req.OrgAndSubmissionInfo.DevelopmentName,
		SubmissionData:        string(propertyDataJSON),
		TotalProperties:       len(req.PropertyDtl),
		Status:                "submitted",
	}

	if err := s.db.Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to store rental submission: %w", err)
	}

	// Return successful response
	return &models.RentalSubmissionResponse{
		ReturnCode: 0,
		Data: &models.RentalSubmissionData{
			RefNo: refNo,
		},
	}, nil
}

// generateRefNo generates a unique reference number for the submission
func (s *RentalService) generateRefNo() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("RNT%s", timestamp)
}

// Admin CRUD operations for rental submission records

// CreateRentalSubmissionRecord creates a new rental submission record
func (s *RentalService) CreateRentalSubmissionRecord(record *models.RentalSubmissionRecord) error {
	return s.db.Create(record).Error
}

// GetRentalSubmissionRecords retrieves rental submission records with pagination
func (s *RentalService) GetRentalSubmissionRecords(offset, limit int) ([]models.RentalSubmissionRecord, int64, error) {
	var records []models.RentalSubmissionRecord
	var total int64

	// Get total count
	if err := s.db.Model(&models.RentalSubmissionRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get records with pagination
	if err := s.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetRentalSubmissionRecordByID retrieves a rental submission record by ID
func (s *RentalService) GetRentalSubmissionRecordByID(id uint) (*models.RentalSubmissionRecord, error) {
	var record models.RentalSubmissionRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GetRentalSubmissionRecordByRefNo retrieves a rental submission record by reference number
func (s *RentalService) GetRentalSubmissionRecordByRefNo(refNo string) (*models.RentalSubmissionRecord, error) {
	var record models.RentalSubmissionRecord
	if err := s.db.Where("ref_no = ?", refNo).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// UpdateRentalSubmissionRecord updates a rental submission record
func (s *RentalService) UpdateRentalSubmissionRecord(id uint, record *models.RentalSubmissionRecord) error {
	return s.db.Model(&models.RentalSubmissionRecord{}).Where("id = ?", id).Updates(record).Error
}

// DeleteRentalSubmissionRecord deletes a rental submission record
func (s *RentalService) DeleteRentalSubmissionRecord(id uint) error {
	return s.db.Delete(&models.RentalSubmissionRecord{}, id).Error
}
