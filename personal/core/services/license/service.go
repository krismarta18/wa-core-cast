package license

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"wacast/core/utils"
)

var (
	ErrInvalidLicense = errors.New("invalid license key")
	ErrExpiredLicense = errors.New("license has expired")
	ErrHWIDMismatch  = errors.New("license not valid for this machine")
)

// Secret constants for encryption. Change these for production!
const (
	licenseKeySalt = "WBC_WACAST_SECURE_SALT_2026"
)

func (s *Service) getLicensePath() string {
	return utils.GetDataPath(".wacast.lic")
}

type LicenseInfo struct {
	HWID       string    `json:"hwid"`
	ExpiryDate time.Time `json:"expiry_date"`
	IsActive   bool      `json:"is_active"`
	IsExpired  bool      `json:"is_expired"`
}

type Service struct {
	currentHWID string
}

func NewService() *Service {
	hwid, _ := utils.GetHWID()
	return &Service{
		currentHWID: hwid,
	}
}

// GetStatus checks the current license status on the machine
func (s *Service) GetStatus() (*LicenseInfo, error) {
	key, err := os.ReadFile(s.getLicensePath())
	if err != nil {
		return &LicenseInfo{HWID: s.currentHWID, IsActive: false}, nil
	}

	info, err := s.ValidateKey(string(key))
	if err != nil {
		return &LicenseInfo{HWID: s.currentHWID, IsActive: false}, err
	}

	return info, nil
}

// Activate saves a new license key to the disk
func (s *Service) Activate(key string) error {
	info, err := s.ValidateKey(key)
	if err != nil {
		return err
	}

	if info.HWID != s.currentHWID {
		return ErrHWIDMismatch
	}

	err = os.WriteFile(s.getLicensePath(), []byte(strings.TrimSpace(key)), 0600)
	if err != nil {
		return fmt.Errorf("failed to save license file: %w", err)
	}

	utils.Info("License activated successfully", zap.String("hwid", info.HWID), zap.Time("expiry", info.ExpiryDate))
	return nil
}

// ValidateKey decrypts and verifies a serial key string
func (s *Service) ValidateKey(serialKey string) (*LicenseInfo, error) {
	serialKey = strings.TrimSpace(serialKey)
	if serialKey == "" {
		return nil, ErrInvalidLicense
	}

	// 1. Decode from Base32
	data, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(serialKey)
	if err != nil {
		return nil, ErrInvalidLicense
	}

	// 2. Decrypt using AES
	decrypted, err := decrypt(data, licenseKeySalt)
	if err != nil {
		return nil, ErrInvalidLicense
	}

	// 3. Parse fields: HWID|EXPIRY|SIGNATURE
	parts := strings.Split(string(decrypted), "|")
	if len(parts) != 3 {
		return nil, ErrInvalidLicense
	}

	hwid := parts[0]
	expiryStr := parts[1]
	providedSignature := parts[2]

	// 4. Verify Signature (Strict check)
	// Signature is SHA256(HWID + "|" + EXPIRY + "|" + SALT)
	expectedSigData := fmt.Sprintf("%s|%s|%s", hwid, expiryStr, licenseKeySalt)
	h := sha256.New()
	h.Write([]byte(expectedSigData))
	expectedSignature := fmt.Sprintf("%x", h.Sum(nil))

	// Important: The provided signature might have trailing null bytes or padding from CFB
	// But since we split by "|", the third part contains EVERYTHING after the last "|"
	// If someone added "QQQQQQ" to the Base32 string, it would change the decrypted signature part.
	if providedSignature != expectedSignature {
		return nil, ErrInvalidLicense
	}

	expiryUnix, err := time.Parse(time.RFC3339, expiryStr)
	if err != nil {
		return nil, ErrInvalidLicense
	}

	isExpired := time.Now().After(expiryUnix)

	return &LicenseInfo{
		HWID:       hwid,
		ExpiryDate: expiryUnix,
		IsActive:   hwid == s.currentHWID,
		IsExpired:  isExpired,
	}, nil
}

// Encryption helpers
func decrypt(data []byte, key string) ([]byte, error) {
	hash := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
