package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BroadcastCampaign represents a broadcast campaign
type BroadcastCampaign struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	DeviceID        uuid.UUID  `json:"device_id"`
	NameBroadcast   string     `json:"name_broadcast"`
	TotalRecipients int32      `json:"total_recipients"`
	ProcessedCount  int32      `json:"processed_count"`
	ScheduledAt     *time.Time `json:"scheduled_at"`
	Status          int32      `json:"status"` // 0: draft, 1: scheduled, 2: running, 3: completed
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// TableName returns the table name
func (BroadcastCampaign) TableName() string {
	return "broadcast_campaigns"
}

// BroadcastMessage represents a message in a campaign
type BroadcastMessage struct {
	ID          uuid.UUID       `json:"id"`
	CampaignID  uuid.UUID       `json:"campaign_id"`
	MessageType int32           `json:"message_type"`
	MessageText string          `json:"message_text"`
	MediaUrl    string          `json:"media_url"`
	ButtonData  json.RawMessage `json:"button_data"`
	CreatedAt   *time.Time      `json:"created_at,omitempty"`
}

// TableName returns the table name
func (BroadcastMessage) TableName() string {
	return "broadcast_messages"
}

// BroadcastRecipient represents a recipient in a broadcast
type BroadcastRecipient struct {
	ID            uuid.UUID  `json:"id"`
	CampaignID    uuid.UUID  `json:"campaign_id"`
	GroupsID      *uuid.UUID `json:"groups_id"`
	ContactID     *uuid.UUID `json:"contact_id"`
	Status        int32      `json:"status"` // 0: pending, 1: sent, 2: failed
	SentAt        *time.Time `json:"sent_at"`
	ErrorMessages string     `json:"error_messages"`
	RetryCount    int32      `json:"retry_count"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}

// TableName returns the table name
func (BroadcastRecipient) TableName() string {
	return "broadcast_recipients"
}

// CreateBroadcastCampaignRequest is the request struct for creating a campaign
type CreateBroadcastCampaignRequest struct {
	DeviceID      uuid.UUID  `json:"device_id" binding:"required"`
	NameBroadcast string     `json:"name_broadcast" binding:"required"`
	ScheduledAt   *time.Time `json:"scheduled_at"`
}

// BroadcastCampaignResponse is the response struct
type BroadcastCampaignResponse struct {
	ID              uuid.UUID `json:"id"`
	NameBroadcast   string    `json:"name_broadcast"`
	TotalRecipients int32     `json:"total_recipients"`
	ProcessedCount  int32     `json:"processed_count"`
	Status          int32     `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// ToResponse converts BroadcastCampaign to response
func (b *BroadcastCampaign) ToResponse() *BroadcastCampaignResponse {
	return &BroadcastCampaignResponse{
		ID:              b.ID,
		NameBroadcast:   b.NameBroadcast,
		TotalRecipients: b.TotalRecipients,
		ProcessedCount:  b.ProcessedCount,
		Status:          b.Status,
		CreatedAt:       b.CreatedAt,
	}
}
