package database

import (
	"fmt"
	"time"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateAutoResponseKeyword creates an auto response keyword rule
func (d *Database) CreateAutoResponseKeyword(autoResp *models.AutoResponseKeyword) error {
	query := `
		INSERT INTO auto_response_keywords (id, user_id, device_id, keyword, match_type, response_text, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	_, err := d.Exec(query,
		autoResp.ID, autoResp.UserID, autoResp.DeviceID, autoResp.Keyword,
		autoResp.MatchType, autoResp.ResponseText, autoResp.IsActive, now, now,
	)

	if err != nil {
		utils.Error("Failed to create auto response keyword", zap.Error(err))
		return err
	}

	return nil
}

// GetAutoResponseKeywordByID retrieves an auto response keyword by ID
func (d *Database) GetAutoResponseKeywordByID(respID uuid.UUID) (*models.AutoResponseKeyword, error) {
	query := `
		SELECT id, user_id, device_id, keyword, match_type, response_text, is_active, created_at, updated_at
		FROM auto_response_keywords
		WHERE id = $1
	`

	resp := &models.AutoResponseKeyword{}
	err := d.QueryRow(query, respID).Scan(
		&resp.ID, &resp.UserID, &resp.DeviceID, &resp.Keyword,
		&resp.MatchType, &resp.ResponseText, &resp.IsActive, &resp.CreatedAt, &resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetAutoResponseKeywordsByDeviceID retrieves active auto response keywords for a device
func (d *Database) GetAutoResponseKeywordsByDeviceID(deviceID uuid.UUID) ([]models.AutoResponseKeyword, error) {
	query := `
		SELECT id, user_id, device_id, keyword, match_type, response_text, is_active, created_at, updated_at
		FROM auto_response_keywords
		WHERE device_id = $1 AND is_active = true
	`

	rows, err := d.Query(query, deviceID)
	if err != nil {
		utils.Error("Failed to get auto response keywords", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	responses := []models.AutoResponseKeyword{}
	for rows.Next() {
		resp := models.AutoResponseKeyword{}
		err := rows.Scan(
			&resp.ID, &resp.UserID, &resp.DeviceID, &resp.Keyword,
			&resp.MatchType, &resp.ResponseText, &resp.IsActive, &resp.CreatedAt, &resp.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan response", zap.Error(err))
			continue
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// UpdateAutoResponseKeyword updates an auto response keyword
func (d *Database) UpdateAutoResponseKeyword(respID uuid.UUID, update *models.UpdateAutoResponseKeywordRequest) error {
	query := `UPDATE auto_response_keywords SET `
	args := []interface{}{}
	argCount := 1

	if update.Keyword != nil {
		query += fmt.Sprintf("keyword = $%d", argCount)
		args = append(args, *update.Keyword)
		argCount++
	}

	if update.MatchType != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("match_type = $%d", argCount)
		args = append(args, *update.MatchType)
		argCount++
	}

	if update.ResponseText != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("response_text = $%d", argCount)
		args = append(args, *update.ResponseText)
		argCount++
	}

	if update.IsActive != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("is_active = $%d", argCount)
		args = append(args, *update.IsActive)
		argCount++
	}

	if argCount > 1 {
		query += ", "
	}
	query += fmt.Sprintf("updated_at = $%d WHERE id = $%d", argCount, argCount+1)
	args = append(args, time.Now(), respID)

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update auto response keyword", zap.Error(err))
		return err
	}

	return nil
}

// DeleteAutoResponseKeyword deletes an auto response keyword
func (d *Database) DeleteAutoResponseKeyword(respID uuid.UUID) error {
	query := `DELETE FROM auto_response_keywords WHERE id = $1`

	_, err := d.Exec(query, respID)
	if err != nil {
		utils.Error("Failed to delete auto response keyword", zap.Error(err))
		return err
	}

	return nil
}

// CreateWebhook creates a webhook configuration
func (d *Database) CreateWebhook(webhook *models.Webhook) error {
	query := `
		INSERT INTO webhooks (id, user_id, device_id, webhook_url, secret_key_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	_, err := d.Exec(query,
		webhook.ID, webhook.UserID, webhook.DeviceID, webhook.WebhookUrl,
		webhook.SecretKeyHash, webhook.IsActive, now, now,
	)

	if err != nil {
		utils.Error("Failed to create webhook", zap.Error(err))
		return err
	}

	return nil
}

// GetWebhookByID retrieves a webhook by ID
func (d *Database) GetWebhookByID(webhookID uuid.UUID) (*models.Webhook, error) {
	query := `
		SELECT id, user_id, device_id, webhook_url, secret_key_hash, is_active, created_at, updated_at
		FROM webhooks
		WHERE id = $1
	`

	webhook := &models.Webhook{}
	err := d.QueryRow(query, webhookID).Scan(
		&webhook.ID, &webhook.UserID, &webhook.DeviceID, &webhook.WebhookUrl,
		&webhook.SecretKeyHash, &webhook.IsActive, &webhook.CreatedAt, &webhook.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return webhook, nil
}

// GetWebhooksByDeviceID retrieves webhooks for a device
func (d *Database) GetWebhooksByDeviceID(deviceID uuid.UUID) ([]models.Webhook, error) {
	query := `
		SELECT id, user_id, device_id, webhook_url, secret_key_hash, is_active, created_at, updated_at
		FROM webhooks
		WHERE device_id = $1
	`

	rows, err := d.Query(query, deviceID)
	if err != nil {
		utils.Error("Failed to get webhooks", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	webhooks := []models.Webhook{}
	for rows.Next() {
		webhook := models.Webhook{}
		err := rows.Scan(
			&webhook.ID, &webhook.UserID, &webhook.DeviceID, &webhook.WebhookUrl,
			&webhook.SecretKeyHash, &webhook.IsActive, &webhook.CreatedAt, &webhook.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan webhook", zap.Error(err))
			continue
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// UpdateWebhook updates a webhook
func (d *Database) UpdateWebhook(webhookID uuid.UUID, webhookURL, secretKeyHash string) error {
	query := `UPDATE webhooks SET webhook_url = $1, secret_key_hash = $2, updated_at = $3 WHERE id = $4`

	_, err := d.Exec(query, webhookURL, secretKeyHash, time.Now(), webhookID)
	if err != nil {
		utils.Error("Failed to update webhook", zap.Error(err))
		return err
	}

	return nil
}

// DeleteWebhook deletes a webhook
func (d *Database) DeleteWebhook(webhookID uuid.UUID) error {
	query := `DELETE FROM webhooks WHERE id = $1`

	_, err := d.Exec(query, webhookID)
	if err != nil {
		utils.Error("Failed to delete webhook", zap.Error(err))
		return err
	}

	return nil
}

// CreateAPILog logs an API call
func (d *Database) CreateAPILog(log *models.APILog) error {
	query := `
		INSERT INTO api_logs (id, user_id, device_id, endpoint, method, status_code, req_body, response_body, ip_address, user_agent, duration_ms, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := d.Exec(query,
		log.ID, log.UserID, log.DeviceID, log.Endpoint, log.Method, log.StatusCode,
		log.ReqBody, log.ResponseBody, log.IPAddress, log.UserAgent, log.DurationMs, time.Now(),
	)

	if err != nil {
		utils.Error("Failed to create API log", zap.Error(err))
		return err
	}

	return nil
}

// GetAPILogsByUserID retrieves API logs for a user
func (d *Database) GetAPILogsByUserID(userID uuid.UUID, limit, offset int) ([]models.APILog, error) {
	query := `
		SELECT id, user_id, device_id, endpoint, method, status_code, req_body, response_body,
		       ip_address, user_agent, duration_ms, created_at
		FROM api_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, userID, limit, offset)
	if err != nil {
		utils.Error("Failed to get API logs", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	logs := []models.APILog{}
	for rows.Next() {
		log := models.APILog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.DeviceID, &log.Endpoint, &log.Method, &log.StatusCode,
			&log.ReqBody, &log.ResponseBody, &log.IPAddress, &log.UserAgent, &log.DurationMs, &log.CreatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan log", zap.Error(err))
			continue
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetSystemSetting retrieves a system setting
func (d *Database) GetSystemSetting(key string) (*models.SystemSetting, error) {
	query := `SELECT id, key, value, description, created_at FROM system_settings WHERE key = $1`

	setting := &models.SystemSetting{}
	err := d.QueryRow(query, key).Scan(
		&setting.ID, &setting.Key, &setting.Value, &setting.Description, &setting.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return setting, nil
}

// UpdateSystemSetting updates a system setting
func (d *Database) UpdateSystemSetting(key, value string) error {
	query := `UPDATE system_settings SET value = $1, updated_at = $2 WHERE key = $3`

	_, err := d.Exec(query, value, time.Now(), key)
	if err != nil {
		utils.Error("Failed to update system setting", zap.Error(err))
		return err
	}

	return nil
}


