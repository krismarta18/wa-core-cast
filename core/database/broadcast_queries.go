package database

import (
	"time"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateBroadcastCampaign creates a new broadcast campaign
func (d *Database) CreateBroadcastCampaign(campaign *models.BroadcastCampaign) error {
	query := `
		INSERT INTO broadcast_campaigns (id, user_id, device_id, name_broadcast, 
			total_recipients, processed_count, scheduled_at, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := d.Exec(query,
		campaign.ID, campaign.UserID, campaign.DeviceID, campaign.NameBroadcast,
		campaign.TotalRecipients, campaign.ProcessedCount, campaign.ScheduledAt,
		campaign.Status, time.Now(),
	)

	if err != nil {
		utils.Error("Failed to create broadcast campaign", zap.Error(err))
		return err
	}

	return nil
}

// GetBroadcastCampaignByID retrieves a broadcast campaign
func (d *Database) GetBroadcastCampaignByID(campaignID uuid.UUID) (*models.BroadcastCampaign, error) {
	query := `
		SELECT id, user_id, device_id, name_broadcast, total_recipients, processed_count,
			scheduled_at, status, created_at, deleted_at
		FROM broadcast_campaigns
		WHERE id = $1 AND deleted_at IS NULL
	`

	campaign := &models.BroadcastCampaign{}
	err := d.QueryRow(query, campaignID).Scan(
		&campaign.ID, &campaign.UserID, &campaign.DeviceID, &campaign.NameBroadcast,
		&campaign.TotalRecipients, &campaign.ProcessedCount, &campaign.ScheduledAt,
		&campaign.Status, &campaign.CreatedAt, &campaign.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return campaign, nil
}

// GetBroadcastCampaignsByUserID retrieves campaigns for a user
func (d *Database) GetBroadcastCampaignsByUserID(userID uuid.UUID, limit, offset int) ([]models.BroadcastCampaign, error) {
	query := `
		SELECT id, user_id, device_id, name_broadcast, total_recipients, processed_count,
			scheduled_at, status, created_at, deleted_at
		FROM broadcast_campaigns
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, userID, limit, offset)
	if err != nil {
		utils.Error("Failed to get broadcast campaigns", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	campaigns := []models.BroadcastCampaign{}
	for rows.Next() {
		campaign := models.BroadcastCampaign{}
		err := rows.Scan(
			&campaign.ID, &campaign.UserID, &campaign.DeviceID, &campaign.NameBroadcast,
			&campaign.TotalRecipients, &campaign.ProcessedCount, &campaign.ScheduledAt,
			&campaign.Status, &campaign.CreatedAt, &campaign.DeletedAt,
		)
		if err != nil {
			utils.Error("Failed to scan campaign", zap.Error(err))
			continue
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

// UpdateBroadcastCampaignStatus updates campaign status
func (d *Database) UpdateBroadcastCampaignStatus(campaignID uuid.UUID, status int32) error {
	query := `UPDATE broadcast_campaigns SET status = $1 WHERE id = $2`

	_, err := d.Exec(query, status, campaignID)
	if err != nil {
		utils.Error("Failed to update campaign status", zap.Error(err))
		return err
	}

	return nil
}

// UpdateBroadcastCampaignProgress updates processed count
func (d *Database) UpdateBroadcastCampaignProgress(campaignID uuid.UUID, processedCount int32) error {
	query := `UPDATE broadcast_campaigns SET processed_count = $1 WHERE id = $2`

	_, err := d.Exec(query, processedCount, campaignID)
	if err != nil {
		utils.Error("Failed to update campaign progress", zap.Error(err))
		return err
	}

	return nil
}

// CreateBroadcastMessage creates a broadcast message
func (d *Database) CreateBroadcastMessage(message *models.BroadcastMessage) error {
	query := `
		INSERT INTO broadcast_messages (id, campaign_id, message_type, message_text, 
			media_url, button_data)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := d.Exec(query,
		message.ID, message.CampaignID, message.MessageType, message.MessageText,
		message.MediaUrl, message.ButtonData,
	)

	if err != nil {
		utils.Error("Failed to create broadcast message", zap.Error(err))
		return err
	}

	return nil
}

// GetBroadcastMessagesByCampaignID retrieves messages for a campaign
func (d *Database) GetBroadcastMessagesByCampaignID(campaignID uuid.UUID) ([]models.BroadcastMessage, error) {
	query := `
		SELECT id, campaign_id, message_type, message_text, media_url, button_data
		FROM broadcast_messages
		WHERE campaign_id = $1
	`

	rows, err := d.Query(query, campaignID)
	if err != nil {
		utils.Error("Failed to get broadcast messages", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	messages := []models.BroadcastMessage{}
	for rows.Next() {
		message := models.BroadcastMessage{}
		err := rows.Scan(
			&message.ID, &message.CampaignID, &message.MessageType, &message.MessageText,
			&message.MediaUrl, &message.ButtonData,
		)
		if err != nil {
			utils.Error("Failed to scan message", zap.Error(err))
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// CreateBroadcastRecipient creates a broadcast recipient
func (d *Database) CreateBroadcastRecipient(recipient *models.BroadcastRecipient) error {
	query := `
		INSERT INTO broadcast_recipients (id, campaign_id, groups_id, contact_id, 
			status, retry_count)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := d.Exec(query,
		recipient.ID, recipient.CampaignID, recipient.GroupsID, recipient.ContactID,
		recipient.Status, recipient.RetryCount,
	)

	if err != nil {
		utils.Error("Failed to create broadcast recipient", zap.Error(err))
		return err
	}

	return nil
}

// GetBroadcastRecipientsByCampaignID retrieves recipients for a campaign
func (d *Database) GetBroadcastRecipientsByCampaignID(campaignID uuid.UUID, limit, offset int) ([]models.BroadcastRecipient, error) {
	query := `
		SELECT id, campaign_id, groups_id, contact_id, status, sent_at, error_messages, retry_count
		FROM broadcast_recipients
		WHERE campaign_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, campaignID, limit, offset)
	if err != nil {
		utils.Error("Failed to get broadcast recipients", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	recipients := []models.BroadcastRecipient{}
	for rows.Next() {
		recipient := models.BroadcastRecipient{}
		err := rows.Scan(
			&recipient.ID, &recipient.CampaignID, &recipient.GroupsID, &recipient.ContactID,
			&recipient.Status, &recipient.SentAt, &recipient.ErrorMessages, &recipient.RetryCount,
		)
		if err != nil {
			utils.Error("Failed to scan recipient", zap.Error(err))
			continue
		}
		recipients = append(recipients, recipient)
	}

	return recipients, nil
}

// UpdateBroadcastRecipientStatus updates recipient status
func (d *Database) UpdateBroadcastRecipientStatus(recipientID uuid.UUID, status int32, errorMsg string) error {
	query := `UPDATE broadcast_recipients SET status = $1, error_messages = $2, sent_at = $3 WHERE id = $4`

	_, err := d.Exec(query, status, errorMsg, time.Now(), recipientID)
	if err != nil {
		utils.Error("Failed to update recipient status", zap.Error(err))
		return err
	}

	return nil
}

// GetPendingBroadcastRecipients retrieves pending recipients
func (d *Database) GetPendingBroadcastRecipients(campaignID uuid.UUID, limit int) ([]models.BroadcastRecipient, error) {
	query := `
		SELECT id, campaign_id, groups_id, contact_id, status, sent_at, error_messages, retry_count
		FROM broadcast_recipients
		WHERE campaign_id = $1 AND status = 0
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := d.Query(query, campaignID, limit)
	if err != nil {
		utils.Error("Failed to get pending recipients", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	recipients := []models.BroadcastRecipient{}
	for rows.Next() {
		recipient := models.BroadcastRecipient{}
		err := rows.Scan(
			&recipient.ID, &recipient.CampaignID, &recipient.GroupsID, &recipient.ContactID,
			&recipient.Status, &recipient.SentAt, &recipient.ErrorMessages, &recipient.RetryCount,
		)
		if err != nil {
			utils.Error("Failed to scan recipient", zap.Error(err))
			continue
		}
		recipients = append(recipients, recipient)
	}

	return recipients, nil
}

// DeleteBroadcastCampaign soft deletes a campaign
func (d *Database) DeleteBroadcastCampaign(campaignID uuid.UUID) error {
	query := `UPDATE broadcast_campaigns SET deleted_at = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), campaignID)
	if err != nil {
		utils.Error("Failed to delete campaign", zap.Error(err))
		return err
	}

	return nil
}

// CountBroadcastRecipientsByStatus counts recipients by status
func (d *Database) CountBroadcastRecipientsByStatus(campaignID uuid.UUID, status int32) (int64, error) {
	query := `SELECT COUNT(*) FROM broadcast_recipients WHERE campaign_id = $1 AND status = $2`

	var count int64
	err := d.QueryRow(query, campaignID, status).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
