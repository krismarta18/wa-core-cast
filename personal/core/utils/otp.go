package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

const otpDigits = 6

// GenerateOTP generates a cryptographically secure 6-digit OTP string.
func GenerateOTP() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}
	return fmt.Sprintf("%0*d", otpDigits, n.Int64()), nil
}

// GenerateSecureToken generates a cryptographically secure URL-safe token.
func GenerateSecureToken(byteLength int) (string, error) {
	buf := make([]byte, byteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
