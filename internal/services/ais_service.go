package services

import (
	"context"
	"strings"

	"api-iras/internal/models"
)

type AISService struct {
	// Add database or external API client here if needed in the future
}

func NewAISService() *AISService {
	return &AISService{}
}

// SearchOrganization performs organization search in AIS
func (s *AISService) SearchOrganization(ctx context.Context, req models.AISorgSearchRequest) (*models.AISorgSearchResponse, error) {
	// Validate client ID
	if strings.TrimSpace(req.ClientID) == "" {
		return &models.AISorgSearchResponse{
			ReturnCode: 40,
			Info: &models.AISorgInfo{
				Message:     "Invalid client ID",
				MessageCode: "40001",
				FieldInfoList: []models.AISorgFieldError{
					{
						Field:   "clientID",
						Message: "Client ID is required and cannot be empty",
					},
				},
			},
		}, nil
	}

	// Validate organization ID
	if strings.TrimSpace(req.OrganizationID) == "" {
		return &models.AISorgSearchResponse{
			ReturnCode: 40,
			Info: &models.AISorgInfo{
				Message:     "Invalid organization ID",
				MessageCode: "40002",
				FieldInfoList: []models.AISorgFieldError{
					{
						Field:   "organizationID",
						Message: "Organization ID is required and cannot be empty",
					},
				},
			},
		}, nil
	}

	// Validate basis year (should be a valid year)
	if req.BasisYear < 1900 || req.BasisYear > 2100 {
		return &models.AISorgSearchResponse{
			ReturnCode: 40,
			Info: &models.AISorgInfo{
				Message:     "Invalid basis year",
				MessageCode: "40003",
				FieldInfoList: []models.AISorgFieldError{
					{
						Field:   "basisYear",
						Message: "Basis year must be between 1900 and 2100",
					},
				},
			},
		}, nil
	}

	// For demo purposes, simulate organization search logic
	// In real implementation, this would call external IRAS API
	organizationInAIS := s.checkOrganizationInAIS(req.ClientID, req.OrganizationID, req.BasisYear)

	// Return successful response
	return &models.AISorgSearchResponse{
		ReturnCode: 10,
		Data: &models.AISorgData{
			OrganizationInAIS: organizationInAIS,
		},
	}, nil
}

// checkOrganizationInAIS simulates checking if organization exists in AIS
// In real implementation, this would integrate with actual IRAS API
func (s *AISService) checkOrganizationInAIS(clientID, organizationID string, basisYear int) string {
	// TODO: In real implementation, use clientID and basisYear for API call
	_ = clientID  // Suppress unused parameter warning
	_ = basisYear // Suppress unused parameter warning
	
	// Demo logic: Return "Y" for specific test cases, "N" for others
	testOrganizations := map[string]bool{
		"4396029847797760": true,  // Test organization exists
		"1234567890123456": true,  // Another test organization
		"9999999999999999": false, // Test organization that doesn't exist
	}

	if exists, found := testOrganizations[organizationID]; found && exists {
		return "Y"
	}

	// Default to "N" for unknown organizations
	return "N"
}
