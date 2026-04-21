package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageTemplate represents a reusable message template
type MessageTemplate struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	Name        string          `json:"name"`
	Category    string          `json:"category"` // marketing/utility/authentication
	Content     string          `json:"content"`
	Variables   json.RawMessage `json:"variables"` // array of variable names
	MediaURL    *string         `json:"media_url,omitempty"`
	Language    string          `json:"language"`
	IsApproved  bool            `json:"is_approved"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *time.Time      `json:"deleted_at,omitempty"`
}

// TableName returns the table name
func (MessageTemplate) TableName() string {
	return "message_templates"
}

// ScheduledMessage represents a message scheduled to be sent at a future time
type ScheduledMessage struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	DeviceID       uuid.UUID  `json:"device_id"`
	TemplateID     *uuid.UUID `json:"template_id,omitempty"`
	Name           string     `json:"name"`
	MessageContent string     `json:"message_content"`
	ScheduledAt    time.Time  `json:"scheduled_at"`
	Status         string     `json:"status"` // pending/processing/completed/failed/cancelled
	SentCount      int        `json:"sent_count"`
	FailedCount    int        `json:"failed_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (ScheduledMessage) TableName() string {
	return "scheduled_messages"
}

// ScheduledMessageRecipient represents a recipient in a scheduled message
type ScheduledMessageRecipient struct {
	ID                 uuid.UUID  `json:"id"`
	ScheduledMessageID uuid.UUID  `json:"scheduled_message_id"`
	PhoneNumber        string     `json:"phone_number"`
	ContactID          *uuid.UUID `json:"contact_id,omitempty"`
	Status             string     `json:"status"` // pending/sent/failed
	SentAt             *time.Time `json:"sent_at,omitempty"`
	ErrorMessage       *string    `json:"error_message,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

// TableName returns the table name
func (ScheduledMessageRecipient) TableName() string {
	return "scheduled_message_recipients"
}

// AutoResponseLog logs each triggered auto-response
type AutoResponseLog struct {
	ID        uuid.UUID `json:"id"`
	KeywordID uuid.UUID `json:"keyword_id"`
	DeviceID  uuid.UUID `json:"device_id"`
	Sender    string    `json:"sender"`
	Trigger   string    `json:"trigger"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName returns the table name
func (AutoResponseLog) TableName() string {
	return "auto_response_logs"
}

// CreateMessageTemplateRequest is the request struct for creating a template
type CreateMessageTemplateRequest struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required,oneof=marketing utility authentication"`
	Content  string  `json:"content" binding:"required"`
	Language string  `json:"language"`
	MediaURL *string `json:"media_url"`
}

// MessageTemplateResponse is the response struct
type MessageTemplateResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	Content    string    `json:"content"`
	Language   string    `json:"language"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToResponse converts MessageTemplate to response
func (t *MessageTemplate) ToResponse() *MessageTemplateResponse {
	return &MessageTemplateResponse{
		ID:         t.ID,
		Name:       t.Name,
		Category:   t.Category,
		Content:    t.Content,
		Language:   t.Language,
		IsApproved: t.IsApproved,
		CreatedAt:  t.CreatedAt,
	}
}
