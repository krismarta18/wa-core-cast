package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// jwtHeader is the fixed base64url-encoded header for HS256 JWTs.
var jwtHeader = base64url(mustMarshal(map[string]string{"alg": "HS256", "typ": "JWT"}))

// JWTClaims holds the claims embedded in every JWT issued by this service.
type JWTClaims struct {
	// Standard JWT claims
	Subject   string `json:"sub"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`

	// Application-specific claims
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone"`
	FullName    string `json:"name"`
}

// GenerateJWT signs a new HS256 JWT for the given user and returns the token string.
func GenerateJWT(userID, phoneNumber, fullName, secret string, expiryHours int) (string, error) {
	if secret == "" {
		return "", errors.New("JWT secret must not be empty")
	}

	now := time.Now()
	claims := JWTClaims{
		Subject:     userID,
		IssuedAt:    now.Unix(),
		ExpiresAt:   now.Add(time.Duration(expiryHours) * time.Hour).Unix(),
		UserID:      userID,
		PhoneNumber: phoneNumber,
		FullName:    fullName,
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JWT claims: %w", err)
	}

	encodedPayload := base64url(payload)
	signingInput := jwtHeader + "." + encodedPayload
	sig := sign(signingInput, secret)

	return signingInput + "." + sig, nil
}

// ValidateJWT parses and validates a JWT string, returning its claims on success.
func ValidateJWT(tokenStr, secret string) (*JWTClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("malformed JWT: expected 3 parts")
	}

	// Verify signature
	signingInput := parts[0] + "." + parts[1]
	if expected := sign(signingInput, secret); !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return nil, errors.New("JWT signature is invalid")
	}

	// Decode and unmarshal payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWT claims: %w", err)
	}

	// Check expiry
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("JWT has expired")
	}

	return &claims, nil
}

// ─── private helpers ──────────────────────────────────────────────────────────

func sign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return base64url(mac.Sum(nil))
}

func base64url(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

