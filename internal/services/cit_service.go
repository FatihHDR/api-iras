package services

import (
	"api-iras/internal/models"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type CITService struct {
	db *gorm.DB
}

func NewCITService(db *gorm.DB) *CITService {
	return &CITService{db: db}
}

// ConvertFormCS performs CIT conversion based on the provided ID
func (s *CITService) ConvertFormCS(req *models.CITConversionRequest) (*models.CITConversionResponse, error) {
	// Validate request ID
	if req.ID <= 0 {
		return &models.CITConversionResponse{
			ReturnCode: 40,
			Info: &models.CITConversionInfo{
				Message:     "Invalid ID",
				MessageCode: 40001,
				FieldInfoList: []models.CITConversionFieldError{
					{
						Field:   "id",
						Message: "ID must be greater than 0",
					},
				},
			},
		}, nil
	}

	// Check if conversion already exists for this ID
	var existingRecord models.CITConversionRecord
	err := s.db.Where("request_id = ?", req.ID).First(&existingRecord).Error
	if err == nil {
		// Return existing conversion
		return &models.CITConversionResponse{
			ReturnCode: 10,
			Data: &models.CITConversionData{
				ConversionID:     existingRecord.ConversionID,
				Status:           existingRecord.Status,
				ConversionDate:   existingRecord.ConversionDate,
				ProcessedBy:      existingRecord.ProcessedBy,
				ConversionResult: existingRecord.ConversionResult,
			},
		}, nil
	}

	// Generate conversion ID
	conversionID := s.generateConversionID()

	// Create new conversion record
	conversionRecord := &models.CITConversionRecord{
		ConversionID:     conversionID,
		RequestID:        req.ID,
		Status:           "completed",
		ConversionDate:   time.Now().Format("2006-01-02 15:04:05"),
		ProcessedBy:      "IRAS_CIT_SYSTEM",
		ConversionResult: fmt.Sprintf("Form CS conversion completed for ID: %d", req.ID),
		ClientID:         strconv.FormatInt(req.ID, 10),
	}

	// Save to database
	if err := s.db.Create(conversionRecord).Error; err != nil {
		return &models.CITConversionResponse{
			ReturnCode: 50,
			Info: &models.CITConversionInfo{
				Message:     "Internal server error",
				MessageCode: 50001,
			},
		}, fmt.Errorf("failed to save conversion record: %w", err)
	}

	// Return successful response
	return &models.CITConversionResponse{
		ReturnCode: 10,
		Data: &models.CITConversionData{
			ConversionID:     conversionRecord.ConversionID,
			Status:           conversionRecord.Status,
			ConversionDate:   conversionRecord.ConversionDate,
			ProcessedBy:      conversionRecord.ProcessedBy,
			ConversionResult: conversionRecord.ConversionResult,
		},
	}, nil
}

// generateConversionID generates a unique conversion ID
func (s *CITService) generateConversionID() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("CIT%s", timestamp)
}

// Admin CRUD operations for CIT conversion records

// CreateCITConversionRecord creates a new CIT conversion record
func (s *CITService) CreateCITConversionRecord(record *models.CITConversionRecord) error {
	return s.db.Create(record).Error
}

// GetCITConversionRecords retrieves CIT conversion records with pagination
func (s *CITService) GetCITConversionRecords(offset, limit int) ([]models.CITConversionRecord, int64, error) {
	var records []models.CITConversionRecord
	var total int64

	// Get total count
	if err := s.db.Model(&models.CITConversionRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get records with pagination
	if err := s.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetCITConversionRecordByID retrieves a CIT conversion record by ID
func (s *CITService) GetCITConversionRecordByID(id uint) (*models.CITConversionRecord, error) {
	var record models.CITConversionRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GetCITConversionRecordByConversionID retrieves a CIT conversion record by conversion ID
func (s *CITService) GetCITConversionRecordByConversionID(conversionID string) (*models.CITConversionRecord, error) {
	var record models.CITConversionRecord
	if err := s.db.Where("conversion_id = ?", conversionID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// GetCITConversionRecordByRequestID retrieves a CIT conversion record by request ID
func (s *CITService) GetCITConversionRecordByRequestID(requestID int64) (*models.CITConversionRecord, error) {
	var record models.CITConversionRecord
	if err := s.db.Where("request_id = ?", requestID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// UpdateCITConversionRecord updates a CIT conversion record
func (s *CITService) UpdateCITConversionRecord(id uint, record *models.CITConversionRecord) error {
	return s.db.Model(&models.CITConversionRecord{}).Where("id = ?", id).Updates(record).Error
}

// DeleteCITConversionRecord deletes a CIT conversion record
func (s *CITService) DeleteCITConversionRecord(id uint) error {
	return s.db.Delete(&models.CITConversionRecord{}, id).Error
}
