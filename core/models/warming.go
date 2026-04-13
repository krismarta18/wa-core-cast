package models

import (
	"time"

	"github.com/google/uuid"
)

// WarmingSession status constants
const (
	WarmingStatusPending  = "pending"
	WarmingStatusSent     = "sent"
	WarmingStatusReplied  = "replied"
	WarmingStatusFailed   = "failed"
)

// WarmingPool represents account warming pool settings
type WarmingPool struct {
	ID               uuid.UUID  `json:"id"`
	DeviceID         uuid.UUID  `json:"device_id"`
	Intensity        int        `json:"intensity"`
	DailyLimit       int        `json:"daily_limit"`
	MessageSendToday int        `json:"message_send_today"`
	IsActive         bool       `json:"is_active"`
	NextActionAt     *time.Time `json:"next_action_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// TableName returns the table name
func (WarmingPool) TableName() string {
	return "warming_pool"
}

// WarmingSession represents a warming session
type WarmingSession struct {
	ID               uuid.UUID  `json:"id"`
	DeviceID         uuid.UUID  `json:"device_id"`
	TargetPhone      string     `json:"target_phone"`
	MessageSent      string     `json:"message_sent"`
	ResponseReceived *string    `json:"response_received,omitempty"`
	Status           string     `json:"status"` // pending/sent/replied/failed
	CreatedAt        time.Time  `json:"created_at"`
}

// TableName returns the table name
func (WarmingSession) TableName() string {
	return "warming_sessions"
}

// CreateWarmingPoolRequest is the request struct
type CreateWarmingPoolRequest struct {
	DeviceID   uuid.UUID `json:"device_id" binding:"required"`
	Intensity  int       `json:"intensity" binding:"required"`
	DailyLimit int       `json:"daily_limit" binding:"required"`
}

// UpdateWarmingPoolRequest is the request struct
type UpdateWarmingPoolRequest struct {
	Intensity  *int  `json:"intensity"`
	DailyLimit *int  `json:"daily_limit"`
	IsActive   *bool `json:"is_active"`
}

// WarmingPoolResponse is the response struct
type WarmingPoolResponse struct {
	ID               uuid.UUID `json:"id"`
	DeviceID         uuid.UUID `json:"device_id"`
	Intensity        int       `json:"intensity"`
	DailyLimit       int       `json:"daily_limit"`
	MessageSendToday int       `json:"message_send_today"`
	IsActive         bool      `json:"is_active"`
}

// ToResponse converts WarmingPool to response
func (w *WarmingPool) ToResponse() *WarmingPoolResponse {
	return &WarmingPoolResponse{
		ID:               w.ID,
		DeviceID:         w.DeviceID,
		Intensity:        w.Intensity,
		DailyLimit:       w.DailyLimit,
		MessageSendToday: w.MessageSendToday,
		IsActive:         w.IsActive,
	}
}
