package database

import (
	"time"

	"github.com/google/uuid"
	"wacast/core/models"
	"wacast/core/utils"
	"go.uber.org/zap"
)

// CreateMessage logs a message
func (d *Database) CreateMessage(message *models.Message) error {
	query := `
		INSERT INTO messages (id, device_id, direction, receipt_number, message_type, 
			content, status_message, error_log, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := d.Exec(query,
		message.ID, message.DeviceID, message.Direction, message.ReceiptNumber,
		message.MessageType, message.Content, message.StatusMessage, message.ErrorLog,
		time.Now(),
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
		SELECT id, device_id, direction, receipt_number, message_type, content, 
			status_message, error_log, created_at
		FROM messages
		WHERE id = $1
	`

	message := &models.Message{}
	err := d.QueryRow(query, messageID).Scan(
		&message.ID, &message.DeviceID, &message.Direction, &message.ReceiptNumber,
		&message.MessageType, &message.Content, &message.StatusMessage, &message.ErrorLog,
		&message.CreatedAt,
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
		SELECT id, device_id, direction, receipt_number, message_type, content, 
			status_message, error_log, created_at
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
			&message.ID, &message.DeviceID, &message.Direction, &message.ReceiptNumber,
			&message.MessageType, &message.Content, &message.StatusMessage, &message.ErrorLog,
			&message.CreatedAt,
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
func (d *Database) UpdateMessageStatus(messageID uuid.UUID, status int32) error {
	query := `UPDATE messages SET status_message = $1 WHERE id = $2`

	_, err := d.Exec(query, status, messageID)
	if err != nil {
		utils.Error("Failed to update message status", zap.Error(err))
		return err
	}

	return nil
}

// UpdateMessageStatusWithError updates message status and error log
func (d *Database) UpdateMessageStatusWithError(messageID uuid.UUID, status int32, errorLog string) error {
	query := `UPDATE messages SET status_message = $1, error_log = $2 WHERE id = $3`

	_, err := d.Exec(query, status, errorLog, messageID)
	if err != nil {
		utils.Error("Failed to update message status", zap.Error(err))
		return err
	}

	return nil
}

// GetPendingMessages retrieves all pending messages
func (d *Database) GetPendingMessages() ([]models.Message, error) {
	query := `
		SELECT id, device_id, direction, receipt_number, message_type, content, 
			status_message, error_log, created_at
		FROM messages
		WHERE status_message = 0
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
			&message.ID, &message.DeviceID, &message.Direction, &message.ReceiptNumber,
			&message.MessageType, &message.Content, &message.StatusMessage, &message.ErrorLog,
			&message.CreatedAt,
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
func (d *Database) GetMessagesByStatusAndDevice(deviceID uuid.UUID, status int32, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT id, device_id, direction, receipt_number, message_type, content, 
			status_message, error_log, created_at
		FROM messages
		WHERE device_id = $1 AND status_message = $2
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
			&message.ID, &message.DeviceID, &message.Direction, &message.ReceiptNumber,
			&message.MessageType, &message.Content, &message.StatusMessage, &message.ErrorLog,
			&message.CreatedAt,
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
		SELECT id, device_id, direction, receipt_number, message_type, content, 
			status_message, error_log, created_at
		FROM messages
		WHERE device_id = $1 AND direction = 'IN'
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
			&message.ID, &message.DeviceID, &message.Direction, &message.ReceiptNumber,
			&message.MessageType, &message.Content, &message.StatusMessage, &message.ErrorLog,
			&message.CreatedAt,
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
		SELECT id, device_id, direction, receipt_number, message_type, content, 
			status_message, error_log, created_at
		FROM messages
		WHERE device_id = $1 AND direction = 'OUT'
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
			&message.ID, &message.DeviceID, &message.Direction, &message.ReceiptNumber,
			&message.MessageType, &message.Content, &message.StatusMessage, &message.ErrorLog,
			&message.CreatedAt,
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
func (d *Database) CountMessagesByStatus(deviceID uuid.UUID, status int32) (int64, error) {
	query := `SELECT COUNT(*) FROM messages WHERE device_id = $1 AND status_message = $2`

	var count int64
	err := d.QueryRow(query, deviceID, status).Scan(&count)
	if err != nil {
		utils.Error("Failed to count messages", zap.Error(err))
		return 0, err
	}

	return count, nil
}

// GetMessageByReceiptNumber retrieves a message by receipt number
func (d *Database) GetMessageByReceiptNumber(receiptNumber string) (*models.Message, error) {
	query := `
		SELECT id, device_id, direction, receipt_number, message_type, content, 
			status_message, error_log, created_at
		FROM messages
		WHERE receipt_number = $1
	`

	message := &models.Message{}
	err := d.QueryRow(query, receiptNumber).Scan(
		&message.ID, &message.DeviceID, &message.Direction, &message.ReceiptNumber,
		&message.MessageType, &message.Content, &message.StatusMessage, &message.ErrorLog,
		&message.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return message, nil
}
