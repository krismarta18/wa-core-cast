package database

import (
	"fmt"
	"time"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateUser creates a new user
func (d *Database) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, phone_number, full_name, email, company_name, timezone,
			is_verified, is_banned, is_api_enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	_, err := d.Exec(query,
		user.ID, user.PhoneNumber, user.FullName, user.Email, user.CompanyName, user.Timezone,
		user.IsVerified, user.IsBanned, user.IsAPIEnabled, now, now,
	)

	if err != nil {
		utils.Error("Failed to create user", zap.Error(err), zap.String("phone", user.PhoneNumber))
		return err
	}

	utils.Debug("User created", zap.String("user_id", user.ID.String()), zap.String("phone", user.PhoneNumber))
	return nil
}

// GetUserByID retrieves a user by ID
func (d *Database) GetUserByID(userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, phone_number, full_name, email, company_name, timezone,
			is_verified, is_banned, is_api_enabled, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := d.QueryRow(query, userID).Scan(
		&user.ID, &user.PhoneNumber, &user.FullName, &user.Email, &user.CompanyName, &user.Timezone,
		&user.IsVerified, &user.IsBanned, &user.IsAPIEnabled,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
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
		SELECT id, phone_number, full_name, email, company_name, timezone,
			is_verified, is_banned, is_api_enabled, created_at, updated_at, last_login_at
		FROM users
		WHERE phone_number = $1
	`

	user := &models.User{}
	err := d.QueryRow(query, phone).Scan(
		&user.ID, &user.PhoneNumber, &user.FullName, &user.Email, &user.CompanyName, &user.Timezone,
		&user.IsVerified, &user.IsBanned, &user.IsAPIEnabled,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
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

	if update.FullName != nil {
		query += fmt.Sprintf("full_name = $%d", argCount)
		args = append(args, *update.FullName)
		argCount++
	}

	if update.Email != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("email = $%d", argCount)
		args = append(args, *update.Email)
		argCount++
	}

	if update.CompanyName != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("company_name = $%d", argCount)
		args = append(args, *update.CompanyName)
		argCount++
	}

	if update.Timezone != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("timezone = $%d", argCount)
		args = append(args, *update.Timezone)
		argCount++
	}

	if update.IsVerified != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("is_verified = $%d", argCount)
		args = append(args, *update.IsVerified)
		argCount++
	}

	if update.IsBanned != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("is_banned = $%d", argCount)
		args = append(args, *update.IsBanned)
		argCount++
	}

	if update.IsAPIEnabled != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("is_api_enabled = $%d", argCount)
		args = append(args, *update.IsAPIEnabled)
		argCount++
	}

	if argCount > 1 {
		query += ", "
	}
	query += fmt.Sprintf("updated_at = $%d WHERE id = $%d", argCount, argCount+1)
	args = append(args, time.Now(), userID)

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update user", zap.Error(err), zap.String("user_id", userID.String()))
		return err
	}

	return nil
}

// UpdateUserLastLogin updates the last login timestamp
func (d *Database) UpdateUserLastLogin(userID uuid.UUID) error {
	now := time.Now()
	query := `UPDATE users SET last_login_at = $1, updated_at = $1 WHERE id = $2`

	_, err := d.Exec(query, now, userID)
	if err != nil {
		utils.Error("Failed to update user last login", zap.Error(err))
		return err
	}

	return nil
}

// DeleteUser soft deletes a user (sets is_banned = true)
func (d *Database) DeleteUser(userID uuid.UUID) error {
	query := `UPDATE users SET is_banned = true, updated_at = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), userID)
	if err != nil {
		utils.Error("Failed to delete user", zap.Error(err), zap.String("user_id", userID.String()))
		return err
	}

	return nil
}

// GetUserCount returns total non-banned user count
func (d *Database) GetUserCount() (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE is_banned = false`

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
		SELECT id, phone_number, full_name, email, company_name, timezone,
			is_verified, is_banned, is_api_enabled, created_at, updated_at, last_login_at
		FROM users
		WHERE is_banned = false
		ORDER BY created_at DESC
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
			&user.ID, &user.PhoneNumber, &user.FullName, &user.Email, &user.CompanyName, &user.Timezone,
			&user.IsVerified, &user.IsBanned, &user.IsAPIEnabled,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
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
	query := `UPDATE users SET is_verified = true, updated_at = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), userID)
	if err != nil {
		utils.Error("Failed to verify user", zap.Error(err))
		return err
	}

	return nil
}

// BanUser bans or unbans a user
func (d *Database) BanUser(userID uuid.UUID, isBanned bool) error {
	query := `UPDATE users SET is_banned = $1, updated_at = $2 WHERE id = $3`

	_, err := d.Exec(query, isBanned, time.Now(), userID)
	if err != nil {
		utils.Error("Failed to ban/unban user", zap.Error(err))
		return err
	}

	return nil
}

