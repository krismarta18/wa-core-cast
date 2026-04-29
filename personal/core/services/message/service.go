package message

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"wacast/core/database"
	"wacast/core/services/analytics"
	"wacast/core/services/billing"
	"wacast/core/services/integration"
	"wacast/core/services/session"
	"wacast/core/utils"
)

// Service implements the message service interface
type Service struct {
	mu                  sync.RWMutex
	db                  *database.Database
	store               MessageStore
	config              *MessageQueueConfig
	sessionService      *session.Service
	deliveryCallbacks   []DeliveryCallback
	receiveCallbacks    []ReceiveCallback
	processing          bool
	processorTicker     *time.Ticker
	done                chan struct{}
	metrics             *ServiceMetrics
	analyticsService    *analytics.Service
	billingService      *billing.Service
	integrationService  *integration.Service
}

// ServiceMetrics holds message service statistics
type ServiceMetrics struct {
	mu                sync.RWMutex
	TotalSent         int64
	TotalReceived     int64
	TotalFailed       int64
	CurrentPending    int64
	AverageLatency    float64
	SuccessfulRetries int64
}

// DefaultQueueConfig returns default configuration with human-like anti-bot delays
func DefaultQueueConfig() *MessageQueueConfig {
	return &MessageQueueConfig{
		MaxRetries:         3,
		RetryDelayBase:     5 * time.Second,
		MaxRetryDelay:      5 * time.Minute,
		BatchSize:          10,              // Smaller batch for more natural pacing
		ProcessInterval:    3 * time.Second, // Check queue every 3 seconds
		MaxConcurrentSends: 1,              // Sequential per device (anti-bot)
		MinSendDelay:       2 * time.Second, // Min gap between messages
		MaxSendDelay:       8 * time.Second, // Max gap between messages
		SimulateTyping:     true,            // Add typing delay for text
		TypingSpeedCPM:     300,             // ~300 chars/min (human speed)
		AntiBotEnabled:     true,            // Enabled by default in Personal Pro
		RandomSuffixLength: 4,               // 4 chars suffix by default
	}
}

// NewService creates a new message service
func NewService(db *database.Database, sessionService *session.Service, analyticsService *analytics.Service, billingService *billing.Service, integrationService *integration.Service, config *MessageQueueConfig) *Service {
	if config == nil {
		config = DefaultQueueConfig()
	}

	store := NewDatabaseMessageStore(db)

	svc := &Service{
		db:              db,
		store:           store,
		config:          config,
		sessionService:  sessionService,
		analyticsService: analyticsService,
		billingService:   billingService,
		integrationService: integrationService,
		deliveryCallbacks: make([]DeliveryCallback, 0),
		receiveCallbacks:  make([]ReceiveCallback, 0),
		done:             make(chan struct{}),
		metrics: &ServiceMetrics{
			TotalSent:      0,
			TotalReceived:  0,
			TotalFailed:    0,
			CurrentPending: 0,
		},
	}

	// Register receipt callback
	sessionService.RegisterReceiptCallback(func(whatsappMsgID string, newStatus int) {
		internalID, err := store.GetDBIDByWhatsappID(whatsappMsgID)
		if err != nil {
			return
		}
		
		msg, err := store.GetQueuedMessage(internalID)
		if err != nil || msg == nil {
			return
		}
		
		deviceID, _ := uuid.Parse(msg.DeviceID)
		userID, _ := sessionService.GetUserID(msg.DeviceID)
		
		if svc.analyticsService != nil {
			svc.analyticsService.RecordDelivery(userID, deviceID)
		}
		
		// Trigger Webhook
		if svc.integrationService != nil {
			svc.integrationService.TriggerWebhook(userID, "message.status_updated", map[string]interface{}{
				"message_id": internalID,
				"status":     newStatus,
				"whatsapp_id": whatsappMsgID,
			})
		}

		if err := store.UpdateQueuedMessageStatus(internalID, MessageStatus(newStatus), nil); err != nil {
			utils.Error("Failed to update message status from receipt",
				zap.String("wa_msg_id", whatsappMsgID),
				zap.Int("status", newStatus),
				zap.Error(err),
			)
		} else {
			// Record delivery/read analytics
			if newStatus == 2 || newStatus == 3 {
				if msg, err := store.GetQueuedMessage(internalID); err == nil {
					uID, err := sessionService.GetUserID(msg.DeviceID)
					if err == nil {
						dID, _ := uuid.Parse(msg.DeviceID)
						_ = analyticsService.RecordDelivery(uID, dID)
					}
				}
			}
		}
	})

	// Register message callback
	sessionService.RegisterMessageCallback(func(deviceID string, evt *session.MessageReceivedEvent) {
		rm := &ReceivedMessage{
			ID:          evt.MessageID,
			DeviceID:    evt.DeviceID,
			FromJID:     evt.FromJID,
			ContentType: evt.ContentType,
			Content:     evt.Content,
			MessageID:   evt.MessageID,
			Timestamp:   evt.Timestamp,
		}
		if evt.IsGroup && evt.GroupJID != "" {
			rm.GroupJID = &evt.GroupJID
		}

		uID, err := sessionService.GetUserID(deviceID)
		if err == nil {
			rm.UserID = uID
		}

		if err := svc.ReceiveMessage(rm); err != nil {
			utils.Error("Failed to handle incoming message callback",
				zap.String("device_id", deviceID),
				zap.String("msg_id", evt.MessageID),
				zap.Error(err),
			)
		}
	})

	return svc
}

