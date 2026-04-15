package message

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"wacast/core/database"
	"wacast/core/utils"
)

// DatabaseMessageStore implementation using PostgreSQL
type DatabaseMessageStore struct {
	db *database.Database
}

// NewDatabaseMessageStore creates a new message store
func NewDatabaseMessageStore(db *database.Database) *DatabaseMessageStore {
	return &DatabaseMessageStore{db: db}
}

// statusToInt maps MessageStatus to int stored in status_message column
// StatusPending=0, StatusSent=1, StatusDelivered=2, StatusRead=3, StatusFailed=4
func statusToInt(s MessageStatus) int { return int(s) }

// intToStatus maps int from status_message column to MessageStatus
func intToStatus(i int) MessageStatus { return MessageStatus(i) }

// contentTypeToInt maps content type string to int stored in message_type column
func contentTypeToInt(ct string) int {
	switch ct {
	case "image":
		return 1
	case "document":
		return 2
	case "audio":
		return 3
	case "video":
		return 4
	default:
		return 0 // text
	}
}

// intToContentType maps int message_type column back to string
func intToContentType(i int) string {
	switch i {
	case 1:
		return "image"
	case 2:
		return "document"
	case 3:
		return "audio"
	case 4:
		return "video"
	default:
		return "text"
	}
}

// EnqueueMessage inserts a pending outgoing message into the messages table.
func (dms *DatabaseMessageStore) EnqueueMessage(qm *QueuedMessage) error {
	query := `
		INSERT INTO messages (
			id, device_id, direction, receipt_number, message_type, content,
			status_message, error_log, priority, retry_count, max_retries,
			scheduled_for, broadcast_id, media_url, caption, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17
		)
	`
	now := time.Now()

	_, err := dms.db.Exec(query,
		qm.ID, qm.DeviceID, "OUT", qm.TargetJID, contentTypeToInt(qm.ContentType), qm.Content,
		statusToInt(qm.Status), qm.ErrorLog, qm.Priority, 0, qm.MaxRetries,
		qm.ScheduledFor, qm.BroadcastID, qm.MediaURL, qm.Caption, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to enqueue message: %w", err)
	}

	utils.Debug("Message queued to DB",
		zap.String("message_id", qm.ID),
		zap.String("device_id", qm.DeviceID),
		zap.String("target_jid", qm.TargetJID),
		zap.String("content_type", qm.ContentType),
	)
	return nil
}

