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
	ID              string        `json:"id"`
	DeviceID        string        `json:"device_id"`
	TargetJID       string        `json:"target_jid"`  // Recipient JID
	GroupID         *string       `json:"group_id,omitempty"` // Optional group ID for group messages
	Content         string        `json:"content"`
	ContentType     string        `json:"content_type"` // text, image, document, etc
	MediaURL        *string       `json:"media_url,omitempty"` // URL to media if applicable
	Caption         *string       `json:"caption,omitempty"` // Caption for media
	Status          MessageStatus `json:"status"`
	RetryCount      int           `json:"retry_count"`
	MaxRetries      int           `json:"max_retries"`
	LastRetryAt     *time.Time    `json:"last_retry_at,omitempty"`
	Priority        int           `json:"priority"` // 1-5, 5 being highest
	ScheduledFor    *time.Time    `json:"scheduled_for,omitempty"` // For scheduled messages
	BroadcastID     *string       `json:"broadcast_id,omitempty"` // For broadcast messages
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	ErrorLog        *string       `json:"error_log,omitempty"` // Last error message
}

// ReceivedMessage represents an incoming message from WhatsApp
type ReceivedMessage struct {
	ID              string  `json:"id"`
	DeviceID        string  `json:"device_id"`
	FromJID         string  `json:"from_jid"`
	GroupJID        *string `json:"group_jid,omitempty"`
	ContentType     string  `json:"content_type"`
	Content         string  `json:"content"`
	MessageID       string  `json:"message_id"` // WhatsApp message ID
	ReceiptNumber   *string `json:"receipt_number,omitempty"`
	Timestamp       int64   `json:"timestamp"`
	IsGroup         bool    `json:"is_group"`
	SenderName      *string `json:"sender_name,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
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

	// Anti-bot: random delay between each message send
	MinSendDelay    time.Duration // Minimum delay between messages (e.g. 1s)
	MaxSendDelay    time.Duration // Maximum delay between messages (e.g. 5s)

	// Anti-bot: simulate human typing before sending text messages
	// Delay is calculated as: len(content) / TypingSpeedCPM * 60 seconds
	SimulateTyping  bool          // If true, add typing delay proportional to message length
	TypingSpeedCPM  int           // Characters per minute (human avg: 200-400)
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

	// WhatsApp message ID tracking (for receipt matching)
	UpdateWhatsappMessageID(internalID, whatsappMsgID string) error
	GetDBIDByWhatsappID(whatsappMsgID string) (string, error)

	// Cleanup
	DeleteOldMessages(beforeDate time.Time) error
	ClearFailedMessages(beforeDate time.Time) error

	// Scheduled and History
	GetScheduledMessages(deviceID string) ([]*QueuedMessage, error)
	GetMessageHistory(deviceID string, limit int) ([]*QueuedMessage, error)
	DeleteQueuedMessage(messageID string) error
}

// DeliveryCallback is called when message delivery status changes
type DeliveryCallback func(*MessageStatusUpdate)

// ReceiveCallback is called when new message is received
type ReceiveCallback func(*ReceivedMessage)

// ServiceInterface defines the message service contract
type ServiceInterface interface {
	// Sending
	SendMessage(ctx context.Context, deviceID string, targetJID string, content string, groupID *string, broadcastID *string) (string, error)
	SendMessageWithMedia(ctx context.Context, deviceID string, targetJID string, mediaURL string, contentType string, caption *string, broadcastID *string) (string, error)
	SendScheduledMessage(ctx context.Context, deviceID string, targetJID string, content string, scheduledFor time.Time, mediaURL *string, contentType string, caption *string, broadcastID *string) (string, error)

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

	// Scheduled and History
	ListScheduledMessages(deviceID string) ([]*QueuedMessage, error)
	ListMessageHistory(deviceID string, limit int) ([]*QueuedMessage, error)
	CancelScheduledMessage(messageID string) error
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
