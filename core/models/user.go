package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id"`
	Phone        string     `json:"phone"`
	NamaLengkap  string     `json:"nama_lengkap"`
	IsVerify     bool       `json:"is_verify"`
	OTPCode      string     `json:"otp_code"`
	OTPExpired   time.Time  `json:"otp_expired"`
	IDSubscribed uuid.UUID  `json:"id_subscribed"`
	MaxDevice    int32      `json:"max_device"`
	IsBan        bool       `json:"is_ban"`
	IsAPI        bool       `json:"is_api"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// TableName returns the table name
func (User) TableName() string {
	return "users"
}

// CreateUserRequest is the request struct for creating a user
type CreateUserRequest struct {
	Phone       string     `json:"phone" binding:"required"`
	NamaLengkap string     `json:"nama_lengkap" binding:"required"`
	IsAPI       bool       `json:"is_api"`
	IDSubscribed uuid.UUID `json:"id_subscribed"`
}

// UpdateUserRequest is the request struct for updating a user
type UpdateUserRequest struct {
	NamaLengkap *string `json:"nama_lengkap"`
	OTPCode     *string `json:"otp_code"`
	OTPExpired  *time.Time `json:"otp_expired"`
	IsVerify    *bool   `json:"is_verify"`
	IsBan       *bool   `json:"is_ban"`
}

// UserResponse is the response struct for user endpoints
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Phone        string    `json:"phone"`
	NamaLengkap  string    `json:"nama_lengkap"`
	IsVerify     bool      `json:"is_verify"`
	MaxDevice    int32     `json:"max_device"`
	IsBan        bool      `json:"is_ban"`
	IsAPI        bool      `json:"is_api"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Phone:       u.Phone,
		NamaLengkap: u.NamaLengkap,
		IsVerify:    u.IsVerify,
		MaxDevice:   u.MaxDevice,
		IsBan:       u.IsBan,
		IsAPI:       u.IsAPI,
	}
}
