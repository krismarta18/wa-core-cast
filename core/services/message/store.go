package message

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"wacast/core/database"
	"wacast/core/utils"
)

// MessageStore implementation using in-memory storage
type DatabaseMessageStore struct {
	db       *database.Database
	mu       sync.RWMutex
	messages map[string]*QueuedMessage
	received map[string]*ReceivedMessage
}

// NewDatabaseMessageStore creates a new message store
func NewDatabaseMessageStore(db *database.Database) *DatabaseMessageStore {
	return &DatabaseMessageStore{
		db:       db,
		messages: make(map[string]*QueuedMessage),
		received: make(map[string]*ReceivedMessage),
	}
}

// EnqueueMessage adds a message to the queue
func (dms *DatabaseMessageStore) EnqueueMessage(qm *QueuedMessage) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	dms.messages[qm.ID] = qm

	utils.Debug("Message queued",
		zap.String("message_id", qm.ID),
		zap.String("device_id", qm.DeviceID),
		zap.String("target_jid", qm.TargetJID),
	)

	return nil
}

// DequeueMessages retrieves messages ready to be sent
func (dms *DatabaseMessageStore) DequeueMessages(deviceID string, limit int) ([]*QueuedMessage, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	var messages []*QueuedMessage
	for _, msg := range dms.messages {
		if msg.DeviceID == deviceID && msg.Status == StatusPending {
			messages = append(messages, msg)
			if len(messages) >= limit {
				break
			}
		}
	}

	return messages, nil
}

// GetQueuedMessage retrieves a specific queued message
func (dms *DatabaseMessageStore) GetQueuedMessage(messageID string) (*QueuedMessage, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	msg, exists := dms.messages[messageID]
	if !exists {
		return nil, fmt.Errorf("message not found")
	}

	return msg, nil
}

// UpdateQueuedMessageStatus updates a queued message status
func (dms *DatabaseMessageStore) UpdateQueuedMessageStatus(messageID string, status MessageStatus, errorMsg *string) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	msg, exists := dms.messages[messageID]
	if !exists {
		return fmt.Errorf("message not found")
	}

	msg.Status = status
	msg.ErrorLog = errorMsg
	msg.UpdatedAt = time.Now()

	return nil
}

// UpdateQueuedMessageRetry updates retry count
func (dms *DatabaseMessageStore) UpdateQueuedMessageRetry(messageID string, retryCount int, lastRetryAt *time.Time) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	msg, exists := dms.messages[messageID]
	if !exists {
		return fmt.Errorf("message not found")
	}

	msg.RetryCount = retryCount
	msg.LastRetryAt = lastRetryAt

	return nil
}

// MarkMessageSent marks a message as sent
func (dms *DatabaseMessageStore) MarkMessageSent(messageID string) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	msg, exists := dms.messages[messageID]
	if !exists {
		return fmt.Errorf("message not found")
	}

	msg.Status = StatusSent
	msg.UpdatedAt = time.Now()

	return nil
}

// GetFailedMessages retrieves failed messages
func (dms *DatabaseMessageStore) GetFailedMessages(limit int) ([]*QueuedMessage, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	var messages []*QueuedMessage
	for _, msg := range dms.messages {
		if msg.Status == StatusFailed {
			messages = append(messages, msg)
			if len(messages) >= limit {
				break
			}
		}
	}

	return messages, nil
}

// SaveReceivedMessage saves an incoming message
func (dms *DatabaseMessageStore) SaveReceivedMessage(rm *ReceivedMessage) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	dms.received[rm.ID] = rm

	utils.Debug("Message received and saved",
		zap.String("message_id", rm.ID),
		zap.String("from_jid", rm.FromJID),
	)

	return nil
}

// GetMessageByID retrieves a received message
func (dms *DatabaseMessageStore) GetMessageByID(messageID string) (*ReceivedMessage, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	msg, exists := dms.received[messageID]
	if !exists {
		return nil, fmt.Errorf("message not found")
	}

	return msg, nil
}

// GetMessagesByDevice retrieves messages for a device
func (dms *DatabaseMessageStore) GetMessagesByDevice(deviceID string, limit int, offset int) ([]*ReceivedMessage, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	var messages []*ReceivedMessage
	for _, msg := range dms.received {
		if msg.DeviceID == deviceID {
			messages = append(messages, msg)
		}
	}

	// Apply pagination
	if offset >= len(messages) {
		return []*ReceivedMessage{}, nil
	}

	end := offset + limit
	if end > len(messages) {
		end = len(messages)
	}

	return messages[offset:end], nil
}

// GetMessagesByJID retrieves messages by JID
func (dms *DatabaseMessageStore) GetMessagesByJID(jid string, limit int, offset int) ([]*ReceivedMessage, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	var messages []*ReceivedMessage
	for _, msg := range dms.received {
		if msg.FromJID == jid {
			messages = append(messages, msg)
		}
	}

	// Apply pagination
	if offset >= len(messages) {
		return []*ReceivedMessage{}, nil
	}

	end := offset + limit
	if end > len(messages) {
		end = len(messages)
	}

	return messages[offset:end], nil
}

// UpdateMessageStatus updates message status
func (dms *DatabaseMessageStore) UpdateMessageStatus(messageID string, status MessageStatus) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	msg, exists := dms.messages[messageID]
	if !exists {
		return fmt.Errorf("message not found")
	}

	msg.Status = status
	msg.UpdatedAt = time.Now()

	return nil
}

// GetMessageStatus gets message status
func (dms *DatabaseMessageStore) GetMessageStatus(messageID string) (MessageStatus, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	msg, exists := dms.messages[messageID]
	if !exists {
		return StatusFailed, fmt.Errorf("message not found")
	}

	return msg.Status, nil
}

// GetPendingMessages gets pending messages
func (dms *DatabaseMessageStore) GetPendingMessages(deviceID string) ([]*QueuedMessage, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	var messages []*QueuedMessage
	for _, msg := range dms.messages {
		if msg.DeviceID == deviceID && msg.Status == StatusPending {
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

// CountByStatus counts messages by status
func (dms *DatabaseMessageStore) CountByStatus(deviceID string, status MessageStatus) (int, error) {
	dms.mu.RLock()
	defer dms.mu.RUnlock()

	count := 0
	for _, msg := range dms.messages {
		if msg.DeviceID == deviceID && msg.Status == status {
			count++
		}
	}

	return count, nil
}

// DeleteOldMessages deletes old messages
func (dms *DatabaseMessageStore) DeleteOldMessages(beforeDate time.Time) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	for id, msg := range dms.messages {
		if msg.CreatedAt.Before(beforeDate) {
			delete(dms.messages, id)
		}
	}

	return nil
}

// ClearFailedMessages clears failed messages
func (dms *DatabaseMessageStore) ClearFailedMessages(beforeDate time.Time) error {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	for id, msg := range dms.messages {
		if msg.Status == StatusFailed && msg.CreatedAt.Before(beforeDate) {
			delete(dms.messages, id)
		}
	}

	return nil
}
