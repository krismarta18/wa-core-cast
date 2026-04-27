package session

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	qrcode "github.com/skip2/go-qrcode"

	"wacast/core/database"
	"wacast/core/models"
	"wacast/core/services/billing"
	"wacast/core/utils"
	importStore "go.mau.fi/whatsmeow/store"
)

// ReceiptCallback is called when a WhatsApp delivery/read receipt is received.
// Parameters: (whatsappMessageID string, newStatus int)
// Status int: 1=sent, 2=delivered, 3=read, 4=failed (matches MessageStatus in message package)
type ReceiptCallback func(whatsappMessageID string, newStatus int)

// MessageCallback is called when a new WhatsApp message is received.
type MessageCallback func(deviceID string, event *MessageReceivedEvent)

// Service implements SessionServiceInterface
type Service struct {
	mu                 sync.RWMutex
	sessions           map[string]*WhatsAppSession
	qrCodes            map[string]string                                    // Store QR code strings
	qrCodeImages       map[string][]byte                                    // ✅ Store QR code PNG images (base64)
	db                 *database.Database
	billingService     *billing.Service
	encryptionKey      string
	maxSessions        int
	sessionTimeout     int
	qrCodeCallbacks    map[string]func(*QRCodeEvent)                        // Callbacks
	statusCallbacks    map[string]func(*ConnectionStatusEvent)
	messageHandlers    map[string]func(*MessageReceivedEvent)
	receiptCallbacks   []ReceiptCallback                                    // ✅ Global receipt callbacks
	messageCallbacks   []MessageCallback                                    // ✅ Global message callbacks
	onQRUpdate         func(deviceID, qrCode string, status int)            // ✅ Callback for WebSocket
	waStoreContainer   *sqlstore.Container                                  // ✅ Shared store for all sessions
}

// NewService creates a new session service
func NewService(db *database.Database, billingService *billing.Service, encryptionKey string, maxSessions int, timeout int) *Service {
	// ✅ Initialize shared store container and run migrations once at startup
	storeContainer := sqlstore.NewWithDB(db.GetConnection(), "postgres", waLog.Noop)
	err := storeContainer.Upgrade(context.Background())
	if err != nil {
		utils.Error("Failed to upgrade whatsmeow sqlstore", zap.Error(err))
	}

	// ✅ Set consistent device identification (Basic props to avoid undefined constants)
	importStore.DeviceProps.Os = proto.String("Windows")

	return &Service{
		db:             db,
		billingService: billingService,
		encryptionKey:  encryptionKey,
		maxSessions:        maxSessions,
		sessionTimeout:     timeout,
		sessions:           make(map[string]*WhatsAppSession),
		qrCodes:            make(map[string]string),
		qrCodeImages:       make(map[string][]byte),
		qrCodeCallbacks:    make(map[string]func(*QRCodeEvent)),
		statusCallbacks:    make(map[string]func(*ConnectionStatusEvent)),
		messageHandlers:    make(map[string]func(*MessageReceivedEvent)),
		receiptCallbacks:   make([]ReceiptCallback, 0),
		messageCallbacks:   make([]MessageCallback, 0),
		waStoreContainer:   storeContainer,
	}
}

// RegisterMessageCallback registers a callback for incoming messages.
func (s *Service) RegisterMessageCallback(fn MessageCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageCallbacks = append(s.messageCallbacks, fn)
}

// RegisterReceiptCallback registers a callback that is invoked whenever a delivery/read receipt arrives.
func (s *Service) RegisterReceiptCallback(fn ReceiptCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.receiptCallbacks = append(s.receiptCallbacks, fn)
}

