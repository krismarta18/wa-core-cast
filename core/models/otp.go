package models

import (
	"time"

	"github.com/google/uuid"
)

// OTP context constants
const (
	OTPContextLogin    = "login"
	OTPContextRegister = "register"
	OTPContextReset    = "reset"
)

// OTPVerification represents a pending OTP verification record
type OTPVerification struct {
	ID           uuid.UUID  `json:"id"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	PhoneNumber  string     `json:"phone_number"`
	Context      string     `json:"context"` // login/register/reset
	OTPCode      string     `json:"-"`
	AttemptCount int        `json:"attempt_count"`
	ExpiresAt    time.Time  `json:"expires_at"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// TableName returns the table name
func (OTPVerification) TableName() string {
	return "otp_verifications"
}

// IsExpired checks if this OTP has expired
func (o *OTPVerification) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsUsed checks if this OTP has already been verified
func (o *OTPVerification) IsUsed() bool {
	return o.VerifiedAt != nil
}

// UserSession represents an active user session
type UserSession struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	SessionTokenHash string     `json:"-"` // never expose
	RefreshTokenHash *string    `json:"-"` // never expose
	IPAddress        *string    `json:"ip_address,omitempty"`
	UserAgent        *string    `json:"user_agent,omitempty"`
	LastActiveAt     time.Time  `json:"last_active_at"`
	ExpiresAt        time.Time  `json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// TableName returns the table name
func (UserSession) TableName() string {
	return "user_sessions"
}

// IsExpired checks if this session has expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
