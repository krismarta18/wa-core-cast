package database

import (
	"time"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateAutoResponse creates an auto response rule
func (d *Database) CreateAutoResponse(autoResp *models.AutoResponse) error {
	query := `
		INSERT INTO auto_response (id, device_id, keyword, response_text, is_active)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := d.Exec(query,
		autoResp.ID, autoResp.DeviceID, autoResp.Keyword, autoResp.ResponseText, autoResp.IsActive,
	)

	if err != nil {
		utils.Error("Failed to create auto response", zap.Error(err))
		return err
	}

	return nil
}

// GetAutoResponseByID retrieves an auto response by ID
func (d *Database) GetAutoResponseByID(respID uuid.UUID) (*models.AutoResponse, error) {
	query := `
		SELECT id, device_id, keyword, response_text, is_active
		FROM auto_response
		WHERE id = $1
	`

	resp := &models.AutoResponse{}
	err := d.QueryRow(query, respID).Scan(
		&resp.ID, &resp.DeviceID, &resp.Keyword, &resp.ResponseText, &resp.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetAutoResponsesByDeviceID retrieves auto responses for a device
func (d *Database) GetAutoResponsesByDeviceID(deviceID uuid.UUID) ([]models.AutoResponse, error) {
	query := `
		SELECT id, device_id, keyword, response_text, is_active
		FROM auto_response
		WHERE device_id = $1 AND is_active = true
	`

	rows, err := d.Query(query, deviceID)
	if err != nil {
		utils.Error("Failed to get auto responses", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	responses := []models.AutoResponse{}
	for rows.Next() {
		resp := models.AutoResponse{}
		err := rows.Scan(
			&resp.ID, &resp.DeviceID, &resp.Keyword, &resp.ResponseText, &resp.IsActive,
		)
		if err != nil {
			utils.Error("Failed to scan response", zap.Error(err))
			continue
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// UpdateAutoResponse updates an auto response
func (d *Database) UpdateAutoResponse(respID uuid.UUID, update *models.UpdateAutoResponseRequest) error {
	query := `UPDATE auto_response SET `
	args := []interface{}{}
	argCount := 1

	if update.Keyword != nil {
		query += `keyword = $1`
		args = append(args, *update.Keyword)
		argCount++
	}

	if update.ResponseText != nil {
		if argCount > 1 {
			query += ", "
		}
		query += `response_text = $` + string(rune(argCount))
		args = append(args, *update.ResponseText)
		argCount++
	}

	if update.IsActive != nil {
		if argCount > 1 {
			query += ", "
		}
		query += `is_active = $` + string(rune(argCount))
		args = append(args, *update.IsActive)
		argCount++
	}

	query += ` WHERE id = $` + string(rune(argCount))
	args = append(args, respID)

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update auto response", zap.Error(err))
		return err
	}

	return nil
}

// DeleteAutoResponse deletes an auto response
func (d *Database) DeleteAutoResponse(respID uuid.UUID) error {
	query := `DELETE FROM auto_response WHERE id = $1`

	_, err := d.Exec(query, respID)
	if err != nil {
		utils.Error("Failed to delete auto response", zap.Error(err))
		return err
	}

	return nil
}

// CreateWebhook creates a webhook configuration
func (d *Database) CreateWebhook(webhook *models.Webhook) error {
	query := `
		INSERT INTO webhooks (id, device_id, webhook_url, secret_key)
		VALUES ($1, $2, $3, $4)
	`

	_, err := d.Exec(query,
		webhook.ID, webhook.DeviceID, webhook.WebhookUrl, webhook.SecretKey,
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
		SELECT id, device_id, webhook_url, secret_key
		FROM webhooks
		WHERE id = $1
	`

	webhook := &models.Webhook{}
	err := d.QueryRow(query, webhookID).Scan(
		&webhook.ID, &webhook.DeviceID, &webhook.WebhookUrl, &webhook.SecretKey,
	)

	if err != nil {
		return nil, err
	}

	return webhook, nil
}

// GetWebhooksByDeviceID retrieves webhooks for a device
func (d *Database) GetWebhooksByDeviceID(deviceID uuid.UUID) ([]models.Webhook, error) {
	query := `
		SELECT id, device_id, webhook_url, secret_key
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
			&webhook.ID, &webhook.DeviceID, &webhook.WebhookUrl, &webhook.SecretKey,
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
func (d *Database) UpdateWebhook(webhookID uuid.UUID, webhookURL, secretKey string) error {
	query := `UPDATE webhooks SET webhook_url = $1, secret_key = $2 WHERE id = $3`

	_, err := d.Exec(query, webhookURL, secretKey, webhookID)
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
		INSERT INTO api_logs (id, user_id, endpoint, req_body, response_body, created_at, ip_address, device_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := d.Exec(query,
		log.ID, log.UserID, log.Endpoint, log.ReqBody, log.ResponseBody, time.Now(), log.IPAddress, log.DeviceID,
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
		SELECT id, user_id, endpoint, req_body, response_body, created_at, ip_address, device_id
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
			&log.ID, &log.UserID, &log.Endpoint, &log.ReqBody, &log.ResponseBody, &log.CreatedAt, &log.IPAddress, &log.DeviceID,
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
	query := `SELECT id, keys, value, description, created_at FROM system_settings WHERE keys = $1`

	setting := &models.SystemSetting{}
	err := d.QueryRow(query, key).Scan(
		&setting.ID, &setting.Keys, &setting.Value, &setting.Description, &setting.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return setting, nil
}

// UpdateSystemSetting updates a system setting
func (d *Database) UpdateSystemSetting(key, value string) error {
	query := `UPDATE system_settings SET value = $1 WHERE keys = $2`

	_, err := d.Exec(query, value, key)
	if err != nil {
		utils.Error("Failed to update system setting", zap.Error(err))
		return err
	}

	return nil
}
