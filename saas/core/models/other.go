package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// APILog represents an API request log
type APILog struct {
	ID           uuid.UUID       `json:"id"`
	UserID       *uuid.UUID      `json:"user_id,omitempty"`
	DeviceID     *uuid.UUID      `json:"device_id,omitempty"`
	Endpoint     string          `json:"endpoint"`
	Method       string          `json:"method"`
	StatusCode   int             `json:"status_code"`
	ReqBody      json.RawMessage `json:"req_body"`
	ResponseBody json.RawMessage `json:"response_body"`
	IPAddress    string          `json:"ip_address"`
	UserAgent    *string         `json:"user_agent,omitempty"`
	DurationMs   *int            `json:"duration_ms,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

// TableName returns the table name
func (APILog) TableName() string {
	return "api_logs"
}

// AutoResponseKeyword represents an auto-response keyword rule
type AutoResponseKeyword struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	DeviceID     *uuid.UUID `json:"device_id,omitempty"`
	Keyword      string     `json:"keyword"`
	MatchType    string     `json:"match_type"` // exact/contains/starts_with/ends_with/regex
	ResponseText string     `json:"response_text"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (AutoResponseKeyword) TableName() string {
	return "auto_response_keywords"
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	DeviceID       *uuid.UUID `json:"device_id,omitempty"`
	WebhookUrl     string     `json:"webhook_url"`
	SecretKeyHash  string     `json:"-"` // never expose in JSON
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (Webhook) TableName() string {
	return "webhooks"
}

// Lookup represents a lookup value
type Lookup struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TableName returns the table name
func (Lookup) TableName() string {
	return "lookup"
}

// SystemSetting represents a system setting
type SystemSetting struct {
	ID          int       `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name
func (SystemSetting) TableName() string {
	return "system_settings"
}

// CreateAutoResponseKeywordRequest is the request struct
type CreateAutoResponseKeywordRequest struct {
	DeviceID     *uuid.UUID `json:"device_id"`
	Keyword      string     `json:"keyword" binding:"required"`
	MatchType    string     `json:"match_type" binding:"required,oneof=exact contains starts_with ends_with regex"`
	ResponseText string     `json:"response_text" binding:"required"`
}

// UpdateAutoResponseKeywordRequest is the request struct
type UpdateAutoResponseKeywordRequest struct {
	Keyword      *string `json:"keyword"`
	MatchType    *string `json:"match_type"`
	ResponseText *string `json:"response_text"`
	IsActive     *bool   `json:"is_active"`
}

// CreateWebhookRequest is the request struct
type CreateWebhookRequest struct {
	DeviceID   *uuid.UUID `json:"device_id"`
	WebhookUrl string     `json:"webhook_url" binding:"required,url"`
	SecretKey  string     `json:"secret_key" binding:"required"`
}

// AutoResponseKeywordResponse is the response struct
type AutoResponseKeywordResponse struct {
	ID           uuid.UUID  `json:"id"`
	DeviceID     *uuid.UUID `json:"device_id,omitempty"`
	Keyword      string     `json:"keyword"`
	MatchType    string     `json:"match_type"`
	ResponseText string     `json:"response_text"`
	IsActive     bool       `json:"is_active"`
}

// ToResponse converts AutoResponseKeyword to response
func (a *AutoResponseKeyword) ToResponse() *AutoResponseKeywordResponse {
	return &AutoResponseKeywordResponse{
		ID:           a.ID,
		DeviceID:     a.DeviceID,
		Keyword:      a.Keyword,
		MatchType:    a.MatchType,
		ResponseText: a.ResponseText,
		IsActive:     a.IsActive,
	}
}

// WebhookResponse is the response struct
type WebhookResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	DeviceID   *uuid.UUID `json:"device_id,omitempty"`
	WebhookUrl string     `json:"webhook_url"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ToResponse converts Webhook to response
func (w *Webhook) ToResponse() *WebhookResponse {
	return &WebhookResponse{
		ID:         w.ID,
		UserID:     w.UserID,
		DeviceID:   w.DeviceID,
		WebhookUrl: w.WebhookUrl,
		IsActive:   w.IsActive,
		CreatedAt:  w.CreatedAt,
	}
}
