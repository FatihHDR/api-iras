package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateEmail validates email format using regex
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// FormatDateTime formats time to string
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// CleanString removes extra spaces and converts to lowercase
func CleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Paginate calculates pagination parameters
func Paginate(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return offset, limit
}

// Response helper functions
func SuccessResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}
}

func ErrorResponse(message string, err error) map[string]interface{} {
	response := map[string]interface{}{
		"success": false,
		"message": message,
	}
	if err != nil {
		response["error"] = err.Error()
	}
	return response
}

// Database helper functions
func BuildSearchQuery(search string) string {
	if search == "" {
		return ""
	}
	return fmt.Sprintf("%%%s%%", strings.ToLower(search))
}
