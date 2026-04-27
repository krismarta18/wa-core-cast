package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"wacast/core/database"
	"wacast/core/models"
	"wacast/core/utils"
)

// Sentinel errors — handlers inspect these to choose the right HTTP status code.
var (
	ErrPhoneAlreadyRegistered = errors.New("phone number is already registered")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserBanned             = errors.New("user is banned")
	ErrOTPNotFound            = errors.New("no active OTP found — request a new one")
	ErrOTPExpired             = errors.New("OTP has expired — request a new one")
	ErrOTPInvalid             = errors.New("invalid OTP code")
	ErrTooManyAttempts        = errors.New("too many failed OTP attempts — request a new one")
	ErrSessionNotFound        = errors.New("session not found or already revoked")
	ErrRefreshTokenInvalid    = errors.New("refresh token is invalid or expired")
	ErrInvalidPassword        = errors.New("invalid password")
)

const (
	otpExpiryMinutes = 5
	otpMaxAttempts   = 5
	refreshTokenBytes = 32
)

// VerifyOTPResult contains the JWT and user returned after a successful OTP verification.
type VerifyOTPResult struct {
	AccessToken      string
	RefreshToken     string
	ExpiresIn        int // seconds
	RefreshExpiresIn int // seconds
	User             *models.User
}

// Service handles all authentication business logic.
type Service struct {
	db                    *database.Database
	jwtSecret             string
	jwtExpiryHours        int
	jwtRefreshExpiryHours int
}

