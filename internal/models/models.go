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

// User model
type User struct {
	BaseModel
	Name     string `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Email    string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password string `json:"-" gorm:"not null" validate:"required,min=6"`
	Role     string `json:"role" gorm:"default:user" validate:"oneof=admin user"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

// Product model (contoh untuk demo)
type Product struct {
	BaseModel
	Name        string  `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Description string  `json:"description" gorm:"type:text"`
	Price       float64 `json:"price" gorm:"not null" validate:"required,min=0"`
	Stock       int     `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	CategoryID  uint    `json:"category_id" gorm:"not null"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
	UserID      uint    `json:"user_id" gorm:"not null"`
	User        User    `json:"user" gorm:"foreignKey:UserID"`
}

// Category model
type Category struct {
	BaseModel
	Name        string    `json:"name" gorm:"uniqueIndex;not null" validate:"required,min=2,max=50"`
	Description string    `json:"description" gorm:"type:text"`
	Products    []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

// Response models
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

// Request models
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"min=0"`
	CategoryID  uint    `json:"category_id" validate:"required"`
}
