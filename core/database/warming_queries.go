package database

import (
	"fmt"
	"time"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateWarmingPool creates a warming pool entry
func (d *Database) CreateWarmingPool(pool *models.WarmingPool) error {
	query := `
		INSERT INTO warming_pool (id, device_id, intensity, daily_limit, message_send_today, 
			is_active, next_action_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := d.Exec(query,
		pool.ID, pool.DeviceID, pool.Intensity, pool.DailyLimit, pool.MessageSendToday,
		pool.IsActive, pool.NextActionAt,
	)

	if err != nil {
		utils.Error("Failed to create warming pool", zap.Error(err))
		return err
	}

	return nil
}

// GetWarmingPoolByDeviceID retrieves warming pool for a device
func (d *Database) GetWarmingPoolByDeviceID(deviceID uuid.UUID) (*models.WarmingPool, error) {
	query := `
		SELECT id, device_id, intensity, daily_limit, message_send_today, is_active, next_action_at
		FROM warming_pool
		WHERE device_id = $1
	`

	pool := &models.WarmingPool{}
	err := d.QueryRow(query, deviceID).Scan(
		&pool.ID, &pool.DeviceID, &pool.Intensity, &pool.DailyLimit, &pool.MessageSendToday,
		&pool.IsActive, &pool.NextActionAt,
	)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

// UpdateWarmingPool updates warming pool settings
func (d *Database) UpdateWarmingPool(poolID uuid.UUID, update *models.UpdateWarmingPoolRequest) error {
	query := `UPDATE warming_pool SET `
	args := []interface{}{}
	argCount := 1

	if update.Intensity != nil {
		query += fmt.Sprintf(`intensity = $%d`, argCount)
		args = append(args, *update.Intensity)
		argCount++
	}

	if update.DailyLimit != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf(`daily_limit = $%d`, argCount)
		args = append(args, *update.DailyLimit)
		argCount++
	}

	if update.IsActive != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf(`is_active = $%d`, argCount)
		args = append(args, *update.IsActive)
		argCount++
	}

	query += fmt.Sprintf(` WHERE id = $%d`, argCount)
	args = append(args, poolID)

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update warming pool", zap.Error(err))
		return err
	}

	return nil
}

// UpdateWarmingPoolMessageCount increments daily message count
func (d *Database) UpdateWarmingPoolMessageCount(poolID uuid.UUID, count int32) error {
	query := `UPDATE warming_pool SET message_send_today = message_send_today + $1 WHERE id = $2`

	_, err := d.Exec(query, count, poolID)
	if err != nil {
		utils.Error("Failed to update message count", zap.Error(err))
		return err
	}

	return nil
}

// ResetWarmingPoolDailyCount resets daily message count
func (d *Database) ResetWarmingPoolDailyCount(poolID uuid.UUID) error {
	query := `UPDATE warming_pool SET message_send_today = 0 WHERE id = $1`

	_, err := d.Exec(query, poolID)
	if err != nil {
		utils.Error("Failed to reset daily count", zap.Error(err))
		return err
	}

	return nil
}

// UpdateWarmingPoolNextAction updates next action timestamp
func (d *Database) UpdateWarmingPoolNextAction(poolID uuid.UUID, nextActionAt time.Time) error {
	query := `UPDATE warming_pool SET next_action_at = $1 WHERE id = $2`

	_, err := d.Exec(query, nextActionAt, poolID)
	if err != nil {
		utils.Error("Failed to update next action", zap.Error(err))
		return err
	}

	return nil
}

// CreateWarmingSession creates a warming session
func (d *Database) CreateWarmingSession(session *models.WarmingSession) error {
	query := `
		INSERT INTO warming_sessions (id, device_id, target_phone, message_sent, response_received, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := d.Exec(query,
		session.ID, session.DeviceID, session.TargetPhone, session.MessageSent,
		session.ResponseReceived, session.Status,
	)

	if err != nil {
		utils.Error("Failed to create warming session", zap.Error(err))
		return err
	}

	return nil
}

// GetWarmingSessionsByDeviceID retrieves warming sessions for a device
func (d *Database) GetWarmingSessionsByDeviceID(deviceID uuid.UUID, limit, offset int) ([]models.WarmingSession, error) {
	query := `
		SELECT id, device_id, target_phone, message_sent, response_received, status, created_at
		FROM warming_sessions
		WHERE device_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, deviceID, limit, offset)
	if err != nil {
		utils.Error("Failed to get warming sessions", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	sessions := []models.WarmingSession{}
	for rows.Next() {
		session := models.WarmingSession{}
		err := rows.Scan(
			&session.ID, &session.DeviceID, &session.TargetPhone, &session.MessageSent,
			&session.ResponseReceived, &session.Status, &session.CreatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan session", zap.Error(err))
			continue
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// UpdateWarmingSessionStatus updates warming session status
func (d *Database) UpdateWarmingSessionStatus(sessionID uuid.UUID, status int32, responseReceived string) error {
	query := `UPDATE warming_sessions SET status = $1, response_received = $2 WHERE id = $3`

	_, err := d.Exec(query, status, responseReceived, sessionID)
	if err != nil {
		utils.Error("Failed to update warming session", zap.Error(err))
		return err
	}

	return nil
}

// GetActiveWarmingPools retrieves all active warming pools
func (d *Database) GetActiveWarmingPools() ([]models.WarmingPool, error) {
	query := `
		SELECT id, device_id, intensity, daily_limit, message_send_today, is_active, next_action_at
		FROM warming_pool
		WHERE is_active = true
		ORDER BY next_action_at ASC
	`

	rows, err := d.Query(query)
	if err != nil {
		utils.Error("Failed to get active warming pools", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	pools := []models.WarmingPool{}
	for rows.Next() {
		pool := models.WarmingPool{}
		err := rows.Scan(
			&pool.ID, &pool.DeviceID, &pool.Intensity, &pool.DailyLimit, &pool.MessageSendToday,
			&pool.IsActive, &pool.NextActionAt,
		)
		if err != nil {
			utils.Error("Failed to scan pool", zap.Error(err))
			continue
		}
		pools = append(pools, pool)
	}

	return pools, nil
}

// DeactivateWarmingPool deactivates a warming pool
func (d *Database) DeactivateWarmingPool(poolID uuid.UUID) error {
	query := `UPDATE warming_pool SET is_active = false WHERE id = $1`

	_, err := d.Exec(query, poolID)
	if err != nil {
		utils.Error("Failed to deactivate warming pool", zap.Error(err))
		return err
	}

	return nil
}

// ActivateWarmingPool activates a warming pool
func (d *Database) ActivateWarmingPool(poolID uuid.UUID) error {
	query := `UPDATE warming_pool SET is_active = true WHERE id = $1`

	_, err := d.Exec(query, poolID)
	if err != nil {
		utils.Error("Failed to activate warming pool", zap.Error(err))
		return err
	}

	return nil
}