// scanQueuedMessage scans a row into QueuedMessage using actual table columns
func scanQueuedMessage(scan func(...interface{}) error) (*QueuedMessage, error) {
	msg := &QueuedMessage{}
	var statusInt, msgTypeInt, priority, retryCount, maxRetries int
	var errLog *string

	err := scan(
		&msg.ID, &msg.DeviceID, &msg.TargetJID, &msg.Content,
		&msgTypeInt, &statusInt, &errLog,
		&priority, &retryCount, &maxRetries,
		&msg.ScheduledFor, &msg.BroadcastID, &msg.MediaURL, &msg.Caption, &msg.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	msg.Status = intToStatus(statusInt)
	msg.ErrorLog = errLog
	msg.Priority = priority
	msg.RetryCount = retryCount
	msg.MaxRetries = maxRetries
	msg.ContentType = intToContentType(msgTypeInt)
	return msg, nil
}

const queuedMessageSelect = `
	SELECT id, device_id, receipt_number, content, message_type, status_message,
	       error_log, priority, retry_count, max_retries, scheduled_for, broadcast_id, media_url, caption, created_at
	FROM messages
`

// DequeueMessages retrieves pending outgoing messages for a device
func (dms *DatabaseMessageStore) DequeueMessages(deviceID string, limit int) ([]*QueuedMessage, error) {
	query := queuedMessageSelect + `
		WHERE device_id = $1 AND status_message = $2 AND direction = 'OUT'
		  AND (scheduled_for IS NULL OR scheduled_for <= NOW())
		ORDER BY priority DESC, created_at ASC
		LIMIT $3
	`
	rows, err := dms.db.Query(query, deviceID, statusToInt(StatusPending), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*QueuedMessage
	for rows.Next() {
		msg, err := scanQueuedMessage(rows.Scan)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// GetQueuedMessage retrieves a specific queued outgoing message by ID
func (dms *DatabaseMessageStore) GetQueuedMessage(messageID string) (*QueuedMessage, error) {
	query := queuedMessageSelect + `WHERE id = $1 AND direction = 'OUT'`
	row := dms.db.QueryRow(query, messageID)
	msg, err := scanQueuedMessage(row.Scan)
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}
	return msg, nil
}

// UpdateQueuedMessageStatus updates a message's status and optionally the error log
func (dms *DatabaseMessageStore) UpdateQueuedMessageStatus(messageID string, status MessageStatus, errorMsg *string) error {
	query := `UPDATE messages SET status_message = $1, error_log = COALESCE($2, error_log), updated_at = $3 WHERE id = $4`
	_, err := dms.db.Exec(query, statusToInt(status), errorMsg, time.Now(), messageID)
	return err
}

// UpdateQueuedMessageRetry updates the retry count and last retry time
func (dms *DatabaseMessageStore) UpdateQueuedMessageRetry(messageID string, retryCount int, lastRetryAt *time.Time) error {
	query := `UPDATE messages SET retry_count = $1, updated_at = $2 WHERE id = $3`
	_, err := dms.db.Exec(query, retryCount, time.Now(), messageID)
	return err
}

// MarkMessageSent marks a message as sent
func (dms *DatabaseMessageStore) MarkMessageSent(messageID string) error {
	query := `UPDATE messages SET status_message = $1, updated_at = $2 WHERE id = $3`
	_, err := dms.db.Exec(query, statusToInt(StatusSent), time.Now(), messageID)
	return err
}

// GetFailedMessages retrieves failed outgoing messages
func (dms *DatabaseMessageStore) GetFailedMessages(limit int) ([]*QueuedMessage, error) {
	query := queuedMessageSelect + `
		WHERE status_message = $1 AND direction = 'OUT'
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := dms.db.Query(query, statusToInt(StatusFailed), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*QueuedMessage
	for rows.Next() {
		msg, err := scanQueuedMessage(rows.Scan)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// SaveReceivedMessage saves an incoming WhatsApp message to the database
func (dms *DatabaseMessageStore) SaveReceivedMessage(rm *ReceivedMessage) error {
	query := `
		INSERT INTO messages (
			id, device_id, direction, receipt_number, message_type, content,
			status_message, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	now := time.Now()
	ts := time.Unix(rm.Timestamp, 0)
	if rm.Timestamp == 0 {
		ts = now
	}

	_, err := dms.db.Exec(query,
		rm.ID, rm.DeviceID, "IN", rm.FromJID, 0, rm.Content,
		statusToInt(StatusDelivered), ts, now,
	)
	if err != nil {
		return fmt.Errorf("failed to save received message: %w", err)
	}

	utils.Debug("Incoming message saved to DB",
		zap.String("message_id", rm.ID),
		zap.String("from_jid", rm.FromJID),
	)
	return nil
}

// GetMessageByID retrieves an incoming message by ID
func (dms *DatabaseMessageStore) GetMessageByID(messageID string) (*ReceivedMessage, error) {
	query := `
		SELECT id, device_id, receipt_number, content, message_type, created_at
		FROM messages WHERE id = $1 AND direction = 'IN'
	`
	msg := &ReceivedMessage{}
	var cAt time.Time
	var msgTypeInt int

	err := dms.db.QueryRow(query, messageID).Scan(
		&msg.ID, &msg.DeviceID, &msg.FromJID, &msg.Content, &msgTypeInt, &cAt,
	)
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}
	msg.CreatedAt = cAt
	msg.Timestamp = cAt.Unix()
	msg.ContentType = "text"
	return msg, nil
}

// GetMessagesByDevice retrieves incoming messages for a device with pagination
func (dms *DatabaseMessageStore) GetMessagesByDevice(deviceID string, limit int, offset int) ([]*ReceivedMessage, error) {
	query := `
		SELECT id, device_id, receipt_number, content, message_type, created_at
		FROM messages WHERE device_id = $1 AND direction = 'IN'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	return dms.scanReceivedMessages(query, deviceID, limit, offset)
}

// GetMessagesByJID retrieves incoming messages by sender JID
func (dms *DatabaseMessageStore) GetMessagesByJID(jid string, limit int, offset int) ([]*ReceivedMessage, error) {
	query := `
		SELECT id, device_id, receipt_number, content, message_type, created_at
		FROM messages WHERE receipt_number = $1 AND direction = 'IN'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	return dms.scanReceivedMessages(query, jid, limit, offset)
}

func (dms *DatabaseMessageStore) scanReceivedMessages(query string, args ...interface{}) ([]*ReceivedMessage, error) {
	rows, err := dms.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*ReceivedMessage
	for rows.Next() {
		msg := &ReceivedMessage{}
		var cAt time.Time
		var msgTypeInt int
		if err := rows.Scan(&msg.ID, &msg.DeviceID, &msg.FromJID, &msg.Content, &msgTypeInt, &cAt); err != nil {
			return nil, err
		}
		msg.CreatedAt = cAt
		msg.Timestamp = cAt.Unix()
		msg.ContentType = "text"
		messages = append(messages, msg)
	}
	return messages, nil
}

// UpdateMessageStatus updates message status (alias for UpdateQueuedMessageStatus)
func (dms *DatabaseMessageStore) UpdateMessageStatus(messageID string, status MessageStatus) error {
	return dms.UpdateQueuedMessageStatus(messageID, status, nil)
}

// GetMessageStatus retrieves the current status of a message
func (dms *DatabaseMessageStore) GetMessageStatus(messageID string) (MessageStatus, error) {
	var statusInt int
	err := dms.db.QueryRow("SELECT status_message FROM messages WHERE id = $1", messageID).Scan(&statusInt)
	if err != nil {
		return StatusFailed, fmt.Errorf("message not found: %w", err)
	}
	return intToStatus(statusInt), nil
}

// GetPendingMessages retrieves all pending outgoing messages for a device
func (dms *DatabaseMessageStore) GetPendingMessages(deviceID string) ([]*QueuedMessage, error) {
	return dms.DequeueMessages(deviceID, 1000)
}

// CountByStatus counts messages by status for a device
func (dms *DatabaseMessageStore) CountByStatus(deviceID string, status MessageStatus) (int, error) {
	var count int
	err := dms.db.QueryRow(
		"SELECT COUNT(*) FROM messages WHERE device_id = $1 AND status_message = $2",
		deviceID, statusToInt(status),
	).Scan(&count)
	return count, err
}

// DeleteOldMessages deletes messages older than the given date
func (dms *DatabaseMessageStore) DeleteOldMessages(beforeDate time.Time) error {
	_, err := dms.db.Exec("DELETE FROM messages WHERE created_at < $1", beforeDate)
	return err
}

// ClearFailedMessages deletes failed messages older than the given date
func (dms *DatabaseMessageStore) ClearFailedMessages(beforeDate time.Time) error {
	_, err := dms.db.Exec(
		"DELETE FROM messages WHERE status_message = $1 AND created_at < $2",
		statusToInt(StatusFailed), beforeDate,
	)
	return err
}

// GetScheduledMessages retrieves pending scheduled messages for a device
func (dms *DatabaseMessageStore) GetScheduledMessages(deviceID string) ([]*QueuedMessage, error) {
	query := queuedMessageSelect + `
		WHERE device_id = $1 AND status_message = $2 AND direction = 'OUT'
		  AND scheduled_for > NOW()
		ORDER BY scheduled_for ASC
	`
	rows, err := dms.db.Query(query, deviceID, statusToInt(StatusPending))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*QueuedMessage
	for rows.Next() {
		msg, err := scanQueuedMessage(rows.Scan)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// GetMessageHistory retrieves the history of sent or failed messages for a device
func (dms *DatabaseMessageStore) GetMessageHistory(deviceID string, limit int) ([]*QueuedMessage, error) {
	query := queuedMessageSelect + `
		WHERE device_id = $1 AND direction = 'OUT' AND status_message != $2
		ORDER BY created_at DESC
		LIMIT $3
	`
	// status_message != 0 (not pending anymore)
	rows, err := dms.db.Query(query, deviceID, statusToInt(StatusPending), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*QueuedMessage
	for rows.Next() {
		msg, err := scanQueuedMessage(rows.Scan)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// DeleteQueuedMessage deletes a queued message
func (dms *DatabaseMessageStore) DeleteQueuedMessage(messageID string) error {
	_, err := dms.db.Exec("DELETE FROM messages WHERE id = $1", messageID)
	return err
}

// UpdateWhatsappMessageID stores the WhatsApp-assigned message ID against our internal DB UUID.
// Called right after SendMessage succeeds so we can match receipts later.
func (dms *DatabaseMessageStore) UpdateWhatsappMessageID(internalID, whatsappMsgID string) error {
	query := `UPDATE messages SET whatsapp_message_id = $1, status_message = $2, updated_at = $3 WHERE id = $4`
	_, err := dms.db.Exec(query, whatsappMsgID, statusToInt(StatusSent), time.Now(), internalID)
	return err
}

// GetDBIDByWhatsappID returns the internal UUID of a message given a WhatsApp message ID.
// Used by the receipt callback to update status_message when a receipt arrives.
func (dms *DatabaseMessageStore) GetDBIDByWhatsappID(whatsappMsgID string) (string, error) {
	var id string
	err := dms.db.QueryRow(
		"SELECT id FROM messages WHERE whatsapp_message_id = $1 LIMIT 1",
		whatsappMsgID,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("message with whatsapp_id %q not found: %w", whatsappMsgID, err)
	}
	return id, nil
}
