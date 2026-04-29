package models

import (
	"time"

	"github.com/google/uuid"
)

// Device status constants
const (
	DeviceStatusConnected    = "connected"
	DeviceStatusDisconnected = "disconnected"
	DeviceStatusPendingQR    = "pending_qr"
	DeviceStatusBanned       = "banned"
)

// Device represents a WhatsApp device/session
type Device struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	UniqueName      string     `json:"unique_name"`
	DisplayName     string     `json:"display_name"`
	PhoneNumber     string     `json:"phone_number"`
	Status          string     `json:"status"` // connected/disconnected/pending_qr/banned
	LastSeenAt      *time.Time `json:"last_seen_at,omitempty"`
	ConnectedSince  *time.Time `json:"connected_since,omitempty"`
	Platform        *string    `json:"platform,omitempty"`
	WaVersion       *string    `json:"wa_version,omitempty"`
	BatteryLevel    *int       `json:"battery_level,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	IsWarming       bool       `json:"is_warming"`
	WarmingUntil    *time.Time `json:"warming_until,omitempty"`
}

// TableName returns the table name
func (Device) TableName() string {
	return "devices"
}

// IsConnected checks if device is connected
func (d *Device) IsConnected() bool {
	return d.Status == DeviceStatusConnected
}

// CreateDeviceRequest is the request struct for creating a device
type CreateDeviceRequest struct {
	UniqueName  string `json:"unique_name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

// UpdateDeviceRequest is the request struct for updating a device
type UpdateDeviceRequest struct {
	DisplayName *string `json:"display_name"`
	PhoneNumber *string `json:"phone_number"`
}

// UpdateDeviceStatusRequest is the request struct for updating device status
type UpdateDeviceStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=connected disconnected pending_qr banned"`
}

// DeviceResponse is the response struct for device endpoints
type DeviceResponse struct {
	ID             uuid.UUID  `json:"id"`
	UniqueName     string     `json:"unique_name"`
	DisplayName    string     `json:"display_name"`
	PhoneNumber    string     `json:"phone_number"`
	Status         string     `json:"status"`
	LastSeenAt     *time.Time `json:"last_seen_at,omitempty"`
	ConnectedSince *time.Time `json:"connected_since,omitempty"`
	Platform       *string    `json:"platform,omitempty"`
	WaVersion      *string    `json:"wa_version,omitempty"`
	BatteryLevel   *int       `json:"battery_level,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	IsWarming      bool       `json:"is_warming"`
	WarmingUntil   *time.Time `json:"warming_until,omitempty"`
}

// ToResponse converts Device to DeviceResponse
func (d *Device) ToResponse() *DeviceResponse {
	return &DeviceResponse{
		ID:             d.ID,
		UniqueName:     d.UniqueName,
		DisplayName:    d.DisplayName,
		PhoneNumber:    d.PhoneNumber,
		Status:         d.Status,
		LastSeenAt:     d.LastSeenAt,
		ConnectedSince: d.ConnectedSince,
		Platform:       d.Platform,
		WaVersion:      d.WaVersion,
		BatteryLevel:   d.BatteryLevel,
		CreatedAt:      d.CreatedAt,
	}
}
