package models

import (
	"time"

	"github.com/google/uuid"
)

// Message status constants
const (
	MessageStatusPending   = "pending"
	MessageStatusSent      = "sent"
	MessageStatusDelivered = "delivered"
	MessageStatusRead      = "read"
	MessageStatusFailed    = "failed"
)

// Message type constants
const (
	MessageTypeText     = "text"
	MessageTypeImage    = "image"
	MessageTypeDocument = "document"
	MessageTypeAudio    = "audio"
	MessageTypeVideo    = "video"
	MessageTypeLocation = "location"
	MessageTypeContact  = "contact"
	MessageTypeSticker  = "sticker"
)

// Message represents a WhatsApp message (inbound or outbound)
type Message struct {
	ID                   uuid.UUID  `json:"id"`
	UserID               uuid.UUID  `json:"user_id"`
	DeviceID             uuid.UUID  `json:"device_id"`
	TemplateID           *uuid.UUID `json:"template_id,omitempty"`
	BroadcastID          *uuid.UUID `json:"broadcast_id,omitempty"`
	ScheduledMessageID   *uuid.UUID `json:"scheduled_message_id,omitempty"`
	Direction            string     `json:"direction"` // inbound / outbound
	RecipientPhone       string     `json:"recipient_phone"`
	SenderPhone          string     `json:"sender_phone"`
	MessageType          string     `json:"message_type"`
	Content              string     `json:"content"`
	MediaURL             *string    `json:"media_url,omitempty"`
	Status               string     `json:"status"`
	WhatsappMessageID    *string    `json:"whatsapp_message_id,omitempty"`
	ErrorLog             *string    `json:"error_log,omitempty"`
	SentAt               *time.Time `json:"sent_at,omitempty"`
	DeliveredAt          *time.Time `json:"delivered_at,omitempty"`
	ReadAt               *time.Time `json:"read_at,omitempty"`
	FailedAt             *time.Time `json:"failed_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (Message) TableName() string {
	return "messages"
}

// SendMessageRequest is the request struct for sending a message
type SendMessageRequest struct {
	DeviceID    uuid.UUID  `json:"device_id" binding:"required"`
	RecipientPhone string  `json:"recipient_phone" binding:"required"`
	Content     string     `json:"content" binding:"required"`
	MessageType string     `json:"message_type"`
	MediaURL    *string    `json:"media_url"`
	TemplateID  *uuid.UUID `json:"template_id"`
}

// MessageResponse is the response struct for message endpoints
type MessageResponse struct {
	ID             uuid.UUID  `json:"id"`
	DeviceID       uuid.UUID  `json:"device_id"`
	Direction      string     `json:"direction"`
	RecipientPhone string     `json:"recipient_phone"`
	SenderPhone    string     `json:"sender_phone"`
	MessageType    string     `json:"message_type"`
	Content        string     `json:"content"`
	Status         string     `json:"status"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	FailedAt       *time.Time `json:"failed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ToResponse converts Message to MessageResponse
func (m *Message) ToResponse() *MessageResponse {
	return &MessageResponse{
		ID:             m.ID,
		DeviceID:       m.DeviceID,
		Direction:      m.Direction,
		RecipientPhone: m.RecipientPhone,
		SenderPhone:    m.SenderPhone,
		MessageType:    m.MessageType,
		Content:        m.Content,
		Status:         m.Status,
		SentAt:         m.SentAt,
		DeliveredAt:    m.DeliveredAt,
		ReadAt:         m.ReadAt,
		FailedAt:       m.FailedAt,
		CreatedAt:      m.CreatedAt,
	}
}
