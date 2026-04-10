package models

import (
	"time"

	"github.com/google/uuid"
)

// Device represents a WhatsApp device/session
type Device struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	UniqueName string    `json:"unique_name"`
	NameDevice string    `json:"name_device"`
	Phone      string    `json:"phone"`
	Status     int32     `json:"status"` // 0: inactive, 1: active, 2: disconnect
	LastSeen   *time.Time `json:"last_seen"`
	SessionData []byte   `json:"session_data"` // Encrypted whatsmeow session
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// TableName returns the table name
func (Device) TableName() string {
	return "devices"
}

// IsActive checks if device is active
func (d *Device) IsActive() bool {
	return d.Status == 1
}

// CreateDeviceRequest is the request struct for creating a device
type CreateDeviceRequest struct {
	UniqueName string `json:"unique_name" binding:"required"`
	NameDevice string `json:"name_device" binding:"required"`
}

// UpdateDeviceRequest is the request struct for updating a device
type UpdateDeviceRequest struct {
	NameDevice *string `json:"name_device"`
	Phone      *string `json:"phone"`
}

// UpdateDeviceStatusRequest is the request struct for updating device status
type UpdateDeviceStatusRequest struct {
	Status int32 `json:"status" binding:"required"`
}

// DeviceResponse is the response struct for device endpoints
type DeviceResponse struct {
	ID         uuid.UUID `json:"id"`
	UniqueName string    `json:"unique_name"`
	NameDevice string    `json:"name_device"`
	Phone      string    `json:"phone"`
	Status     int32     `json:"status"`
	LastSeen   *time.Time `json:"last_seen"`
}

// ToResponse converts Device to DeviceResponse
func (d *Device) ToResponse() *DeviceResponse {
	return &DeviceResponse{
		ID:         d.ID,
		UniqueName: d.UniqueName,
		NameDevice: d.NameDevice,
		Phone:      d.Phone,
		Status:     d.Status,
		LastSeen:   d.LastSeen,
	}
}