// StartSession initializes and connects a new WhatsApp session
func (s *Service) StartSession(ctx context.Context, cfg *SessionConfig) error {
	// --- RECONNECT LOGIC START ---
	// 1. Check if ANY device already has this phone number and is connected
	if oldDevice, err := s.db.GetDeviceByPhone(cfg.Phone); err == nil && oldDevice != nil {
		oldID := oldDevice.ID.String()
		utils.Info("Force reconnect: Phone already active in another session. Stopping old session.", 
			zap.String("phone", cfg.Phone), 
			zap.String("old_device_id", oldID),
			zap.String("new_device_id", cfg.DeviceID),
		)
		// We stop it before proceeding. Logout = true to unlink from WA.
		_ = s.StopSession(ctx, oldID) 
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sessions) >= s.maxSessions {
		return fmt.Errorf("maximum sessions reached: %d", s.maxSessions)
	}

	// 2. Check if a session with the SAME DeviceID already exists in memory
	if _, exists := s.sessions[cfg.DeviceID]; exists {
		utils.Info("Session already exists in memory for this ID. Stopping before restart.", zap.String("device_id", cfg.DeviceID))
		s.stopSessionLocked(ctx, cfg.DeviceID, true)
	}
	// --- RECONNECT LOGIC END ---

	// Create session wrapper - initially PENDING status (waiting for QR scan)
	session := &WhatsAppSession{
		ID:               cfg.DeviceID,
		Client:           nil,            // ✅ Initialize as nil, will be created in goroutine
		Status:           SessionPending, // ✅ Start as pending - waiting for QR
		LastActivity:     time.Now().Unix(),
		EnableReceiptAck: true,
		Config:           cfg,
	}

	s.sessions[cfg.DeviceID] = session

	// Generate Device Record in Database
	devID, errDev := uuid.Parse(cfg.DeviceID)
	if errDev != nil {
		utils.Error("Failed to parse DeviceID as UUID", zap.Error(errDev), zap.String("device_id", cfg.DeviceID))
	} else {
		uID, errUser := uuid.Parse(cfg.UserID)
		if errUser != nil {
			utils.Error("Failed to parse UserID as UUID", zap.Error(errUser), zap.String("user_id", cfg.UserID))
		} else {
			displayName := cfg.DisplayName
			if displayName == "" {
				displayName = "WhatsApp Device"
			}
			device := &models.Device{
				ID:          devID,
				UserID:      uID,
				UniqueName:  cfg.DeviceID,
				DisplayName: displayName,
				PhoneNumber: cfg.Phone,
				Status:      models.DeviceStatusPendingQR,
			}
			if errDB := s.db.CreateDevice(device); errDB != nil {
				utils.Error("Failed to insert device into database", zap.Error(errDB))
			} else {
				utils.Info("Device record successfully inserted into database", zap.String("device_id", cfg.DeviceID))
			}
		}
	}

	// ========================================================================
	// Initialize WhatsApp Web connection - Proper whatsmeow implementation
	// ========================================================================
	
	go func() {
		deviceID := cfg.DeviceID

		// ✅ Create new device for this session from the SHARED container
		deviceStore := s.waStoreContainer.NewDevice()
		
		// ✅ Create whatsmeow client with proper store and logger
		client := whatsmeow.NewClient(deviceStore, waLog.Noop)
		
		// Update session with actual client
		s.mu.Lock()
		session.Client = client
		s.mu.Unlock()
		
		s.setupEventHandlersAndConnect(deviceID, client)
	}()

	utils.Info("Session initiated (connecting to WhatsApp Web)",
		zap.String("device_id", cfg.DeviceID),
		zap.String("phone", cfg.Phone),
	)

	return nil
}

