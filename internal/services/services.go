package services

import (
	"api-iras/internal/models"
	"errors"
	"math"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (s *UserService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

func (s *UserService) GetUsers(page, limit int) (*models.PaginationResponse, error) {
	var users []models.User
	var total int64

	// Count total records
	s.db.Model(&models.User{}).Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Get users with pagination
	err := s.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginationResponse{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Product Service
type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.db.Create(product).Error
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := s.db.Preload("Category").Preload("User").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) GetProducts(page, limit int) (*models.PaginationResponse, error) {
	var products []models.Product
	var total int64

	// Count total records
	s.db.Model(&models.Product{}).Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Get products with pagination and preload associations
	err := s.db.Preload("Category").Preload("User").Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginationResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *ProductService) UpdateProduct(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.db.Delete(&models.Product{}, id).Error
}

// Category Service
type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	return s.db.Create(category).Error
}

func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := s.db.Preload("Products").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := s.db.Find(&categories).Error
	return categories, err
}

func (s *CategoryService) UpdateCategory(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error
}

func (s *CategoryService) DeleteCategory(id uint) error {
	// Check if category has products
	var count int64
	s.db.Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("cannot delete category with existing products")
	}
	return s.db.Delete(&models.Category{}, id).Error
}
