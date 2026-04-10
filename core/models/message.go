package models

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a WhatsApp message (in or out)
type Message struct {
	ID             uuid.UUID `json:"id"`
	DeviceID       uuid.UUID `json:"device_id"`
	Direction      string    `json:"direction"` // IN or OUT
	ReceiptNumber  string    `json:"receipt_number"`
	MessageType    int32     `json:"message_type"`
	Content        string    `json:"content"`
	StatusMessage  int32     `json:"status_message"` // 0: pending, 1: sent, 2: delivered, 3: read, 4: failed
	ErrorLog       string    `json:"error_log"`
	CreatedAt      time.Time `json:"created_at"`
}

// TableName returns the table name
func (Message) TableName() string {
	return "messages"
}

// SendMessageRequest is the request struct for sending a message
type SendMessageRequest struct {
	DeviceID    uuid.UUID `json:"device_id" binding:"required"`
	Phone       string    `json:"phone" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	MessageType int32     `json:"message_type"`
}

// MessageResponse is the response struct for message endpoints
type MessageResponse struct {
	ID            uuid.UUID `json:"id"`
	Direction     string    `json:"direction"`
	ReceiptNumber string    `json:"receipt_number"`
	MessageType   int32     `json:"message_type"`
	Content       string    `json:"content"`
	StatusMessage int32     `json:"status_message"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToResponse converts Message to MessageResponse
func (m *Message) ToResponse() *MessageResponse {
	return &MessageResponse{
		ID:            m.ID,
		Direction:     m.Direction,
		ReceiptNumber: m.ReceiptNumber,
		MessageType:   m.MessageType,
		Content:       m.Content,
		StatusMessage: m.StatusMessage,
		CreatedAt:     m.CreatedAt,
	}
}

// GetStatusText returns human-readable status
func (m *Message) GetStatusText() string {
	switch m.StatusMessage {
	case 0:
		return "pending"
	case 1:
		return "sent"
	case 2:
		return "delivered"
	case 3:
		return "read"
	case 4:
		return "failed"
	default:
		return "unknown"
	}
}
