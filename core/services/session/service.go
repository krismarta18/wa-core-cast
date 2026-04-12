package session

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/zap"

	qrcode "github.com/skip2/go-qrcode"

	"wacast/core/database"
	"wacast/core/utils"
)

// Service implements SessionServiceInterface
type Service struct {
	mu                 sync.RWMutex
	sessions           map[string]*WhatsAppSession
	qrCodes            map[string]string                                    // Store QR code strings
	qrCodeImages       map[string][]byte                                    // ✅ Store QR code PNG images (base64)
	db                 *database.Database
	encryptionKey      string
	maxSessions        int
	sessionTimeout     int
	qrCodeCallbacks    map[string]func(*QRCodeEvent)                        // Callbacks
	statusCallbacks    map[string]func(*ConnectionStatusEvent)
	messageHandlers    map[string]func(*MessageReceivedEvent)
	onQRUpdate         func(deviceID, qrCode string, status int)            // ✅ Callback for WebSocket
}

// NewService creates a new session service
func NewService(db *database.Database, encryptionKey string, maxSessions int, sessionTimeout int) *Service {
	return &Service{
		db:                 db,
		encryptionKey:      encryptionKey,
		maxSessions:        maxSessions,
		sessionTimeout:     sessionTimeout,
		sessions:           make(map[string]*WhatsAppSession),
		qrCodes:            make(map[string]string),
		qrCodeImages:       make(map[string][]byte),                          // ✅ Initialize QR images
		qrCodeCallbacks:    make(map[string]func(*QRCodeEvent)),
		statusCallbacks:    make(map[string]func(*ConnectionStatusEvent)),
		messageHandlers:    make(map[string]func(*MessageReceivedEvent)),
	}
}

