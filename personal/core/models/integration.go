package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// APIKey represents an external API key for access
type APIKey struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Name       string     `json:"name"`
	KeyHash    string     `json:"-"`
	Prefix     string     `json:"prefix"` // first few chars to show user
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"-"`
}

// TableName returns the table name
func (APIKey) TableName() string {
	return "api_keys"
}

// WebhookSettings represents webhook configuration for a user
type WebhookSettings struct {
	UserID        uuid.UUID       `json:"user_id"`
	URL           string          `json:"url"`
	Secret        string          `json:"secret"`
	IsActive      bool            `json:"is_active"`
	EnabledEvents json.RawMessage `json:"enabled_events"` // array of strings
	UpdatedAt     time.Time       `json:"updated_at"`
}

// TableName returns the table name
func (WebhookSettings) TableName() string {
	return "webhook_settings"
}

// APIKeyResponse is the response struct for API Keys
type APIKeyResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Key        string     `json:"key,omitempty"` // Only populated on creation
	Prefix     string     `json:"prefix"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
