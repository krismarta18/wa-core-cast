package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// APIKey represents a user's API key for programmatic access
type APIKey struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Name        string     `json:"name"`
	KeyHash     string     `json:"-"` // never expose in JSON
	KeyPrefix   string     `json:"key_prefix"` // first 8 chars for display
	Permissions []string   `json:"permissions"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (APIKey) TableName() string {
	return "api_keys"
}

// WebhookEventSubscription defines which events a webhook listens to
type WebhookEventSubscription struct {
	ID        uuid.UUID `json:"id"`
	WebhookID uuid.UUID `json:"webhook_id"`
	EventType string    `json:"event_type"` // message.received/message.sent/device.status/etc
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName returns the table name
func (WebhookEventSubscription) TableName() string {
	return "webhook_event_subscriptions"
}

// WebhookDelivery records each webhook delivery attempt
type WebhookDelivery struct {
	ID             uuid.UUID  `json:"id"`
	WebhookID      uuid.UUID  `json:"webhook_id"`
	EventType      string     `json:"event_type"`
	Payload        string     `json:"payload"`
	StatusCode     *int       `json:"status_code,omitempty"`
	ResponseBody   *string    `json:"response_body,omitempty"`
	AttemptCount   int        `json:"attempt_count"`
	IsSuccess      bool       `json:"is_success"`
	NextRetryAt    *time.Time `json:"next_retry_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (WebhookDelivery) TableName() string {
	return "webhook_deliveries"
}

// CreateAPIKeyRequest is the request struct for creating an API key
type CreateAPIKeyRequest struct {
	Name        string     `json:"name" binding:"required"`
	Permissions []string   `json:"permissions"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// APIKeyResponse is the response struct
type APIKeyResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Key        string     `json:"key,omitempty"`    // for personal compatibility
	Prefix     string     `json:"prefix,omitempty"` // for personal compatibility
	KeyPrefix  string     `json:"key_prefix"`
	RawKey     *string    `json:"raw_key,omitempty"`
	IsActive   bool       `json:"is_active"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type WebhookSettings struct {
	UserID        uuid.UUID       `json:"user_id"`
	URL           string          `json:"url"`
	Secret        string          `json:"secret"`
	IsActive      bool            `json:"is_active"`
	EnabledEvents json.RawMessage `json:"enabled_events"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ToResponse converts APIKey to response
func (k *APIKey) ToResponse() *APIKeyResponse {
	return &APIKeyResponse{
		ID:         k.ID,
		Name:       k.Name,
		KeyPrefix:  k.KeyPrefix,
		IsActive:   k.IsActive,
		ExpiresAt:  k.ExpiresAt,
		LastUsedAt: k.LastUsedAt,
		CreatedAt:  k.CreatedAt,
	}
}
