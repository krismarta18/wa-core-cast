package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID  `json:"id"`
	PhoneNumber   string     `json:"phone_number"`
	FullName      string     `json:"full_name"`
	Email         *string    `json:"email,omitempty"`
	CompanyName   *string    `json:"company_name,omitempty"`
	Timezone      string     `json:"timezone"`
	IsVerified    bool       `json:"is_verified"`
	IsBanned      bool       `json:"is_banned"`
	IsAPIEnabled  bool       `json:"is_api_enabled"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

// TableName returns the table name
func (User) TableName() string {
	return "users"
}

// CreateUserRequest is the request struct for creating a user
type CreateUserRequest struct {
	PhoneNumber string  `json:"phone_number" binding:"required"`
	FullName    string  `json:"full_name" binding:"required"`
	Email       *string `json:"email"`
	CompanyName *string `json:"company_name"`
	Timezone    string  `json:"timezone"`
}

// UpdateUserRequest is the request struct for updating a user
type UpdateUserRequest struct {
	FullName     *string `json:"full_name"`
	Email        *string `json:"email"`
	CompanyName  *string `json:"company_name"`
	Timezone     *string `json:"timezone"`
	IsVerified   *bool   `json:"is_verified"`
	IsBanned     *bool   `json:"is_banned"`
	IsAPIEnabled *bool   `json:"is_api_enabled"`
}

// UserResponse is the response struct for user endpoints
type UserResponse struct {
	ID           uuid.UUID  `json:"id"`
	PhoneNumber  string     `json:"phone_number"`
	FullName     string     `json:"full_name"`
	Email        *string    `json:"email,omitempty"`
	CompanyName  *string    `json:"company_name,omitempty"`
	Timezone     string     `json:"timezone"`
	IsVerified   bool       `json:"is_verified"`
	IsBanned     bool       `json:"is_banned"`
	IsAPIEnabled bool       `json:"is_api_enabled"`
	CreatedAt    time.Time  `json:"created_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		PhoneNumber:  u.PhoneNumber,
		FullName:     u.FullName,
		Email:        u.Email,
		CompanyName:  u.CompanyName,
		Timezone:     u.Timezone,
		IsVerified:   u.IsVerified,
		IsBanned:     u.IsBanned,
		IsAPIEnabled: u.IsAPIEnabled,
		CreatedAt:    u.CreatedAt,
		LastLoginAt:  u.LastLoginAt,
	}
}