// setupEventHandlersAndConnect centralizes the WhatsApp events and connects to the websocket.
func (s *Service) setupEventHandlersAndConnect(deviceID string, client *whatsmeow.Client) {
	// Set up event handlers BEFORE connecting
	// This ensures we capture all events including QR code
	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
			case *events.PairSuccess:
				// ✅ Pairing completed instantly!
				utils.Info("WhatsApp pairing successful", zap.String("device_id", deviceID))
				
				s.mu.Lock()
				if sess, exists := s.sessions[deviceID]; exists {
					sess.Status = SessionActive
					sess.LastActivity = time.Now().Unix()
				}
				onQRUpdate := s.onQRUpdate
				s.mu.Unlock()
				
				// Update database status asynchronously to avoid blocking event dispatcher
				go func() {
					if parsedID, err := uuid.Parse(deviceID); err == nil {
						_ = s.db.UpdateDeviceStatus(parsedID, models.DeviceStatusConnected)
						
						if client.Store != nil && client.Store.ID != nil {
							actualPhone := client.Store.ID.User
							updatePhone := &models.UpdateDeviceRequest{
								PhoneNumber: &actualPhone,
							}
							_ = s.db.UpdateDeviceInfo(parsedID, updatePhone)
						}
					}
				}()
				
				// Force immediate WebSocket notification so UI updates instantly
				if onQRUpdate != nil {
					go onQRUpdate(deviceID, "", int(SessionActive))
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
				
				// Update status in Database asynchronously
				go func() {
					if parsedID, err := uuid.Parse(deviceID); err == nil {
						_ = s.db.UpdateDeviceStatus(parsedID, models.DeviceStatusConnected)
						
						// Ensure we save the REAL WhatsApp number to the device's phone_number
						// This acts as the bridge for multi-session restoration!
						if client.Store != nil && client.Store.ID != nil {
							actualPhone := client.Store.ID.User
							updatePhone := &models.UpdateDeviceRequest{
								PhoneNumber: &actualPhone,
							}
							_ = s.db.UpdateDeviceInfo(parsedID, updatePhone)
						}
					}
				}()

				// ✅ Notify WebSocket clients about connection success
				s.mu.RLock()
				onQRUpdate := s.onQRUpdate
				s.mu.RUnlock()
				if onQRUpdate != nil {
					go onQRUpdate(deviceID, "", int(SessionActive))
				}
				
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
				
				// Update status in Database asynchronously
				go func() {
					if parsedID, err := uuid.Parse(deviceID); err == nil {
						_ = s.db.UpdateDeviceStatus(parsedID, models.DeviceStatusDisconnected)
					}
				}()

				// ✅ Notify WebSocket clients about disconnection
				s.mu.RLock()
				onQRUpdate := s.onQRUpdate
				s.mu.RUnlock()
				if onQRUpdate != nil {
					go onQRUpdate(deviceID, "", int(SessionPending))
				}
				
			case *events.ConnectFailure:
				// ❌ Connection failed
				utils.Error("Connection failure",
					zap.String("device_id", deviceID),
					zap.String("reason", v.Message),
				)

			case *events.LoggedOut:
				// ❌ Device unlinked by user from WhatsApp App
				utils.Warn("Device unlinked from WhatsApp app",
					zap.String("device_id", deviceID),
					zap.String("reason", v.Reason.String()),
				)
				s.mu.Lock()
				if sess, exists := s.sessions[deviceID]; exists {
					sess.Status = SessionInactive // Set to inactive because session is invalid
				}
				s.mu.Unlock()
				
				// Update status in Database asynchronously
				go func() {
					if parsedID, err := uuid.Parse(deviceID); err == nil {
						_ = s.db.UpdateDeviceStatus(parsedID, models.DeviceStatusDisconnected)
					}
				}()

				// ✅ Notify WebSocket clients about logout
				s.mu.RLock()
				onQRUpdate := s.onQRUpdate
				s.mu.RUnlock()
				if onQRUpdate != nil {
					go onQRUpdate(deviceID, "", int(SessionInactive))
				}

			case *events.Receipt:
				// ✅ Message delivery / read receipt
				// v.Type == types.ReceiptTypeDelivered → pesan diterima di HP penerima
				// v.Type == types.ReceiptTypeRead     → pesan sudah dibaca
				var newStatus int
				switch v.Type {
				case types.ReceiptTypeDelivered:
					newStatus = 2 // StatusDelivered
				case types.ReceiptTypeRead, types.ReceiptTypeReadSelf:
					newStatus = 3 // StatusRead
				default:
					newStatus = 0 // Unknown, skip
				}

				if newStatus > 0 {
					s.mu.RLock()
					callbacks := s.receiptCallbacks
					s.mu.RUnlock()
					for _, msgID := range v.MessageIDs {
						utils.Debug("Receipt received",
							zap.String("device_id", deviceID),
							zap.String("wa_msg_id", msgID),
							zap.Int("status", newStatus),
						)
						for _, cb := range callbacks {
							go cb(msgID, newStatus)
						}
					}
				}

			case *events.Message:
				// Avoid processing our own messages sent from other devices (if needed)
				if v.Info.IsFromMe {
					return
				}
				// Log raw message info
				utils.Debug("Raw WhatsApp message received",
					zap.String("device_id", deviceID),
					zap.String("sender", v.Info.Sender.String()),
					zap.String("chat", v.Info.Chat.String()),
				)

				// Extract Message Text Content
				content := ""
				if v.Message.GetConversation() != "" {
					content = v.Message.GetConversation()
				} else if v.Message.GetExtendedTextMessage().GetText() != "" {
					content = v.Message.GetExtendedTextMessage().GetText()
				} else if v.Message.GetImageMessage().GetCaption() != "" {
					content = v.Message.GetImageMessage().GetCaption()
				} else if v.Message.GetVideoMessage().GetCaption() != "" {
					content = v.Message.GetVideoMessage().GetCaption()
				} else if v.Message.GetButtonsResponseMessage().GetSelectedDisplayText() != "" {
					content = v.Message.GetButtonsResponseMessage().GetSelectedDisplayText()
				} else if v.Message.GetTemplateButtonReplyMessage().GetSelectedDisplayText() != "" {
					content = v.Message.GetTemplateButtonReplyMessage().GetSelectedDisplayText()
				}

				// If we found text content, trigger callbacks
				if content != "" {
					utils.Info("Incoming text message",
						zap.String("device_id", deviceID),
						zap.String("from", v.Info.Sender.String()),
						zap.String("content", content),
					)
					event := &MessageReceivedEvent{
						DeviceID:    deviceID,
						FromJID:     v.Info.Sender.String(),
						GroupJID:    v.Info.Chat.String(),
						MessageID:   v.Info.ID,
						Content:     content,
						ContentType: "text",
						Timestamp:   v.Info.Timestamp.Unix(),
						IsGroup:     v.Info.IsGroup,
					}
					
					if !v.Info.IsGroup {
						event.GroupJID = ""
					}

					s.mu.RLock()
					callbacks := s.messageCallbacks
					s.mu.RUnlock()

					for _, cb := range callbacks {
						go cb(deviceID, event)
					}
				}
			}
		})
		
	// ✅ Get QR Channel for auto-refreshing QR codes if not logged in
	var qrChan <-chan whatsmeow.QRChannelItem
	if client.Store.ID == nil {
		qrChan, _ = client.GetQRChannel(context.Background())
	}

	// ✅ Start listening to QR channel for auto-refreshing QR codes (must be before Connect)
	if qrChan != nil {
		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					qrString := evt.Code
					utils.Info("New QR code received from WhatsApp channel", zap.String("device_id", deviceID))
					
					s.mu.Lock()
					s.qrCodes[deviceID] = qrString
					onQRUpdate := s.onQRUpdate
					s.mu.Unlock()
					
					// Generate PNG image from the real QR code
					qrImage, err := s.GenerateQRCodeImage(qrString)
					if err != nil {
						utils.Error("Failed to generate QR code image", zap.String("device_id", deviceID), zap.Error(err))
					} else {
						s.mu.Lock()
						s.qrCodeImages[deviceID] = qrImage
						s.mu.Unlock()
					}
					
					// Notify WebSocket clients about QR code update
					if onQRUpdate != nil {
						go onQRUpdate(deviceID, qrString, int(SessionPending))
					}
				} else {
					utils.Info("QR channel event", zap.String("device_id", deviceID), zap.String("event", evt.Event))
				}
			}
		}()
	}

	// ✅ Connect to WhatsApp Web
	utils.Info("Attempting to connect to WhatsApp Web",
		zap.String("device_id", deviceID),
	)
	
	err := client.Connect()
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
		return
	}
}

