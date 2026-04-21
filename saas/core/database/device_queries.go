package database

import (
	"fmt"
	"time"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateDevice creates a new device/session
func (d *Database) CreateDevice(device *models.Device) error {
	query := `
		INSERT INTO devices (id, user_id, unique_name, display_name, phone_number, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	_, err := d.Exec(query,
		device.ID, device.UserID, device.UniqueName, device.DisplayName,
		device.PhoneNumber, device.Status, now, now,
	)

	if err != nil {
		utils.Error("Failed to create device", zap.Error(err), zap.String("device", device.DisplayName))
		return err
	}

	utils.Debug("Device created", zap.String("device_id", device.ID.String()), zap.String("phone", device.PhoneNumber))
	return nil
}

// GetDeviceByID retrieves a device by ID
func (d *Database) GetDeviceByID(deviceID uuid.UUID) (*models.Device, error) {
	query := `
		SELECT id, user_id, unique_name, display_name, phone_number, status,
		       last_seen_at, connected_since, platform, wa_version, battery_level,
		       created_at, updated_at
		FROM devices
		WHERE id = $1
	`

	device := &models.Device{}
	err := d.QueryRow(query, deviceID).Scan(
		&device.ID, &device.UserID, &device.UniqueName, &device.DisplayName,
		&device.PhoneNumber, &device.Status, &device.LastSeenAt, &device.ConnectedSince,
		&device.Platform, &device.WaVersion, &device.BatteryLevel,
		&device.CreatedAt, &device.UpdatedAt,
	)

	if err != nil {
		utils.Debug("Device not found", zap.String("device_id", deviceID.String()))
		return nil, err
	}

	return device, nil
}

// GetDevicesByUserID retrieves all devices for a user
func (d *Database) GetDevicesByUserID(userID uuid.UUID) ([]models.Device, error) {
	query := `
		SELECT id, user_id, unique_name, display_name, phone_number, status,
		       last_seen_at, connected_since, platform, wa_version, battery_level,
		       created_at, updated_at
		FROM devices
		WHERE user_id = $1 AND status != 'banned'
		ORDER BY created_at DESC
	`

	rows, err := d.Query(query, userID)
	if err != nil {
		utils.Error("Failed to get devices", zap.Error(err), zap.String("user_id", userID.String()))
		return nil, err
	}
	defer rows.Close()

	devices := []models.Device{}
	for rows.Next() {
		device := models.Device{}
		err := rows.Scan(
			&device.ID, &device.UserID, &device.UniqueName, &device.DisplayName,
			&device.PhoneNumber, &device.Status, &device.LastSeenAt, &device.ConnectedSince,
			&device.Platform, &device.WaVersion, &device.BatteryLevel,
			&device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan device", zap.Error(err))
			continue
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// GetActiveDevices retrieves all connected devices
func (d *Database) GetActiveDevices() ([]models.Device, error) {
	query := `
		SELECT id, user_id, unique_name, display_name, phone_number, status,
		       last_seen_at, connected_since, platform, wa_version, battery_level,
		       created_at, updated_at
		FROM devices
		WHERE status = 'connected'
		ORDER BY last_seen_at DESC
	`

	rows, err := d.Query(query)
	if err != nil {
		utils.Error("Failed to get active devices", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	devices := []models.Device{}
	for rows.Next() {
		device := models.Device{}
		err := rows.Scan(
			&device.ID, &device.UserID, &device.UniqueName, &device.DisplayName,
			&device.PhoneNumber, &device.Status, &device.LastSeenAt, &device.ConnectedSince,
			&device.Platform, &device.WaVersion, &device.BatteryLevel,
			&device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			utils.Error("Failed to scan device", zap.Error(err))
			continue
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// UpdateDeviceStatus updates device status
func (d *Database) UpdateDeviceStatus(deviceID uuid.UUID, status string) error {
	query := `UPDATE devices SET status = $1, last_seen_at = $2, updated_at = $2 WHERE id = $3`

	_, err := d.Exec(query, status, time.Now(), deviceID)
	if err != nil {
		utils.Error("Failed to update device status", zap.Error(err))
		return err
	}

	return nil
}

// UpdateDeviceInfo updates device information
func (d *Database) UpdateDeviceInfo(deviceID uuid.UUID, update *models.UpdateDeviceRequest) error {
	query := `UPDATE devices SET `
	args := []interface{}{}
	argCount := 1

	if update.DisplayName != nil {
		query += fmt.Sprintf("display_name = $%d", argCount)
		args = append(args, *update.DisplayName)
		argCount++
	}

	if update.PhoneNumber != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("phone_number = $%d", argCount)
		args = append(args, *update.PhoneNumber)
		argCount++
	}

	if argCount > 1 {
		query += ", "
	}
	query += fmt.Sprintf("updated_at = $%d WHERE id = $%d", argCount, argCount+1)
	args = append(args, time.Now(), deviceID)

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update device", zap.Error(err))
		return err
	}

	return nil
}

// DeleteDevice soft deletes a device (sets status to banned)
func (d *Database) DeleteDevice(deviceID uuid.UUID) error {
	query := `UPDATE devices SET status = 'banned', updated_at = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), deviceID)
	if err != nil {
		utils.Error("Failed to delete device", zap.Error(err))
		return err
	}

	return nil
}

// GetDeviceByPhone retrieves a device by phone number
func (d *Database) GetDeviceByPhone(phone string) (*models.Device, error) {
	query := `
		SELECT id, user_id, unique_name, display_name, phone_number, status,
		       last_seen_at, connected_since, platform, wa_version, battery_level,
		       created_at, updated_at
		FROM devices
		WHERE phone_number = $1 AND status = 'connected'
	`

	device := &models.Device{}
	err := d.QueryRow(query, phone).Scan(
		&device.ID, &device.UserID, &device.UniqueName, &device.DisplayName,
		&device.PhoneNumber, &device.Status, &device.LastSeenAt, &device.ConnectedSince,
		&device.Platform, &device.WaVersion, &device.BatteryLevel,
		&device.CreatedAt, &device.UpdatedAt,
	)

	if err != nil {
		utils.Debug("Device not found by phone", zap.String("phone", phone))
		return nil, err
	}

	return device, nil
}

// CountUserDevices counts non-banned devices for a user
func (d *Database) CountUserDevices(userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM devices
		WHERE user_id = $1
		  AND COALESCE(status::text, '') NOT IN ('banned', '3')
	`

	var count int
	err := d.QueryRow(query, userID).Scan(&count)
	if err != nil {
		utils.Error("Failed to count user devices", zap.Error(err))
		return 0, err
	}

	return count, nil
}

// UpdateDeviceLastSeen updates last seen timestamp
func (d *Database) UpdateDeviceLastSeen(deviceID uuid.UUID) error {
	query := `UPDATE devices SET last_seen_at = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), deviceID)
	if err != nil {
		utils.Error("Failed to update device last seen", zap.Error(err))
		return err
	}

	return nil
}

