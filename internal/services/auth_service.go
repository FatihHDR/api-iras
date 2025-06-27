package services

import (
	"api-iras/internal/models"
	"api-iras/pkg/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Register creates a new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if username already exists
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // Default role
		IsActive: true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates user and returns user info
func (s *AuthService) Login(req *models.LoginRequest) (*models.User, error) {
	var user models.User

	// Find user by username or email
	err := s.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves user by username
func (s *AuthService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

// UpdateUser updates user information
func (s *AuthService) UpdateUser(id uint, updates map[string]interface{}) error {
	// If password is being updated, hash it
	if password, exists := updates["password"]; exists {
		if passwordStr, ok := password.(string); ok {
			hashedPassword, err := utils.HashPassword(passwordStr)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}
			updates["password"] = hashedPassword
		}
	}

	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// DeactivateUser deactivates a user account
func (s *AuthService) DeactivateUser(id uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", false).Error
}

// GetAllUsers retrieves all users with pagination (admin only)
func (s *AuthService) GetAllUsers(page, limit int) (*models.PaginationResponse, error) {
	var users []models.User
	var total int64

	// Count total records
	s.db.Model(&models.User{}).Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Get users with pagination
	err := s.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &models.PaginationResponse{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
