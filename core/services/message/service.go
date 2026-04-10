package message

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"wacast/core/database"
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

// DefaultQueueConfig returns default configuration
func DefaultQueueConfig() *MessageQueueConfig {
	return &MessageQueueConfig{
		MaxRetries:         3,
		RetryDelayBase:     5 * time.Second,
		MaxRetryDelay:      5 * time.Minute,
		BatchSize:          50,
		ProcessInterval:    2 * time.Second,
		MaxConcurrentSends: 5,
	}
}

// NewService creates a new message service
func NewService(db *database.Database, sessionService *session.Service, config *MessageQueueConfig) *Service {
	if config == nil {
		config = DefaultQueueConfig()
	}

	store := NewDatabaseMessageStore(db)

	svc := &Service{
		db:              db,
		store:           store,
		config:          config,
		sessionService:  sessionService,
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
func (s *Service) SendMessage(ctx context.Context, deviceID string, targetJID string, content string, groupID *string) (string, error) {
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
func (s *Service) SendMessageWithMedia(ctx context.Context, deviceID string, targetJID string, mediaURL string, contentType string, caption *string) (string, error) {
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
func (s *Service) SendScheduledMessage(ctx context.Context, deviceID string, targetJID string, content string, scheduledFor time.Time) (string, error) {
	if !s.sessionService.IsSessionActive(deviceID) {
		return "", fmt.Errorf("session not active for device %s", deviceID)
	}

	messageID := uuid.New().String()

	queuedMsg := &QueuedMessage{
		ID:          messageID,
		DeviceID:    deviceID,
		TargetJID:   targetJID,
		Content:     content,
		ContentType: "text",
		Status:      StatusPending,
		ScheduledFor: &scheduledFor,
		MaxRetries:  s.config.MaxRetries,
		Priority:    1, // Lower priority for scheduled
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.store.EnqueueMessage(queuedMsg); err != nil {
		return "", fmt.Errorf("failed to queue scheduled message: %w", err)
	}

	utils.Info("Scheduled message queued",
		zap.String("message_id", messageID),
		zap.Time("scheduled_for", scheduledFor),
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

// processDeviceQueue processes messages for a specific device
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

	semaphore := make(chan struct{}, s.config.MaxConcurrentSends)

	for _, msg := range messages {
		// Skip scheduled messages that aren't ready
		if msg.ScheduledFor != nil && time.Now().Before(*msg.ScheduledFor) {
			continue
		}

		semaphore <- struct{}{}
		go s.sendQueuedMessage(deviceID, msg, semaphore)
	}
}

// sendQueuedMessage attempts to send a queued message
func (s *Service) sendQueuedMessage(deviceID string, msg *QueuedMessage, semaphore chan struct{}) {
	defer func() { <-semaphore }()

	session := s.sessionService.GetSession(deviceID)
	if session == nil {
		utils.Warn("Session not found for message delivery",
			zap.String("device_id", deviceID),
			zap.String("message_id", msg.ID),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Send the message via session service
	_, err := s.sessionService.SendMessage(ctx, deviceID, msg.TargetJID, msg.Content, msg.GroupID)
	if err != nil {
		s.handleSendFailure(msg, err)
		return
	}

	// Mark as sent if no error
	if err := s.store.MarkMessageSent(msg.ID); err != nil {
		utils.Error("Failed to mark message as sent",
			zap.String("message_id", msg.ID),
			zap.Error(err),
		)
	}

	s.metrics.recordMessageSent()

	utils.Debug("Message sent successfully",
		zap.String("message_id", msg.ID),
		zap.String("target_jid", msg.TargetJID),
	)
}

// handleSendFailure handles a failed message send attempt
func (s *Service) handleSendFailure(msg *QueuedMessage, err error) {
	msg.RetryCount++
	msg.LastRetryAt = func() *time.Time { t := time.Now(); return &t }()

	if msg.RetryCount >= msg.MaxRetries {
		// Max retries exceeded, mark as failed
		errMsg := fmt.Sprintf("Failed after %d retries: %v", msg.MaxRetries, err)
		if updateErr := s.store.UpdateQueuedMessageStatus(msg.ID, StatusFailed, &errMsg); updateErr != nil {
			utils.Error("Failed to update message status",
				zap.String("message_id", msg.ID),
				zap.Error(updateErr),
			)
		}

		s.metrics.recordMessageFailed()

		utils.Warn("Message failed - max retries exceeded",
			zap.String("message_id", msg.ID),
			zap.Error(err),
		)

		return
	}

	// Update retry count
	if updateErr := s.store.UpdateQueuedMessageRetry(msg.ID, msg.RetryCount, msg.LastRetryAt); updateErr != nil {
		utils.Error("Failed to update retry count",
			zap.String("message_id", msg.ID),
			zap.Error(updateErr),
		)
	}

	utils.Debug("Message send failed, will retry",
		zap.String("message_id", msg.ID),
		zap.Int("retry_count", msg.RetryCount),
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