// StopSession disconnects and removes a session
func (s *Service) StopSession(ctx context.Context, deviceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.stopSessionLocked(ctx, deviceID, true) // Default to Logout for "Putuskan"
}

// stopSessionLocked performs the actual cleanup logic.
// CALLER MUST HOLD s.mu LOCK.
func (s *Service) stopSessionLocked(ctx context.Context, deviceID string, logout bool) error {
	session, exists := s.sessions[deviceID]
	if !exists {
		// If not in memory, check if it exists in database to at least update status there
		if parsedID, err := uuid.Parse(deviceID); err == nil {
			// Check if device exists in DB
			if _, err := s.db.GetDeviceByID(parsedID); err == nil {
				utils.Info("Session not in memory, updating DB status only", zap.String("device_id", deviceID))
				_ = s.db.UpdateDeviceStatus(parsedID, models.DeviceStatusDisconnected)
				return nil
			}
		}
		return fmt.Errorf("session not found: %s", deviceID)
	}

	// 1. Unlink session from WhatsApp (Logout) or just Disconnect
	if session.Client != nil {
		if logout && session.Client.IsConnected() {
			utils.Info("Logging out WhatsApp session", zap.String("device_id", deviceID))
			if err := session.Client.Logout(ctx); err != nil {
				utils.Warn("Failed to logout session, falling back to disconnect",
					zap.String("device_id", deviceID),
					zap.Error(err),
				)
				session.Client.Disconnect()
			}
		} else {
			session.Client.Disconnect()
		}
	}

	// 2. Remove from manager
	delete(s.sessions, deviceID)

	// 3. Update status in Database immediately
	if parsedID, err := uuid.Parse(deviceID); err == nil {
		_ = s.db.UpdateDeviceStatus(parsedID, models.DeviceStatusDisconnected)
	}

	utils.Info("Session stopped/unlinked",
		zap.String("device_id", deviceID),
		zap.Bool("was_logout", logout),
	)

	return nil
}

