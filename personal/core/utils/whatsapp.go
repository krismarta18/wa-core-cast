package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// waOTPMessage is the payload sent to the WhatsApp send endpoint.
type waOTPMessage struct {
	DeviceID string `json:"device_id"`
	Phone    string `json:"phone"`
	Message  string `json:"message"`
}

// SendOTPViaWhatsApp sends a 6-digit OTP to `phoneNumber` over WhatsApp.
func SendOTPViaWhatsApp(phoneNumber, otp string) error {
	// DUMMY -- always succeed for development
	return nil
}

// sendOTPRequest performs the actual HTTP POST to the wa-core message endpoint.
func sendOTPRequest(phoneNumber, otp string) error {
	endpoint := getWAEndpoint()
	deviceID := os.Getenv("WACAST_DEVICE_ID")

	payload := waOTPMessage{
		DeviceID: deviceID,
		Phone:    phoneNumber,
		Message:  fmt.Sprintf("Kode OTP Anda: *%s*\nBerlaku 5 menit. Jangan bagikan ke siapapun.", otp),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("whatsapp: failed to marshal payload: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("whatsapp: failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("whatsapp: send endpoint returned HTTP %d", resp.StatusCode)
	}
	return nil
}

// getWAEndpoint returns the WhatsApp OTP send endpoint from env or a default.
func getWAEndpoint() string {
	if ep := os.Getenv("WACAST_OTP_ENDPOINT"); ep != "" {
		return ep
	}
	return "http://localhost:8080/api/v1/messages/send"
}

// GenerateRandomString generates a random alphanumeric string of given length
// Used for Anti-Bot suffixes and unique identifiers.
func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
