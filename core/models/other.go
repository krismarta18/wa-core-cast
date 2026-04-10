package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// APILog represents an API request log
type APILog struct {
	ID           uuid.UUID       `json:"id"`
	UserID       *uuid.UUID      `json:"user_id"`
	Endpoint     string          `json:"endpoint"`
	ReqBody      json.RawMessage `json:"req_body"`
	ResponseBody json.RawMessage `json:"response_body"`
	CreatedAt    string          `json:"created_at"`
	IPAddress    string          `json:"ip_address"`
	DeviceID     *uuid.UUID      `json:"device_id"`
}

// TableName returns the table name
func (APILog) TableName() string {
	return "api_logs"
}

// AutoResponse represents an auto response rule
type AutoResponse struct {
	ID           uuid.UUID `json:"id"`
	DeviceID     uuid.UUID `json:"device_id"`
	Keyword      string    `json:"keyword"`
	ResponseText string    `json:"response_text"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}

// TableName returns the table name
func (AutoResponse) TableName() string {
	return "auto_response"
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID        uuid.UUID `json:"id"`
	DeviceID  uuid.UUID `json:"device_id"`
	WebhookUrl string   `json:"webhook_url"`
	SecretKey string   `json:"secret_key"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// TableName returns the table name
func (Webhook) TableName() string {
	return "webhooks"
}

// Lookup represents a lookup value
type Lookup struct {
	ID     int32  `json:"id"`
	Keys   string `json:"keys"`
	Values string `json:"values"`
}

// TableName returns the table name
func (Lookup) TableName() string {
	return "lookup"
}

// SystemSetting represents a system setting
type SystemSetting struct {
	ID          int32      `json:"id"`
	Keys        string     `json:"keys"`
	Value       string     `json:"value"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

// TableName returns the table name
func (SystemSetting) TableName() string {
	return "system_settings"
}

// CreateAutoResponseRequest is the request struct
type CreateAutoResponseRequest struct {
	DeviceID     uuid.UUID `json:"device_id" binding:"required"`
	Keyword      string    `json:"keyword" binding:"required"`
	ResponseText string    `json:"response_text" binding:"required"`
}

// UpdateAutoResponseRequest is the request struct
type UpdateAutoResponseRequest struct {
	Keyword      *string `json:"keyword"`
	ResponseText *string `json:"response_text"`
	IsActive     *bool   `json:"is_active"`
}

// CreateWebhookRequest is the request struct
type CreateWebhookRequest struct {
	DeviceID   uuid.UUID `json:"device_id" binding:"required"`
	WebhookUrl string    `json:"webhook_url" binding:"required"`
	SecretKey  string    `json:"secret_key" binding:"required"`
}

// AutoResponseResponse is the response struct
type AutoResponseResponse struct {
	ID           uuid.UUID `json:"id"`
	Keyword      string    `json:"keyword"`
	ResponseText string    `json:"response_text"`
	IsActive     bool      `json:"is_active"`
}

// ToResponse converts AutoResponse to response
func (a *AutoResponse) ToResponse() *AutoResponseResponse {
	return &AutoResponseResponse{
		ID:           a.ID,
		Keyword:      a.Keyword,
		ResponseText: a.ResponseText,
		IsActive:     a.IsActive,
	}
}

// WebhookResponse is the response struct
type WebhookResponse struct {
	ID        uuid.UUID `json:"id"`
	DeviceID  uuid.UUID `json:"device_id"`
	WebhookUrl string   `json:"webhook_url"`
}

// ToResponse converts Webhook to response
func (w *Webhook) ToResponse() *WebhookResponse {
	return &WebhookResponse{
		ID:        w.ID,
		DeviceID:  w.DeviceID,
		WebhookUrl: w.WebhookUrl,
	}
}