// DeleteSession stops the session and removes the device from the database
func (s *Service) DeleteSession(ctx context.Context, deviceID string) error {
	// 1. Try to stop the session if it's running in memory
	_ = s.StopSession(ctx, deviceID)

	// 2. Delete from database (mark as banned)
	parsedID, err := uuid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("invalid device id: %w", err)
	}

	if err := s.db.DeleteDevice(parsedID); err != nil {
		return fmt.Errorf("failed to delete device from database: %w", err)
	}

	utils.Info("Device deleted successfully", zap.String("device_id", deviceID))
	return nil
}

// RestoreSession restores a previous session from encrypted data
func (s *Service) RestoreSession(ctx context.Context, deviceID string, sessionData *SessionData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sessions) >= s.maxSessions {
		return fmt.Errorf("maximum sessions reached: %d", s.maxSessions)
	}

	// Create whatsmeow client with SHARED store
	deviceStore := s.waStoreContainer.NewDevice()
	client := whatsmeow.NewClient(deviceStore, waLog.Noop)

	// Create session wrapper
	session := &WhatsAppSession{
		ID:               deviceID,
		Client:           client,
		Status:           SessionActive,
		LastActivity:     time.Now().Unix(),
		EnableReceiptAck: true,
		Config: &SessionConfig{
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

// GetAllSessions returns all sessions for a user, merging DB data and in-memory state
func (s *Service) GetAllSessions(userID string) ([]*WhatsAppSession, error) {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	// Fetch devices from database
	dbDevices, err := s.db.GetDevicesByUserID(uID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch devices from db: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var sessions []*WhatsAppSession
	for _, dev := range dbDevices {
		deviceID := dev.ID.String()

		// Check if we have an in-memory session
		if sess, exists := s.sessions[deviceID]; exists {
			sessions = append(sessions, sess)
			continue
		}

		// Otherwise, create a representative session object from DB data
		status := SessionInactive
		switch dev.Status {
		case models.DeviceStatusConnected:
			status = SessionActive
		case models.DeviceStatusPendingQR:
			status = SessionPending
		}

		sessions = append(sessions, &WhatsAppSession{
			ID:     deviceID,
			Status: status,
			Config: &SessionConfig{
				DeviceID:    deviceID,
				UserID:      dev.UserID.String(),
				Phone:       dev.PhoneNumber,
				DisplayName: dev.DisplayName,
			},
		})
	}

	return sessions, nil
}

// GetUserID retrieves the owner UserID of a given device session.
func (s *Service) GetUserID(deviceID string) (uuid.UUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[deviceID]
	if !exists || session.Config == nil {
		return uuid.Nil, fmt.Errorf("session not found: %s", deviceID)
	}

	return uuid.Parse(session.Config.UserID)
}

// RestorePreviousSessions restores previous sessions from database
func (s *Service) RestorePreviousSessions(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	utils.Info("Restoring previous sessions...")

	// 1. Fetch Whatsmeow device stores that have logged in Data from SHARED container
	waDevices, err := s.waStoreContainer.GetAllDevices(context.Background())
	if err != nil {
		utils.Error("Failed to get whatsmeow devices from database", zap.Error(err))
		return err
	}

	// 3. Fetch active/connected backend Devices
	dbDevices, err := s.db.GetActiveDevices()
	if err != nil {
		utils.Error("Failed to get backend devices from database", zap.Error(err))
		return err
	}

	// 4. Match and Restore
	restoredCount := 0
	for _, dbDevice := range dbDevices {
		var matchedStore *importStore.Device
		
		// Bridging logic: Check if WhatsApp phone number matches our db device phone number!
		for _, wDevice := range waDevices {
			if wDevice.ID != nil && wDevice.ID.User == dbDevice.PhoneNumber {
				matchedStore = wDevice
				break
			}
		}

		if matchedStore == nil {
			utils.Warn("Device missing in whatsmeow storage. Moving to disconnected.", 
				zap.String("device_id", dbDevice.ID.String()),
				zap.String("phone", dbDevice.PhoneNumber),
			)
			_ = s.db.UpdateDeviceStatus(dbDevice.ID, models.DeviceStatusDisconnected)
			continue
		}

		// Recreate configuration block
		cfg := &SessionConfig{
			DeviceID:       dbDevice.ID.String(),
			UserID:         dbDevice.UserID.String(),
			Phone:          dbDevice.PhoneNumber,
			DisplayName:    dbDevice.DisplayName,
			EncryptionKey:  s.encryptionKey,
			SessionTimeout: s.sessionTimeout,
		}

		// Allocate new whatsmeow client utilizing the restored Storage State
		client := whatsmeow.NewClient(matchedStore, waLog.Noop)

		// Create Session Memory Instance
		session := &WhatsAppSession{
			ID:               dbDevice.ID.String(),
			Client:           client,
			Status:           SessionActive, // Resume active
			LastActivity:     time.Now().Unix(),
			EnableReceiptAck: true,
			Config:           cfg,
		}

		s.sessions[cfg.DeviceID] = session

		// Perform WebSocket Reconnection asynchronously 
		go s.setupEventHandlersAndConnect(cfg.DeviceID, client)
		
		restoredCount++
	}

	utils.Info("Successfully triggered session restoration", zap.Int("total_restored", restoredCount))
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

// SendMessage sends a text message via the active WhatsApp session using whatsmeow
func (s *Service) SendMessage(ctx context.Context, deviceID, targetJID, content string, groupID *string) (string, error) {
	session := s.GetSession(deviceID)
	if session == nil {
		return "", fmt.Errorf("session not found: %s", deviceID)
	}

	if session.Status != SessionActive {
		return "", fmt.Errorf("session not active: %s", deviceID)
	}

	if session.Client == nil {
		return "", fmt.Errorf("whatsapp client not initialized for device: %s", deviceID)
	}

	if !session.Client.IsConnected() {
		return "", fmt.Errorf("whatsapp client not connected for device: %s", deviceID)
	}

	// Parse target JID
	// targetJID can be phone number like "6285887373722" or full JID "6285887373722@s.whatsapp.net"
	var recipient types.JID
	var err error

	if strings.Contains(targetJID, "@") {
		// Already a full JID
		recipient, err = types.ParseJID(targetJID)
		if err != nil {
			return "", fmt.Errorf("invalid JID format %q: %w", targetJID, err)
		}
	} else {
		// Phone number only — append @s.whatsapp.net for personal chat
		// or use groupID for group messages
		if groupID != nil && *groupID != "" {
			recipient, err = types.ParseJID(*groupID + "@g.us")
			if err != nil {
				return "", fmt.Errorf("invalid group JID %q: %w", *groupID, err)
			}
		} else {
			recipient = types.NewJID(targetJID, types.DefaultUserServer)
		}
	}

	// Build text message proto
	msg := &waE2E.Message{
		Conversation: proto.String(content),
	}

	// Send via whatsmeow
	resp, err := session.Client.SendMessage(ctx, recipient, msg)
	if err != nil {
		utils.Error("Failed to send WhatsApp message",
			zap.String("device_id", deviceID),
			zap.String("target_jid", targetJID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	utils.Info("WhatsApp message sent successfully",
		zap.String("device_id", deviceID),
		zap.String("target_jid", targetJID),
		zap.String("message_id", resp.ID),
	)

	return resp.ID, nil
}

// SendMessageWithMedia downloads media from mediaURL, uploads to WhatsApp, and sends to targetJID.
// contentType: "image" | "document" | "audio" | "video"
func (s *Service) SendMessageWithMedia(ctx context.Context, deviceID, targetJID, mediaURL, contentType string, caption *string) (string, error) {
	session := s.GetSession(deviceID)
	if session == nil {
		return "", fmt.Errorf("session not found: %s", deviceID)
	}
	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("whatsapp client not connected for device: %s", deviceID)
	}

	// 1. Parse recipient JID
	var recipient types.JID
	var err error
	if strings.Contains(targetJID, "@") {
		recipient, err = types.ParseJID(targetJID)
		if err != nil {
			return "", fmt.Errorf("invalid JID %q: %w", targetJID, err)
		}
	} else {
		recipient = types.NewJID(targetJID, types.DefaultUserServer)
	}

	// 2. Download media bytes from URL
	resp2, err := httpGet(mediaURL)
	if err != nil {
		return "", fmt.Errorf("failed to download media from %q: %w", mediaURL, err)
	}
	defer resp2.Body.Close()

	mediaBytes, err := io.ReadAll(resp2.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read media body: %w", err)
	}

	mimeType := resp2.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 3. Determine whatsmeow media type & build proto message
	var waMsg *waE2E.Message

	captionStr := ""
	if caption != nil {
		captionStr = *caption
	}

	switch contentType {
	case "image":
		uploaded, err := session.Client.Upload(ctx, mediaBytes, whatsmeow.MediaImage)
		if err != nil {
			return "", fmt.Errorf("failed to upload image: %w", err)
		}
		waMsg = &waE2E.Message{
			ImageMessage: &waE2E.ImageMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(mediaBytes))),
				Mimetype:      proto.String(mimeType),
				Caption:       proto.String(captionStr),
			},
		}

	case "video":
		uploaded, err := session.Client.Upload(ctx, mediaBytes, whatsmeow.MediaVideo)
		if err != nil {
			return "", fmt.Errorf("failed to upload video: %w", err)
		}
		waMsg = &waE2E.Message{
			VideoMessage: &waE2E.VideoMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(mediaBytes))),
				Mimetype:      proto.String(mimeType),
				Caption:       proto.String(captionStr),
			},
		}

	case "audio":
		uploaded, err := session.Client.Upload(ctx, mediaBytes, whatsmeow.MediaAudio)
		if err != nil {
			return "", fmt.Errorf("failed to upload audio: %w", err)
		}
		waMsg = &waE2E.Message{
			AudioMessage: &waE2E.AudioMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(mediaBytes))),
				Mimetype:      proto.String(mimeType),
			},
		}

	default: // "document" and anything else
		uploaded, err := session.Client.Upload(ctx, mediaBytes, whatsmeow.MediaDocument)
		if err != nil {
			return "", fmt.Errorf("failed to upload document: %w", err)
		}
		waMsg = &waE2E.Message{
			DocumentMessage: &waE2E.DocumentMessage{
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(mediaBytes))),
				Mimetype:      proto.String(mimeType),
				Caption:       proto.String(captionStr),
			},
		}
	}

	// 4. Send message
	sendResp, err := session.Client.SendMessage(ctx, recipient, waMsg)
	if err != nil {
		utils.Error("Failed to send WhatsApp media message",
			zap.String("device_id", deviceID),
			zap.String("target_jid", targetJID),
			zap.String("content_type", contentType),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to send media message: %w", err)
	}

	utils.Info("WhatsApp media message sent successfully",
		zap.String("device_id", deviceID),
		zap.String("target_jid", targetJID),
		zap.String("content_type", contentType),
		zap.String("message_id", sendResp.ID),
	)

	return sendResp.ID, nil
}

// httpGet is a simple HTTP GET helper (abstracted for testability)
var httpGet = func(url string) (*http.Response, error) {
	return http.Get(url) //nolint:noctx
}
