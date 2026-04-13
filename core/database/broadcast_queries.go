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
		INSERT INTO broadcast_campaigns (id, user_id, device_id, template_id, name,
			message_content, total_recipients, success_count, failed_count, delay_seconds,
			scheduled_at, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	now := time.Now()
	_, err := d.Exec(query,
		campaign.ID, campaign.UserID, campaign.DeviceID, campaign.TemplateID, campaign.Name,
		campaign.MessageContent, campaign.TotalRecipients, campaign.SuccessCount, campaign.FailedCount,
		campaign.DelaySeconds, campaign.ScheduledAt, campaign.Status, now, now,
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
		SELECT id, user_id, device_id, template_id, name, message_content,
			total_recipients, success_count, failed_count, delay_seconds,
			scheduled_at, started_at, completed_at, status, created_at, updated_at, deleted_at
		FROM broadcast_campaigns
		WHERE id = $1 AND deleted_at IS NULL
	`

	campaign := &models.BroadcastCampaign{}
	err := d.QueryRow(query, campaignID).Scan(
		&campaign.ID, &campaign.UserID, &campaign.DeviceID, &campaign.TemplateID, &campaign.Name,
		&campaign.MessageContent, &campaign.TotalRecipients, &campaign.SuccessCount, &campaign.FailedCount,
		&campaign.DelaySeconds, &campaign.ScheduledAt, &campaign.StartedAt, &campaign.CompletedAt,
		&campaign.Status, &campaign.CreatedAt, &campaign.UpdatedAt, &campaign.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return campaign, nil
}

// GetBroadcastCampaignsByUserID retrieves campaigns for a user
func (d *Database) GetBroadcastCampaignsByUserID(userID uuid.UUID, limit, offset int) ([]models.BroadcastCampaign, error) {
	query := `
		SELECT id, user_id, device_id, template_id, name, message_content,
			total_recipients, success_count, failed_count, delay_seconds,
			scheduled_at, started_at, completed_at, status, created_at, updated_at, deleted_at
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
			&campaign.ID, &campaign.UserID, &campaign.DeviceID, &campaign.TemplateID, &campaign.Name,
			&campaign.MessageContent, &campaign.TotalRecipients, &campaign.SuccessCount, &campaign.FailedCount,
			&campaign.DelaySeconds, &campaign.ScheduledAt, &campaign.StartedAt, &campaign.CompletedAt,
			&campaign.Status, &campaign.CreatedAt, &campaign.UpdatedAt, &campaign.DeletedAt,
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
func (d *Database) UpdateBroadcastCampaignStatus(campaignID uuid.UUID, status string) error {
	query := `UPDATE broadcast_campaigns SET status = $1, updated_at = $2 WHERE id = $3`

	_, err := d.Exec(query, status, time.Now(), campaignID)
	if err != nil {
		utils.Error("Failed to update campaign status", zap.Error(err))
		return err
	}

	return nil
}

// UpdateBroadcastCampaignProgress updates success/failed counts
func (d *Database) UpdateBroadcastCampaignProgress(campaignID uuid.UUID, successCount, failedCount int) error {
	query := `UPDATE broadcast_campaigns SET success_count = $1, failed_count = $2, updated_at = $3 WHERE id = $4`

	_, err := d.Exec(query, successCount, failedCount, time.Now(), campaignID)
	if err != nil {
		utils.Error("Failed to update campaign progress", zap.Error(err))
		return err
	}

	return nil
}

// CreateBroadcastRecipient creates a broadcast recipient
func (d *Database) CreateBroadcastRecipient(recipient *models.BroadcastRecipient) error {
	query := `
		INSERT INTO broadcast_recipients (id, campaign_id, group_id, contact_id,
			phone_number, status, retry_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := d.Exec(query,
		recipient.ID, recipient.CampaignID, recipient.GroupID, recipient.ContactID,
		recipient.PhoneNumber, recipient.Status, recipient.RetryCount, time.Now(),
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
		SELECT id, campaign_id, group_id, contact_id, phone_number, status,
			sent_at, failed_at, error_message, retry_count, created_at
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
			&recipient.ID, &recipient.CampaignID, &recipient.GroupID, &recipient.ContactID,
			&recipient.PhoneNumber, &recipient.Status, &recipient.SentAt, &recipient.FailedAt,
			&recipient.ErrorMessage, &recipient.RetryCount, &recipient.CreatedAt,
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
func (d *Database) UpdateBroadcastRecipientStatus(recipientID uuid.UUID, status string, errMsg *string) error {
	var sentAt *time.Time
	var failedAt *time.Time
	now := time.Now()

	if status == models.BroadcastStatusCompleted {
		sentAt = &now
	} else if status == models.BroadcastStatusFailed {
		failedAt = &now
	}

	query := `
		UPDATE broadcast_recipients
		SET status = $1, sent_at = $2, failed_at = $3, error_message = $4
		WHERE id = $5
	`

	_, err := d.Exec(query, status, sentAt, failedAt, errMsg, recipientID)
	if err != nil {
		utils.Error("Failed to update recipient status", zap.Error(err))
		return err
	}

	return nil
}

// GetPendingBroadcastRecipients retrieves pending recipients
func (d *Database) GetPendingBroadcastRecipients(campaignID uuid.UUID, limit int) ([]models.BroadcastRecipient, error) {
	query := `
		SELECT id, campaign_id, group_id, contact_id, phone_number, status,
			sent_at, failed_at, error_message, retry_count, created_at
		FROM broadcast_recipients
		WHERE campaign_id = $1 AND status = 'pending'
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
			&recipient.ID, &recipient.CampaignID, &recipient.GroupID, &recipient.ContactID,
			&recipient.PhoneNumber, &recipient.Status, &recipient.SentAt, &recipient.FailedAt,
			&recipient.ErrorMessage, &recipient.RetryCount, &recipient.CreatedAt,
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
func (d *Database) CountBroadcastRecipientsByStatus(campaignID uuid.UUID, status string) (int64, error) {
	query := `SELECT COUNT(*) FROM broadcast_recipients WHERE campaign_id = $1 AND status = $2`

	var count int64
	err := d.QueryRow(query, campaignID, status).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
