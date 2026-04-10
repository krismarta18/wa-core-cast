package message

import (
	"context"
	"time"
)

// MessageStatus represents the state of a message
type MessageStatus int

const (
	StatusPending   MessageStatus = 0 // Queued, not yet sent
	StatusSent      MessageStatus = 1 // Sent to WhatsApp server
	StatusDelivered MessageStatus = 2 // Delivered to recipient device
	StatusRead      MessageStatus = 3 // Read by recipient
	StatusFailed    MessageStatus = 4 // Failed to send
)

// MessageDirection represents whether message is incoming or outgoing
type MessageDirection int

const (
	DirectionOut MessageDirection = 0 // Outgoing (sent by us)
	DirectionIn  MessageDirection = 1 // Incoming (received from contact)
)

// QueuedMessage represents a message in the outgoing queue
type QueuedMessage struct {
	ID              string
	DeviceID        string
	TargetJID       string  // Recipient JID
	GroupID         *string // Optional group ID for group messages
	Content         string
	ContentType     string // text, image, document, etc
	MediaURL        *string // URL to media if applicable
	Caption         *string // Caption for media
	Status          MessageStatus
	RetryCount      int
	MaxRetries      int
	LastRetryAt     *time.Time
	Priority        int // 1-5, 5 being highest
	ScheduledFor    *time.Time // For scheduled messages
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ErrorLog        *string // Last error message
}

// ReceivedMessage represents an incoming message from WhatsApp
type ReceivedMessage struct {
	ID              string
	DeviceID        string
	FromJID         string
	GroupJID        *string
	ContentType     string
	Content         string
	MessageID       string // WhatsApp message ID
	ReceiptNumber   *string
	Timestamp       int64
	IsGroup         bool
	SenderName      *string
	CreatedAt       time.Time
}

// MessageStatusUpdate represents a status change event
type MessageStatusUpdate struct {
	MessageID    string
	DeviceID     string
	OldStatus    MessageStatus
	NewStatus    MessageStatus
	Timestamp    int64
	ErrorMessage *string
}

// MessageQueueConfig holds queue configuration
type MessageQueueConfig struct {
	MaxRetries         int           // Maximum retry attempts
	RetryDelayBase     time.Duration // Base delay for exponential backoff
	MaxRetryDelay      time.Duration // Maximum delay between retries
	BatchSize          int           // Messages to process per batch
	ProcessInterval    time.Duration // Interval between processing cycles
	MaxConcurrentSends int           // Max concurrent sends per device
}

// MessageQueue stores messages waiting to be sent
type MessageQueue struct {
	ID         string
	DeviceID   string
	MessageID  string
	Content    string
	TargetJID  string
	RetryCount int
	MaxRetries int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// MessageStore interface for message persistence
type MessageStore interface {
	// Queue operations
	EnqueueMessage(qm *QueuedMessage) error
	DequeueMessages(deviceID string, limit int) ([]*QueuedMessage, error)
	GetQueuedMessage(messageID string) (*QueuedMessage, error)
	UpdateQueuedMessageStatus(messageID string, status MessageStatus, errorMsg *string) error
	UpdateQueuedMessageRetry(messageID string, retryCount int, lastRetryAt *time.Time) error
	MarkMessageSent(messageID string) error
	GetFailedMessages(limit int) ([]*QueuedMessage, error)

	// Message operations
	SaveReceivedMessage(rm *ReceivedMessage) error
	GetMessageByID(messageID string) (*ReceivedMessage, error)
	GetMessagesByDevice(deviceID string, limit int, offset int) ([]*ReceivedMessage, error)
	GetMessagesByJID(jid string, limit int, offset int) ([]*ReceivedMessage, error)

	// Status tracking
	UpdateMessageStatus(messageID string, status MessageStatus) error
	GetMessageStatus(messageID string) (MessageStatus, error)
	GetPendingMessages(deviceID string) ([]*QueuedMessage, error)
	CountByStatus(deviceID string, status MessageStatus) (int, error)

	// Cleanup
	DeleteOldMessages(beforeDate time.Time) error
	ClearFailedMessages(beforeDate time.Time) error
}

// DeliveryCallback is called when message delivery status changes
type DeliveryCallback func(*MessageStatusUpdate)

// ReceiveCallback is called when new message is received
type ReceiveCallback func(*ReceivedMessage)

// ServiceInterface defines the message service contract
type ServiceInterface interface {
	// Sending
	SendMessage(ctx context.Context, deviceID string, targetJID string, content string, groupID *string) (string, error)
	SendMessageWithMedia(ctx context.Context, deviceID string, targetJID string, mediaURL string, contentType string, caption *string) (string, error)
	SendScheduledMessage(ctx context.Context, deviceID string, targetJID string, content string, scheduledFor time.Time) (string, error)

	// Receiving
	ReceiveMessage(rm *ReceivedMessage) error

	// Status tracking
	GetMessageStatus(messageID string) (MessageStatus, error)
	UpdateMessageStatus(messageID string, status MessageStatus, errorMsg *string) error
	GetFailedMessages(limit int) ([]*QueuedMessage, error)

	// Queue
	ProcessQueue() error
	GetQueueStats() map[string]interface{}

	// Callbacks
	RegisterDeliveryCallback(fn DeliveryCallback)
	RegisterReceiveCallback(fn ReceiveCallback)

	// Cleanup
	Cleanup() error
}

// StatisticsFields for message metrics
type StatisticsFields struct {
	TotalSent      int
	TotalReceived  int
	TotalFailed    int
	AverageLatency float64 // Milliseconds
	PendingCount   int
	FailedCount    int
}

// MessageStore implementation helpers
type MessageStoreError string

const (
	ErrMessageNotFound     MessageStoreError = "message not found"
	ErrInvalidStatus      MessageStoreError = "invalid status"
	ErrQueueFull          MessageStoreError = "queue is full"
	ErrDatabaseError      MessageStoreError = "database error"
	ErrMessageAlreadySent MessageStoreError = "message already sent"
)
