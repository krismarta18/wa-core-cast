package models

import (
	"time"

	"github.com/google/uuid"
)

// WarmingPool represents account warming pool settings
type WarmingPool struct {
	ID                uuid.UUID  `json:"id"`
	DeviceID          uuid.UUID  `json:"device_id"`
	Intensity         int32      `json:"intensity"`
	DailyLimit        int32      `json:"daily_limit"`
	MessageSendToday  int32      `json:"message_send_today"`
	IsActive          bool       `json:"is_active"`
	NextActionAt      *time.Time `json:"next_action_at"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
}

// TableName returns the table name
func (WarmingPool) TableName() string {
	return "warming_pool"
}

// WarmingSession represents a warming session
type WarmingSession struct {
	ID              uuid.UUID `json:"id"`
	DeviceID        uuid.UUID `json:"device_id"`
	TargetPhone     string    `json:"target_phone"`
	MessageSent     string    `json:"message_sent"`
	ResponseReceived string   `json:"response_received"`
	Status          int32     `json:"status"` // 0: pending, 1: sent, 2: response_received, 3: failed
	CreatedAt       *time.Time `json:"created_at,omitempty"`
}

// TableName returns the table name
func (WarmingSession) TableName() string {
	return "warming_sessions"
}

// CreateWarmingPoolRequest is the request struct
type CreateWarmingPoolRequest struct {
	DeviceID   uuid.UUID `json:"device_id" binding:"required"`
	Intensity  int32     `json:"intensity" binding:"required"`
	DailyLimit int32     `json:"daily_limit" binding:"required"`
}

// UpdateWarmingPoolRequest is the request struct
type UpdateWarmingPoolRequest struct {
	Intensity  *int32 `json:"intensity"`
	DailyLimit *int32 `json:"daily_limit"`
	IsActive   *bool  `json:"is_active"`
}

// WarmingPoolResponse is the response struct
type WarmingPoolResponse struct {
	ID               uuid.UUID `json:"id"`
	DeviceID         uuid.UUID `json:"device_id"`
	Intensity        int32     `json:"intensity"`
	DailyLimit       int32     `json:"daily_limit"`
	MessageSendToday int32     `json:"message_send_today"`
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