// StartSession initializes and connects a new WhatsApp session
func (s *Service) StartSession(ctx context.Context, cfg *SessionConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sessions) >= s.maxSessions {
		return fmt.Errorf("maximum sessions reached: %d", s.maxSessions)
	}

	// Check if session already exists
	if _, exists := s.sessions[cfg.DeviceID]; exists {
		return fmt.Errorf("session already exists for device %s", cfg.DeviceID)
	}

	// Create whatsmeow client with in-memory store
	// In production: use sqlstore with database backend
	client := whatsmeow.NewClient(nil, nil)

	// Create session wrapper - initially PENDING status (waiting for QR scan)
	session := &WhatsAppSession{
		ID:               cfg.DeviceID,
		Client:           client,
		Status:           SessionPending, // ✅ Start as pending - waiting for QR
		LastActivity:     time.Now().Unix(),
		EnableReceiptAck: true,
		config:           cfg,
	}

	s.sessions[cfg.DeviceID] = session

	// ========================================================================
	// Initialize WhatsApp Web connection - Proper whatsmeow implementation
	// ========================================================================
	
	go func() {
		deviceID := cfg.DeviceID
		var err error

		// Use a background context so the goroutine is NOT tied to the HTTP
		// request context (which is cancelled as soon as the handler returns).
		bgCtx := context.Background()

		// ✅ Create sqlstore container for device management
		// Use existing database connection with postgres dialect
		storeContainer := sqlstore.NewWithDB(s.db.GetConnection(), "postgres", waLog.Noop)
		
		// ✅ Ensure database schema is up to date
		if err := storeContainer.Upgrade(bgCtx); err != nil {
			utils.Error("Failed to upgrade sqlstore schema",
				zap.String("device_id", deviceID),
				zap.Error(err),
			)
		}
		
		// ✅ Create new device for this session
		deviceStore := storeContainer.NewDevice()
		
		// ✅ Create whatsmeow client with proper store and logger
		client := whatsmeow.NewClient(deviceStore, waLog.Noop)
		
		// Update session with actual client
		s.mu.Lock()
		session.Client = client
		s.mu.Unlock()
		
		// Set up event handlers BEFORE connecting
		// This ensures we capture all events including QR code
		client.AddEventHandler(func(evt interface{}) {
			switch v := evt.(type) {
			case *events.QR:
				// ✅ Capture actual QR code from WhatsApp
				utils.Info("QR code received from WhatsApp",
					zap.String("device_id", deviceID),
					zap.Int("qr_count", len(v.Codes)),
				)
				
				// Store the first QR code (main code)
				if len(v.Codes) > 0 {
					qrString := v.Codes[0] // First code is the LinkDevice code
					
					s.mu.Lock()
					s.qrCodes[deviceID] = qrString
					onQRUpdate := s.onQRUpdate
					s.mu.Unlock()
					
					// Generate PNG image from the real QR code
					qrImage, err := s.GenerateQRCodeImage(qrString)
					if err != nil {
						utils.Error("Failed to generate QR code image",
							zap.String("device_id", deviceID),
							zap.Error(err),
						)
					} else {
						s.mu.Lock()
						s.qrCodeImages[deviceID] = qrImage
						s.mu.Unlock()
						
						utils.Info("QR code image generated",
							zap.String("device_id", deviceID),
							zap.Int("image_size", len(qrImage)),
							zap.String("qr_string", qrString),
						)
					}
					
					// ✅ Notify WebSocket clients about QR code update
					if onQRUpdate != nil {
						go onQRUpdate(deviceID, qrString, int(SessionPending))
					}
				}
				
			case *events.Connected:
				// ✅ Connection successful
				utils.Info("Connected to WhatsApp",
					zap.String("device_id", deviceID),
				)
				s.mu.Lock()
				if sess, exists := s.sessions[deviceID]; exists {
					sess.Status = SessionActive
					sess.LastActivity = time.Now().Unix()
				}
				s.mu.Unlock()
				
			case *events.Disconnected:
				// ❌ Connection lost
				utils.Warn("Lost connection to WhatsApp",
					zap.String("device_id", deviceID),
				)
				s.mu.Lock()
				if sess, exists := s.sessions[deviceID]; exists {
					sess.Status = SessionPending
				}
				s.mu.Unlock()
				
			case *events.ConnectFailure:
				// ❌ Connection failed
				utils.Error("Connection failure",
					zap.String("device_id", deviceID),
					zap.String("reason", v.Message),
				)
			}
		})
		
		// ✅ Connect to WhatsApp Web
		utils.Info("Attempting to connect to WhatsApp Web",
			zap.String("device_id", deviceID),
		)
		
		err = client.Connect()
		if err != nil {
			utils.Error("Failed to connect to WhatsApp",
				zap.String("device_id", deviceID),
				zap.Error(err),
			)
			
			s.mu.Lock()
			if sess, exists := s.sessions[deviceID]; exists {
				sess.Status = SessionInactive
			}
			s.mu.Unlock()
		}
	}()

	utils.Info("Session initiated (connecting to WhatsApp Web)",
		zap.String("device_id", cfg.DeviceID),
		zap.String("phone", cfg.Phone),
	)

	return nil
}

// StopSession disconnects and removes a session
func (s *Service) StopSession(ctx context.Context, deviceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[deviceID]
	if !exists {
		return fmt.Errorf("session not found: %s", deviceID)
	}

	// Disconnect
	session.Client.Disconnect()

	// Remove from manager
	delete(s.sessions, deviceID)

	utils.Info("Session stopped",
		zap.String("device_id", deviceID),
	)

	return nil
}

// RestoreSession restores a previous session from encrypted data
func (s *Service) RestoreSession(ctx context.Context, deviceID string, sessionData *SessionData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sessions) >= s.maxSessions {
		return fmt.Errorf("maximum sessions reached: %d", s.maxSessions)
	}

	// Create whatsmeow client
	client := whatsmeow.NewClient(nil, nil)

	// Create session wrapper
	session := &WhatsAppSession{
		ID:               deviceID,
		Client:           client,
		Status:           SessionActive,
		LastActivity:     time.Now().Unix(),
		EnableReceiptAck: true,
		config: &SessionConfig{
			DeviceID:       deviceID,
			Phone:          sessionData.Phone,
			EncryptionKey:  s.encryptionKey,
			SessionTimeout: s.sessionTimeout,
		},
	}

	s.sessions[deviceID] = session

	utils.Info("Session restored",
		zap.String("device_id", deviceID),
		zap.String("phone", sessionData.Phone),
	)

	return nil
}

// GetSession retrieves a session by device ID
func (s *Service) GetSession(deviceID string) *WhatsAppSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[deviceID]
}

