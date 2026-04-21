package models

import (
	"time"

	"github.com/google/uuid"
)

// Broadcast campaign status constants
const (
	BroadcastStatusDraft     = "draft"
	BroadcastStatusQueued    = "queued"
	BroadcastStatusSending   = "sending"
	BroadcastStatusCompleted = "completed"
	BroadcastStatusFailed    = "failed"
	BroadcastStatusCancelled = "cancelled"
)

// BroadcastCampaign represents a broadcast campaign
type BroadcastCampaign struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	DeviceID        uuid.UUID  `json:"device_id"`
	TemplateID      *uuid.UUID `json:"template_id,omitempty"`
	Name            string     `json:"name"`
	MessageContent  *string    `json:"message_content,omitempty"`
	TotalRecipients int        `json:"total_recipients"`
	SuccessCount    int        `json:"success_count"`
	FailedCount     int        `json:"failed_count"`
	DelaySeconds    int        `json:"delay_seconds"`
	ScheduledAt     *time.Time `json:"scheduled_at,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	Status          string     `json:"status"` // draft/queued/sending/completed/failed/cancelled
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// TableName returns the table name
func (BroadcastCampaign) TableName() string {
	return "broadcast_campaigns"
}

// BroadcastRecipient represents a recipient in a broadcast
type BroadcastRecipient struct {
	ID          uuid.UUID  `json:"id"`
	CampaignID  uuid.UUID  `json:"campaign_id"`
	GroupID     *uuid.UUID `json:"group_id,omitempty"`
	ContactID   *uuid.UUID `json:"contact_id,omitempty"`
	PhoneNumber string     `json:"phone_number"`
	Status      string     `json:"status"` // pending/sent/failed
	SentAt      *time.Time `json:"sent_at,omitempty"`
	FailedAt    *time.Time `json:"failed_at,omitempty"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	RetryCount  int        `json:"retry_count"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TableName returns the table name
func (BroadcastRecipient) TableName() string {
	return "broadcast_recipients"
}

// CreateBroadcastCampaignRequest is the request struct for creating a campaign
type CreateBroadcastCampaignRequest struct {
	DeviceID       uuid.UUID  `json:"device_id" binding:"required"`
	Name           string     `json:"name" binding:"required"`
	TemplateID     *uuid.UUID `json:"template_id"`
	MessageContent *string    `json:"message_content"`
	DelaySeconds   int        `json:"delay_seconds"`
	ScheduledAt    *time.Time `json:"scheduled_at"`
}

// BroadcastCampaignResponse is the response struct
type BroadcastCampaignResponse struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	TotalRecipients int        `json:"total_recipients"`
	SuccessCount    int        `json:"success_count"`
	FailedCount     int        `json:"failed_count"`
	Status          string     `json:"status"`
	ScheduledAt     *time.Time `json:"scheduled_at,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ToResponse converts BroadcastCampaign to response
func (b *BroadcastCampaign) ToResponse() *BroadcastCampaignResponse {
	return &BroadcastCampaignResponse{
		ID:              b.ID,
		Name:            b.Name,
		TotalRecipients: b.TotalRecipients,
		SuccessCount:    b.SuccessCount,
		FailedCount:     b.FailedCount,
		Status:          b.Status,
		ScheduledAt:     b.ScheduledAt,
		StartedAt:       b.StartedAt,
		CompletedAt:     b.CompletedAt,
		CreatedAt:       b.CreatedAt,
	}
}
