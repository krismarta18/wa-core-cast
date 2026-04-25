package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	licenseKeySalt = "WBC_WACAST_SECURE_SALT_2026"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run license_generator.go <HWID> <Duration> <Unit: s|d|y>")
		fmt.Println("Example 10 Seconds: go run license_generator.go 1A2B3D4D 10 s")
		fmt.Println("Example 7 Days:     go run license_generator.go 1A2B3C4D 7 d")
		fmt.Println("Example 1 Year:     go run license_generator.go 1A2B3C4D 1 y")
		return
	}

	hwid := strings.ToUpper(strings.TrimSpace(os.Args[1]))
	if hwid == "" {
		fmt.Println("Error: HWID cannot be empty")
		return
	}
	
	duration := 1
	unit := "y"

	if len(os.Args) > 2 {
		fmt.Sscanf(os.Args[2], "%d", &duration)
	}
	if len(os.Args) > 3 {
		unit = strings.ToLower(os.Args[3])
	}

	var expiry time.Time
	unitLabel := ""
	now := time.Now().UTC()
	
	if unit == "s" {
		expiry = now.Add(time.Duration(duration) * time.Second)
		unitLabel = "Detik"
	} else if unit == "d" {
		expiry = now.AddDate(0, 0, duration)
		unitLabel = "Hari"
	} else {
		expiry = now.AddDate(duration, 0, 0)
		unitLabel = "Tahun"
	}
	
	// Create payload: HWID|EXPIRY|SIGNATURE
	payload := fmt.Sprintf("%s|%s|SIGNATURE", hwid, expiry.Format(time.RFC3339))
	
	encrypted, err := encrypt([]byte(payload), licenseKeySalt)
	if err != nil {
		fmt.Printf("Error encrypting: %v\n", err)
		return
	}

	serialKey := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(encrypted)

	fmt.Println("\n========================================")
	fmt.Println("       WACAST LICENSE GENERATOR         ")
	fmt.Println("========================================")
	fmt.Println("HWID:        ", hwid)
	fmt.Println("Expiry:      ", expiry.Format("02 Jan 2006"))
	fmt.Println("Masa Aktif:  ", duration, unitLabel)
	fmt.Println("----------------------------------------")
	fmt.Println("SERIAL KEY:")
	fmt.Println(serialKey)
	fmt.Println("========================================\n")
}

func encrypt(data []byte, key string) ([]byte, error) {
	hash := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(data))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, data)

	return append(iv, ciphertext...), nil
}