// NewService creates a new auth service.
func NewService(db *database.Database, jwtSecret string, jwtExpiryHours, jwtRefreshExpiryHours int) *Service {
	return &Service{
		db:                    db,
		jwtSecret:             jwtSecret,
		jwtExpiryHours:        jwtExpiryHours,
		jwtRefreshExpiryHours: jwtRefreshExpiryHours,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Register — creates a new unverified user and sends OTP
// ─────────────────────────────────────────────────────────────────────────────

// Register creates or re-uses a user record for phoneNumber and dispatches an OTP.
// Returns ErrPhoneAlreadyRegistered if the phone is already verified.
func (s *Service) Register(ctx context.Context, phoneNumber, fullName string) error {
	existing, err := getUserByPhone(ctx, s.db, phoneNumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("register: query user: %w", err)
	}

	if existing != nil {
		if existing.IsVerified {
			return ErrPhoneAlreadyRegistered
		}
		// User exists but not yet verified — re-send OTP
		return s.sendOTP(ctx, s.db, existing.ID, phoneNumber, models.OTPContextRegister)
	}

	// New user — insert into users table
	user, err := insertUser(ctx, s.db, phoneNumber, fullName)
	if err != nil {
		return fmt.Errorf("register: insert user: %w", err)
	}

	return s.sendOTP(ctx, s.db, user.ID, phoneNumber, models.OTPContextRegister)
}

// ─────────────────────────────────────────────────────────────────────────────
// RequestOTP — sends an OTP to an existing verified user (login)
// ─────────────────────────────────────────────────────────────────────────────

// RequestOTP looks up a verified user by phone and dispatches a login OTP.
func (s *Service) RequestOTP(ctx context.Context, phoneNumber string) error {
	user, err := getUserByPhone(ctx, s.db, phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("request-otp: query user: %w", err)
	}

	if user.IsBanned {
		return ErrUserBanned
	}

	return s.sendOTP(ctx, s.db, user.ID, phoneNumber, models.OTPContextLogin)
}

// ─────────────────────────────────────────────────────────────────────────────
// VerifyOTP — validates OTP and returns a signed JWT
// ─────────────────────────────────────────────────────────────────────────────

// VerifyOTP validates the supplied OTP code for phoneNumber.
// On success it returns a signed JWT and the resolved user.
func (s *Service) VerifyOTP(ctx context.Context, phoneNumber, otpCode, ipAddress, userAgent string) (*VerifyOTPResult, error) {
	otp, err := getActiveOTP(ctx, s.db, phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOTPNotFound
		}
		return nil, fmt.Errorf("verify-otp: query otp: %w", err)
	}

	// Increment attempt counter first (before checking, so even reading counts)
	attempts, err := incrementOTPAttempts(ctx, s.db, otp.ID)
	if err != nil {
		return nil, fmt.Errorf("verify-otp: increment attempts: %w", err)
	}

	if attempts > otpMaxAttempts {
		return nil, ErrTooManyAttempts
	}

	if otp.OTPCode != otpCode {
		return nil, ErrOTPInvalid
	}

	// Mark OTP as used
	if err := markOTPVerified(ctx, s.db, otp.ID); err != nil {
		return nil, fmt.Errorf("verify-otp: mark verified: %w", err)
	}

	// Update user: set is_verified = true, last_login_at = NOW()
	user, err := markUserVerifiedAndLogin(ctx, s.db, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("verify-otp: update user: %w", err)
	}

	result, accessExpiresAt, refreshExpiresAt, err := s.issueSessionTokens(user)
	if err != nil {
		return nil, fmt.Errorf("verify-otp: issue session tokens: %w", err)
	}

	if err := createUserSession(ctx, s.db, user.ID, result.AccessToken, result.RefreshToken, ipAddress, userAgent, accessExpiresAt, refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("verify-otp: create session: %w", err)
	}

	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// LoginWithPassword — validates master password and returns a signed JWT
// ─────────────────────────────────────────────────────────────────────────────

func (s *Service) LoginWithPassword(ctx context.Context, password, ipAddress, userAgent string) (*VerifyOTPResult, error) {
	passwordFile := "./admin.password"
	var adminPassword string

	// Read admin password from file
	data, err := os.ReadFile(passwordFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default password file if it doesn't exist
			adminPassword = "admin"
			_ = os.WriteFile(passwordFile, []byte(adminPassword), 0600)
			utils.Info("Created default admin.password file")
		} else {
			return nil, fmt.Errorf("failed to read password file: %w", err)
		}
	} else {
		adminPassword = string(data)
	}

	if password != adminPassword {
		return nil, ErrInvalidPassword
	}

	// Get or create the master admin user (using a fixed phone number for personal mode)
	adminPhone := "0000000000"
	user, err := getUserByPhone(ctx, s.db, adminPhone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create it
			user, err = insertUser(ctx, s.db, adminPhone, "Admin Access")
			if err != nil {
				return nil, fmt.Errorf("login-password: create admin: %w", err)
			}
		} else {
			return nil, fmt.Errorf("login-password: query admin: %w", err)
		}
	}

	// Double check if name is correct, if not update it
	if user.FullName != "Admin Access" {
		updateReq := models.UpdateUserRequest{
			FullName: func(s string) *string { return &s }("Admin Access"),
		}
		user, _ = s.UpdateProfile(ctx, user.ID.String(), updateReq)
	}

	if user.IsBanned {
		return nil, ErrUserBanned
	}

	// Update last login
	user, err = markUserVerifiedAndLogin(ctx, s.db, user.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("login-password: update user: %w", err)
	}

	result, accessExpiresAt, refreshExpiresAt, err := s.issueSessionTokens(user)
	if err != nil {
		return nil, fmt.Errorf("login-password: issue session tokens: %w", err)
	}

	if err := createUserSession(ctx, s.db, user.ID, result.AccessToken, result.RefreshToken, ipAddress, userAgent, accessExpiresAt, refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("login-password: create session: %w", err)
	}

	return result, nil
}

// RefreshSession rotates an existing refresh token and issues a new access token pair.
func (s *Service) RefreshSession(ctx context.Context, refreshToken, ipAddress, userAgent string) (*VerifyOTPResult, error) {
	session, err := getSessionByRefreshToken(ctx, s.db, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRefreshTokenInvalid
		}
		return nil, fmt.Errorf("refresh-session: query session: %w", err)
	}

	user, err := getUserByID(ctx, s.db, session.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("refresh-session: query user: %w", err)
	}

	result, accessExpiresAt, refreshExpiresAt, err := s.issueSessionTokens(user)
	if err != nil {
		return nil, fmt.Errorf("refresh-session: issue session tokens: %w", err)
	}

	if err := rotateUserSession(ctx, s.db, session.ID, result.AccessToken, result.RefreshToken, ipAddress, userAgent, accessExpiresAt, refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("refresh-session: rotate session: %w", err)
	}

	return result, nil
}

