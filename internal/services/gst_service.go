package services

import (
	"api-iras/internal/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type GSTService struct {
	db *gorm.DB
}

func NewGSTService(db *gorm.DB) *GSTService {
	return &GSTService{db: db}
}

// SearchGSTRegistered performs GST registration lookup
func (s *GSTService) SearchGSTRegistered(req *models.GSTRequest) (*models.GSTResponse, error) {
	// Validate client ID (basic validation)
	if strings.TrimSpace(req.ClientID) == "" {
		return &models.GSTResponse{
			ReturnCode: 40,
			Info: &models.GSTInfo{
				Message:     "Invalid client ID",
				MessageCode: 40001,
				FieldInfoList: []models.GSTFieldError{
					{
						Field:   "clientID",
						Message: "Client ID is required",
					},
				},
			},
		}, nil
	}

	// Validate registration ID
	if strings.TrimSpace(req.RegID) == "" {
		return &models.GSTResponse{
			ReturnCode: 40,
			Info: &models.GSTInfo{
				Message:     "Invalid registration ID",
				MessageCode: 40002,
				FieldInfoList: []models.GSTFieldError{
					{
						Field:   "regID",
						Message: "Registration ID is required",
					},
				},
			},
		}, nil
	}

	// Search for GST registration in database
	var gstReg models.GSTRegistration
	err := s.db.Where("registration_id = ? AND client_id = ?", req.RegID, req.ClientID).First(&gstReg).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return not found response
			return &models.GSTResponse{
				ReturnCode: 20,
				Info: &models.GSTInfo{
					Message:     "GST registration not found",
					MessageCode: 20001,
				},
			}, nil
		}
		// Return server error
		return &models.GSTResponse{
			ReturnCode: 50,
			Info: &models.GSTInfo{
				Message:     "Internal server error",
				MessageCode: 50001,
			},
		}, err
	}

	// Return successful response
	return &models.GSTResponse{
		ReturnCode: 10,
		Data: &models.GSTData{
			Name:                  gstReg.Name,
			GSTRegistrationNumber: gstReg.GSTRegistrationNumber,
			RegistrationID:        gstReg.RegistrationID,
			RegisteredFrom:        gstReg.RegisteredFrom,
			RegisteredTo:          gstReg.RegisteredTo,
			Remarks:               gstReg.Remarks,
			Status:                gstReg.Status,
		},
	}, nil
}

// CreateGSTRegistration creates a new GST registration record (for admin/setup purposes)
func (s *GSTService) CreateGSTRegistration(gstReg *models.GSTRegistration) error {
	return s.db.Create(gstReg).Error
}

// GetGSTRegistrationByID gets GST registration by ID
func (s *GSTService) GetGSTRegistrationByID(id uint) (*models.GSTRegistration, error) {
	var gstReg models.GSTRegistration
	err := s.db.First(&gstReg, id).Error
	if err != nil {
		return nil, err
	}
	return &gstReg, nil
}

// UpdateGSTRegistration updates GST registration
func (s *GSTService) UpdateGSTRegistration(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.GSTRegistration{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteGSTRegistration deletes GST registration
func (s *GSTService) DeleteGSTRegistration(id uint) error {
	return s.db.Delete(&models.GSTRegistration{}, id).Error
}

// GetAllGSTRegistrations gets all GST registrations for admin purposes
func (s *GSTService) GetAllGSTRegistrations(page, limit int) ([]models.GSTRegistration, int64, error) {
	var gstRegs []models.GSTRegistration
	var total int64

	// Count total records
	s.db.Model(&models.GSTRegistration{}).Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Get records with pagination
	err := s.db.Offset(offset).Limit(limit).Find(&gstRegs).Error
	if err != nil {
		return nil, 0, err
	}

	return gstRegs, total, nil
}
