package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"wacast/core/models"
	"wacast/core/utils"
	"go.uber.org/zap"
)

// CreateDevice creates a new device/session
func (d *Database) CreateDevice(device *models.Device) error {
	query := `
		INSERT INTO devices (id, user_id, unique_name, name_device, phone, status, session_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := d.Exec(query,
		device.ID, device.UserID, device.UniqueName, device.NameDevice, 
		device.Phone, device.Status, device.SessionData,
	)

	if err != nil {
		utils.Error("Failed to create device", zap.Error(err), zap.String("device", device.NameDevice))
		return err
	}

	utils.Debug("Device created", zap.String("device_id", device.ID.String()), zap.String("phone", device.Phone))
	return nil
}

// GetDeviceByID retrieves a device by ID
func (d *Database) GetDeviceByID(deviceID uuid.UUID) (*models.Device, error) {
	query := `
		SELECT id, user_id, unique_name, name_device, phone, status, last_seen, session_data
		FROM devices
		WHERE id = $1
	`

	device := &models.Device{}
	err := d.QueryRow(query, deviceID).Scan(
		&device.ID, &device.UserID, &device.UniqueName, &device.NameDevice,
		&device.Phone, &device.Status, &device.LastSeen, &device.SessionData,
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
		SELECT id, user_id, unique_name, name_device, phone, status, last_seen, session_data
		FROM devices
		WHERE user_id = $1 AND status != 2
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
			&device.ID, &device.UserID, &device.UniqueName, &device.NameDevice,
			&device.Phone, &device.Status, &device.LastSeen, &device.SessionData,
		)
		if err != nil {
			utils.Error("Failed to scan device", zap.Error(err))
			continue
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// GetActiveDevices retrieves all active devices
func (d *Database) GetActiveDevices() ([]models.Device, error) {
	query := `
		SELECT id, user_id, unique_name, name_device, phone, status, last_seen, session_data
		FROM devices
		WHERE status = 1
		ORDER BY last_seen DESC
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
			&device.ID, &device.UserID, &device.UniqueName, &device.NameDevice,
			&device.Phone, &device.Status, &device.LastSeen, &device.SessionData,
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
func (d *Database) UpdateDeviceStatus(deviceID uuid.UUID, status int32) error {
	query := `UPDATE devices SET status = $1, last_seen = $2 WHERE id = $3`

	_, err := d.Exec(query, status, time.Now(), deviceID)
	if err != nil {
		utils.Error("Failed to update device status", zap.Error(err))
		return err
	}

	return nil
}

// UpdateDeviceSessionData updates device session data
func (d *Database) UpdateDeviceSessionData(deviceID uuid.UUID, sessionData []byte) error {
	query := `UPDATE devices SET session_data = $1, last_seen = $2 WHERE id = $3`

	_, err := d.Exec(query, sessionData, time.Now(), deviceID)
	if err != nil {
		utils.Error("Failed to update device session data", zap.Error(err))
		return err
	}

	return nil
}

// UpdateDeviceInfo updates device information
func (d *Database) UpdateDeviceInfo(deviceID uuid.UUID, update *models.UpdateDeviceRequest) error {
	query := `UPDATE devices SET `
	args := []interface{}{}
	argCount := 1

	if update.NameDevice != nil {
		query += fmt.Sprintf("name_device = $%d", argCount)
		args = append(args, *update.NameDevice)
		argCount++
	}

	if update.Phone != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("phone = $%d", argCount)
		args = append(args, *update.Phone)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, deviceID)

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update device", zap.Error(err))
		return err
	}

	return nil
}

// DeleteDevice soft deletes a device (sets status = 2)
func (d *Database) DeleteDevice(deviceID uuid.UUID) error {
	query := `UPDATE devices SET status = 2 WHERE id = $1`

	_, err := d.Exec(query, deviceID)
	if err != nil {
		utils.Error("Failed to delete device", zap.Error(err))
		return err
	}

	return nil
}

// GetDeviceByPhone retrieves a device by phone number
func (d *Database) GetDeviceByPhone(phone string) (*models.Device, error) {
	query := `
		SELECT id, user_id, unique_name, name_device, phone, status, last_seen, session_data
		FROM devices
		WHERE phone = $1 AND status = 1
	`

	device := &models.Device{}
	err := d.QueryRow(query, phone).Scan(
		&device.ID, &device.UserID, &device.UniqueName, &device.NameDevice,
		&device.Phone, &device.Status, &device.LastSeen, &device.SessionData,
	)

	if err != nil {
		utils.Debug("Device not found by phone", zap.String("phone", phone))
		return nil, err
	}

	return device, nil
}

// CountUserDevices counts active devices for a user
func (d *Database) CountUserDevices(userID uuid.UUID) (int32, error) {
	query := `SELECT COUNT(*) FROM devices WHERE user_id = $1 AND status != 2`

	var count int32
	err := d.QueryRow(query, userID).Scan(&count)
	if err != nil {
		utils.Error("Failed to count user devices", zap.Error(err))
		return 0, err
	}

	return count, nil
}

// UpdateDeviceLastSeen updates last seen timestamp
func (d *Database) UpdateDeviceLastSeen(deviceID uuid.UUID) error {
	query := `UPDATE devices SET last_seen = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), deviceID)
	if err != nil {
		utils.Error("Failed to update device last seen", zap.Error(err))
		return err
	}

	return nil
}
