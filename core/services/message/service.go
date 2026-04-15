package message

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"wacast/core/database"
	"wacast/core/services/analytics"
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
		MinSendDelay:       1 * time.Second, // Min gap between messages
		MaxSendDelay:       5 * time.Second, // Max gap between messages
		SimulateTyping:     true,            // Add typing delay for text
		TypingSpeedCPM:     300,             // ~300 chars/min (human speed)
	}
}

// NewService creates a new message service
func NewService(db *database.Database, sessionService *session.Service, analyticsService *analytics.Service, config *MessageQueueConfig) *Service {
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

	// Register receipt callback: whenever WA sends a delivery/read receipt,
	// look up the internal DB UUID by whatsapp_message_id and update status_message.
	sessionService.RegisterReceiptCallback(func(whatsappMsgID string, newStatus int) {
		internalID, err := store.GetDBIDByWhatsappID(whatsappMsgID)
		if err != nil {
			// Message may not be from this service (broadcast, etc.) — ignore silently
			return
		}
		if err := store.UpdateQueuedMessageStatus(internalID, MessageStatus(newStatus), nil); err != nil {
			utils.Error("Failed to update message status from receipt",
				zap.String("wa_msg_id", whatsappMsgID),
				zap.Int("status", newStatus),
				zap.Error(err),
			)
		} else {
			utils.Debug("Message status updated from receipt",
				zap.String("internal_id", internalID),
				zap.Int("status", newStatus),
			)
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


	return svc
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

	utils.Info("Message service started",
		zap.Duration("process_interval", s.config.ProcessInterval),
		zap.Int("max_retries", s.config.MaxRetries),
		zap.Int("batch_size", s.config.BatchSize),
	)

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

	utils.Info("Message service stopped")
	return nil
}

// SendMessage queues a text message for delivery
func (s *Service) SendMessage(ctx context.Context, deviceID string, targetJID string, content string, groupID *string, broadcastID *string) (string, error) {
	if !s.sessionService.IsSessionActive(deviceID) {
		return "", fmt.Errorf("session not active for device %s", deviceID)
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

	utils.Info("Message queued",
		zap.String("message_id", messageID),
		zap.String("device_id", deviceID),
		zap.String("target_jid", targetJID),
	)

	return messageID, nil
}

// SendMessageWithMedia queues a message with media attachment
func (s *Service) SendMessageWithMedia(ctx context.Context, deviceID string, targetJID string, mediaURL string, contentType string, caption *string, broadcastID *string) (string, error) {
	if !s.sessionService.IsSessionActive(deviceID) {
		return "", fmt.Errorf("session not active for device %s", deviceID)
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

	utils.Info("Media message queued",
		zap.String("message_id", messageID),
		zap.String("media_type", contentType),
	)

	return messageID, nil
}

// SendScheduledMessage queues a message to be sent at a specific time
func (s *Service) SendScheduledMessage(ctx context.Context, deviceID string, targetJID string, content string, scheduledFor time.Time, mediaURL *string, contentType string, caption *string, broadcastID *string) (string, error) {
	if !s.sessionService.IsSessionActive(deviceID) {
		return "", fmt.Errorf("session not active for device %s", deviceID)
	}

	messageID := uuid.New().String()

	// If it's a media message, ContentType should be the media type (image, video, etc)
	// Otherwise it defaults to "text"
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
		Priority:     1, // Lower priority for scheduled
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.store.EnqueueMessage(queuedMsg); err != nil {
		return "", fmt.Errorf("failed to queue scheduled message: %w", err)
	}

	utils.Info("Scheduled message queued",
		zap.String("message_id", messageID),
		zap.Time("scheduled_for", scheduledFor),
		zap.String("media_type", finalContentType),
	)

	return messageID, nil
}

// ReceiveMessage processes an incoming message
func (s *Service) ReceiveMessage(rm *ReceivedMessage) error {
	if err := s.store.SaveReceivedMessage(rm); err != nil {
		return fmt.Errorf("failed to save received message: %w", err)
	}

	s.metrics.recordMessageReceived()

	// Trigger receive callbacks
	for _, callback := range s.receiveCallbacks {
		go callback(rm)
	}

	return nil
}

// GetMessageStatus retrieves the status of a sent message
func (s *Service) GetMessageStatus(messageID string) (MessageStatus, error) {
	return s.store.GetMessageStatus(messageID)
}

// UpdateMessageStatus updates message status (called from session service events)
func (s *Service) UpdateMessageStatus(messageID string, status MessageStatus, errorMsg *string) error {
	oldStatus, err := s.store.GetMessageStatus(messageID)
	if err != nil {
		return err
	}

	if err := s.store.UpdateQueuedMessageStatus(messageID, status, errorMsg); err != nil {
		return err
	}

	// Trigger delivery callbacks
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

// GetFailedMessages retrieves messages that failed to send
func (s *Service) GetFailedMessages(limit int) ([]*QueuedMessage, error) {
	return s.store.GetFailedMessages(limit)
}

// ProcessQueue processes queued messages for all devices
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

// processDeviceQueue processes messages for a specific device.
// Messages are sent SEQUENTIALLY with random delays to mimic human behaviour.
func (s *Service) processDeviceQueue(deviceID string) {
	messages, err := s.store.DequeueMessages(deviceID, s.config.BatchSize)
	if err != nil {
		utils.Error("Failed to dequeue messages",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return
	}

	if len(messages) == 0 {
		return
	}

	for i, msg := range messages {
		// Skip scheduled messages that aren't ready
		if msg.ScheduledFor != nil && time.Now().Before(*msg.ScheduledFor) {
			continue
		}

		// 1. Anti-bot: simulate typing delay for text messages
		if s.config.SimulateTyping && msg.ContentType == "text" && len(msg.Content) > 0 && s.config.TypingSpeedCPM > 0 {
			typingSeconds := float64(len(msg.Content)) / float64(s.config.TypingSpeedCPM) * 60.0
			// Cap at 10 seconds maximum typing delay
			if typingSeconds > 10 {
				typingSeconds = 10
			}
			// Add ±20% jitter on typing time
			jitter := (rand.Float64()*0.4 - 0.2) * typingSeconds
			typingDuration := time.Duration((typingSeconds+jitter)*1000) * time.Millisecond
			utils.Debug("Simulating typing delay",
				zap.String("device_id", deviceID),
				zap.Duration("typing_delay", typingDuration),
			)
			time.Sleep(typingDuration)
		}

		// Send synchronously (sequential, no goroutine)
		s.sendQueuedMessageSync(deviceID, msg)

		// 2. Anti-bot: random delay between messages (skip after last message)
		if i < len(messages)-1 && s.config.MaxSendDelay > 0 {
			minMs := int64(s.config.MinSendDelay / time.Millisecond)
			maxMs := int64(s.config.MaxSendDelay / time.Millisecond)
			randMs := minMs + rand.Int63n(maxMs-minMs+1)
			delay := time.Duration(randMs) * time.Millisecond
			utils.Debug("Anti-bot delay between messages",
				zap.String("device_id", deviceID),
				zap.Duration("delay", delay),
			)
			time.Sleep(delay)
		}
	}
}

// sendQueuedMessageSync sends a single queued message synchronously.
func (s *Service) sendQueuedMessageSync(deviceID string, msg *QueuedMessage) {
	session := s.sessionService.GetSession(deviceID)
	if session == nil {
		utils.Warn("Session not found for message delivery",
			zap.String("device_id", deviceID),
			zap.String("message_id", msg.ID),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// ⚠️  CLAIM the message first: mark as 'sent' (status=1) in DB BEFORE sending.
	// This prevents the next queue tick from picking up the same message again
	// while the upload/send is still in progress (race condition).
	if err := s.store.MarkMessageSent(msg.ID); err != nil {
		utils.Error("Failed to claim message before send — skipping to avoid double-send",
			zap.String("message_id", msg.ID),
			zap.Error(err),
		)
		return
	}

	// Route to correct send method based on content type
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
	default: // "text" and anything else
		waMessageID, sendErr = s.sessionService.SendMessage(
			ctx, deviceID, msg.TargetJID, msg.Content, msg.GroupID,
		)
	}

	if sendErr != nil {
		// Revert status back to pending so it can be retried
		s.handleSendFailure(msg, sendErr)
		return
	}

	// Persist WA-assigned message ID (status stays at sent=1 already set above)
	if waMessageID != "" {
		if err := s.store.UpdateWhatsappMessageID(msg.ID, waMessageID); err != nil {
			utils.Error("Failed to save whatsapp_message_id",
				zap.String("message_id", msg.ID),
				zap.String("wa_message_id", waMessageID),
				zap.Error(err),
			)
		}
	}

	utils.Info("Message sent successfully",
		zap.String("message_id", msg.ID),
		zap.String("wa_message_id", waMessageID),
		zap.String("content_type", msg.ContentType),
		zap.String("target_jid", msg.TargetJID),
	)
	// Record success analytics
	if uID, err := s.sessionService.GetUserID(deviceID); err == nil {
		mID, _ := uuid.Parse(msg.ID)
		dID, _ := uuid.Parse(deviceID)
		_ = s.analyticsService.RecordSent(uID, dID)
		_ = mID 
	}
}


// handleSendFailure handles a failed message send attempt.
// If retries remain, reverts status back to pending so the queue picks it up again.
func (s *Service) handleSendFailure(msg *QueuedMessage, err error) {
	msg.RetryCount++
	msg.LastRetryAt = func() *time.Time { t := time.Now(); return &t }()

	if msg.RetryCount >= msg.MaxRetries {
		// Max retries exceeded, mark as failed permanently
		errMsg := fmt.Sprintf("Failed after %d retries: %v", msg.MaxRetries, err)
		if updateErr := s.store.UpdateQueuedMessageStatus(msg.ID, StatusFailed, &errMsg); updateErr != nil {
			utils.Error("Failed to update message status to failed",
				zap.String("message_id", msg.ID),
				zap.Error(updateErr),
			)
		}

		s.metrics.recordMessageFailed()

		utils.Warn("Message failed - max retries exceeded",
			zap.String("message_id", msg.ID),
			zap.Error(err),
		)

		// Record failure analytics
		if uID, errS := s.sessionService.GetUserID(msg.DeviceID); errS == nil {
			mID, _ := uuid.Parse(msg.ID)
			dID, _ := uuid.Parse(msg.DeviceID)
			_ = s.analyticsService.RecordFailure(uID, dID, mID, msg.TargetJID, "Send Error", err.Error())
		}

		return
	}

	// Revert status back to pending so queue processor picks it up on next tick
	errMsg := err.Error()
	if updateErr := s.store.UpdateQueuedMessageStatus(msg.ID, StatusPending, &errMsg); updateErr != nil {
		utils.Error("Failed to revert message status to pending",
			zap.String("message_id", msg.ID),
			zap.Error(updateErr),
		)
	}

	// Update retry count in DB
	if updateErr := s.store.UpdateQueuedMessageRetry(msg.ID, msg.RetryCount, msg.LastRetryAt); updateErr != nil {
		utils.Error("Failed to update retry count",
			zap.String("message_id", msg.ID),
			zap.Error(updateErr),
		)
	}

	utils.Warn("Message send failed, will retry",
		zap.String("message_id", msg.ID),
		zap.Int("retry_count", msg.RetryCount),
		zap.Int("max_retries", msg.MaxRetries),
		zap.Error(err),
	)
}

// RegisterDeliveryCallback registers a callback for delivery status updates
func (s *Service) RegisterDeliveryCallback(fn DeliveryCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deliveryCallbacks = append(s.deliveryCallbacks, fn)
}

// RegisterReceiveCallback registers a callback for incoming messages
func (s *Service) RegisterReceiveCallback(fn ReceiveCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.receiveCallbacks = append(s.receiveCallbacks, fn)
}

// GetQueueStats returns queue statistics
func (s *Service) GetQueueStats() map[string]interface{} {
	s.metrics.mu.RLock()
	defer s.metrics.mu.RUnlock()

	return map[string]interface{}{
		"total_sent":     s.metrics.TotalSent,
		"total_received": s.metrics.TotalReceived,
		"total_failed":   s.metrics.TotalFailed,
		"pending":        s.metrics.CurrentPending,
		"avg_latency_ms": s.metrics.AverageLatency,
	}
}

// Cleanup performs cleanup operations
func (s *Service) Cleanup() error {
	// Delete messages older than 30 days
	cutoff := time.Now().AddDate(0, 0, -30)
	if err := s.store.DeleteOldMessages(cutoff); err != nil {
		utils.Warn("Failed to cleanup old messages", zap.Error(err))
	}

	return nil
}

// ListScheduledMessages returns a list of pending scheduled messages for a device
func (s *Service) ListScheduledMessages(deviceID string) ([]*QueuedMessage, error) {
	return s.store.GetScheduledMessages(deviceID)
}

// ListMessageHistory returns a list of sent or failed messages for a device
func (s *Service) ListMessageHistory(deviceID string, limit int) ([]*QueuedMessage, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.store.GetMessageHistory(deviceID, limit)
}

// CancelScheduledMessage removes a pending scheduled message from the queue
func (s *Service) CancelScheduledMessage(messageID string) error {
	// We only allow deleting messages that are still pending
	msg, err := s.store.GetQueuedMessage(messageID)
	if err != nil {
		return err
	}

	if msg.Status != StatusPending {
		return fmt.Errorf("cannot cancel message with status %v", msg.Status)
	}

	return s.store.DeleteQueuedMessage(messageID)
}

// processLoop runs the message processing loop
func (s *Service) processLoop() {
	for {
		select {
		case <-s.processorTicker.C:
			if err := s.ProcessQueue(); err != nil {
				utils.Debug("Error processing queue", zap.Error(err))
			}
		case <-s.done:
			return
		}
	}
}

// Metrics helper functions

func (sm *ServiceMetrics) recordMessageQueued() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.CurrentPending++
}

func (sm *ServiceMetrics) recordMessageSent() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.TotalSent++
	if sm.CurrentPending > 0 {
		sm.CurrentPending--
	}
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

func (sm *ServiceMetrics) recordStatusUpdate(status MessageStatus) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	// Track status updates for metrics
}