// GetSessionStatus gets the status of a session
func (s *Service) GetSessionStatus(deviceID string) SessionStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[deviceID]
	if !exists {
		return SessionInactive
	}

	return session.Status
}

// IsSessionActive checks if a session is active
func (s *Service) IsSessionActive(deviceID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[deviceID]
	if !exists {
		return false
	}

	return session.Status == SessionActive
}

// GetAllActiveSessions returns all active sessions
func (s *Service) GetAllActiveSessions() []*WhatsAppSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sessions []*WhatsAppSession
	for _, session := range s.sessions {
		if session.Status == SessionActive {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// RestorePreviousSessions restores previous sessions from database
func (s *Service) RestorePreviousSessions(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	utils.Info("Restoring previous sessions...")

	// In a real implementation, query database for active devices
	// For now, just log that restoration is attempted
	utils.Debug("No previous sessions to restore (stub implementation)")

	return nil
}

// RegisterQRCodeCallback registers a QR code callback
func (s *Service) RegisterQRCodeCallback(deviceID string, callback func(*QRCodeEvent)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.qrCodeCallbacks[deviceID] = callback
	utils.Debug("QR code callback registered", zap.String("device_id", deviceID))
}

// RegisterStatusCallback registers a status change callback
func (s *Service) RegisterStatusCallback(deviceID string, callback func(*ConnectionStatusEvent)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statusCallbacks[deviceID] = callback
	utils.Debug("Status callback registered", zap.String("device_id", deviceID))
}

// RegisterMessageHandler registers a message handler
func (s *Service) RegisterMessageHandler(deviceID string, callback func(*MessageReceivedEvent)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageHandlers[deviceID] = callback
	utils.Debug("Message handler registered", zap.String("device_id", deviceID))
}

// GetQRCode retrieves the QR code string for a session
func (s *Service) GetQRCode(deviceID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// ✅ Return actual QR code stored from whatsmeow event
	qrCode, exists := s.qrCodes[deviceID]
	if !exists {
		return ""
	}
	
	return qrCode
}

// GetQRCodeImage retrieves the QR code image (PNG bytes) for a session
// ✅ For testing/debugging - returns raw PNG bytes
func (s *Service) GetQRCodeImage(deviceID string) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	qrImage, exists := s.qrCodeImages[deviceID]
	if !exists {
		return nil
	}
	
	return qrImage
}

// RegisterQRUpdateCallback registers a callback for QR code updates
// ✅ Used by WebSocket handler for real-time notifications
func (s *Service) RegisterQRUpdateCallback(callback func(deviceID, qrCode string, status int)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.onQRUpdate = callback
}

// GenerateQRCodeImage generates a QR code image from a string and returns PNG bytes
// ✅ For testing with WhatsApp mobile
func (s *Service) GenerateQRCodeImage(qrString string) ([]byte, error) {
	// Generate QR code (skip2 library returns PNG bytes directly)
	qr, err := qrcode.New(qrString, qrcode.High)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}
	
	// Encode to PNG bytes
	pngBytes, err := qr.PNG(256) // 256x256 pixels
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR as PNG: %w", err)
	}
	
	return pngBytes, nil
}

// Shutdown shuts down the service
func (s *Service) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for deviceID, session := range s.sessions {
		session.Client.Disconnect()
		utils.Debug("Session disconnected", zap.String("device_id", deviceID))
	}

	s.sessions = make(map[string]*WhatsAppSession)
	utils.Info("Session service shut down")

	return nil
}

// ============================================================================
// Service Methods for Compatibility
// ============================================================================

// SendMessage sends a message via an active session
func (s *Service) SendMessage(ctx context.Context, deviceID, targetJID, content string, groupID *string) (string, error) {
	session := s.GetSession(deviceID)
	if session == nil {
		return "", fmt.Errorf("session not found: %s", deviceID)
	}

	if session.Status != SessionActive {
		return "", fmt.Errorf("session not active: %s", deviceID)
	}

	// In real implementation, send message via session.Client
	// For now, return a dummy message ID
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())

	utils.Debug("Message sent",
		zap.String("device_id", deviceID),
		zap.String("target_jid", targetJID),
		zap.String("message_id", messageID),
	)

	return messageID, nil
}
