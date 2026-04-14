package session

import (
	"context"
	"sync"

	"go.mau.fi/whatsmeow"
)

// SessionStatus represents the current state of a WhatsApp session
type SessionStatus int

const (
	SessionInactive SessionStatus = 0 // Device not connected
	SessionActive   SessionStatus = 1 // Device connected and ready
	SessionPending  SessionStatus = 2 // QR code waiting to be scanned
)

// QRCodeEvent is emitted when QR code is ready to be scanned
type QRCodeEvent struct {
	DeviceID  string
	QRCode    []byte // Binary QR code data
	QRCodeURL string // Data URL for easy display
}

// ConnectionStatusEvent is emitted when connection status changes
type ConnectionStatusEvent struct {
	DeviceID string
	Status   SessionStatus
	Error    string // Optional error message
}

// MessageStatusEvent is emitted when message status changes
type MessageStatusEvent struct {
	DeviceID      string
	MessageID     string
	Status        string // pending, sent, delivered, read, failed
	Timestamp     int64
	ParticipantJID string // For group message status
}

// MessageReceivedEvent is emitted when a message is received
type MessageReceivedEvent struct {
	DeviceID    string
	FromJID     string
	GroupJID    string // Empty for direct messages
	MessageID   string
	Content     string
	ContentType string // text, image, document, audio, video, contact, location, etc
	Timestamp   int64
	IsGroup     bool
}

// SessionData stores encrypted session information
type SessionData struct {
	DeviceID       string
	Phone          string
	WID            string // WhatsApp ID (phone@s.whatsapp.net)
	EncryptedData  []byte // Encrypted session data
	BackupToken    []byte // Backup token for recovery
	ClientToken    []byte // Client token
	ServerToken    []byte // Server token
	IsConnected    bool
	LastSeen       int64
	PushedName     string
}

// SessionConfig holds configuration for session initialization
type SessionConfig struct {
	DeviceID        string
	UserID          string
	Phone           string
	DisplayName     string
	EncryptionKey   string // 32-byte hex string
	SessionTimeout  int    // Seconds before auto-disconnect
	ReconnectLimit  int    // Maximum reconnection attempts
	reconnectCount  int
}

// WhatsAppSession wraps a whatsmeow client connection
type WhatsAppSession struct {
	mu               sync.RWMutex
	ID               string
	Client           *whatsmeow.Client
	Status           SessionStatus
	LastActivity     int64
	EnableReceiptAck bool
	config           *SessionConfig
}

// SessionManager manages multiple WhatsApp sessions
type SessionManager struct {
	mu              sync.RWMutex
	sessions        map[string]*WhatsAppSession      // deviceID -> session
	configs         map[string]*SessionConfig         // deviceID -> config
	qrCodeCallbacks map[string]func(*QRCodeEvent)    // deviceID -> qr handler
	statusCallbacks map[string]func(*ConnectionStatusEvent) // deviceID -> status handler
	messageHandlers map[string]func(*MessageReceivedEvent)  // deviceID -> message handler
	db              interface{} // *database.Database
	encryptionKey   string
	maxSessions     int
	sessionTimeout  int
}

// EventHandler is the interface for event handling
type EventHandler interface {
	OnQRCode(event *QRCodeEvent)
	OnConnectionStatusChange(event *ConnectionStatusEvent)
	OnMessageStatus(event *MessageStatusEvent)
	OnMessageReceived(event *MessageReceivedEvent)
}

// ServiceInterface defines the session service contract
type ServiceInterface interface {
	// Session lifecycle
	StartSession(ctx context.Context, cfg *SessionConfig) error
	StopSession(ctx context.Context, deviceID string) error
	RestoreSession(ctx context.Context, deviceID string, sessionData *SessionData) error
	
	// Session queries
	GetSession(deviceID string) *WhatsAppSession
	GetAllActiveSessions() []*WhatsAppSession
	GetSessionStatus(deviceID string) SessionStatus
	IsSessionActive(deviceID string) bool
	
	// Message operations
	SendMessage(ctx context.Context, deviceID string, targetJID string, message string, isGroup bool) (string, error)
	SendMessageWithMedia(ctx context.Context, deviceID string, targetJID string, mediaURL string, mediaType string, caption string) (string, error)
	
	// QR code & setup
	RegisterQRCodeCallback(deviceID string, handler func(*QRCodeEvent))
	RegisterStatusCallback(deviceID string, handler func(*ConnectionStatusEvent))
	RegisterMessageHandler(deviceID string, handler func(*MessageReceivedEvent))
	
	// Auto-restore
	RestorePreviousSessions(ctx context.Context) error
	
	// Cleanup
	Shutdown(ctx context.Context) error
}

// ReceiptAckConfig for message receipt configuration
type ReceiptAckConfig struct {
	EnableReadReceipts   bool
	EnableDeliveryStatus bool
}