// UpdateConfig updates the message queue configuration dynamically
func (s *Service) UpdateConfig(config *MessageQueueConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = config
}

// GetConfig returns a copy of the current configuration
func (s *Service) GetConfig() MessageQueueConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.config
}

// Start begins processing the message queue
func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.processing {
		return fmt.Errorf("message service already running")
	}

	s.processing = true
	s.processorTicker = time.NewTicker(s.config.ProcessInterval)

	go s.processLoop()

	return nil
}

// Stop gracefully stops the message service
func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.processing {
		return fmt.Errorf("message service not running")
	}

	s.processing = false
	close(s.done)

	if s.processorTicker != nil {
		s.processorTicker.Stop()
	}

	return nil
}

func (s *Service) SendMessage(ctx context.Context, deviceID string, targetJID string, content string, groupID *string, broadcastID *string) (string, error) {
	return s.sendMessageInternal(ctx, deviceID, targetJID, content, groupID, broadcastID, false)
}

// SendInternalMessage sends a message bypassing the warming lockdown (used by Warming service)
func (s *Service) SendInternalMessage(ctx context.Context, deviceID string, targetJID string, content string) (string, error) {
	return s.sendMessageInternal(ctx, deviceID, targetJID, content, nil, nil, true)
}

func (s *Service) sendMessageInternal(ctx context.Context, deviceID string, targetJID string, content string, groupID *string, broadcastID *string, bypassLockdown bool) (string, error) {
	if !s.sessionService.IsSessionActive(deviceID) {
		return "", fmt.Errorf("session not active for device %s", deviceID)
	}

	// Check for warming lockdown if not bypassed
	if !bypassLockdown {
		parsedID, _ := uuid.Parse(deviceID)
		dev, err := s.db.GetDeviceByID(parsedID)
		if err == nil && dev != nil {
			if dev.IsWarming && dev.WarmingUntil != nil && time.Now().Before(*dev.WarmingUntil) {
				return "", fmt.Errorf("nomor sedang dalam mode Warming (Lockdown). Mohon tunggu hingga sesi selesai.")
			}
		}
	}

	messageID := uuid.New().String()

	queuedMsg := &QueuedMessage{
		ID:          messageID,
		DeviceID:    deviceID,
		TargetJID:   targetJID,
		GroupID:     groupID,
		Content:     content,
		ContentType: "text",
		BroadcastID: broadcastID,
		Status:      StatusPending,
		MaxRetries:  s.config.MaxRetries,
		Priority:    3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.store.EnqueueMessage(queuedMsg); err != nil {
		return "", fmt.Errorf("failed to queue message: %w", err)
	}

	s.metrics.recordMessageQueued()
	return messageID, nil
}

func (s *Service) SendMessageWithMedia(ctx context.Context, deviceID string, targetJID string, mediaURL string, contentType string, caption *string, broadcastID *string) (string, error) {
	if !s.sessionService.IsSessionActive(deviceID) {
		return "", fmt.Errorf("session not active for device %s", deviceID)
	}

	// Check for warming lockdown
	parsedID, _ := uuid.Parse(deviceID)
	dev, err := s.db.GetDeviceByID(parsedID)
	if err == nil && dev != nil {
		if dev.IsWarming && dev.WarmingUntil != nil && time.Now().Before(*dev.WarmingUntil) {
			return "", fmt.Errorf("nomor sedang dalam mode Warming (Lockdown). Mohon tunggu hingga sesi selesai.")
		}
	}

	messageID := uuid.New().String()

	queuedMsg := &QueuedMessage{
		ID:          messageID,
		DeviceID:    deviceID,
		TargetJID:   targetJID,
		MediaURL:    &mediaURL,
		Caption:     caption,
		ContentType: contentType,
		BroadcastID: broadcastID,
		Status:      StatusPending,
		MaxRetries:  s.config.MaxRetries,
		Priority:    2,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.store.EnqueueMessage(queuedMsg); err != nil {
		return "", fmt.Errorf("failed to queue media message: %w", err)
	}

	s.metrics.recordMessageQueued()
	return messageID, nil
}

func (s *Service) SendScheduledMessage(ctx context.Context, deviceID string, targetJID string, content string, scheduledFor time.Time, mediaURL *string, contentType string, caption *string, broadcastID *string) (string, error) {
	if !s.sessionService.IsSessionActive(deviceID) {
		return "", fmt.Errorf("session not active for device %s", deviceID)
	}

	// Check for warming lockdown
	parsedID, _ := uuid.Parse(deviceID)
	dev, err := s.db.GetDeviceByID(parsedID)
	if err == nil && dev != nil {
		if dev.IsWarming && dev.WarmingUntil != nil && time.Now().Before(*dev.WarmingUntil) {
			return "", fmt.Errorf("nomor sedang dalam mode Warming (Lockdown). Mohon tunggu hingga sesi selesai.")
		}
	}

	messageID := uuid.New().String()
	finalContentType := "text"
	if contentType != "" {
		finalContentType = contentType
	}

	queuedMsg := &QueuedMessage{
		ID:           messageID,
		DeviceID:     deviceID,
		TargetJID:    targetJID,
		Content:      content,
		ContentType:  finalContentType,
		MediaURL:     mediaURL,
		Caption:      caption,
		Status:       StatusPending,
		BroadcastID:  broadcastID,
		ScheduledFor: &scheduledFor,
		MaxRetries:   s.config.MaxRetries,
		Priority:     1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.store.EnqueueMessage(queuedMsg); err != nil {
		return "", fmt.Errorf("failed to queue scheduled message: %w", err)
	}

	return messageID, nil
}

func (s *Service) ReceiveMessage(rm *ReceivedMessage) error {
	if err := s.store.SaveReceivedMessage(rm); err != nil {
		return fmt.Errorf("failed to save received message: %w", err)
	}

	s.metrics.recordMessageReceived()

	for _, callback := range s.receiveCallbacks {
		go callback(rm)
	}

	return nil
}

func (s *Service) GetMessageStatus(messageID string) (MessageStatus, error) {
	return s.store.GetMessageStatus(messageID)
}

func (s *Service) UpdateMessageStatus(messageID string, status MessageStatus, errorMsg *string) error {
	oldStatus, err := s.store.GetMessageStatus(messageID)
	if err != nil {
		return err
	}

	if err := s.store.UpdateQueuedMessageStatus(messageID, status, errorMsg); err != nil {
		return err
	}

	update := &MessageStatusUpdate{
		MessageID:    messageID,
		OldStatus:    oldStatus,
		NewStatus:    status,
		Timestamp:    time.Now().Unix(),
		ErrorMessage: errorMsg,
	}

	for _, callback := range s.deliveryCallbacks {
		go callback(update)
	}

	s.metrics.recordStatusUpdate(status)
	return nil
}

func (s *Service) GetFailedMessages(limit int) ([]*QueuedMessage, error) {
	return s.store.GetFailedMessages(limit)
}

func (s *Service) ProcessQueue() error {
	s.mu.RLock()
	if !s.processing {
		s.mu.RUnlock()
		return fmt.Errorf("message service not running")
	}
	s.mu.RUnlock()

	sessions := s.sessionService.GetAllActiveSessions()
	if len(sessions) == 0 {
		return nil
	}

	for _, session := range sessions {
		go s.processDeviceQueue(session.ID)
	}

	return nil
}

func (s *Service) processDeviceQueue(deviceID string) {
	messages, err := s.store.DequeueMessages(deviceID, s.config.BatchSize)
	if err != nil {
		return
	}

	if len(messages) == 0 {
		return
	}

	for i, msg := range messages {
		if msg.ScheduledFor != nil && time.Now().Before(*msg.ScheduledFor) {
			continue
		}

		if s.config.SimulateTyping && msg.ContentType == "text" && len(msg.Content) > 0 && s.config.TypingSpeedCPM > 0 {
			typingSeconds := float64(len(msg.Content)) / float64(s.config.TypingSpeedCPM) * 60.0
			if typingSeconds > 10 {
				typingSeconds = 10
			}
			jitter := (rand.Float64()*0.4 - 0.2) * typingSeconds
			typingDuration := time.Duration((typingSeconds+jitter)*1000) * time.Millisecond
			time.Sleep(typingDuration)
		}

		s.sendQueuedMessageSync(deviceID, msg)

		if i < len(messages)-1 && s.config.MaxSendDelay > 0 {
			minMs := int64(s.config.MinSendDelay / time.Millisecond)
			maxMs := int64(s.config.MaxSendDelay / time.Millisecond)
			randMs := minMs + rand.Int63n(maxMs-minMs+1)
			
			// Personal Pro: Add additional jitter if anti-bot is enabled
			if s.config.AntiBotEnabled {
				jitter := rand.Int63n(maxMs / 4) // Additional ±25% noise
				randMs += jitter
			}

			delay := time.Duration(randMs) * time.Millisecond
			time.Sleep(delay)
		}
	}
}

func (s *Service) sendQueuedMessageSync(deviceID string, msg *QueuedMessage) {
	session := s.sessionService.GetSession(deviceID)
	if session == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := s.store.MarkMessageSent(msg.ID); err != nil {
		return
	}

	var waMessageID string
	var sendErr error

	switch msg.ContentType {
	case "image", "video", "audio", "document":
		if msg.MediaURL == nil || *msg.MediaURL == "" {
			sendErr = fmt.Errorf("media message missing media_url")
		} else {
			waMessageID, sendErr = s.sessionService.SendMessageWithMedia(
				ctx, deviceID, msg.TargetJID, *msg.MediaURL, msg.ContentType, msg.Caption,
			)
		}
	default:
		content := msg.Content
		// Personal Pro: Enhanced Anti-bot suffix
		if s.config.AntiBotEnabled && s.config.RandomSuffixLength > 0 && msg.ContentType == "text" {
			suffix := "\n\n[" + utils.GenerateRandomString(s.config.RandomSuffixLength) + "]"
			content += suffix
		}

		waMessageID, sendErr = s.sessionService.SendMessage(
			ctx, deviceID, msg.TargetJID, content, msg.GroupID,
		)
	}

	if sendErr != nil {
		s.handleSendFailure(msg, sendErr)
		return
	}

	if waMessageID != "" {
		_ = s.store.UpdateWhatsappMessageID(msg.ID, waMessageID)
	}

	if uID, err := s.sessionService.GetUserID(deviceID); err == nil {
		dID, _ := uuid.Parse(deviceID)
		_ = s.analyticsService.RecordSent(uID, dID)
	}

	s.metrics.recordMessageSent()
}

func (s *Service) handleSendFailure(msg *QueuedMessage, err error) {
	msg.RetryCount++
	msg.LastRetryAt = func() *time.Time { t := time.Now(); return &t }()

	if msg.RetryCount >= msg.MaxRetries {
		errMsg := fmt.Sprintf("Failed after %d retries: %v", msg.MaxRetries, err)
		_ = s.store.UpdateQueuedMessageStatus(msg.ID, StatusFailed, &errMsg)
		s.metrics.recordMessageFailed()

		if uID, errS := s.sessionService.GetUserID(msg.DeviceID); errS == nil {
			mID, _ := uuid.Parse(msg.ID)
			dID, _ := uuid.Parse(msg.DeviceID)
			_ = s.analyticsService.RecordFailure(uID, dID, mID, msg.TargetJID, "Send Error", err.Error())
		}
		return
	}

	errMsg := err.Error()
	_ = s.store.UpdateQueuedMessageStatus(msg.ID, StatusPending, &errMsg)
	_ = s.store.UpdateQueuedMessageRetry(msg.ID, msg.RetryCount, msg.LastRetryAt)
}

// Personal Pro Features: Bulk CSV Processing

func (s *Service) ProcessBulkCSV(ctx context.Context, deviceID string, r io.Reader) (int, error) {
	reader := csv.NewReader(r)
	headers, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	phoneIdx := -1
	messageIdx := -1
	for i, h := range headers {
		hLower := strings.ToLower(strings.TrimSpace(h))
		if hLower == "phone" || hLower == "nomor" || hLower == "penerima" {
			phoneIdx = i
		}
		if hLower == "message" || hLower == "pesan" || hLower == "content" {
			messageIdx = i
		}
	}

	if phoneIdx == -1 {
		return 0, fmt.Errorf("column 'phone' not found in CSV")
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		phone := strings.TrimSpace(record[phoneIdx])
		if phone == "" {
			continue
		}

		// Prepare message content with placeholder replacement
		var finalMessage string
		if messageIdx != -1 {
			finalMessage = record[messageIdx]
		}
		
		// Universal Placeholder Replacement: [ColumnName]
		for i, h := range headers {
			placeholder := "[" + h + "]"
			finalMessage = strings.ReplaceAll(finalMessage, placeholder, record[i])
		}

		if finalMessage == "" {
			continue
		}

		// Enqueue the message
		_, err = s.SendMessage(ctx, deviceID, phone, finalMessage, nil, nil)
		if err == nil {
			count++
		}
	}

	return count, nil
}

// Personal Pro Features: Export Logs CSV

func (s *Service) ExportLogsCSV(userID uuid.UUID) ([]byte, error) {
	logs, err := s.store.GetGlobalMessageLogs(userID.String(), 5000, 0)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write Headers
	_ = writer.Write([]string{"ID", "Device ID", "Recipient", "Content", "Status", "Error", "Created At"})

	for _, log := range logs {
		statusStr := "Pending"
		switch log.Status {
		case StatusSent:
			statusStr = "Sent"
		case StatusDelivered:
			statusStr = "Delivered"
		case StatusRead:
			statusStr = "Read"
		case StatusFailed:
			statusStr = "Failed"
		}

		errLog := ""
		if log.ErrorLog != nil {
			errLog = *log.ErrorLog
		}

		_ = writer.Write([]string{
			log.ID,
			log.DeviceID,
			log.TargetJID,
			log.Content,
			statusStr,
			errLog,
			log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	writer.Flush()
	return buf.Bytes(), nil
}

// Callbacks and Cleanup remains same

func (s *Service) RegisterDeliveryCallback(fn DeliveryCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deliveryCallbacks = append(s.deliveryCallbacks, fn)
}

func (s *Service) RegisterReceiveCallback(fn ReceiveCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.receiveCallbacks = append(s.receiveCallbacks, fn)
}

func (s *Service) GetQueueStats() map[string]interface{} {
	s.metrics.mu.RLock()
	defer s.metrics.mu.RUnlock()
	return map[string]interface{}{
		"total_sent":     s.metrics.TotalSent,
		"total_received": s.metrics.TotalReceived,
		"avg_latency_ms": s.metrics.AverageLatency,
	}
}

func (s *Service) Cleanup() error {
	cutoff := time.Now().AddDate(0, 0, -30)
	return s.store.DeleteOldMessages(cutoff)
}

func (s *Service) ListScheduledMessages(deviceID string) ([]*QueuedMessage, error) {
	return s.store.GetScheduledMessages(deviceID)
}

func (s *Service) ListMessageHistory(deviceID string, limit int) ([]*QueuedMessage, error) {
	return s.store.GetMessageHistory(deviceID, limit)
}

func (s *Service) ListGlobalMessageLogs(userID uuid.UUID, limit, offset int) ([]*QueuedMessage, error) {
	return s.store.GetGlobalMessageLogs(userID.String(), limit, offset)
}

func (s *Service) CancelScheduledMessage(messageID string) error {
	msg, err := s.store.GetQueuedMessage(messageID)
	if err != nil {
		return err
	}
	if msg.Status != StatusPending {
		return fmt.Errorf("cannot cancel message with status %v", msg.Status)
	}
	return s.store.DeleteQueuedMessage(messageID)
}

func (s *Service) processLoop() {
	for {
		select {
		case <-s.processorTicker.C:
			_ = s.ProcessQueue()
		case <-s.done:
			return
		}
	}
}

func (sm *ServiceMetrics) recordMessageQueued() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.CurrentPending++
}

func (sm *ServiceMetrics) recordMessageReceived() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.TotalReceived++
}

func (sm *ServiceMetrics) recordMessageFailed() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.TotalFailed++
	if sm.CurrentPending > 0 {
		sm.CurrentPending--
	}
}

func (sm *ServiceMetrics) recordMessageSent() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.TotalSent++
	if sm.CurrentPending > 0 {
		sm.CurrentPending--
	}
}

func (sm *ServiceMetrics) recordStatusUpdate(status MessageStatus) {
	// Optional trace
}