func (s *Service) issueSessionTokens(user *models.User) (*VerifyOTPResult, time.Time, time.Time, error) {
	accessToken, err := utils.GenerateJWT(
		user.ID.String(),
		user.PhoneNumber,
		user.FullName,
		s.jwtSecret,
		s.jwtExpiryHours,
	)
	if err != nil {
		return nil, time.Time{}, time.Time{}, err
	}

	refreshToken, err := utils.GenerateSecureToken(refreshTokenBytes)
	if err != nil {
		return nil, time.Time{}, time.Time{}, err
	}

	accessExpiresAt := time.Now().UTC().Add(time.Duration(s.jwtExpiryHours) * time.Hour)
	refreshExpiresAt := time.Now().UTC().Add(time.Duration(s.jwtRefreshExpiryHours) * time.Hour)

	return &VerifyOTPResult{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        s.jwtExpiryHours * 3600,
		RefreshExpiresIn: s.jwtRefreshExpiryHours * 3600,
		User:             user,
	}, accessExpiresAt, refreshExpiresAt, nil
}

// ValidateSession ensures the supplied access token is still active in user_sessions.
func (s *Service) ValidateSession(ctx context.Context, accessToken string) error {
	active, err := touchActiveSession(ctx, s.db, accessToken)
	if err != nil {
		return fmt.Errorf("validate-session: %w", err)
	}
	if !active {
		return ErrSessionNotFound
	}
	return nil
}

