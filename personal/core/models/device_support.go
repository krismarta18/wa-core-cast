package models

import (
	"time"

	"github.com/google/uuid"
)

// DeviceSession stores the whatsmeow session data for a device
type DeviceSession struct {
	ID          uuid.UUID `json:"id"`
	DeviceID    uuid.UUID `json:"device_id"`
	SessionData []byte    `json:"-"` // encrypted session blob, never expose
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name
func (DeviceSession) TableName() string {
	return "device_sessions"
}

// DeviceQRCode holds a pending QR code for device pairing
type DeviceQRCode struct {
	ID        uuid.UUID `json:"id"`
	DeviceID  uuid.UUID `json:"device_id"`
	QRCode    string    `json:"qr_code"` // base64 or raw string
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName returns the table name
func (DeviceQRCode) TableName() string {
	return "device_qr_codes"
}

// IsExpired checks if the QR code has expired
func (q *DeviceQRCode) IsExpired() bool {
	return time.Now().After(q.ExpiresAt)
}

// DeviceMetrics holds periodic metric snapshots for a device
type DeviceMetrics struct {
	ID              uuid.UUID `json:"id"`
	DeviceID        uuid.UUID `json:"device_id"`
	BatteryLevel    *int      `json:"battery_level,omitempty"`
	MessagesSent    int       `json:"messages_sent"`
	MessagesReceived int      `json:"messages_received"`
	Uptime          int       `json:"uptime"` // seconds
	RecordedAt      time.Time `json:"recorded_at"`
}

// TableName returns the table name
func (DeviceMetrics) TableName() string {
	return "device_metrics"
}
