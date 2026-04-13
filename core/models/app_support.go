package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OnboardingProgress tracks a user's onboarding checklist progress
type OnboardingProgress struct {
	ID                   uuid.UUID `json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	StepAddDevice        bool      `json:"step_add_device"`
	StepSendFirstMessage bool      `json:"step_send_first_message"`
	StepSetupTemplate    bool      `json:"step_setup_template"`
	StepConfigureWebhook bool      `json:"step_configure_webhook"`
	IsCompleted          bool      `json:"is_completed"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// TableName returns the table name
func (OnboardingProgress) TableName() string {
	return "onboarding_progress"
}

// NotificationSetting holds a user's notification preferences
type NotificationSetting struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	Channel     string          `json:"channel"` // email/push/sms
	EventType   string          `json:"event_type"`
	IsEnabled   bool            `json:"is_enabled"`
	Extra       json.RawMessage `json:"extra"` // channel-specific config
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// TableName returns the table name
func (NotificationSetting) TableName() string {
	return "notification_settings"
}

// Notification represents a notification delivered to a user
type Notification struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Data      json.RawMessage `json:"data"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// TableName returns the table name
func (Notification) TableName() string {
	return "notifications"
}

// AuditLog records user and system actions for accountability
type AuditLog struct {
	ID         uuid.UUID       `json:"id"`
	UserID     *uuid.UUID      `json:"user_id,omitempty"`
	Action     string          `json:"action"`
	Resource   string          `json:"resource"`
	ResourceID *uuid.UUID      `json:"resource_id,omitempty"`
	OldData    json.RawMessage `json:"old_data"`
	NewData    json.RawMessage `json:"new_data"`
	IPAddress  string          `json:"ip_address"`
	UserAgent  *string         `json:"user_agent,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

// TableName returns the table name
func (AuditLog) TableName() string {
	return "audit_logs"
}
