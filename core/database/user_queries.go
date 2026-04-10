package database

import (
	"fmt"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateUser creates a new user
func (d *Database) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, phone, nama_lengkap, is_verify, otp_code, otp_expired, 
			id_subscribed, max_device, is_ban, is_api)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := d.Exec(query,
		user.ID, user.Phone, user.NamaLengkap, user.IsVerify, user.OTPCode,
		user.OTPExpired, user.IDSubscribed, user.MaxDevice, user.IsBan, user.IsAPI,
	)

	if err != nil {
		utils.Error("Failed to create user", zap.Error(err), zap.String("phone", user.Phone))
		return err
	}

	utils.Debug("User created", zap.String("user_id", user.ID.String()), zap.String("phone", user.Phone))
	return nil
}

// GetUserByID retrieves a user by ID
func (d *Database) GetUserByID(userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, phone, nama_lengkap, is_verify, otp_code, otp_expired, 
			id_subscribed, max_device, is_ban, is_api
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := d.QueryRow(query, userID).Scan(
		&user.ID, &user.Phone, &user.NamaLengkap, &user.IsVerify, &user.OTPCode,
		&user.OTPExpired, &user.IDSubscribed, &user.MaxDevice, &user.IsBan, &user.IsAPI,
	)

	if err != nil {
		utils.Debug("User not found", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}

	return user, nil
}

// GetUserByPhone retrieves a user by phone number
func (d *Database) GetUserByPhone(phone string) (*models.User, error) {
	query := `
		SELECT id, phone, nama_lengkap, is_verify, otp_code, otp_expired, 
			id_subscribed, max_device, is_ban, is_api
		FROM users
		WHERE phone = $1
	`

	user := &models.User{}
	err := d.QueryRow(query, phone).Scan(
		&user.ID, &user.Phone, &user.NamaLengkap, &user.IsVerify, &user.OTPCode,
		&user.OTPExpired, &user.IDSubscribed, &user.MaxDevice, &user.IsBan, &user.IsAPI,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user information
func (d *Database) UpdateUser(userID uuid.UUID, update *models.UpdateUserRequest) error {
	query := `UPDATE users SET `
	args := []interface{}{}
	argCount := 1

	if update.NamaLengkap != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("nama_lengkap = $%d", argCount)
		args = append(args, *update.NamaLengkap)
		argCount++
	}

	if update.OTPCode != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("otp_code = $%d", argCount)
		args = append(args, *update.OTPCode)
		argCount++
	}

	if update.OTPExpired != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("otp_expired = $%d", argCount)
		args = append(args, *update.OTPExpired)
		argCount++
	}

	if update.IsVerify != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("is_verify = $%d", argCount)
		args = append(args, *update.IsVerify)
		argCount++
	}

	if update.IsBan != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("is_ban = $%d", argCount)
		args = append(args, *update.IsBan)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, userID)

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update user", zap.Error(err), zap.String("user_id", userID.String()))
		return err
	}

	return nil
}

// DeleteUser soft deletes a user (sets is_ban = true)
func (d *Database) DeleteUser(userID uuid.UUID) error {
	query := `UPDATE users SET is_ban = true WHERE id = $1`

	_, err := d.Exec(query, userID)
	if err != nil {
		utils.Error("Failed to delete user", zap.Error(err), zap.String("user_id", userID.String()))
		return err
	}

	return nil
}

// GetUserCount returns total user count
func (d *Database) GetUserCount() (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE is_ban = false`

	var count int64
	err := d.QueryRow(query).Scan(&count)
	if err != nil {
		utils.Error("Failed to get user count", zap.Error(err))
		return 0, err
	}

	return count, nil
}

// GetAllUsers retrieves all active users with pagination
func (d *Database) GetAllUsers(limit, offset int) ([]models.User, error) {
	query := `
		SELECT id, phone, nama_lengkap, is_verify, otp_code, otp_expired, 
			id_subscribed, max_device, is_ban, is_api
		FROM users
		WHERE is_ban = false
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := d.Query(query, limit, offset)
	if err != nil {
		utils.Error("Failed to get users", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(
			&user.ID, &user.Phone, &user.NamaLengkap, &user.IsVerify, &user.OTPCode,
			&user.OTPExpired, &user.IDSubscribed, &user.MaxDevice, &user.IsBan, &user.IsAPI,
		)
		if err != nil {
			utils.Error("Failed to scan user", zap.Error(err))
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// VerifyUser marks user as verified
func (d *Database) VerifyUser(userID uuid.UUID) error {
	query := `UPDATE users SET is_verify = true, otp_code = NULL WHERE id = $1`

	_, err := d.Exec(query, userID)
	if err != nil {
		utils.Error("Failed to verify user", zap.Error(err))
		return err
	}

	return nil
}

// BanUser marks user as banned
func (d *Database) BanUser(userID uuid.UUID, isBan bool) error {
	query := `UPDATE users SET is_ban = $1 WHERE id = $2`

	_, err := d.Exec(query, isBan, userID)
	if err != nil {
		utils.Error("Failed to ban/unban user", zap.Error(err))
		return err
	}

	return nil
}

// UpdateUserMaxDevice updates the max device count for a user
func (d *Database) UpdateUserMaxDevice(userID uuid.UUID, maxDevice int32) error {
	query := `UPDATE users SET max_device = $1 WHERE id = $2`

	_, err := d.Exec(query, maxDevice, userID)
	if err != nil {
		utils.Error("Failed to update user max device", zap.Error(err))
		return err
	}

	return nil
}
