package database

import (
	"time"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const messageSelectColumns = `
	id, user_id, device_id, template_id, broadcast_id, scheduled_message_id,
	direction, recipient_phone, sender_phone, message_type, content, media_url,
	status, whatsapp_message_id, error_log, sent_at, delivered_at, read_at,
	failed_at, created_at, updated_at
`

// CreateMessage logs a message
func (d *Database) CreateMessage(message *models.Message) error {
	query := `
		INSERT INTO messages (
			id, user_id, device_id, template_id, broadcast_id, scheduled_message_id,
			direction, recipient_phone, sender_phone, message_type, content, media_url,
			status, whatsapp_message_id, error_log, sent_at, delivered_at, read_at,
			failed_at, created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18,
			$19, $20, $21
		)
	`

	now := time.Now()
	_, err := d.Exec(query,
		message.ID, message.UserID, message.DeviceID, message.TemplateID, message.BroadcastID, message.ScheduledMessageID,
		message.Direction, message.RecipientPhone, message.SenderPhone, message.MessageType, message.Content, message.MediaURL,
		message.Status, message.WhatsappMessageID, message.ErrorLog, message.SentAt, message.DeliveredAt, message.ReadAt,
		message.FailedAt, now, now,
	)
	if err != nil {
		utils.Error("Failed to create message", zap.Error(err))
		return err
	}

	return nil
}

// GetMessageByID retrieves a message by ID
func (d *Database) GetMessageByID(messageID uuid.UUID) (*models.Message, error) {
	query := `
		SELECT ` + messageSelectColumns + `
		FROM messages
		WHERE id = $1
	`

	message := &models.Message{}
	err := d.QueryRow(query, messageID).Scan(
		&message.ID, &message.UserID, &message.DeviceID, &message.TemplateID, &message.BroadcastID, &message.ScheduledMessageID,
		&message.Direction, &message.RecipientPhone, &message.SenderPhone, &message.MessageType, &message.Content, &message.MediaURL,
		&message.Status, &message.WhatsappMessageID, &message.ErrorLog, &message.SentAt, &message.DeliveredAt, &message.ReadAt,
		&message.FailedAt, &message.CreatedAt, &message.UpdatedAt,
	)
	if err != nil {
		utils.Debug("Message not found", zap.String("message_id", messageID.String()))
		return nil, err
	}

	return message, nil
}

// GetMessagesByDeviceID retrieves messages for a device with pagination
func (d *Database) GetMessagesByDeviceID(deviceID uuid.UUID, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT ` + messageSelectColumns + `
		FROM messages
		WHERE device_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, deviceID, limit, offset)
	if err != nil {
		utils.Error("Failed to get messages", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		message := models.Message{}
		err := rows.Scan(
			&message.ID, &message.UserID, &message.DeviceID, &message.TemplateID, &message.BroadcastID, &message.ScheduledMessageID,
			&message.Direction, &message.RecipientPhone, &message.SenderPhone, &message.MessageType, &message.Content, &message.MediaURL,
			&message.Status, &message.WhatsappMessageID, &message.ErrorLog, &message.SentAt, &message.DeliveredAt, &message.ReadAt,
			&message.FailedAt, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan message", zap.Error(err))
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// UpdateMessageStatus updates message status
func (d *Database) UpdateMessageStatus(messageID uuid.UUID, status string) error {
	now := time.Now()
	query := `
		UPDATE messages
		SET status = $1,
			sent_at = CASE WHEN $1 = 'sent' THEN $2 ELSE sent_at END,
			delivered_at = CASE WHEN $1 = 'delivered' THEN $2 ELSE delivered_at END,
			read_at = CASE WHEN $1 = 'read' THEN $2 ELSE read_at END,
			failed_at = CASE WHEN $1 = 'failed' THEN $2 ELSE failed_at END,
			updated_at = $2
		WHERE id = $3
	`

	_, err := d.Exec(query, status, now, messageID)
	if err != nil {
		utils.Error("Failed to update message status", zap.Error(err))
		return err
	}

	return nil
}

// UpdateMessageStatusWithError updates message status and error log
func (d *Database) UpdateMessageStatusWithError(messageID uuid.UUID, status string, errorLog string) error {
	now := time.Now()
	query := `
		UPDATE messages
		SET status = $1,
			error_log = $2,
			failed_at = CASE WHEN $1 = 'failed' THEN $3 ELSE failed_at END,
			updated_at = $3
		WHERE id = $4
	`

	_, err := d.Exec(query, status, errorLog, now, messageID)
	if err != nil {
		utils.Error("Failed to update message status", zap.Error(err))
		return err
	}

	return nil
}

// GetPendingMessages retrieves all pending messages
func (d *Database) GetPendingMessages() ([]models.Message, error) {
	query := `
		SELECT ` + messageSelectColumns + `
		FROM messages
		WHERE status = 'pending'
		ORDER BY created_at ASC
	`

	rows, err := d.Query(query)
	if err != nil {
		utils.Error("Failed to get pending messages", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		message := models.Message{}
		err := rows.Scan(
			&message.ID, &message.UserID, &message.DeviceID, &message.TemplateID, &message.BroadcastID, &message.ScheduledMessageID,
			&message.Direction, &message.RecipientPhone, &message.SenderPhone, &message.MessageType, &message.Content, &message.MediaURL,
			&message.Status, &message.WhatsappMessageID, &message.ErrorLog, &message.SentAt, &message.DeliveredAt, &message.ReadAt,
			&message.FailedAt, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan message", zap.Error(err))
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetMessagesByStatusAndDevice retrieves messages by status and device
func (d *Database) GetMessagesByStatusAndDevice(deviceID uuid.UUID, status string, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT ` + messageSelectColumns + `
		FROM messages
		WHERE device_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := d.Query(query, deviceID, status, limit, offset)
	if err != nil {
		utils.Error("Failed to get messages by status", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		message := models.Message{}
		err := rows.Scan(
			&message.ID, &message.UserID, &message.DeviceID, &message.TemplateID, &message.BroadcastID, &message.ScheduledMessageID,
			&message.Direction, &message.RecipientPhone, &message.SenderPhone, &message.MessageType, &message.Content, &message.MediaURL,
			&message.Status, &message.WhatsappMessageID, &message.ErrorLog, &message.SentAt, &message.DeliveredAt, &message.ReadAt,
			&message.FailedAt, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan message", zap.Error(err))
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetIncomingMessages retrieves incoming messages for a device
func (d *Database) GetIncomingMessages(deviceID uuid.UUID, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT ` + messageSelectColumns + `
		FROM messages
		WHERE device_id = $1 AND direction = 'inbound'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, deviceID, limit, offset)
	if err != nil {
		utils.Error("Failed to get incoming messages", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		message := models.Message{}
		err := rows.Scan(
			&message.ID, &message.UserID, &message.DeviceID, &message.TemplateID, &message.BroadcastID, &message.ScheduledMessageID,
			&message.Direction, &message.RecipientPhone, &message.SenderPhone, &message.MessageType, &message.Content, &message.MediaURL,
			&message.Status, &message.WhatsappMessageID, &message.ErrorLog, &message.SentAt, &message.DeliveredAt, &message.ReadAt,
			&message.FailedAt, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan message", zap.Error(err))
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetOutgoingMessages retrieves outgoing messages for a device
func (d *Database) GetOutgoingMessages(deviceID uuid.UUID, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT ` + messageSelectColumns + `
		FROM messages
		WHERE device_id = $1 AND direction = 'outbound'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, deviceID, limit, offset)
	if err != nil {
		utils.Error("Failed to get outgoing messages", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		message := models.Message{}
		err := rows.Scan(
			&message.ID, &message.UserID, &message.DeviceID, &message.TemplateID, &message.BroadcastID, &message.ScheduledMessageID,
			&message.Direction, &message.RecipientPhone, &message.SenderPhone, &message.MessageType, &message.Content, &message.MediaURL,
			&message.Status, &message.WhatsappMessageID, &message.ErrorLog, &message.SentAt, &message.DeliveredAt, &message.ReadAt,
			&message.FailedAt, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan message", zap.Error(err))
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// CountMessagesByStatus counts messages by status for a device
func (d *Database) CountMessagesByStatus(deviceID uuid.UUID, status string) (int64, error) {
	query := `SELECT COUNT(*) FROM messages WHERE device_id = $1 AND status = $2`

	var count int64
	err := d.QueryRow(query, deviceID, status).Scan(&count)
	if err != nil {
		utils.Error("Failed to count messages", zap.Error(err))
		return 0, err
	}

	return count, nil
}

// GetMessageByReceiptNumber retrieves a message by WhatsApp message ID.
func (d *Database) GetMessageByReceiptNumber(receiptNumber string) (*models.Message, error) {
	query := `
		SELECT ` + messageSelectColumns + `
		FROM messages
		WHERE whatsapp_message_id = $1
	`

	message := &models.Message{}
	err := d.QueryRow(query, receiptNumber).Scan(
		&message.ID, &message.UserID, &message.DeviceID, &message.TemplateID, &message.BroadcastID, &message.ScheduledMessageID,
		&message.Direction, &message.RecipientPhone, &message.SenderPhone, &message.MessageType, &message.Content, &message.MediaURL,
		&message.Status, &message.WhatsappMessageID, &message.ErrorLog, &message.SentAt, &message.DeliveredAt, &message.ReadAt,
		&message.FailedAt, &message.CreatedAt, &message.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return message, nil
}
// GetMessageCountToday counts outbound messages sent by a user today
func (d *Database) GetMessageCountToday(userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages
		WHERE user_id = $1
		  AND direction = 'outbound'
		  AND created_at >= CURRENT_DATE
	`

	var count int
	err := d.QueryRow(query, userID).Scan(&count)
	if err != nil {
		utils.Error("Failed to count messages today", zap.Error(err), zap.String("user_id", userID.String()))
		return 0, err
	}

	return count, nil
}
