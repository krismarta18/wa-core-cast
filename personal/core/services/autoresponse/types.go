package autoresponse

import (
	"time"

	"github.com/google/uuid"
)

// Keyword represents an auto-response keyword configured by the user
type Keyword struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	DeviceID     *uuid.UUID `json:"device_id" db:"device_id"` // null means applies to all devices
	Keyword      string    `json:"keyword" db:"keyword"`
	MatchType    string    `json:"match_type" db:"match_type"` // exact/contains/starts_with/ends_with/regex
	ResponseText string    `json:"response_text" db:"response_text"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateKeywordRequest is the payload to create a new keyword
type CreateKeywordRequest struct {
	DeviceID     *string `json:"device_id"`
	Keyword      string  `json:"keyword" binding:"required"`
	MatchType    string  `json:"match_type"` // Optional, default handled in logic
	ResponseText string  `json:"response_text" binding:"required"`
}

// UpdateKeywordRequest is the payload to update a keyword
type UpdateKeywordRequest struct {
	Keyword      *string `json:"keyword"`
	MatchType    *string `json:"match_type"`
	ResponseText *string `json:"response_text"`
	IsActive     *bool   `json:"is_active"`
}

// MessageTemplate represents a reusable message template
type MessageTemplate struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Category  string    `json:"category" db:"category"`
	Content   string    `json:"content" db:"content"`
	UsedCount int       `json:"used_count" db:"used_count"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTemplateRequest is the payload to create a new template
type CreateTemplateRequest struct {
	Name     string `json:"name" binding:"required"`
	Category string `json:"category" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// UpdateTemplateRequest is the payload to update a template
type UpdateTemplateRequest struct {
	Name     *string `json:"name"`
	Category *string `json:"category"`
	Content  *string `json:"content"`
}