// Logout revokes the persisted session associated with the current bearer token.
func (s *Service) Logout(ctx context.Context, accessToken string) error {
	revoked, err := revokeSession(ctx, s.db, accessToken)
	if err != nil {
		return fmt.Errorf("logout: revoke session: %w", err)
	}
	if !revoked {
		return ErrSessionNotFound
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// GetUser — fetch user by ID (used by /auth/me)
// ─────────────────────────────────────────────────────────────────────────────

// GetUser returns the user identified by userID.
func (s *Service) GetUser(ctx context.Context, userID string) (*models.User, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := getUserByID(ctx, s.db, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get-user: %w", err)
	}
	return user, nil
}

// UpdateProfile updates the profile information of the user identified by userID.
func (s *Service) UpdateProfile(ctx context.Context, userID string, req models.UpdateUserRequest) (*models.User, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	existing, err := getUserByID(ctx, s.db, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("update-profile: check user: %w", err)
	}

	if req.FullName != nil {
		existing.FullName = *req.FullName
	}
	if req.Email != nil {
		existing.Email = req.Email
	}
	if req.CompanyName != nil {
		existing.CompanyName = req.CompanyName
	}
	if req.Timezone != nil {
		existing.Timezone = *req.Timezone
	}

	user, err := updateUserProfile(ctx, s.db, existing)
	if err != nil {
		return nil, fmt.Errorf("update-profile: execute update: %w", err)
	}

	return user, nil
}

// UpdatePassword verifies oldPassword and updates it to newPassword.
// If the user is the Master Admin (phone 0000000000), it updates the admin.password file.
func (s *Service) UpdatePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := getUserByID(ctx, s.db, uid)
	if err != nil {
		return fmt.Errorf("update-password: check user: %w", err)
	}

	// 1. Handle Admin Password Update
	if user.PhoneNumber == "0000000000" {
		passwordFile := "./admin.password"
		data, err := os.ReadFile(passwordFile)
		if err != nil {
			return fmt.Errorf("failed to read admin password file: %w", err)
		}

		currentAdminPassword := string(data)
		if oldPassword != currentAdminPassword {
			return ErrInvalidPassword
		}

		// Update file
		if err := os.WriteFile(passwordFile, []byte(newPassword), 0600); err != nil {
			return fmt.Errorf("failed to update admin password file: %w", err)
		}

		utils.Info("Admin password updated successfully", zap.String("user_id", userID))
		return nil
	}

	// 2. Regular user password update logic (if any) could go here
	return errors.New("password update only supported for admin in personal mode")
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

// sendOTP invalidates any pending OTPs, generates a new one, stores it, and
// dispatches it via WhatsApp.
func (s *Service) sendOTP(ctx context.Context, db *database.Database, userID uuid.UUID, phoneNumber, otpContext string) error {
	// Invalidate all pending OTPs for this phone
	if err := invalidatePendingOTPs(ctx, db, phoneNumber); err != nil {
		return fmt.Errorf("sendOTP: invalidate old otps: %w", err)
	}

	// Generate plaintext OTP
	code, err := utils.GenerateOTP()
	if err != nil {
		return fmt.Errorf("sendOTP: generate code: %w", err)
	}

	expiresAt := time.Now().UTC().Add(otpExpiryMinutes * time.Minute)

	if err := insertOTP(ctx, db, userID, phoneNumber, otpContext, code, expiresAt); err != nil {
		return fmt.Errorf("sendOTP: insert otp: %w", err)
	}

	// Send via WhatsApp (dummy in current mode)
	if err := utils.SendOTPViaWhatsApp(phoneNumber, code); err != nil {
		// Log but do not fail — the OTP is already stored
		utils.Warn("sendOTP: WhatsApp delivery failed",
			zap.String("phone", phoneNumber),
			zap.Error(err),
		)
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// DB query functions
// ─────────────────────────────────────────────────────────────────────────────

const userSelectCols = `
	id, phone_number, full_name, email, company_name, timezone,
	is_verified, is_banned, is_api_enabled, created_at, updated_at, last_login_at
`

type authSessionRecord struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func scanUser(row interface {
	Scan(dest ...interface{}) error
}) (*models.User, error) {
	u := &models.User{}
	err := row.Scan(
		&u.ID, &u.PhoneNumber, &u.FullName, &u.Email, &u.CompanyName, &u.Timezone,
		&u.IsVerified, &u.IsBanned, &u.IsAPIEnabled, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func getUserByPhone(ctx context.Context, db *database.Database, phone string) (*models.User, error) {
	q := `SELECT ` + userSelectCols + ` FROM users WHERE phone_number = $1`
	return scanUser(db.QueryRowContext(ctx, q, phone))
}

func getUserByID(ctx context.Context, db *database.Database, id uuid.UUID) (*models.User, error) {
	q := `SELECT ` + userSelectCols + ` FROM users WHERE id = $1`
	return scanUser(db.QueryRowContext(ctx, q, id))
}

func insertUser(ctx context.Context, db *database.Database, phoneNumber, fullName string) (*models.User, error) {
	q := `
		INSERT INTO users (phone_number, full_name)
		VALUES ($1, $2)
		RETURNING ` + userSelectCols
	return scanUser(db.QueryRowContext(ctx, q, phoneNumber, fullName))
}

func markUserVerifiedAndLogin(ctx context.Context, db *database.Database, phoneNumber string) (*models.User, error) {
	q := `
		UPDATE users
		SET is_verified = true, last_login_at = NOW(), updated_at = NOW()
		WHERE phone_number = $1
		RETURNING ` + userSelectCols
	return scanUser(db.QueryRowContext(ctx, q, phoneNumber))
}

func updateUserProfile(ctx context.Context, db *database.Database, user *models.User) (*models.User, error) {
	q := `
		UPDATE users
		SET full_name = $1, email = $2, company_name = $3, timezone = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING ` + userSelectCols
	return scanUser(db.QueryRowContext(ctx, q, user.FullName, user.Email, user.CompanyName, user.Timezone, user.ID))
}

func invalidatePendingOTPs(ctx context.Context, db *database.Database, phoneNumber string) error {
	_, err := db.ExecContext(ctx,
		`UPDATE otp_verifications SET verified_at = NOW()
		 WHERE phone_number = $1 AND verified_at IS NULL`,
		phoneNumber,
	)
	return err
}

func insertOTP(ctx context.Context, db *database.Database, userID uuid.UUID, phoneNumber, otpContext, code string, expiresAt time.Time) error {
	_, err := db.ExecContext(ctx,
		`INSERT INTO otp_verifications (user_id, phone_number, context, otp_code, expires_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		userID, phoneNumber, otpContext, code, expiresAt,
	)
	return err
}

func getActiveOTP(ctx context.Context, db *database.Database, phoneNumber string) (*models.OTPVerification, error) {
	q := `
		SELECT id, user_id, phone_number, context, otp_code, attempt_count, expires_at, verified_at, created_at
		FROM otp_verifications
		WHERE phone_number = $1
		  AND verified_at IS NULL
		  AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1`

	otp := &models.OTPVerification{}
	err := db.QueryRowContext(ctx, q, phoneNumber).Scan(
		&otp.ID, &otp.UserID, &otp.PhoneNumber, &otp.Context,
		&otp.OTPCode, &otp.AttemptCount, &otp.ExpiresAt, &otp.VerifiedAt, &otp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func incrementOTPAttempts(ctx context.Context, db *database.Database, otpID uuid.UUID) (int, error) {
	var attempts int
	err := db.QueryRowContext(ctx,
		`UPDATE otp_verifications SET attempt_count = attempt_count + 1
		 WHERE id = $1 RETURNING attempt_count`,
		otpID,
	).Scan(&attempts)
	return attempts, err
}

func markOTPVerified(ctx context.Context, db *database.Database, otpID uuid.UUID) error {
	_, err := db.ExecContext(ctx,
		`UPDATE otp_verifications SET verified_at = NOW() WHERE id = $1`,
		otpID,
	)
	return err
}

func createUserSession(ctx context.Context, db *database.Database, userID uuid.UUID, accessToken, refreshToken, ipAddress, userAgent string, expiresAt, refreshExpiresAt time.Time) error {
	_, err := db.ExecContext(ctx,
		`INSERT INTO user_sessions (user_id, session_token_hash, refresh_token_hash, ip_address, user_agent, expires_at, refresh_expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID,
		hashToken(accessToken),
		hashToken(refreshToken),
		nullableString(ipAddress),
		nullableString(userAgent),
		expiresAt,
		refreshExpiresAt,
	)
	return err
}

func getSessionByRefreshToken(ctx context.Context, db *database.Database, refreshToken string) (*authSessionRecord, error) {
	record := &authSessionRecord{}
	err := db.QueryRowContext(ctx,
		`SELECT id, user_id
		 FROM user_sessions
		 WHERE refresh_token_hash = $1
		   AND revoked_at IS NULL
		   AND refresh_expires_at > NOW()
		 LIMIT 1`,
		hashToken(refreshToken),
	).Scan(&record.ID, &record.UserID)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func rotateUserSession(ctx context.Context, db *database.Database, sessionID uuid.UUID, accessToken, refreshToken, ipAddress, userAgent string, expiresAt, refreshExpiresAt time.Time) error {
	_, err := db.ExecContext(ctx,
		`UPDATE user_sessions
		 SET session_token_hash = $1,
		     refresh_token_hash = $2,
		     ip_address = $3,
		     user_agent = $4,
		     last_active_at = NOW(),
		     expires_at = $5,
		     refresh_expires_at = $6,
		     revoked_at = NULL
		 WHERE id = $7`,
		hashToken(accessToken),
		hashToken(refreshToken),
		nullableString(ipAddress),
		nullableString(userAgent),
		expiresAt,
		refreshExpiresAt,
		sessionID,
	)
	return err
}

func touchActiveSession(ctx context.Context, db *database.Database, accessToken string) (bool, error) {
	result, err := db.ExecContext(ctx,
		`UPDATE user_sessions
		 SET last_active_at = NOW()
		 WHERE session_token_hash = $1
		   AND revoked_at IS NULL
		   AND expires_at > NOW()`,
		hashToken(accessToken),
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func revokeSession(ctx context.Context, db *database.Database, accessToken string) (bool, error) {
	result, err := db.ExecContext(ctx,
		`UPDATE user_sessions
		 SET revoked_at = NOW(), last_active_at = NOW()
		 WHERE session_token_hash = $1
		   AND revoked_at IS NULL`,
		hashToken(accessToken),
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}
